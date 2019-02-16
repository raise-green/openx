package main

import (
	"fmt"
	"log"

	utils "github.com/YaleOpenLab/openx/utils"
)

func displayHelper(input []string, username string, pwhash string, role string) {
	// display is a  broad command and needs to have a subcommand
	if len(input) == 1 {
		// only display was given, so display help command
		log.Println("<display><balance, profile, projects>")
		return
	}
	subcommand := input[1]
	switch subcommand {
	case "balance":
		if len(input) == 2 {
			log.Println("Calling balances API")
			balances, err := GetBalances(username, pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			PrintBalances(balances)
			break
		}
		subcommand := input[2]
		switch subcommand {
		case "xlm":
			// print xlm balance
			balance, err := GetXLMBalance(username, pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			ColorOutput("BALANCE: "+balance, MagentaColor)
			break
		case "all":
			balances, err := GetBalances(username, pwhash)
			if err != nil {
				log.Println(err)
				break
			}
			PrintBalances(balances)
			break
		default:
			balance, err := GetAssetBalance(username, pwhash, subcommand)
			if err != nil {
				log.Println(err)
				break
			}
			ColorOutput("BALANCE: "+balance, MagentaColor)
			break
			// print asset balance
		}
	case "profile":
		log.Println("Displaying Profile")
		switch role {
		case "investor":
			PrintInvestor(LocalInvestor)
		case "recipient":
			PrintRecipient(LocalRecipient)
		case "contractor":
			PrintEntity(LocalContractor)
		case "originator":
			PrintEntity(LocalOriginator)
		}
		break
	case "projects":
		if len(input) != 4 {
			// only display was given, so display help command
			log.Println("display projects <platform> <preorigin, origin, seed, proposed, open, funded, installed, power, fin>")
			break
		}
		platform := input[2]
		switch platform {
		case "opzones":
			log.Println("OPZONES PLATFORM")
			subsubcommand := input[3]
			switch subsubcommand {
			case "cbonds":
				log.Println("PRINTGING ALL OPEN Construction Bonds")
			case "lucoops":
				log.Println("PRINTGING ALL OPEN Living unit coops")
			}
			break
		case "opensolar":
			subsubcommand := input[3]
			var stage float64
			switch subsubcommand {
			case "preorigin":
				log.Println("Displaying all pre-originated (stage 0) projects")
				stage = 0
				break
			case "origin":
				log.Println("Displaying all originated (stage 1) projects")
				stage = 1
				break
			case "seed":
				log.Println("Displaying all seed (stage 1.5) projects")
				stage = 1.5
				break
			case "proposed":
				log.Println("Displaying all proposed (stage 2) projects")
				stage = 2
				break
			case "open":
				log.Println("Displaying open (stage 3) projects")
				stage = 3
				break
			case "funded":
				log.Println("Displaying funded (stage 4) projects")
				stage = 4
				break
			case "installed":
				log.Println("Displaying installed (stage 5) projects")
				stage = 5
				break
			case "power":
				log.Println("Displaying funded (stage 6) projects")
				stage = 6
				break
			case "fin":
				log.Println("Displaying funded (stage 7) projects")
				stage = 7
				break
			}
			arr, err := RetrieveProject(stage)
			if err != nil {
				log.Println(err)
				break
			}
			PrintProjects(arr)
			break
		}
	} // end of display
}

func exchangeHelper(input []string, username string, pwhash string, seed string) {
	if len(input) != 2 {
		// only display was given, so display help command
		log.Println("<exchange> amount")
		return
	}
	amount, err := utils.StoICheck(input[1])
	if err != nil {
		log.Println(err)
		return
	}
	// convert this to int and check if int
	fmt.Println("Exchanging", amount, "XLM for STABLEUSD")
	response, err := GetStableCoin(username, pwhash, seed, input[1])
	if err != nil {
		log.Println(err)
		return
	}
	if response.Code == 200 {
		ColorOutput("SUCCESSFUL, CHECK BALANCES", GreenColor)
	} else {
		ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Code), GreenColor)
	}
}

func ipfsHelper(input []string, username string, pwhash string) {
	if len(input) != 2 {
		log.Println("<ipfs> string")
		return
	}
	inputString := input[1]
	hashString, err := GetIpfsHash(username, pwhash, inputString)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("IPFS HASH", hashString)
	// end of ipfs
}

func pingHelper() {
	err := PingRpc()
	if err != nil {
		log.Println(err)
		return
	}
}

func sendHelper(input []string, username string, pwhash string) {
	var err error
	if len(input) == 1 {
		log.Println("send <asset>")
		return
	}
	subcommand := input[1]
	switch subcommand {
	case "asset":
		if len(input) != 5 {
			log.Println("send asset <assetName> <destination> <amount>")
			return
		}

		assetName := input[2]
		destination := input[3]
		amount := input[4]

		txhash, err := SendLocalAsset(username, pwhash,
			LocalSeedPwd, assetName, destination, amount)
		if err != nil {
			log.Println(err)
		}
		ColorOutput("TX HASH: "+txhash, MagentaColor)
		// end of asset
	case "xlm":
		if len(input) < 4 {
			log.Println("send xlm <destination> <amount> <<memo>>")
			break
		}
		destination := input[2]
		_, err = utils.StoFWithCheck(input[3])
		if err != nil {
			log.Println(err)
			break
		}
		// send xlm overs
		amount := input[3]
		var memo string
		if len(input) > 4 {
			memo = input[4]
		}
		txhash, err := SendXLM(username, pwhash, LocalSeedPwd, destination, amount, memo)
		if err != nil {
			log.Println(err)
		}
		ColorOutput("TX HASH: "+txhash, MagentaColor)
	}
}

func receiveHelper(input []string, username string, pwhash string) {
	// we can either receive from the faucet or trust issuers to receive assets
	var err error
	if len(input) == 1 {
		log.Println("receive <xlm, asset>")
		return
	}
	subcommand := input[1]
	switch subcommand {
	case "xlm":
		status, err := AskXLM(username, pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("COIN REQUEST SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("COIN REQUEST NOT SUCCESSFUL", RedColor)
		}
		// ask for coins from the faucet
	case "asset":
		if len(input) != 5 {
			log.Println("receive asset <assetName> <issuerPubkey> <limit>")
			break
		}

		assetName := input[2]
		issuerPubkey := input[3]
		_, err = utils.StoFWithCheck(input[4])
		if err != nil {
			log.Println(err)
			break
		}

		limit := input[4]

		status, err := TrustAsset(username, pwhash, assetName, issuerPubkey, limit, LocalSeedPwd)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("COIN REQUEST SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("COIN REQUEST NOT SUCCESSFUL", RedColor)
		}
		break
	} // end of receive
}

func createHelper(input []string, username string, pwhash string, pubkey string) {
	// create enables you to create tokens on stellar that you can excahnge with third parties.
	if len(input) == 1 {
		log.Println("create <asset>")
		return
	}
	subcommand := input[1]
	switch subcommand {
	case "asset":
		// create a new asset
		if len(input) != 3 {
			log.Println("create asset <name>")
			break
		}
		assetName := input[2]
		status, err := CreateAssetInv(username, pwhash, assetName, pubkey)
		if err != nil {
			log.Println(err)
			return
		}
		if status.Code == 200 {
			ColorOutput("INVESTMENT SUCCESSFUL, CHECK EMAIL", GreenColor)
		} else {
			ColorOutput("INVESTMENT NOT SUCCESSFUL", RedColor)
		}
	} // end of create
}

func kycHelper(input []string, username string, pwhash string, inspector bool) {
	var err error
	if !inspector {
		ColorOutput("YOU ARE NOT A KYC INSPECTOR", RedColor)
		return
	}
	if len(input) == 1 {
		log.Println("kyc <auth, view>")
	}
	subcommand := input[1]
	switch subcommand {
	case "auth":
		if len(input) != 3 {
			log.Println("kyc auth <userIndex>")
			break
		}
		_, err = utils.StoICheck(input[1])
		if err != nil {
			log.Println(err)
			break
		}
		status, err := AuthKyc(input[1], username, pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		if status.Code == 200 {
			ColorOutput("USER KYC'D!", GreenColor)
		} else {
			ColorOutput("USER NOT KYC'D", RedColor)
		}
		break
	case "notdone":
		users, err := NotKycView(username, pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		PrintUsers(users)
		// print all the users who have kyc'd
		break
	case "done":
		users, err := KycView(username, pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		PrintUsers(users)
		// print all the users who have kyc'd
		break
	}
	// end of kyc
}
