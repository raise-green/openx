package opensolar

import (
	"fmt"
)

// Contract auctions are specific for public infrastructure and for projects with multiple stakeholders
// where transparency is important. The main point is to avoid corruption and produce positive competition
// among providers. When funding solar project using international investor money, transparency and auditability
// to them is also a crucial aspect since they are often not close to the project, and want to make sure there
// is a tobust due diligence prior to unlocking the funds.

// Different Auctions or Tenders are designed based on the nature of the project.
// In general, the criteria for selection is price, technical quality (eg. hardware), engineering model, development time
// and other perks offered by developers (eg. extra guarantees).

// you need to have a lock in period beyond which contractors can not post what
// stuff they want. now, how do you choose which contractor wins? Ideally,
// the school would want the most stuff but you need to vet which contracts are good
// and not.

// For more trivia, see https://en.wikipedia.org/wiki/Auction
type ContractAuction struct {
	// TODO: this struct isn't used yet as it needs handlers and stuff, but when
	// we move off main.go for testing, this must be used in order to make stuff
	// easier for us.
	AllContracts    []Project
	AllContractors  []Entity
	WinningContract Project
}

// SelectContractBlind selects the winning bid  based on blind auctio nrules
// in a blind auction, the bid with the highest price wins
func SelectContractBlind(arr []Project) (Project, error) {
	var a Project
	if len(arr) == 0 {
		return a, fmt.Errorf("Empty array passed!")
	}
	// array is not empty, min 1 elem
	a = arr[0]
	for _, elem := range arr {
		if elem.TotalValue < a.TotalValue {
			a = elem
			continue
		}
	}
	return a, nil
}

// SelectContractVickrey selects the winning bid based on vickrey auction rules
// in a vickrey auction, the bid with the second highest price wins
func SelectContractVickrey(arr []Project) (Project, error) {
	var winningContract Project
	if len(arr) == 0 {
		return winningContract, fmt.Errorf("Empty array passed!")
	}
	// array is not empty, min 1 elem
	winningContract = arr[0]
	var pos int
	for i, elem := range arr {
		if elem.TotalValue < winningContract.TotalValue {
			winningContract = elem
			pos = i
			continue
		}
	}
	// here we have the highest bidder. Now we need to delete this guy from the array
	// and get the second highest bidder
	// delete a[pos] from arr
	arr = append(arr[:pos], arr[pos+1:]...)
	if len(arr) == 0 {
		// means only one contract was proposed for this project, so fall back to blind auction
		return winningContract, nil
	}
	vickreyPrice := arr[0].TotalValue
	for _, elem := range arr {
		if elem.TotalValue < vickreyPrice {
			vickreyPrice = elem.TotalValue
		}
	}
	// we have the winner, who's elem and we have the price which is vickreyPrice
	// overwrite the winning contractor's contract
	winningContract.TotalValue = vickreyPrice
	return winningContract, winningContract.Save()
}

// SelectContractTime selects the winning contract based on the least time to completion
func SelectContractTime(arr []Project) (Project, error) {
	var a Project
	if len(arr) == 0 {
		return a, fmt.Errorf("Empty array passed!")
	}

	a = arr[0]
	for _, elem := range arr {
		if elem.Years < a.Years {
			a = elem
			continue
		}
	}
	return a, nil
}

// SetAuctionType sets the auction type of a specific project
func (project *Project) SetAuctionType(auctionType string) error {
	switch auctionType {
	case "blind":
		project.AuctionType = "blind"
	case "vickrey":
		project.AuctionType = "vickrey"
	case "english":
		project.AuctionType = "english"
	case "dutch":
		project.AuctionType = "dutch"
	default:
		project.AuctionType = "blind"
	}
	return project.Save()
}