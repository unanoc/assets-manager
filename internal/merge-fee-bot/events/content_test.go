package events

import "testing"

func Test_GetQR(t *testing.T) {
	tests := []struct {
		name       string
		amount     float64
		token      string
		address    string
		memo       string
		wantTwQR   string
		wantFullQR string
	}{
		{
			name:       "QR for TWT token",
			amount:     2000,
			address:    "bnb1tqq9llyr3dyjd559dha6z5r5etk3qfwk07m098",
			memo:       "3395",
			wantTwQR:   "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=https%3A%2F%2Flink.trustwallet.com%2Fsend%3Fcoin%3D714%26address%3Dbnb1tqq9llyr3dyjd559dha6z5r5etk3qfwk07m098%26amount%3D2000%26memo%3D3395",
			wantFullQR: "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=bnb1tqq9llyr3dyjd559dha6z5r5etk3qfwk07m098%3Famount%3D2000%26memo%3D3395",
		},
		{
			name:       "QR for BNB token",
			amount:     5,
			address:    "bnb1tqq9llyr3dyjd559dha6z5r5etk3qfwk07m098",
			memo:       "12",
			wantTwQR:   "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=https%3A%2F%2Flink.trustwallet.com%2Fsend%3Fcoin%3D714%26address%3Dbnb1tqq9llyr3dyjd559dha6z5r5etk3qfwk07m098%26amount%3D5%26memo%3D12",
			wantFullQR: "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=bnb1tqq9llyr3dyjd559dha6z5r5etk3qfwk07m098%3Famount%3D5%26memo%3D12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwQR, gotFullQR := getQR(tt.amount, tt.address, tt.memo)
			if gotTwQR != tt.wantTwQR {
				t.Errorf("getQR() = %v, want %v", gotTwQR, tt.wantTwQR)
			}
			if gotFullQR != tt.wantFullQR {
				t.Errorf("getQR() = %v, want %v", gotFullQR, tt.wantFullQR)
			}
		})
	}
}

func Test_GetMinAmount(t *testing.T) {
	tests := []struct {
		name             string
		tolarancePercent float64
		amount           float64
		want             float64
	}{
		{
			name:             "96 percents for 2000",
			tolarancePercent: 96,
			amount:           2000,
			want:             1920,
		},
		{
			name:             "96 percent for 5",
			tolarancePercent: 96,
			amount:           5,
			want:             4.8,
		},
		{
			name:             "100 percent for 2000",
			tolarancePercent: 100,
			amount:           2000,
			want:             2000,
		},
		{
			name:             "97 percent for 2000",
			tolarancePercent: 97,
			amount:           2000,
			want:             1940,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getMinAmount(tt.tolarancePercent, tt.amount)
			if got != tt.want {
				t.Errorf("value = %v, want %v", got, tt.want)
			}
		})
	}
}
