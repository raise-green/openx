package main

import (
	"fmt"
	"log"
	"strings"

	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

// inputorig.go contains all the relevant emulator commands for the originator
func LoopOrig() error {
	// This loop is exclusive to an originator
	promptColor := color.New(color.FgHiYellow).SprintFunc()
	whiteColor := color.New(color.FgHiWhite).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      promptColor("emulator") + whiteColor("# "),
		HistoryFile: consts.TellerHomeDir + "/history.txt",
		// AutoComplete: lc.NewAutoCompleter(),
	})

	ColorOutput("YOUR SEED IS: "+LocalSeed, RedColor)

	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			log.Println(err)
			return err
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		rl.SaveHistory(msg)

		cmdslice := strings.Fields(msg)
		ColorOutput("entered command: "+msg, YellowColor)

		err = ParseInputOrig(cmdslice)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func ParseInputOrig(input []string) error {
	var err error
	// Various command supported for the recipient
	if len(input) == 0 {
		// shouldn't happen, still
		return fmt.Errorf("Length of input array is zero, quitting!")
	}
	// input is greater than length 1 which means we can parse according to the command given
	command := input[0]
	switch command {
	case "help":
		fmt.Println("LIST OF SUPPORTED COMMANDS: ")
		fmt.Println("ping, display, exchange, ipfs, create, send, receive, propose, " +
			"preoriginate, myproposed, addcollateral, myoriginated, mypreoriginated")
	case "ping":
		pingHelper()
	case "display":
		displayHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, "originator")
	case "exchange":
		exchangeHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, LocalSeed)
	case "ipfs":
		ipfsHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
	case "send":
		sendHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
	case "receive":
		receiveHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
	case "create":
		createHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, LocalOriginator.U.PublicKey)
	case "kyc":
		kycHelper(input, LocalOriginator.U.Username, LocalOriginator.U.Pwhash, LocalOriginator.U.Inspector)
	// Originator only functions
	case "propose":
		fmt.Println("Proposing a contract can be done only through the opensolar webui" +
			"since that involves document verification")
		break
	case "preoriginate":
		fmt.Println("Pre originating a contract can be done only through the opensolar webui" +
			"since that involves document verification")
		break
		// end of propose
	case "addcollateral":
		if len(input) != 3 {
			log.Println("<addcollateral> collateral amount")
			break
		}

		_, err = utils.StoFWithCheck(input[2])
		if err != nil {
			log.Println(err)
			break
		}

		collateral := input[1]
		amount := input[2]

		response, err := AddCollateral(LocalOriginator.U.Username, LocalOriginator.U.Pwhash, collateral, amount)
		if err != nil {
			log.Println(err)
			break
		}

		if response.Code == 200 {
			ColorOutput("SUCCESSFULLY ADDED COLLATERAL", GreenColor)
		} else {
			ColorOutput("RESPONSE STATUS: "+utils.ItoS(response.Code), GreenColor)
		}
		break
		// end of addcollateral
	case "myproposed":
		x, err := GetProposedContracts(LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
		// end of myproposed
	case "mypreoriginated":
		x, err := GetPreOriginatedContracts(LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
		// end of myoriginated
	case "myoriginated":
		x, err := GetOriginatedContracts(LocalOriginator.U.Username, LocalOriginator.U.Pwhash)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(x)
		break
		// end of myoriginated
	}
	return nil
}
