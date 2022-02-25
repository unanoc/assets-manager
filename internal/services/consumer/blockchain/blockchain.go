package blockchain

import (
	"fmt"
	"strings"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/go-libs/blockchain/binance/api"
)

const AmountPrecision = 100000000.0

type PaymentStatus struct {
	Paid         bool
	Amount       float64
	Token        string
	Transactions []Tx
}

type Tx struct {
	Hash         string
	Amount       float64
	Token        string
	Date         int64
	Memo         string
	DestAddress  string
	FromAddress  string
	ExplorerLink string
}

func GetPaymentStatus(
	txs []api.Tx, address, memo, token string, startTime, endTime int64, minAmount float64,
) *PaymentStatus {
	if len(txs) == 0 {
		return &PaymentStatus{}
	}

	var sum float64
	transactions := make([]Tx, 0)

	for i := range txs {
		tx := Tx{
			Hash:        txs[i].Hash,
			Amount:      txs[i].Amount / 100000000,
			Token:       txs[i].Asset,
			Date:        txs[i].BlockTime,
			Memo:        txs[i].Memo,
			DestAddress: txs[i].ToAddr,
			FromAddress: txs[i].FromAddr,
			ExplorerLink: fmt.Sprintf("%s/tx/%s",
				config.Default.Clients.Binance.Explorer, txs[i].Hash),
		}

		valid := ValidateTx(&txs[i], address, memo, token, startTime, endTime)
		if !valid {
			continue
		}

		sum += tx.Amount
		transactions = append(transactions, tx)
	}

	return &PaymentStatus{
		Paid:         sum >= minAmount,
		Amount:       sum,
		Token:        token,
		Transactions: transactions,
	}
}

func ValidateTx(tx *api.Tx, address, memo, token string, startTime, endTime int64) bool {
	if tx.Type != "TRANSFER" {
		return false
	}
	if tx.ToAddr != address {
		return false
	}
	if tx.BlockTime < startTime || tx.BlockTime > endTime {
		return false
	}
	if !strings.EqualFold(tx.Memo, memo) {
		return false
	}
	if !strings.EqualFold(tx.Asset, token) {
		return false
	}

	return true
}
