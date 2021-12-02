package order

import (
	"log"
	"sync"
	"time"

	"github.com/PxyUp/binance_bot/services"
	"github.com/adshao/go-binance/v2"
)

type OrderWatcher struct {
	service        *services.Service
	symbol         string
	currentOrderId *int64
	callback       Callback
	interval       time.Duration
	dryRun         bool

	mutex sync.Mutex
}

type Callback func(order *binance.Order, livePrice float64)

func New(symbol string, service *services.Service, callback Callback, dryRun bool, interval time.Duration) *OrderWatcher {
	return &OrderWatcher{
		dryRun:   dryRun,
		symbol:   symbol,
		service:  service,
		callback: callback,
		interval: interval,
	}
}

func (w *OrderWatcher) SetOrderId(orderId *int64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w == nil {
		return
	}
	w.currentOrderId = orderId
}

func (w *OrderWatcher) checkOrder() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.currentOrderId == nil {
		return
	}

	if w.dryRun {
		price, err := w.service.GetPrice(w.symbol)
		if err != nil {
			return
		}
		go func(livePrice float64) {
			w.callback(&binance.Order{OrderID: -1}, livePrice)
		}(price)

		return
	}

	if *w.currentOrderId <= 0 {
		return
	}
	res, err := w.service.GetOrder(w.symbol, *w.currentOrderId)
	if err != nil {
		log.Println("Order watcher return error", err.Error())
		return
	}

	go func(order *binance.Order) {
		w.callback(order, 0)
	}(res)

}

func (w *OrderWatcher) Run() chan struct{} {
	ticker := time.NewTicker(w.interval)
	done := make(chan struct{}, 1)
	log.Println("Run watcher order for", w.symbol)
	go func() {
		for {
			w.checkOrder()
			select {
			case <-done:
				ticker.Stop()
				log.Println("Stop watcher for order")
				return
			case <-ticker.C:
				continue
			}
		}
	}()

	return done
}
