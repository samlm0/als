package config

import (
	"io"
	"log"
	"net/http"
	"os"
)

var Config *ALSConfig
var IsInternalCall bool

type ALSConfig struct {
	ListenHost string `json:"-"`
	ListenPort string `json:"-"`

	Location string `json:"location"`

	PublicIPv4 string `json:"public_ipv4"`
	PublicIPv6 string `json:"public_ipv6"`

	Iperf3StartPort int `json:"-"`
	Iperf3EndPort   int `json:"-"`

	SpeedtestFileList []string `json:"speedtest_files"`

	SponsorMessage string `json:"sponsor_message"`

	FeaturePing          bool `json:"feature_ping"`
	FeatureShell         bool `json:"feature_shell"`
	FeatureLibrespeed    bool `json:"feature_librespeed"`
	FeatureFileSpeedtest bool `json:"feature_filespeedtest"`
	FeatureIperf3        bool `json:"feature_iperf3"`
	FeatureMTR           bool `json:"feature_mtr"`
	FeatureTraceroute    bool `json:"feature_traceroute"`
	FeatureIfaceTraffic  bool `json:"feature_iface_traffic"`
}

func GetDefaultConfig() *ALSConfig {
	defaultConfig := &ALSConfig{
		ListenHost:      "0.0.0.0",
		ListenPort:      "8080",
		Location:        "未设置",
		Iperf3StartPort: 20000,
		Iperf3EndPort:   30000,

		SpeedtestFileList: []string{"1MB", "10MB", "100MB", "1GB", "100GB"},
		PublicIPv4:        "",
		PublicIPv6:        "",

		FeaturePing:          true,
		FeatureShell:         true,
		FeatureLibrespeed:    true,
		FeatureFileSpeedtest: true,
		FeatureIperf3:        true,
		FeatureMTR:           true,
		FeatureTraceroute:    true,
		FeatureIfaceTraffic:  true,
	}

	return defaultConfig
}

func Load() {
	// default config
	Config = GetDefaultConfig()
	LoadFromEnv()
}

func LoadWebConfig() {
	Load()
	log.Default().Println("Loading config for web services...")
	if Config.PublicIPv4 == "" && Config.PublicIPv6 == "" {
		go updatePublicIP()
	}
}

func LoadSponsorMessage() {
	if Config.SponsorMessage == "" {
		return
	}

	log.Default().Println("Loading sponser message...")

	if _, err := os.Stat(Config.SponsorMessage); err == nil {
		content, err := os.ReadFile(Config.SponsorMessage)
		if err == nil {
			Config.SponsorMessage = string(content)
			return
		}
	}

	resp, err := http.Get(Config.SponsorMessage)
	if err == nil {
		content, err := io.ReadAll(resp.Body)
		if err == nil {
			Config.SponsorMessage = string(content)
			return
		}
	}

	log.Default().Panicln("Failed to load sponsor message.")
}
