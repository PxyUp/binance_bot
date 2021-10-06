package currency

import (
	"errors"
	"github.com/PxyUp/binance_bot/order"
	"github.com/PxyUp/binance_bot/patterns"
	"github.com/PxyUp/binance_bot/services"
	"github.com/PxyUp/binance_bot/utils"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
	"golang.org/x/sync/errgroup"
	"log"
	"math"
	"sync"
	"time"
)

type WatcherMode string

const (
	Candle       WatcherMode = "Candle"
	LivePrice    WatcherMode = "LivePrice"
	AccountError int64       = 1010
	BalanceError int64       = -2010
	HardLimit                = 0.006
)

var (
	userErrors       = []int64{AccountError, BalanceError}
	errWrongInterval = errors.New("wrong interval")
)

type Watcher struct {
	name                string
	useSR               bool
	useStopOrder        bool
	stopOrderDiff       float64
	srDiffPercent       float64
	precision           int
	symbol              string
	buyPatterns         [][]int8
	sellPatterns        [][]int8
	interval            time.Duration
	lastPrice           float64
	lastAction          patterns.Action
	orderService        interface{}
	afterCommission     float64
	count               float64
	initialCount        float64
	pricesBuffer        int
	profit              float64
	format              string
	numbersAfter        int
	needBuy             bool
	mode                WatcherMode
	buyCandlesPatterns  [][]*patterns.CandlePattern
	sellCandlesPatterns [][]*patterns.CandlePattern
	dryRun              bool
	candleTimeBuffer    time.Duration
	prices              []*patterns.PriceRow
	klineIntervals      []services.CandleInterval
	orderWatcher        *order.OrderWatcher
	closerOrderWatcher  chan struct{}

	stopOrderId        *int64
	lastStopOrderPrice *float64
	mutex              sync.Mutex
	service            *services.Service
}

func NewCurrencyWatcher(
	name string,
	useSR bool,
	useStopOrder bool,
	stopOrderDiff float64,
	srDiffPercent float64,
	precision int,
	symbol string,
	mode WatcherMode,
	interval time.Duration,
	buyPatterns [][]int8,
	sellPatterns [][]int8,
	buyCandlesPatterns [][]*patterns.CandlePattern,
	sellCandlesPatterns [][]*patterns.CandlePattern,
	commission float64,
	count float64,
	pricesBuffer int,
	format string,
	numbersAfter int,
	candleTimeBuffer time.Duration,
	dryRun bool,
	service *services.Service,
	klineIntervals []services.CandleInterval,
) *Watcher {
	return &Watcher{
		name:                name,
		useSR:               useSR,
		useStopOrder:        useStopOrder,
		stopOrderDiff:       stopOrderDiff,
		srDiffPercent:       srDiffPercent,
		precision:           precision,
		interval:            interval,
		symbol:              symbol,
		buyPatterns:         buyPatterns,
		sellPatterns:        sellPatterns,
		afterCommission:     100 - commission,
		count:               count,
		initialCount:        count,
		pricesBuffer:        pricesBuffer,
		format:              format,
		needBuy:             true,
		numbersAfter:        numbersAfter,
		mode:                mode,
		sellCandlesPatterns: sellCandlesPatterns,
		buyCandlesPatterns:  buyCandlesPatterns,
		dryRun:              dryRun,
		candleTimeBuffer:    candleTimeBuffer,
		service:             service,
		klineIntervals:      klineIntervals,
	}
}

func (w *Watcher) GetName() string {
	return w.name
}

func (w *Watcher) analise() {
	if w.mode == Candle {
		w.candleAnalise()
		return
	}
	w.livePriceAnalyse()
}

func (w *Watcher) candleAnalise() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	price, err := w.service.GetPrice(w.symbol)

	if err != nil {
		log.Println("Can't get live price for", w.symbol, err.Error())
		return
	}

	if len(w.klineIntervals) == 0 {
		log.Println("No intervals set for", w.symbol)
		return
	}

	actions := make([]patterns.Action, len(w.klineIntervals))

	var eg errgroup.Group

	for index, v := range w.klineIntervals {
		i := index
		interval := v
		eg.Go(func() error {
			duration, ok := services.CandleIntervalToDuration[interval]
			if !ok {
				return errWrongInterval
			}
			candles, errCandles := w.service.GetCandles(w.symbol, interval, time.Now().Add(-1*time.Duration(w.pricesBuffer)*duration), time.Now())
			if errCandles != nil {
				return errCandles
			}
			action := patterns.CandleMultiPatterns(candles, w.sellCandlesPatterns, w.buyCandlesPatterns, w.useSR, price, w.srDiffPercent)
			actions[i] = action
			return nil
		})
	}

	err = eg.Wait()

	if err != nil {
		log.Println("Can't get live price for", w.symbol, err.Error())
		return
	}

	w.doAction(utils.GetAction(actions), price)
}

func (w *Watcher) doAction(action patterns.Action, livePrice float64) {
	if action == patterns.Hold {
		w.lastAction = patterns.Hold
		return
	}
	if action == patterns.Clean {
		w.prices = make([]*patterns.PriceRow, 0)
		return
	}

	prevSum := w.count * w.lastPrice
	if w.lastAction != patterns.Buy && action == patterns.Buy && w.needBuy {
		var buyPrice, count float64
		var errBuy error
		if !w.dryRun {
			buyPrice, count, errBuy = w.service.Buy(w.symbol, w.count, w.format)
		} else {
			buyPrice = livePrice
			count = w.count
			errBuy = nil
		}

		if errBuy != nil {
			log.Println("Can't buy", w.symbol, errBuy.Error())
			return
		}
		w.needBuy = false
		w.lastPrice = buyPrice
		w.count = utils.ToFixed(count*w.afterCommission/100, w.numbersAfter)
		w.lastAction = patterns.Buy
		log.Println("Successfully buy", w.symbol, "new count", w.count, "price", w.lastPrice)

		if w.useStopOrder {
			orderPrice := utils.ToFixed(w.lastPrice-w.stopOrderDiff, w.precision)
			var orderResp *binance.CreateOrderResponse
			var errCreateOrder error

			if !w.dryRun {
				orderResp, errCreateOrder = w.service.CreateStopOrder(w.symbol, w.count, w.format, orderPrice)
			} else {
				orderResp = &binance.CreateOrderResponse{OrderID: -1}
				errCreateOrder = nil
			}

			if errCreateOrder != nil {
				log.Println("Can't create for stop order", w.symbol, "for price", orderPrice, "err", errCreateOrder, "retrying")
				apiErr, ok := errCreateOrder.(*common.APIError)
				if ok && utils.Int64InArray(apiErr.Code, userErrors) && apiErr.Message == "Stop price would trigger immediately." {
					newOrderPrice := utils.ToFixed(w.lastPrice-w.stopOrderDiff*2, w.precision)
					orderResp, errCreateOrder = w.service.CreateStopOrder(w.symbol, w.count, w.format, newOrderPrice)
					if errCreateOrder != nil {
						log.Println("Can't create for stop order", w.symbol, "for price", orderPrice, "err", errCreateOrder, "retry not help")
						return
					}
				}

			}

			w.stopOrderId = &orderResp.OrderID
			w.lastStopOrderPrice = &orderPrice
			w.orderWatcher.SetOrderId(w.stopOrderId)
			log.Println("Created stop order", w.symbol, "for price", orderPrice)
		}
		return
	}

	estimatedProfit := livePrice*w.count*w.afterCommission/100 - w.count*w.lastPrice

	if !w.needBuy && w.lastAction != patterns.Sell && (action == patterns.Sell && (estimatedProfit > 0 || (w.lastPrice != 0 && (w.lastPrice-livePrice)/w.lastPrice > HardLimit))) {
		var sellPrice float64
		var errSell error
		if !w.dryRun {
			if w.useStopOrder && w.stopOrderId != nil {
				_, stopOrderErr := w.service.CancelOrder(w.symbol, *w.stopOrderId)

				if stopOrderErr != nil {
					log.Println("Can't cancel stop order", w.symbol, "due error", stopOrderErr.Error())
				} else {
					log.Println("Canceled stop order", w.symbol, "ID", *w.stopOrderId, "dryRUN", w.dryRun)

					sellPrice, _, errSell = w.service.Sell(w.symbol, w.count, w.format)
				}
				w.orderWatcher.SetOrderId(nil)
				w.stopOrderId = nil
			} else {
				sellPrice, _, errSell = w.service.Sell(w.symbol, w.count, w.format)
			}
		} else {
			// with dry run we just get live price
			sellPrice, errSell = w.service.GetPrice(w.symbol)
		}

		if w.dryRun && w.useStopOrder && w.lastStopOrderPrice != nil {
			if sellPrice <= *w.lastStopOrderPrice {
				sellPrice = *w.lastStopOrderPrice
			}
			log.Println("Stop order happens in dry run", w.symbol, "price", sellPrice)
			w.lastStopOrderPrice = nil
			w.orderWatcher.SetOrderId(nil)
			w.stopOrderId = nil
		}

		if errSell != nil && !w.dryRun {
			if !w.useStopOrder {
				log.Println("Can't sell", w.symbol, errSell.Error())
				return
			}
			apiErr, ok := errSell.(*common.APIError)
			if !ok {
				log.Println("Can't sell", w.symbol, "strange error", errSell.Error())
				return
			}

			if utils.Int64InArray(apiErr.Code, userErrors) {
				// It is means order already happend, so need clear current stop order ID
				log.Println("Stop order happens", w.symbol, "price", *w.lastStopOrderPrice)
				sellPrice = *w.lastStopOrderPrice
				w.lastStopOrderPrice = nil
			} else {
				log.Println("Can't sell", w.symbol, "strange error", errSell.Error())
			}
		}
		w.needBuy = true
		diff := w.count*sellPrice*w.afterCommission/100 - prevSum
		w.profit += diff
		w.count = w.initialCount
		log.Println("Successfully sell", w.symbol, "new count", w.count, "profit", diff, "estimation", livePrice*w.count*w.afterCommission/100-w.count*w.lastPrice)
		w.lastAction = patterns.Sell
		w.lastPrice = sellPrice
		return
	}
}

func (w *Watcher) livePriceAnalyse() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	price, err := w.service.GetPrice(w.symbol)
	if err != nil {
		w.prices = append(w.prices, &patterns.PriceRow{
			Date:  time.Now(),
			Price: math.MaxFloat64,
		})
		log.Println("Can't get prices for", w.symbol, err.Error())
		return
	}

	w.prices = append(w.prices, &patterns.PriceRow{
		Date:  time.Now(),
		Price: price,
	})

	if len(w.prices) > w.pricesBuffer {
		w.prices = w.prices[1:]
	}

	action := patterns.LivePriceMultiPattern(w.prices, w.sellPatterns, w.buyPatterns)

	w.doAction(action, price)
}

func (w *Watcher) IsDryRun() bool {
	return w.dryRun
}

func (w *Watcher) GetProfit() float64 {
	return w.profit
}

func (w *Watcher) OrderCallback(order *binance.Order, livePrice float64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	prevSum := w.count * w.lastPrice

	if w.dryRun {
		if w.stopOrderId != nil && order.OrderID == *w.stopOrderId && livePrice <= utils.ToFixed(w.lastPrice-w.stopOrderDiff, w.precision) {
			w.stopOrderId = nil
			w.lastPrice = *w.lastStopOrderPrice
			diff := w.count*w.lastPrice*w.afterCommission/100 - prevSum
			w.lastStopOrderPrice = nil
			w.lastAction = patterns.Sell
			w.needBuy = true
			w.profit += diff
			w.count = w.initialCount
			log.Println("Successfully sell by stop order dry run", w.symbol, "new count", w.count, "profit", diff)
			w.orderWatcher.SetOrderId(nil)
		}
		return
	}
	if w.stopOrderId != nil && order.OrderID == *w.stopOrderId && (order.Status == binance.OrderStatusTypeFilled) {
		w.stopOrderId = nil
		w.lastPrice = *w.lastStopOrderPrice
		diff := w.count*w.lastPrice*w.afterCommission/100 - prevSum
		w.lastStopOrderPrice = nil
		w.lastAction = patterns.Sell
		w.needBuy = true
		w.profit += diff
		w.count = w.initialCount
		log.Println("Successfully sell by stop order", w.symbol, "new count", w.count, "profit", diff)
		w.orderWatcher.SetOrderId(nil)
	}
}

func (w *Watcher) Run() chan struct{} {
	done := make(chan struct{}, 1)
	if w.dryRun {
		log.Println("We dont do any real actions sell/buy for", w.name)
	}

	if w.mode == "" {
		log.Println("Watcher mode not set for", w.symbol)
		return done
	}

	go func() {
		currentSecond := time.Now().Second()
		var waitTime time.Duration = 0
		diffTime := (59*time.Second - w.candleTimeBuffer) - time.Duration(currentSecond)*time.Second

		if diffTime > 0 {
			waitTime = diffTime
			log.Println("We will wait", waitTime, "until second will be", time.Duration(59)*time.Second-w.candleTimeBuffer, "current second is", currentSecond, "symbol", w.symbol)
		}

		time.Sleep(waitTime)

		log.Println("Run watcher for", w.symbol, "mode is", w.mode, "interval", w.interval, "use stopOrder", w.useStopOrder, "use SR", w.useSR)
		watcherInterval := time.Second * 5
		if w.dryRun {
			watcherInterval = time.Second * 2
		}
		if w.useStopOrder {
			ch := order.New(w.symbol, w.service, w.OrderCallback, w.dryRun, watcherInterval)
			w.orderWatcher = ch
			w.closerOrderWatcher = ch.Run()
		} else {
			w.closerOrderWatcher = make(chan struct{}, 1)
		}

		ticker := time.NewTicker(w.interval)
		go func() {
			for {
				w.analise()
				select {
				case <-done:
					w.closerOrderWatcher <- struct{}{}
					ticker.Stop()
					log.Println("Stop watcher for", w.symbol)
					return
				case <-ticker.C:
					continue
				}
			}
		}()
	}()

	return done
}
