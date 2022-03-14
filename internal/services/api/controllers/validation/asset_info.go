package validation

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-go-libs/validation"
	"github.com/trustwallet/assets-go-libs/validation/info"
	"github.com/trustwallet/assets-go-libs/validation/info/external"
	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/go-primitives/address"
	"github.com/trustwallet/go-primitives/coin"
	"github.com/trustwallet/go-primitives/types"
)

func (i *Controller) ValidateAssetInfo(asset AssetInfoRequest) *AssetInfoResponse {
	errors := make([]Error, 0)

	assetModel := mapAssetModel(asset)

	externalTokenInfo, err := external.GetTokenInfo(asset.ID, asset.Type)
	if err != nil {
		log.WithError(err).Debugf("Failed to get token info")
	}

	checks := []func() error{
		func() error { return validateAssetInfoID(asset.ID, asset.Type) },
		func() error { return validateAssetInfoType(asset.Type) },
		func() error { return validateAssetInfoDecimals(asset.Decimals, externalTokenInfo) },
		func() error { return info.ValidateDescription(asset.Description) },
		func() error { return info.ValidateDescriptionWebsite(asset.Description, asset.Website) },
		func() error { return validateAssetInfoExplorer(asset.Explorer, asset.ID, asset.Type) },
		func() error { return info.ValidateStatus(asset.Status) },
		func() error { return validateAssetInfoLinks(asset.Links) },
		func() error { return validateAssetInfoTags(asset.Tags) },
		func() error { return validateAssetInfoHolders(externalTokenInfo) },
		func() error { return info.ValidateAssetRequiredKeys(assetModel) },
	}

	for _, check := range checks {
		if err := check(); err != nil {
			errors = append(errors, Error{Message: err.Error()})
		}
	}

	status := StatusTypeOk
	if len(errors) > 0 {
		status = StatusTypeError
	}

	return &AssetInfoResponse{
		Status: status,
		Errors: errors,
	}
}

func mapAssetModel(asset AssetInfoRequest) info.AssetModel {
	links := make([]info.Link, len(asset.Links))
	for i := range asset.Links {
		links[i] = info.Link{
			Name: &asset.Links[i].Name,
			URL:  &asset.Links[i].URL,
		}
	}

	return info.AssetModel{
		ID:          &asset.ID,
		Name:        &asset.Name,
		Type:        &asset.Type,
		Symbol:      &asset.Symbol,
		Website:     &asset.Website,
		Explorer:    &asset.Explorer,
		Description: &asset.Description,
		Status:      &asset.Status,
		Decimals:    &asset.Decimals,
		Links:       links,
		Tags:        asset.Tags,
	}
}

func validateAssetInfoID(tokenID, tokenType string) error {
	if tokenID == "" {
		return fmt.Errorf("token id cannot be empty")
	}

	if strings.ToUpper(tokenType) == string(types.ERC20) || strings.ToUpper(tokenType) == string(types.BEP20) {
		if !validation.IsEthereumAddress(tokenID) {
			return fmt.Errorf("ID is not a valid Ethereum address")
		}

		checksum, err := address.EIP55Checksum(tokenID)
		if err != nil {
			return fmt.Errorf("failed to get checksum for %s: %w", tokenID, err)
		}

		if checksum != tokenID {
			return fmt.Errorf("id is not in checksum format, should be %s (not %s). "+
				"Please rename it. You may need to rename to a temp name first, "+
				"then to the checksum format, because lowercase-uppercase-only renames "+
				"are often ignored by the Git client or the filesystem", checksum, tokenID)
		}
	}

	return nil
}

func validateAssetInfoType(tokenType string) error {
	if tokenType == "" {
		return fmt.Errorf("type field cannot be empty")
	}

	_, err := types.GetChainFromAssetType(tokenType)
	if err != nil {
		return fmt.Errorf("invalid type field: %w", err)
	}

	return nil
}

func validateAssetInfoDecimals(decimals int, externalTokenInfo *external.TokenInfo) error {
	decimalsMaxValue := config.Default.Validation.Asset.DecimalsMaxValue

	if decimals > decimalsMaxValue {
		return fmt.Errorf("decimals value is invalid: %d (max %d)", decimals, decimalsMaxValue)
	}

	if externalTokenInfo == nil {
		return nil
	}

	if decimals != externalTokenInfo.Decimals {
		return fmt.Errorf("decimals value is incorrect: expected %d instead of %d",
			externalTokenInfo.Decimals, decimals)
	}

	return nil
}

func validateAssetInfoExplorer(explorer, tokenID, tokenType string) error {
	if explorer == "" {
		return fmt.Errorf("explorer field cannot be empty")
	}

	chain, err := types.GetChainFromAssetType(tokenType)
	if err != nil {
		return err
	}

	explorerStandart, err := coin.GetCoinExploreURL(chain, tokenID, tokenType)
	if err != nil {
		return err
	}

	if !strings.EqualFold(explorer, explorerStandart) {
		return fmt.Errorf("explorer field incorrect: should be standard %s, not %s", explorerStandart, explorer)
	}

	return nil
}

func validateAssetInfoLinks(links []Link) error {
	linksMinRequired := config.Default.Validation.Asset.LinksMinRequired

	if len(links) < linksMinRequired {
		return fmt.Errorf("at least %d links are required, %d present. "+
			"Add as many as you can: twitter, github, telegram, reddit, etc", linksMinRequired, len(links))
	}

	log.Debugf("links in request: %v", links)

	infoLinks := make([]info.Link, len(links))
	for i := range links {
		infoLinks[i] = info.Link{
			Name: &links[i].Name,
			URL:  &links[i].URL,
		}
	}

	return info.ValidateLinks(infoLinks)
}

func validateAssetInfoTags(tags []string) error {
	log.WithFields(log.Fields{
		"count": len(tags),
		"tags":  tags,
	}).Debug("tags in request")

	tagsMinRequired := config.Default.Validation.Asset.TagsMinRequired
	if len(tags) < tagsMinRequired {
		return fmt.Errorf("at least %d tag is required", tagsMinRequired)
	}

	tagMap := make(map[string]bool)
	for _, tag := range config.Default.Tags {
		tagMap[tag.ID] = true
	}

	for _, tag := range tags {
		if _, exists := tagMap[tag]; !exists {
			return fmt.Errorf("tag '%s' is not allowed", tag)
		}
	}

	return nil
}

func validateAssetInfoHolders(extTokenInfo *external.TokenInfo) error {
	if extTokenInfo == nil {
		return fmt.Errorf("number of holders not checked: please, check it manually")
	}

	holdersMinRequired := config.Default.Validation.Asset.HoldersMinRequired

	if extTokenInfo.HoldersCount < holdersMinRequired {
		return fmt.Errorf("low token circulation: number of holders is %d, below limit of %d",
			extTokenInfo.HoldersCount, holdersMinRequired)
	}

	return nil
}
