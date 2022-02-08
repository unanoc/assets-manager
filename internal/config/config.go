package config

import (
	"os"
	"time"

	"github.com/trustwallet/go-libs/config/viper"
)

type Configuration struct {
	ServiceName string `mapstructure:"service_name"`
	LogLevel    string `mapstructure:"log_level"`
	Port        int    `mapstructure:"port"`

	Sentry struct {
		DSN        string  `mapstructure:"dsn"`
		SampleRate float32 `mapstructure:"sample_rate"`
	} `mapstructure:"sentry"`

	Rabbitmq struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"rabbitmq"`

	Consumer struct {
		Workers int `mapstructure:"workers"`
	} `mapstructure:"consumer"`

	Metrics struct {
		Path           string        `mapstructure:"path"`
		PushGatewayURL string        `mapstructure:"pushgateway_url"`
		PushInterval   time.Duration `mapstructure:"push_interval"`
	} `mapstructure:"metrics"`

	Clients struct {
		Binance struct {
			DEX      string `mapstructure:"dex"`
			API      string `mapstructure:"api"`
			Explorer string `mapstructure:"explorer"`
		} `mapstructure:"binance"`

		AssetsManager struct {
			API string `mapstructure:"api"`
			App string `mapstructure:"app"`
		} `mapstructure:"assets_manager"`
	} `mapstructure:"clients"`

	Github struct {
		APIURL           string `mapstructure:"api_url"`
		AppID            int64  `mapstructure:"app_id"`
		AppWebhookSecret string `mapstructure:"app_webhook_secret"`
		AppPrivateKey    string `mapstructure:"app_private_key"`
		RepoOwner        string `mapstructure:"repo_owner"`
		RepoName         string `mapstructure:"repo_name"`
		BaseURL          string `mapstructure:"base_url"`
		ClientID         string `mapstructure:"client_id"`
		ClientSecret     string `mapstructure:"client_secret"`
	} `mapstructure:"github"`

	Payment struct {
		Options []struct {
			Amount float64 `mapstructure:"amount"`
			Symbol string  `mapstructure:"symbol"`
			Token  string  `mapstructure:"token"`
		} `mapstructure:"options"`

		Address          string  `mapstructure:"address"`
		SeedPhrase       string  `mapstructure:"seed_phrase"`
		TolerancePercent float64 `mapstructure:"tolerance_percent"`
	} `mapstructure:"payment"`

	Message struct {
		Initial       string `mapstructure:"initial"`
		NotReceived   string `mapstructure:"not_received"`
		Received      string `mapstructure:"received"`
		ReviewCreated string `mapstructure:"review_created"`
		Reviewed      string `mapstructure:"reviewed"`
		Reminder      string `mapstructure:"reminder"`
		ClosingOldPR  string `mapstructure:"closing_old_pr"`
		Burned        string `mapstructure:"burned"`
	} `mapstructure:"message"`

	Label struct {
		Requested string `mapstructure:"requested"`
		Paid      string `mapstructure:"paid"`
	} `mapstructure:"label"`

	UserAccess struct {
		DeleteCommentsFromExternal bool   `mapstructure:"delete_comments_from_external"`
		Collaborators              string `mapstructure:"collaborators"`
		Moderators                 string `mapstructure:"moderators"`
	} `mapstructure:"user_access"`

	Timeout struct {
		MaxAgeClose     time.Duration `mapstructure:"max_age_close"`
		MaxIdleRemind   time.Duration `mapstructure:"max_idle_remind"`
		BackgroundCheck time.Duration `mapstructure:"background_check"`
	} `mapstructure:"timeout"`

	Limitation struct {
		PrFilesNumMax int `mapstructure:"pr_files_num_max"`
	} `mapstructure:"limitation"`

	Validation struct {
		Asset struct {
			DecimalsMaxValue     int `mapstructure:"decimals_max_value"`
			DescriptionMaxLength int `mapstructure:"description_max_length"`
			LinksMinRequired     int `mapstructure:"links_min_required"`
			TagsMinRequired      int `mapstructure:"tags_min_required"`
			HoldersMinRequired   int `mapstructure:"holders_min_required"`
		} `mapstructure:"asset"`
	} `mapstructure:"validation"`

	Tags []struct {
		ID          string `mapstructure:"id"`
		Name        string `mapstructure:"name"`
		Description string `mapstructure:"description"`
	} `mapstructure:"tags"`
}

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
