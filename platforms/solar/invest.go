package solar

import (
	"fmt"
	"log"
	"time"

	assets "github.com/OpenFinancing/openfinancing/assets"
	consts "github.com/OpenFinancing/openfinancing/consts"
	database "github.com/OpenFinancing/openfinancing/database"
	issuer "github.com/OpenFinancing/openfinancing/issuer"
	notif "github.com/OpenFinancing/openfinancing/notif"
	stablecoin "github.com/OpenFinancing/openfinancing/stablecoin"
	utils "github.com/OpenFinancing/openfinancing/utils"
	wallet "github.com/OpenFinancing/openfinancing/wallet"
	xlm "github.com/OpenFinancing/openfinancing/xlm"
	"github.com/stellar/go/build"
)

func RetrieveValues(projIndex int, invIndex int, recpIndex int) (Project, database.Investor, database.Recipient, error) {
	var project Project
	var investor database.Investor
	var recipient database.Recipient
	var err error

	project, err = RetrieveProject(projIndex)
	if err != nil {
		return project, investor, recipient, err
	}

	investor, err = database.RetrieveInvestor(invIndex)
	if err != nil {
		return project, investor, recipient, err
	}

	recipient, err = database.RetrieveRecipient(recpIndex)
	if err != nil {
		return project, investor, recipient, err
	}
	return project, investor, recipient, nil
}

func SendUSDToPlatform(platformSeed string, invSeed string, invAmount string, projIndex int) (string, error) {
	// send stableusd to the platform (not the issuer) since the issuer will be locked
	// and we can't use the funds
	// if we are the issuer, we can burn the stableUSD. but if we are not, then we need
	// to redeem this stableUSD for fiat and hence we need the asset to be usable.
	platformPubkey, err := wallet.ReturnPubkey(platformSeed)
	if err != nil {
		return "", err
	}

	invPubkey, err := wallet.ReturnPubkey(invSeed)
	if err != nil {
		return "", err
	}

	oldPlatformBalance, err := xlm.GetAssetBalance(platformPubkey, stablecoin.Code)
	if err != nil {
		return "", err
	}

	_, txhash, err := assets.SendAsset(stablecoin.Code, stablecoin.PublicKey, platformPubkey, invAmount, invSeed, invPubkey, "Opensolar investment: "+utils.ItoS(projIndex))
	if err != nil {
		log.Println("Sending stableusd to platform failed", platformPubkey, invAmount, invSeed, invPubkey)
		return txhash, err
	}

	log.Println("Sent STABLEUSD to platform, confirmation: ", txhash)
	time.Sleep(5 * time.Second) // wait for a block

	newPlatformBalance, err := xlm.GetAssetBalance(platformPubkey, stablecoin.Code)
	if err != nil {
		return txhash, err
	}

	if utils.StoF(newPlatformBalance)-utils.StoF(oldPlatformBalance) < utils.StoF(invAmount)-1 {
		return txhash, fmt.Errorf("Sent amount doesn't match with investment amount")
	}
	return txhash, nil
}

// this file does not contain any tests associated with it right now. In the future,
// once we have a robust frontend, we can modify the CLI interface to act as a test
// for this file

// InvestInProject invests in a particular solar project given required parameters
func InvestInProject(projIndex int, invIndex int, recpIndex int, invAmount string,
	invSeed string, recpSeed string, platformSeed string) (Project, error) {
	var err error

	project, investor, recipient, err := RetrieveValues(projIndex, invIndex, recpIndex)
	if err != nil {
		return project, err
	}

	if !investor.CanInvest(investor.U.PublicKey, invAmount) {
		log.Println("Investor has less balance than what is required to ivnest in this asset")
		return project, err
	}

	// check if investment amount is greater than or equal to the project requirements
	if utils.StoF(invAmount) > project.Params.TotalValue-project.Params.MoneyRaised {
		return project, fmt.Errorf("User is trying to invest more than what is needed")
	}

	var InvestorAsset build.Asset
	var PaybackAsset build.Asset
	var DebtAsset build.Asset
	// user has decided to invest in a part of the project (don't know if full yet)
	// so if there has been no asset codes assigned yet, we need to create them and
	// assign them here
	// you can retrieve these anywhere since the metadata will most likely be unique
	if project.Params.InvestorAssetCode == "" {
		// this person is the first investor, set the investor asset name and create the
		// issuer that will be created for this particular project
		project.Params.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Params.Metadata) // set the investor asset code
		err = issuer.InitIssuer(project.Params.Index, consts.IssuerSeedPwd)
		if err != nil {
			log.Fatal(err)
		}
		err = issuer.FundIssuer(project.Params.Index, consts.IssuerSeedPwd, platformSeed)
		if err != nil {
			log.Fatal(err)
		}
	}

	stableTxHash, err := SendUSDToPlatform(platformSeed, invSeed, invAmount, project.Params.Index)
	if err != nil {
		return project, err
	}

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(project.Params.Index), consts.IssuerSeedPwd)
	if err != nil {
		return project, err
	}

	// InvAsset is not a native asset, so don't set the native flag
	InvestorAsset = assets.CreateAsset(project.Params.InvestorAssetCode, issuerPubkey)
	// make investor trust the asset, trustlimit is upto the value of the project
	invTrustTxHash, err := assets.TrustAsset(InvestorAsset.Code, issuerPubkey, utils.FtoS(project.Params.TotalValue), investor.U.PublicKey, invSeed)
	if err != nil {
		return project, err
	}

	log.Println("Investor trusted asset: ", InvestorAsset.Code, " tx hash: ", invTrustTxHash)
	_, invAssetTxHash, err := assets.SendAssetFromIssuer(InvestorAsset.Code, investor.U.PublicKey, invAmount, issuerSeed, issuerPubkey)
	if err != nil {
		return project, err
	}

	log.Printf("Sent InvAsset %s to investor %s with txhash %s", InvestorAsset.Code, investor.U.PublicKey, invAssetTxHash)
	// investor asset sent, update project.Params's BalLeft
	fmt.Println("Updating investor to handle invested amounts and assets")
	project.Params.MoneyRaised += utils.StoF(invAmount)
	investor.AmountInvested += utils.StoF(invAmount)
	investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, InvestorAsset.Code)
	// keep note of who all invested in this asset (even though it should be easy
	// to get that from the blockchain)
	err = investor.Save() // save investor creds now that we're done
	if err != nil {
		return project, err
	}
	fmt.Println("Updated investor database")
	// append the investor class to the list of project investors
	// if the same investor has invested twice, he will appear twice
	// can be resolved on the UI side by requiring unique, so not doing that here
	project.ProjectInvestors = append(project.ProjectInvestors, investor)
	if project.Params.MoneyRaised == project.Params.TotalValue {
		// this project covers up the amount nedeed for the project, so set the DebtAssetCode
		// and PaybackAssetCodes, generate them and give to the recipient
		project.Params.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + project.Params.Metadata)
		project.Params.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + project.Params.Metadata)

		DebtAsset = assets.CreateAsset(project.Params.DebtAssetCode, issuerPubkey)
		PaybackAsset = assets.CreateAsset(project.Params.PaybackAssetCode, issuerPubkey)

		pbAmtTrust := utils.ItoS(project.Params.Years * 12 * 2) // two way exchange possible, to account for errors

		recpPbTrustHash, err := assets.TrustAsset(PaybackAsset.Code, issuerPubkey, pbAmtTrust, recipient.U.PublicKey, recpSeed)
		if err != nil {
			return project, err
		}

		log.Println("Recipient Trusts Debt asset: ", DebtAsset.Code, " tx hash: ", recpPbTrustHash)
		_, recpAssetHash, err := assets.SendAssetFromIssuer(PaybackAsset.Code, recipient.U.PublicKey, pbAmtTrust, issuerSeed, issuerPubkey) // same amount as debt
		if err != nil {
			return project, err
		}

		log.Printf("Sent PaybackAsset to recipient %s with txhash %s", recipient.U.PublicKey, recpAssetHash)
		recpDebtTrustHash, err := assets.TrustAsset(DebtAsset.Code, issuerPubkey, utils.FtoS(project.Params.TotalValue*2), recipient.U.PublicKey, recpSeed)
		if err != nil {
			return project, err
		}

		log.Println("Recipient Trusts Payback asset: ", PaybackAsset.Code, " tx hash: ", recpDebtTrustHash)
		_, recpDebtAssetHash, err := assets.SendAssetFromIssuer(DebtAsset.Code, recipient.U.PublicKey, utils.FtoS(project.Params.TotalValue), issuerSeed, issuerPubkey) // same amount as debt
		if err != nil {
			return project, err
		}

		log.Printf("Sent PaybackAsset to recipient %s with txhash %s\n", recipient.U.PublicKey, recpDebtAssetHash)
		project.Params.BalLeft = float64(project.Params.TotalValue)
		recipient.ReceivedSolarProjects = append(recipient.ReceivedSolarProjects, DebtAsset.Code)
		project.ProjectRecipient = recipient // need to udpate project.Params each time recipient is mutated
		// only here does the recipient part change, so update it only here
		err = recipient.Save()
		if err != nil {
			return project, err
		}

		project.Stage = FundedProject // set funded project stage
		err = project.Save()
		if err != nil {
			log.Println("Couldn't insert project")
			return project, err
		}

		fmt.Println("Updated recipient bucket")
		txhash, err := issuer.FreezeIssuer(project.Params.Index, "blah")
		if err != nil {
			return project, err
		}

		log.Printf("Tx hash for freezing issuer is: %s", txhash)
		if recipient.U.Notification {
			notif.SendInvestmentNotifToRecipient(projIndex, recipient.U.Email, recpPbTrustHash, recpAssetHash, recpDebtTrustHash, recpDebtAssetHash)
		}
	}
	// update the project finally now that we have updated other databases
	err = project.Save()
	// send notification emails out
	if investor.U.Notification {
		notif.SendInvestmentNotifToInvestor(projIndex, investor.U.Email, stableTxHash, invTrustTxHash, invAssetTxHash)
	}
	return project, err
}
