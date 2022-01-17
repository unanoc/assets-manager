package blockchain

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/go-libs/blockchain/binance/api"
)

func Test_validateTx(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		memo      string
		token     string
		startTime int64
		endTime   int64
		value     *api.Tx
		want      bool
	}{
		{
			name:      "Tx is valid",
			address:   "correct addr",
			memo:      "correct memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			value: &api.Tx{
				ToAddr:    "correct addr",
				Memo:      "correct memo",
				Asset:     "TWT-8C2",
				BlockTime: 1630757347146,
				Amount:    2000,
				Type:      "TRANSFER",
			},
			want: true,
		},
		{
			name:      "Tx is valid with memo in other cases",
			address:   "correct addr",
			memo:      "correct memo",
			token:     "BNB",
			startTime: 1630187474000,
			endTime:   1630792274000,
			value: &api.Tx{
				ToAddr:    "correct addr",
				Memo:      "CoRrect mEmO",
				Asset:     "BNB",
				BlockTime: 1630757347146,
				Amount:    2000,
				Type:      "TRANSFER",
			},
			want: true,
		},
		{
			name:      "Tx type is invalid",
			address:   "correct addr",
			memo:      "correct memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			value: &api.Tx{
				ToAddr:    "correct addr",
				Memo:      "correct memo",
				Asset:     "TWT-8C2",
				BlockTime: 1630757347146,
				Amount:    2000,
				Type:      "WRONG TYPE",
			},
			want: false,
		},
		{
			name:      "Tx addr is invalid",
			address:   "correct addr",
			memo:      "correct memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			value: &api.Tx{
				ToAddr:    "incorrect addr",
				Memo:      "correct memo",
				Asset:     "TWT-8C2",
				BlockTime: 1630757347146,
				Amount:    2000,
				Type:      "TRANSFER",
			},
			want: false,
		},
		{
			name:      "Tx blocktime is invalid",
			address:   "correct addr",
			memo:      "correct memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			value: &api.Tx{
				ToAddr:    "correct addr",
				Memo:      "correct memo",
				Asset:     "TWT-8C2",
				BlockTime: 1640757347146,
				Amount:    2000,
				Type:      "TRANSFER",
			},
			want: false,
		},
		{
			name:      "Tx memo is invalid",
			address:   "correct addr",
			memo:      "correct memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			value: &api.Tx{
				ToAddr:    "correct addr",
				Memo:      "incorrect memo",
				Asset:     "TWT-8C2",
				BlockTime: 1630757347146,
				Amount:    2000,
				Type:      "TRANSFER",
			},
			want: false,
		},
		{
			name:      "Tx asset is invalid",
			address:   "correct addr",
			memo:      "correct memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			value: &api.Tx{
				ToAddr:    "correct addr",
				Memo:      "correct memo",
				Asset:     "BNB",
				BlockTime: 1630757347146,
				Amount:    2000,
				Type:      "TRANSFER",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateTx(tt.value, tt.address, tt.memo, tt.token, tt.startTime, tt.endTime)
			if got != tt.want {
				t.Errorf("value = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetPaymentStatus(t *testing.T) {
	tests := []struct {
		name      string
		txs       []api.Tx
		address   string
		memo      string
		token     string
		startTime int64
		endTime   int64
		minAmount float64
		want      *PaymentStatus
	}{
		{
			name:      "Empty transaction list",
			txs:       nil,
			address:   "addr",
			memo:      "memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			minAmount: 2000,
			want:      &PaymentStatus{},
		},

		{
			name: "Paid status with 1 transaction",
			txs: []api.Tx{
				{
					Hash:      "hash",
					ToAddr:    "to_addr",
					FromAddr:  "from_addr",
					Memo:      "memo",
					Asset:     "TWT-8C2",
					BlockTime: 1630757347146,
					Amount:    2000 * AmountPrecision,
					Type:      "TRANSFER",
				},
			},
			address:   "to_addr",
			memo:      "memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			minAmount: 2000,
			want: &PaymentStatus{
				Paid:   true,
				Amount: 2000,
				Token:  "TWT-8C2",
				Transactions: []Tx{
					{
						Hash:         "hash",
						Amount:       2000,
						Token:        "TWT-8C2",
						Date:         1630757347146,
						Memo:         "memo",
						DestAddress:  "to_addr",
						FromAddress:  "from_addr",
						ExplorerLink: fmt.Sprintf("%s/tx/%s", config.Default.ClientURLs.BinanceExplorer, "hash"),
					},
				},
			},
		},

		{
			name: "Paid status with 1 transaction, more than needed",
			txs: []api.Tx{
				{
					Hash:      "hash",
					ToAddr:    "to_addr",
					FromAddr:  "from_addr",
					Memo:      "memo",
					Asset:     "TWT-8C2",
					BlockTime: 1630757347146,
					Amount:    2500 * AmountPrecision,
					Type:      "TRANSFER",
				},
			},
			address:   "to_addr",
			memo:      "memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			minAmount: 2000,
			want: &PaymentStatus{
				Paid:   true,
				Amount: 2500,
				Token:  "TWT-8C2",
				Transactions: []Tx{
					{
						Hash:        "hash",
						Amount:      2500,
						Token:       "TWT-8C2",
						Date:        1630757347146,
						Memo:        "memo",
						DestAddress: "to_addr",
						FromAddress: "from_addr",
						ExplorerLink: fmt.Sprintf("%s/tx/%s",
							config.Default.ClientURLs.BinanceExplorer, "hash"),
					},
				},
			},
		},

		{
			name: "Paid status with a few transaction",
			txs: []api.Tx{
				{
					Hash:      "hash1",
					ToAddr:    "to_addr",
					FromAddr:  "from_addr",
					Memo:      "memo",
					Asset:     "TWT-8C2",
					BlockTime: 1630757347146,
					Amount:    1000 * AmountPrecision,
					Type:      "TRANSFER",
				},
				{
					Hash:      "hash2",
					ToAddr:    "to_addr",
					FromAddr:  "from_addr",
					Memo:      "memo",
					Asset:     "TWT-8C2",
					BlockTime: 1630757347147,
					Amount:    500 * AmountPrecision,
					Type:      "TRANSFER",
				},
				{
					Hash:      "hash3",
					ToAddr:    "to_addr",
					FromAddr:  "from_addr",
					Memo:      "memo",
					Asset:     "TWT-8C2",
					BlockTime: 1630757347148,
					Amount:    500 * AmountPrecision,
					Type:      "TRANSFER",
				},
			},
			address:   "to_addr",
			memo:      "memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			minAmount: 2000,
			want: &PaymentStatus{
				Paid:   true,
				Amount: 2000,
				Token:  "TWT-8C2",
				Transactions: []Tx{
					{
						Hash:        "hash1",
						Amount:      1000,
						Token:       "TWT-8C2",
						Date:        1630757347146,
						Memo:        "memo",
						DestAddress: "to_addr",
						FromAddress: "from_addr",
						ExplorerLink: fmt.Sprintf("%s/tx/%s",
							config.Default.ClientURLs.BinanceExplorer, "hash1"),
					},
					{
						Hash:        "hash2",
						Amount:      500,
						Token:       "TWT-8C2",
						Date:        1630757347147,
						Memo:        "memo",
						DestAddress: "to_addr",
						FromAddress: "from_addr",
						ExplorerLink: fmt.Sprintf("%s/tx/%s",
							config.Default.ClientURLs.BinanceExplorer, "hash2"),
					},
					{
						Hash:        "hash3",
						Amount:      500,
						Token:       "TWT-8C2",
						Date:        1630757347148,
						Memo:        "memo",
						DestAddress: "to_addr",
						FromAddress: "from_addr",
						ExplorerLink: fmt.Sprintf("%s/tx/%s",
							config.Default.ClientURLs.BinanceExplorer, "hash3"),
					},
				},
			},
		},

		{
			name: "Not paid status, not enough",
			txs: []api.Tx{
				{
					Hash:      "hash",
					ToAddr:    "to_addr",
					FromAddr:  "from_addr",
					Memo:      "memo",
					Asset:     "TWT-8C2",
					BlockTime: 1630757347146,
					Amount:    1800 * AmountPrecision,
					Type:      "TRANSFER",
				},
			},
			address:   "to_addr",
			memo:      "memo",
			token:     "TWT-8C2",
			startTime: 1630187474000,
			endTime:   1630792274000,
			minAmount: 2000,
			want: &PaymentStatus{
				Paid:   false,
				Amount: 1800,
				Token:  "TWT-8C2",
				Transactions: []Tx{
					{
						Hash:        "hash",
						Amount:      1800,
						Token:       "TWT-8C2",
						Date:        1630757347146,
						Memo:        "memo",
						DestAddress: "to_addr",
						FromAddress: "from_addr",
						ExplorerLink: fmt.Sprintf("%s/tx/%s",
							config.Default.ClientURLs.BinanceExplorer, "hash"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPaymentStatus(tt.txs, tt.address, tt.memo, tt.token, tt.startTime, tt.endTime, tt.minAmount)
			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("value = %v, want %v", got, tt.want)
			}
		})
	}
}
