package config

import (
	"os"
	"time"

	"github.com/trustwallet/go-libs/config/viper"
)

type (
	Configuration struct {
		MergeFeeBot MergeFeeBot `mapstructure:"merge-fee-bot"`
	}

	MergeFeeBot struct {
		ServiceName string      `mapstructure:"service_name"`
		LogLevel    string      `mapstructure:"log_level"`
		ClientsURLs ClientsURLs `mapstructure:"clients_urls"`
		HTTP        HTTP        `mapstructure:"http"`
		Github      Github      `mapstructure:"github"`
		Payment     Payment     `mapstructure:"payment"`
		Message     Message     `mapstructure:"message"`
		Label       Label       `mapstructure:"label"`
		User        User        `mapstructure:"user"`
		Timeout     Timeout     `mapstructure:"timeout"`
		Limitations Limitations `mapstructure:"limitations"`
	}

	ClientsURLs struct {
		BinanceDEX      string `mapstructure:"binance_dex"`
		BinanceAPI      string `mapstructure:"binance_api"`
		BinanceExplorer string `mapstructure:"binance_explorer"`
		BackendAPI      string `mapstructure:"backend_api"`
	}

	HTTP struct {
		Port int `mapstructure:"port"`
	}

	Github struct {
		BaseURL          string `mapstructure:"base_url"`
		AppID            int64  `mapstructure:"app_id"`
		AppWebhookSecret string `mapstructure:"app_webhook_secret"`
		AppPrivateKey    string `mapstructure:"app_private_key"`
		RepoOwner        string `mapstructure:"repo_owner"`
		RepoName         string `mapstructure:"repo_name"`
	}

	Payment struct {
		Assets                 []Asset `mapstructure:"assets"`
		Address                string  `mapstructure:"address"`
		Phrase                 string  `mapstructure:"phrase"`
		AmountTolerancePercent float64 `mapstructure:"amount_tolerance_percent"`
	}

	Asset struct {
		Amount float64 `mapstructure:"amount"`
		Symbol string  `mapstructure:"symbol"`
		Token  string  `mapstructure:"token"`
	}

	Message struct {
		Initial       string `mapstructure:"initial"`
		NotReceived   string `mapstructure:"not_received"`
		Received      string `mapstructure:"received"`
		ReviewCreated string `mapstructure:"review_created"`
		Reviewed      string `mapstructure:"reviewed"`
		Reminder      string `mapstructure:"reminder"`
		ClosingOldPr  string `mapstructure:"closing_old_pr"`
		Burned        string `mapstructure:"burned"`
	}

	Label struct {
		Requested string `mapstructure:"requested"`
		Paid      string `mapstructure:"paid"`
	}

	User struct {
		DeleteCommentsFromExternal bool   `mapstructure:"delete_comments_from_external"`
		Collaborators              string `mapstructure:"collaborators"`
		Moderators                 string `mapstructure:"moderators"`
	}

	Timeout struct {
		MaxAgeClose     time.Duration `mapstructure:"max_age_close"`
		MaxIdleRemind   time.Duration `mapstructure:"max_idle_remind"`
		BackgroundCheck time.Duration `mapstructure:"background_check"`
	}

	Limitations struct {
		PrFiles int `mapstructure:"pr_files"`
	}
)

// Default is a configuration instance.
var Default = Configuration{} // nolint:gochecknoglobals // config must be global

// SetConfig reads a config file and returs an initialized config instance.
func SetConfig() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yml"
	}

	viper.Load(path, &Default)
}
