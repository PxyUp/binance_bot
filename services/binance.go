package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type CandleInterval string

const (
	Min      CandleInterval = "1m"
	ThreeMin CandleInterval = "3m"
	FiveMin  CandleInterval = "5m"
	Hour     CandleInterval = "1h"
	Day      CandleInterval = "1d"
)

var (
	CandleIntervalToDuration = map[CandleInterval]time.Duration{
		Min:      time.Minute,
		ThreeMin: time.Minute * 3,
		FiveMin:  time.Minute * 5,
		Hour:     time.Hour,
		Day:      time.Hour * 24,
	}
)

type Binance struct {
	ApiKey     string  `yaml:"apiKey"`
	Secret     string  `yaml:"secret"`
	Commission float64 `yaml:"commission"`
	Debug      bool    `yaml:"debug"`
}

var (
	ErrWrongCandle   = errors.New("wrong candle")
	ErrTokenNotExist = errors.New("no token")
	Commission       = 0.1
	requestTimeout   = 5 * time.Second
)

type Service struct {
	commission float64
	client     *binance.Client
}

func New(config *Binance) *Service {
	client := binance.NewClient(config.ApiKey, config.Secret)
	client.Debug = config.Debug
	return &Service{
		commission: config.Commission,
		client:     client,
	}
}

func (s *Service) GetCommission() float64 {
	return s.commission
}

func (s *Service) CreateStopOrder(symbol string, count float64, format string, stopPrice float64) (*binance.CreateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.client.NewCreateOrderService().Symbol(symbol).Type(binance.OrderTypeStopLossLimit).Side(binance.SideTypeSell).Price(fmt.Sprintf("%g", stopPrice)).TimeInForce(binance.TimeInForceTypeGTC).StopPrice(fmt.Sprintf("%g", stopPrice)).Quantity(fmt.Sprintf(format, count)).Do(ctx, binance.WithRecvWindow(1000))
}

func (s *Service) GetOrder(symbol string, orderId int64) (*binance.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.client.NewGetOrderService().Symbol(symbol).OrderID(orderId).Do(ctx)
}

func (s *Service) CancelOrder(symbol string, orderId int64) (*binance.CancelOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	return s.client.NewCancelOrderService().Symbol(symbol).OrderID(orderId).Do(ctx)
}

func (s *Service) Sell(symbol string, count float64, format string) (float64, float64, error) {
	log.Println("Try sell", symbol, "count is", fmt.Sprintf(format, count))
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	res, err := s.client.NewCreateOrderService().Symbol(symbol).Type(binance.OrderTypeMarket).Side(binance.SideTypeSell).Quantity(fmt.Sprintf(format, count)).Do(ctx, binance.WithRecvWindow(1000))
	if err != nil {
		return 0, 0, err
	}
	if len(res.Fills) == 0 {
		return 0, 0, errors.New("not filled")
	}
	price, errPrice := strconv.ParseFloat(res.Fills[len(res.Fills)-1].Price, 64)

	if err != errPrice {
		return 0, 0, errPrice
	}

	newCount, errCount := strconv.ParseFloat(res.ExecutedQuantity, 64)
	if err != errCount {
		return 0, 0, errCount
	}

	return price, newCount, nil

}

func (s *Service) GetCandles(symbol string, interval CandleInterval, startTime time.Time, endTime time.Time) ([]*Candle, error) {
	log.Println("Get Candle", symbol, "interval", interval)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	defer cancel()

	res, err := s.client.NewKlinesService().Symbol(symbol).Interval(string(interval)).StartTime(startTime.Unix() * 1000).EndTime(endTime.Unix() * 1000).Do(ctx)

	if err != nil {
		return nil, err
	}
	var candles []*Candle

	for _, v := range res {
		candle := klineToCandle(v)

		if candle == nil {
			return nil, ErrWrongCandle
		}

		candles = append(candles, candle)
	}

	return candles, nil
}

func (s *Service) Buy(symbol string, count float64, format string) (float64, float64, error) {
	log.Println("Try buy", symbol, "count is", fmt.Sprintf(format, count))
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	res, err := s.client.NewCreateOrderService().Symbol(symbol).Type(binance.OrderTypeMarket).Side(binance.SideTypeBuy).Quantity(fmt.Sprintf(format, count)).Do(ctx, binance.WithRecvWindow(1000))

	if err != nil {
		return 0, 0, err
	}

	if len(res.Fills) == 0 {
		return 0, 0, errors.New("not filled")
	}

	price, errPrice := strconv.ParseFloat(res.Fills[len(res.Fills)-1].Price, 64)

	if err != errPrice {
		return 0, 0, errPrice
	}

	newCount, errCount := strconv.ParseFloat(res.ExecutedQuantity, 64)

	if err != errCount {
		return 0, 0, errCount
	}

	return price, newCount, nil
}

func (s *Service) GetPrice(symbol string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	res, err := s.client.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return 0, ErrTokenNotExist
	}

	price, errPrice := strconv.ParseFloat(res[0].Price, 64)

	if errPrice != nil {
		return 0, errPrice
	}

	return price, nil
}
