package pool

import (
	"context"
	"fmt"
	"github.com/PxyUp/binance_bot/config"
	"github.com/PxyUp/binance_bot/currency"
	"github.com/PxyUp/binance_bot/patterns"
	"github.com/PxyUp/binance_bot/services"
	"log"
	"time"
)

type Pool struct {
	watchers []*currency.Watcher
	closers  []chan struct{}
	ctx      context.Context
}

func (p *Pool) UseConfig(config *config.Configs) *Pool {
	for _, symbol := range config.Currencies {
		// We by default use global, but can overwrite
		name := fmt.Sprintf("Symbol: %s; DryRun: %t; UseSR: %t; UseStopOrder: %t, Shorts: %t", symbol.Symbol, symbol.DryRun, symbol.UseSR, symbol.UseStopOrder, symbol.ShortPatterns)
		binConf := config.Binance
		if symbol.Binance != nil {
			binConf = symbol.Binance
		}
		if  symbol.Name != "" {
			name = symbol.Name
		}
		p.watchers = append(p.watchers,
			currency.NewCurrencyWatcher(
				name,
				symbol.UseSR,
				symbol.UseStopOrder,
				symbol.StopOrderDiff,
				symbol.SRDiffPercent,
				symbol.Precision,
				symbol.Symbol,
				symbol.Mode,
				time.Duration(symbol.Interval)*time.Second,
				patterns.AllBuy,
				patterns.AllSell,
				patterns.UseDefaults(true, symbol.ShortPatterns),
				patterns.UseDefaults(false, symbol.ShortPatterns),
				services.Commission,
				symbol.Count,
				symbol.PriceBuffer,
				"%."+fmt.Sprintf("%df", symbol.FloatNumbers),
				symbol.CountNumber,
				time.Duration(symbol.CandleTimeBuffer)*time.Second,
				symbol.DryRun,
				services.New(binConf),
				symbol.KlineIntervals,
			),
		)
	}
	return p
}

func New() *Pool {
	var watchers []*currency.Watcher

	return &Pool{
		watchers: watchers,
	}
}

func (p *Pool) Run() (context.Context, chan struct{}) {
	ctx := context.Background()

	for _, watcher := range p.watchers {
		p.closers = append(p.closers, watcher.Run())
	}

	p.ctx = ctx

	closer := make(chan struct{}, 1)

	go func() {
		<-closer

		for index, _ := range p.watchers {
			p.closers[index] <- struct{}{}
		}
	}()

	return ctx, closer
}

func (p *Pool) GetProfit() float64 {
	var totalProfit float64 = 0

	for _, w := range p.watchers {
		profit := w.GetProfit()
		if !w.IsDryRun() {
			log.Println("Profit from", w.GetName(), "is", profit)
			totalProfit += profit
			continue
		}
		log.Println("Profit from dry RUN", w.GetName(), "is", profit)
	}

	return totalProfit
}
