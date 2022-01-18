package events

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v38/github"

	"github.com/trustwallet/assets-manager/internal/config"
)

const (
	qrGeneratorLink       = "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data="
	deepLinkToTrustWallet = "https://link.trustwallet.com/send?coin="
	monthInSec            = 30 * 86400 * 1000
	tokenID               = 714
)

type PaymentsParams struct {
	Payments []Payment
	User     string
	Address  string
	Phrase   string
	QR       string
}

type Payment struct {
	Amount      float64
	Symbol      string
	Token       string
	MinAmount   float64
	Memo        string
	CreatedTime int64
	EndTime     int64
}

func getPaymentParams(pr *github.PullRequest) *PaymentsParams {
	if pr == nil {
		return &PaymentsParams{}
	}

	createdTime := pr.GetCreatedAt().Unix() * 1000
	memo := strconv.Itoa(pr.GetNumber())
	payments := make([]Payment, len(config.Default.Payment.Options))

	for i := range payments {
		payments[i].Amount = config.Default.Payment.Options[i].Amount
		payments[i].Symbol = config.Default.Payment.Options[i].Symbol
		payments[i].Token = config.Default.Payment.Options[i].Token
		payments[i].MinAmount = getMinAmount(config.Default.Payment.TolerancePercent,
			config.Default.Payment.Options[i].Amount)
		payments[i].Memo = memo
		payments[i].CreatedTime = createdTime
		payments[i].EndTime = createdTime + monthInSec
	}

	qrTw, qrFull := getQR(config.Default.Payment.Options[0].Amount,
		config.Default.Payment.Address, memo)

	return &PaymentsParams{
		Payments: payments,
		User:     pr.GetUser().GetLogin(),
		Address:  config.Default.Payment.Address,
		Phrase:   config.Default.Payment.SeedPhrase,
		QR:       fmt.Sprintf("**QR** code: [Trust]( %s ) | [other wallet]( %s )", qrTw, qrFull),
	}
}

func getMinAmount(tolarancePercent, amount float64) float64 {
	return 0.01 * math.Min(100.0, math.Max(95, tolarancePercent)) * amount
}

func getQR(amount float64, address, memo string) (qrTw, qrFull string) {
	deepLink := fmt.Sprintf("%s%d&address=%s&amount=%d&memo=%s",
		deepLinkToTrustWallet,
		tokenID,
		address,
		int(amount),
		memo,
	)
	uri := fmt.Sprintf("%s?amount=%d&memo=%s", address, int(amount), memo)
	qrTw = fmt.Sprintf("%s%s", qrGeneratorLink, url.QueryEscape(deepLink))
	qrFull = fmt.Sprintf("%s%s", qrGeneratorLink, url.QueryEscape(uri))

	return
}

type contentParams struct {
	PP               *PaymentsParams
	PaidAmount       float64
	PaidSymbol       string
	PaidExplorerLink string
	BurnExplorerLink string
	Moderators       string
}

func substituteDynamicContent(text string, p *contentParams) string {
	m := make(map[string]string)

	if p != nil && p.PP != nil {
		m["$PAY1_AMOUNT"] = strconv.Itoa(int(p.PP.Payments[0].Amount))
		m["$PAY1_SYMBOL"] = p.PP.Payments[0].Symbol
		m["$PAY2_AMOUNT"] = strconv.Itoa(int(p.PP.Payments[1].Amount))
		m["$PAY2_SYMBOL"] = p.PP.Payments[1].Symbol
		m["$PAY1_MEMO"] = p.PP.Payments[0].Memo
		m["$PAY1_ADDRESS"] = p.PP.Address
		m["$QR_CODE"] = p.PP.QR
		m["$USER"] = p.PP.User
	}

	if p != nil {
		m["$PAID_AMOUNT"] = fmt.Sprintf("%.2f", p.PaidAmount)
		m["$PAID_SYMBOL"] = p.PaidSymbol
		m["$PAID_EXPLORER_LINK"] = p.PaidExplorerLink
		m["$BURN_EXPLORER_LINK"] = p.BurnExplorerLink
		m["$MODERATORS"] = getModerators(p.Moderators)
	}

	for k, v := range m {
		text = strings.ReplaceAll(text, k, v)
	}

	return text
}

func getLogoHTML(logoURL string) string {
	return `<span style="padding: 5px; background-color: rgb(32,32,32);"><img src="` + logoURL +
		`" style="max-width: 64px; border-radius: 48%;" width="48" height="48"/></span>`
}

func getModerators(rawModerators string) string {
	moderatorsList := strings.Split(rawModerators, ",")
	if len(moderatorsList) == 0 {
		return ""
	}

	var result string

	for i := range moderatorsList {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("@%s", moderatorsList[i])
	}

	return fmt.Sprintf("%s: please take note.", result)
}
