package main

import (
	"github.com/pkg/errors"
	"log"

	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
	utils "github.com/YaleOpenLab/openx/utils"
)

// newLivingUnitCoop creates a new living unit coop
func newLivingUnitCoop(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, totalAmount float64, typeOfUnit string, monthlyPayment float64,
	title string, location string, description string) (opzones.LivingUnitCoop, error) {
	var coop opzones.LivingUnitCoop
	coop.MaturationDate = mdate
	coop.MemberRights = mrights
	coop.SecurityType = stype
	coop.InterestRate = intrate
	coop.Rating = rating
	coop.BondIssuer = bIssuer
	coop.Underwriter = uWriter
	coop.Title = title
	coop.Location = location
	coop.Description = description
	coop.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllLivingUnitCoops()
	if err != nil {
		return coop, errors.Wrap(err, "could not retrieve all living unit coops")
	}
	coop.Index = len(x) + 1
	coop.UnitsSold = 0
	coop.Amount = totalAmount
	coop.TypeOfUnit = typeOfUnit
	coop.MonthlyPayment = monthlyPayment
	err = coop.Save()
	return coop, err
}

// newConstructionBond returns a New Construction Bond and automatically stores it in the db
func newConstructionBond(mdate string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, unitCost float64, itype string, nUnits int, tax string, recIndex int,
	title string, location string, description string) (opzones.ConstructionBond, error) {
	var cBond opzones.ConstructionBond
	cBond.MaturationDate = mdate
	cBond.SecurityType = stype
	cBond.InterestRate = intrate
	cBond.Rating = rating
	cBond.BondIssuer = bIssuer
	cBond.Underwriter = uWriter
	cBond.Title = title
	cBond.Location = location
	cBond.Description = description
	cBond.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllConstructionBonds()
	if err != nil {
		return cBond, errors.Wrap(err, "could not retrieve all living unit coops")
	}

	cBond.Index = len(x) + 1
	cBond.CostOfUnit = unitCost
	cBond.InstrumentType = itype
	cBond.NoOfUnits = nUnits
	cBond.Tax = tax
	cBond.RecipientIndex = recIndex
	err = cBond.Save()
	return cBond, err
}

// ALL 5 PROJECT DATA WILL BE ADDED HERE FOR THE DEMO
// InsertDummyData inserts sample data
func InsertDummyData() error {
	var err error
	// populate database with dumym data
	var recp database.Recipient
	allRecs, err := database.RetrieveAllRecipients()
	if err != nil {
		log.Fatal(err)
	}
	if len(allRecs) == 0 {
		// there is no recipient right now, so create a dummy recipient
		var err error
		recp, err = database.NewRecipient("martin", "p", "x", "Martin")
		if err != nil {
			log.Fatal(err)
		}
		recp.U.Notification = true
		err = recp.AddEmail("varunramganesh@gmail.com")
		if err != nil {
			log.Fatal(err)
		}
	}

	var inv database.Investor
	allInvs, err := database.RetrieveAllInvestors()
	if err != nil {
		log.Fatal(err)
	}
	if len(allInvs) == 0 {
		var err error
		inv, err = database.NewInvestor("john", "p", "x", "John")
		if err != nil {
			log.Fatal(err)
		}
		err = inv.AddVotingBalance(100000)
		// this function saves as well, so there's no need to save again
		if err != nil {
			log.Fatal(err)
		}
		err = database.AddInspector(inv.U.Index)
		if err != nil {
			log.Fatal(err)
		}
		x, err := database.RetrieveUser(inv.U.Index)
		if err != nil {
			log.Fatal(err)
		}
		inv.U = x
		err = inv.Save()
		if err != nil {
			log.Fatal(err)
		}
		err = x.Authorize(inv.U.Index)
		if err != nil {
			log.Fatal(err)
		}
		inv.U.Notification = true
		err = inv.AddEmail("varunramganesh@gmail.com")
		if err != nil {
			log.Fatal(err)
		}
	}

	originator, err := opensolar.NewOriginator("samuel", "p", "x", "Samuel L. Jackson", "ABC Street, London", "I am an originator")
	if err != nil {
		log.Fatal(err)
	}

	contractor, err := opensolar.NewContractor("sam", "p", "x", "Samuel Jackson", "14 ABC Street London", "This is a competing contractor")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newConstructionBond("Dec 21 2021", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Opportunity Zone Construction", 200, "5% tax for 10 years", 1, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newConstructionBond("Apr 2 2025", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Opportunity Zone Construction", 400, "No tax for 20 years", 1, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newConstructionBond("Jul 9 2029", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Opportunity Zone Construction", 100, "3% Tax for 5 Years", 1, "Osaka Development Project", "Osaka", "This Project seeks to develop cutting edge technologies in Osaka")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newLivingUnitCoop("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Coop Model", 4000, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newLivingUnitCoop("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Coop Model", 1000, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		log.Fatal(err)
	}

	_, err = newLivingUnitCoop("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Coop Model", 2000, "Osaka Development Project", "Osaka", "ODP seeks to develop cutting edge technologies in Osaka and invites investors all around the world to be a part of this new age")
	if err != nil {
		log.Fatal(err)
	}

	_, err = testSolarProject(1, "100 1000 sq.ft homes each with their own private spaces for luxury", 14000, "India Basin, San Francisco",
		0, "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate", "", "", "",
		3, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = testSolarProject(2, "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square", 30000, "Kendall Square, Boston",
		0, "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub", "", "", "",
		5, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = testSolarProject(3, "260 1500 sq.ft homes set in a medieval cathedral style construction", 40000, "Trafalgar Square, London",
		0, "Trafalgar Square is set in the heart of London's financial district, with big banks all over", "", "", "",
		7, recp.U.Index, contractor, originator, 4, 2, "blind")

	if err != nil {
		log.Fatal(err)
	}

	_, err = originator.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1, "blind") // 1 is the idnex for martin
	if err != nil {
		log.Fatal(err)
	}

	// project: Puerto Rico Project
	// STAGE 7 - Puerto Rico
// 	err = createPuertoRicoProject()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
	// project: One Mega Watt Project
	// STAGE 4 - New Hampshire
//	err = createOneMegaWattProject()
//	if err != nil {
//		log.Fatal(err)
//	}
	// project: Ten Kilowatt Project
	// STAGE 8 - Connecticut Homeless Shelter
	err = createTenKiloWattProject()
	if err != nil {
		log.Fatal(err)
	}
	// project: Ten Mega Watt Project
	// STAGE 2 - Puerto Rico Public School Bond
	err = createTenMegaWattProject()
	if err != nil {
		log.Fatal(err)
	}

	err = createOneHundredKiloWattProject()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
