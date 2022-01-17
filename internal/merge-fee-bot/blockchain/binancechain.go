package blockchain

import (
	"fmt"

	"github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/go-libs/blockchain/binance/api"
	"github.com/trustwallet/go-primitives/coin"
)

type Client struct {
	client       client.DexClient
	customClient api.Client
}

func NewClient() *Client {
	keyManager, err := keys.NewMnemonicKeyManager(config.Default.MergeFeeBot.Payment.Phrase)
	if err != nil {
		// return nil, errors.Wrap(err, "failed to create a mnemonic key manager")
		log.Error(err, "failed to create a mnemonic key manager")
	}

	c, err := client.NewDexClient(config.Default.MergeFeeBot.ClientsURLs.BinanceDEX, types.ProdNetwork, keyManager)
	if err != nil {
		// return nil, errors.Wrap(err, "failed to create a dex client to binance api")
		log.Error(err, "failed to create a dex client to binance api")
	}

	return &Client{
		client:       c,
		customClient: api.InitClient(config.Default.MergeFeeBot.ClientsURLs.BinanceAPI, nil),
	}
}

func (c *Client) BurnToken(token string, amount int64) (string, error) {
	if c.client == nil {
		return "", nil
	}

	if token == coin.Coins[coin.BINANCE].Symbol {
		return "", nil
	}

	res, err := c.client.BurnToken(token, amount, true)
	if err != nil {
		return "", errors.Wrap(err, "failed to burn a token")
	}

	log.WithFields(log.Fields{
		"token":  token,
		"amount": amount,
	}).Debugf("tokens has been burned")

	return fmt.Sprintf("%s/tx/%s", config.Default.MergeFeeBot.ClientsURLs.BinanceExplorer, res.Hash), nil
}

func (c *Client) GetTransactionsForAddress(address string) ([]api.Tx, error) {
	var txs []api.Tx
	var err error

	if txs, err = c.customClient.GetTransactionsByAddress(address, 50); err != nil {
		return nil, errors.Wrap(err, "failed to fetch transactions by address")
	}

	return txs, nil
}
