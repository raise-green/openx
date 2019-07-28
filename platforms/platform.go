package platform

import (
	"log"
	"os"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	scan "github.com/Varunram/essentials/scan"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
	"github.com/pkg/errors"
)

// the platform structure is the backend representation of the frontend UI.
// on a very low level, this should just be a pubkey + seed pair. Each platform
// needs to be hosted somewhere, so it is necessary that each platform should have
// its own pubkey and seed pair
// InitializePlatform returns the platform publickey and seed
// We have a new model in which we have a new seed for every project that is
// advertised on the platform. The way this would work is that it sets up the assets,
// and then we freeze the account to freeze issuance. This would mean we would no longer
// be able to transact with the account although people can still send funds to it
// in this case, they would send us back DebtAssets provided they have sufficient
// stableUSD balance. Else they would not be able to trigger payback.

// InitializePlatform starts the platform
func InitializePlatform() error {
	var publicKey string
	var seed string
	var err error

	// now we can be sure we have the directory, check for seed
	if _, err := os.Stat(consts.PlatformSeedFile); !os.IsNotExist(err) {
		// the seed exists
		log.Println("ENTER YOUR PASSWORD TO DECRYPT THE PLATFORM SEED FILE")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return errors.Wrap(err, "couldn't scan raw password")
		}
		consts.PlatformPublicKey, consts.PlatformSeed, err = wallet.RetrieveSeed(consts.PlatformSeedFile, password)
		return err
	}
	// platform doesn't exist or user doesn't have encrypted file. Ask
	log.Println("DO YOU HAVE YOUR RAW SEED? IF SO, ENTER SEED. ELSE ENTER N")
	seed, err = scan.ScanString()
	if err != nil {
		return errors.Wrap(err, "couldn't scan raw string")
	}
	if seed == "N" || seed == "n" {
		// no seed, no file, create new keypair
		// need to pass the password for the  eed file
		log.Println("Enter a password to encrypt your master seed. Please store this in a very safe place. This prompt will not ask to confirm your password")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return err
		}
		publicKey, seed, err = wallet.NewSeedStore(consts.PlatformSeedFile, password)
		if err != nil {
			return errors.Wrap(err, "couldn't retrieve seed")
		}
		err = xlm.GetXLM(publicKey)
		if err != nil {
			return errors.Wrap(err, "error while getting xlm")
		}
	} else {
		// no file, retrieve pukbey
		// user has given us a seed, validate
		log.Println("ENTER A PASSWORD TO DECRYPT YOUR SEED")
		password, err := scan.ScanRawPassword()
		if err != nil {
			return err
		}

		publicKey, err = wallet.ReturnPubkey(seed)
		if err != nil {
			return err
		}

		err = wallet.StoreSeed(seed, password, consts.PlatformSeed)
		if err != nil {
			return err
		}
	}
	_ = xlm.GetXLM(publicKey) // the API request errors out even on success, so
	// don't catch this error
	_, txhash, err := xlm.SetAuthImmutable(seed)
	log.Println("TX HASH FOR SETOPTIONS: ", txhash)
	if err != nil {
		log.Println("ERROR WHILE SETTING OPTIONS")
	}
	// make the platform trust the stablecoin for receiving payments
	txhash, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, "10000000000", seed)
	if err != nil {
		return err
	}

	_, _, err = assets.SendAssetFromIssuer(consts.StablecoinCode, publicKey, "10", consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}

	log.Println("Platform trusts stablecoin: ", txhash)
	consts.PlatformPublicKey = publicKey
	consts.PlatformSeed = seed
	return err
}

// RefillPlatform checks whether the publicKey passed has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask the friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	// check whether the investor has XLM already
	balance, err := xlm.GetNativeBalance(publicKey)
	if err != nil {
		return err
	}
	// balance is in string, convert to int
	balanceF, err := utils.ToFloat(balance)
	if err != nil {
		return err
	}
	log.Println("Platform's balance is: ", balanceF)
	if balanceF < 21 { // 1 to account for fees
		// get coins if balance is this low
		log.Println("Refilling platform balance")
		err := xlm.GetXLM(publicKey)
		// refill platform sufficiently well and interact with a cold wallet that we
		// have previously set earlier to avoid hacks and similar
		if err != nil {
			return err
		}
	}
	return nil
}
