package munibond

import (
	"log"
	"time"

	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	models "github.com/YaleOpenLab/openx/models"
	notif "github.com/YaleOpenLab/openx/notif"
	oracle "github.com/YaleOpenLab/openx/oracle"
	stablecoin "github.com/YaleOpenLab/openx/stablecoin"
	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
	assets "github.com/YaleOpenLab/openx/xlm/assets"
	issuer "github.com/YaleOpenLab/openx/xlm/issuer"
	wallet "github.com/YaleOpenLab/openx/xlm/wallet"
	"github.com/pkg/errors"
)

// MunibondInvest invests in a specific munibond
func MunibondInvest(issuerPath string, invIndex int, invSeed string, invAmount string,
	projIndex int, invAssetCode string, totalValue float64, seedInvestmentFactor float64) error {
	// offer user to exchange xlm for stableusd and invest directly if the user does not have stableusd
	// this should be a menu on the Frontend but here we do this automatically
	var err error

	investor, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve investor from database")
	}

	err = stablecoin.OfferExchange(investor.U.StellarWallet.PublicKey, invSeed, invAmount)
	if err != nil {
		return errors.Wrap(err, "Unable to offer xlm to STABLEUSD excahnge for investor")
	}

	stableTxHash, err := SendUSDToPlatform(invSeed, invAmount, "Opensolar investment: "+utils.ItoS(projIndex))
	if err != nil {
		return errors.Wrap(err, "Unable to send STABLEUSD to platform")
	}

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve seed")
	}

	InvestorAsset := assets.CreateAsset(invAssetCode, issuerPubkey)
	invTrustTxHash, err := assets.TrustAsset(InvestorAsset.Code, issuerPubkey, utils.FtoS(totalValue), invSeed)
	if err != nil {
		return errors.Wrap(err, "Error while trusting investor asset")
	}

	log.Printf("Investor trusts InvAsset %s with txhash %s", InvestorAsset.Code, invTrustTxHash)
	_, invAssetTxHash, err := assets.SendAssetFromIssuer(InvestorAsset.Code, investor.U.StellarWallet.PublicKey, invAmount, issuerSeed, issuerPubkey)
	if err != nil {
		return errors.Wrap(err, "Error while sending out investor asset")
	}

	log.Printf("Sent InvAsset %s to investor %s with txhash %s", InvestorAsset.Code, investor.U.StellarWallet.PublicKey, invAssetTxHash)

	investor.AmountInvested += utils.StoF(invAmount) //  / seedInvestmentFactor -> figure out after demo
	investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, InvestorAsset.Code)
	investor.InvestedSolarProjectsIndices = append(investor.InvestedSolarProjectsIndices, projIndex)
	// keep note of who all invested in this asset (even though it should be easy
	// to get that from the blockchain)
	err = investor.Save()
	if err != nil {
		return err
	}

	if investor.U.Notification {
		notif.SendInvestmentNotifToInvestor(projIndex, investor.U.Email, stableTxHash, invTrustTxHash, invAssetTxHash)
	}
	return nil
}

// MunibondReceive sends assets to the recipient
func MunibondReceive(issuerPath string, recpIndex int, projIndex int, debtAssetId string,
	paybackAssetId string, years int, recpSeed string, totalValue float64, paybackPeriod int) error {

	log.Println("Retrieving recipient")
	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve recipient from database")
	}

	log.Println("Retrieving issuer")
	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.CreatePath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve issuer seed")
	}

	DebtAsset := assets.CreateAsset(debtAssetId, issuerPubkey)
	PaybackAsset := assets.CreateAsset(paybackAssetId, issuerPubkey)

	if years == 0 {
		years = 1
	}
	pbAmtTrust := utils.ItoS(years * 12 * 2) // two way exchange possible, to account for errors

	paybackTrustHash, err := assets.TrustAsset(PaybackAsset.Code, issuerPubkey, pbAmtTrust, recpSeed)
	if err != nil {
		return errors.Wrap(err, "Error while trusting Payback Asset")
	}
	log.Printf("Recipient Trusts Payback asset %s with txhash %s", PaybackAsset.Code, paybackTrustHash)

	_, paybackAssetHash, err := assets.SendAssetFromIssuer(PaybackAsset.Code, recipient.U.StellarWallet.PublicKey, pbAmtTrust, issuerSeed, issuerPubkey) // same amount as debt
	if err != nil {
		return errors.Wrap(err, "Error while sending payback asset from issue")
	}

	log.Printf("Sent PaybackAsset to recipient %s with txhash %s", recipient.U.StellarWallet.PublicKey, paybackAssetHash)
	debtTrustHash, err := assets.TrustAsset(DebtAsset.Code, issuerPubkey, utils.FtoS(totalValue*2), recpSeed)
	if err != nil {
		return errors.Wrap(err, "Error while trusting debt asset")
	}
	log.Printf("Recipient Trusts Debt asset %s with txhash %s", DebtAsset.Code, debtTrustHash)

	_, recpDebtAssetHash, err := assets.SendAssetFromIssuer(DebtAsset.Code, recipient.U.StellarWallet.PublicKey, utils.FtoS(totalValue), issuerSeed, issuerPubkey) // same amount as debt
	if err != nil {
		return errors.Wrap(err, "Error while sending debt asset")
	}

	log.Printf("Sent DebtAsset to recipient %s with txhash %s\n", recipient.U.StellarWallet.PublicKey, recpDebtAssetHash)
	recipient.ReceivedSolarProjects = append(recipient.ReceivedSolarProjects, DebtAsset.Code)
	recipient.ReceivedSolarProjectIndices = append(recipient.ReceivedSolarProjectIndices, projIndex)
	err = recipient.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save recipient")
	}

	txhash, err := issuer.FreezeIssuer(issuerPath, projIndex, "blah")
	if err != nil {
		return errors.Wrap(err, "Error while freezing issuer")
	}

	log.Printf("Tx hash for freezing issuer is: %s", txhash)
	log.Printf("PROJECT %d's INVESTMENT CONFIRMED!", projIndex)

	if recipient.U.Notification {
		notif.SendInvestmentNotifToRecipient(projIndex, recipient.U.Email, paybackTrustHash, paybackAssetHash, debtTrustHash, recpDebtAssetHash)
	}

	go sendPaymentNotif(recipient.U.Index, projIndex, paybackPeriod, recipient.U.Email)
	return nil
}

// sendPaymentNotif sends a notification every payback period to the recipient to
// kindly remind him to payback towards the project
func sendPaymentNotif(recpIndex int, projIndex int, paybackPeriod int, email string) {
	// setup a payback monitoring routine for monitoring if the recipient pays us back on time
	// the recipient must give his email to receive updates
	paybackTimes := 0
	for {

		_, err := database.RetrieveRecipient(recpIndex) // need to retrieve to make sure nothing goes awry
		if err != nil {
			log.Println("Error while retrieving recipient from database", err)
			message := "Error while retrieving your account details, please contact help as soon as you receive this message " + err.Error()
			notif.SendAlertEmail(message, email) // don't catch the error here
			time.Sleep(time.Second * 2 * 604800)
		}

		if paybackTimes == 0 {
			// sleep and bother during the next cycle
			time.Sleep(time.Second * 2 * 604800)
		}

		// PAYBACK TIME!!
		// we don't know if the user has paid, but we send an email anyway
		notif.SendPaybackAlertEmail(projIndex, email)
		// sleep until the next payment is due
		paybackTimes += 1
		log.Println("Sent: ", email, "a notification on payments for payment cycle: ", paybackTimes)
		time.Sleep(2 * time.Duration(paybackPeriod) * time.Second)
	}
}

// MunibondPayback is used by the recipient to pay the platform back. Here, we pay the
// project escrow instead of the platform since it would be responsible for redistribution of funds
func MunibondPayback(issuerPath string, recpIndex int, amount string, recipientSeed string, projIndex int,
	assetName string, projectInvestors []int, totalValue float64, escrowPubkey string) (float64, error) {

	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return -1, errors.Wrap(err, "Error while retrieving recipient from database")
	}

	issuerPubkey, _, err := wallet.RetrieveSeed(issuer.CreatePath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to retrieve issuer seed")
	}

	monthlyBill := oracle.MonthlyBill()
	if err != nil {
		return -1, errors.Wrap(err, "Unable to fetch oracle price, exiting")
	}

	log.Println("Retrieved average price from oracle: ", monthlyBill)
	mBillFloat := utils.StoF(monthlyBill)

	if utils.StoF(amount) < mBillFloat {
		return -1, errors.New("amount paid is less than amount needed. Please refill your main account")
	}

	err = stablecoin.OfferExchange(recipient.U.StellarWallet.PublicKey, recipientSeed, amount)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to offer xlm to STABLEUSD exchange for investor")
	}

	StableBalance, err := xlm.GetAssetBalance(recipient.U.StellarWallet.PublicKey, "STABLEUSD")
	if err != nil || (utils.StoF(StableBalance) < utils.StoF(amount)) {
		return -1, errors.Wrap(err, "You do not have the required stablecoin balance, please refill")
	}

	if projIndex == 4 {
		escrowPubkey = "GB2MLKOMRHXJQRX3ZCBLAZ2PMJOYHMFXQUR525CNCSPKBE643FRQNZBM"
	}
	log.Println("ESCROW PUBKEY: ", escrowPubkey)
	_, stableUSDHash, err := assets.SendAsset(consts.StablecoinCode, consts.StableCoinAddress, escrowPubkey, amount, recipientSeed, "Opensolar payback: "+utils.ItoS(projIndex))
	if err != nil {
		log.Println("ESCROW PUBKEY: ", escrowPubkey)
		return -1, errors.Wrap(err, "Error while sending STABLEUSD back")
	}
	log.Printf("Paid %s back to platform in stableUSD, txhash %s ", amount, stableUSDHash)

	_, debtPaybackHash, err := assets.SendAssetToIssuer(assetName, issuerPubkey, amount, recipientSeed)
	if err != nil {
		return -1, errors.Wrap(err, "Error while sending debt asset back")
	}
	log.Printf("Paid %s back to platform in DebtAsset, txhash %s ", amount, debtPaybackHash)

	ownershipAmt := utils.StoF(amount) - mBillFloat
	ownershipPct := ownershipAmt / totalValue
	if recipient.U.Notification {
		notif.SendPaybackNotifToRecipient(projIndex, recipient.U.Email, stableUSDHash, debtPaybackHash)
	}

	for _, i := range projectInvestors {
		investor, err := database.RetrieveInvestor(i)
		if err != nil {
			log.Println("Error while retrieving investor from list of investors", err)
			continue
		}
		if investor.U.Notification {
			notif.SendPaybackNotifToInvestor(projIndex, investor.U.Email, stableUSDHash, debtPaybackHash)
		}
	}

	return ownershipPct, nil
}

// SendUSDToPlatform is used to send usd back to the platform
func SendUSDToPlatform(invSeed string, invAmount string, memo string) (string, error) {
	return models.SendUSDToPlatform(invSeed, invAmount, memo)
}
