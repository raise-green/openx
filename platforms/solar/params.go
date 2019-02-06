package solar

// TODO: add more parameters here that would help identify a given solar project
type SolarParams struct {
	Index int // an Index to keep quick track of how many projects exist

	PanelSize    string  // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue   float64 // the total money that we need from investors
	Location     string  // where this specific solar panel is located
	MoneyRaised  float64 // total money that has been raised until now
	Years        int     // number of years the recipient is expected to repay the initial investment amount by
	InterestRate float64 // the interest rate provided to potential investors
	Metadata     string  // any other metadata can be stored here

	// once all funds have been raised, we need to set assetCodes
	InvestorAssetCode string
	DebtAssetCode     string
	PaybackAssetCode  string
	SeedAssetCode     string

	BalLeft float64 // denotes the balance left to pay by the party
	Votes   int     // the number of votes towards a proposed contract by investors

	DateInitiated string // date the project was created
	DateFunded    string // date when the project was funded
	DateLastPaid  int64  // int64 since we need comparisons on this one
	// Percentage raised is not stored in the database since that can be calculated by the UI

	// extra params based on data available
	Inverter        string
	ChargeRegulator string
	ControlPanel    string
	CommBox         string
	ACTransfer      string
	SolarCombiner   string
	Batteries       string
	IoTHub          string
}

// these are the reputation values associated with a specific project. For eg if
// a project's total worth is 10000 and everything in the project goes well and
// all entities are satisfied by the outcome, the originator gains 1000 points,
// the contractor gains 3000 points and so on
// MWTODO: get comments on the weights and tweak them if needed
var (
	InvestorWeight   = 0.1
	OriginatorWeight = 0.1
	ContractorWeight = 0.3
	DeveloperWeight  = 0.2
	RecipientWeight  = 0.3
)
