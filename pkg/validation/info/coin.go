package info

import (
	"github.com/trustwallet/assets-manager/pkg/validation"
	"github.com/trustwallet/go-primitives/coin"
)

func ValidateCoin(c CoinModel, chain coin.Coin, addr string, allowedTags []string) error {
	if err := ValidateCoinRequiredKeys(c); err != nil {
		return err
	}

	// All fields validated for nil and can be safety used.
	compErr := validation.NewErrComposite()
	if err := ValidateCoinType(*c.Type); err != nil {
		compErr.Append(err)
	}

	if err := ValidateDecimals(*c.Decimals); err != nil {
		compErr.Append(err)
	}

	if err := ValidateStatus(*c.Status); err != nil {
		compErr.Append(err)
	}

	if err := ValidateTags(c.Tags, allowedTags); err != nil {
		compErr.Append(err)
	}

	if err := ValidateDescription(*c.Description); err != nil {
		compErr.Append(err)
	}

	if err := ValidateDescriptionWebsite(*c.Description, *c.Website); err != nil {
		compErr.Append(err)
	}

	if err := ValidateExplorer(*c.Explorer, *c.Name, chain, addr); err != nil {
		compErr.Append(err)
	}

	if err := ValidateLinks(c.Links); err != nil {
		compErr.Append(err)
	}

	if compErr.Len() > 0 {
		return compErr
	}

	return nil
}