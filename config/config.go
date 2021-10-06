package config

import (
	"github.com/PxyUp/binance_bot/currency"
	"github.com/PxyUp/binance_bot/services"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configs struct {
	Currencies []*Item           `yaml:"currencies"`
	Binance    *services.Binance `yaml:"binance"`
}

type Item struct {
	Name             string                    `yaml:"name"`
	UseSR            bool                      `yaml:"useSR"`
	UseStopOrder     bool                      `yaml:"useStopOrder"`
	StopOrderDiff    float64                   `yaml:"stopOrderDiff"`
	Precision        int                       `yaml:"precision"`
	Binance          *services.Binance         `yaml:"binance"`
	SRDiffPercent    float64                   `yaml:"srDiffPercent"`
	ShortPatterns    bool                      `yaml:"shortPatterns"`
	Symbol           string                    `yaml:"symbol"`
	Count            float64                   `yaml:"count"`
	FloatNumbers     int                       `yaml:"floatNumbers"`
	CountNumber      int                       `yaml:"countNumber"`
	KlineIntervals   []services.CandleInterval `yaml:"klineIntervals"`
	Interval         int                       `yaml:"interval"`
	Mode             currency.WatcherMode      `yaml:"mode"`
	PriceBuffer      int                       `yaml:"priceBuffer"`
	DryRun           bool                      `yaml:"dryRun"`
	CandleTimeBuffer int                       `yaml:"candleTimeBuffer"`
}

type config struct {
	path string
}

func New(path string) *config {
	return &config{
		path: path,
	}
}

func (c *config) Get() (*Configs, error) {
	conf := &Configs{}

	yamlFile, err := ioutil.ReadFile(c.path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, conf)

	if err != nil {
		return nil, err
	}

	return conf, nil
}
