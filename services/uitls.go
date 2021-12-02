package services

import (
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

func klineToCandle(kline *binance.Kline) *Candle {
	openTime := time.Unix(kline.OpenTime/1000, 0)
	closeTIme := time.Unix(kline.CloseTime/1000, 0)
	open, err := strconv.ParseFloat(kline.Open, 64)

	if err != nil {
		return nil
	}

	closeValue, err := strconv.ParseFloat(kline.Close, 64)

	if err != nil {
		return nil
	}

	low, err := strconv.ParseFloat(kline.Low, 64)

	if err != nil {
		return nil
	}

	high, err := strconv.ParseFloat(kline.High, 64)

	if err != nil {
		return nil
	}

	volume, err := strconv.ParseFloat(kline.Volume, 64)

	if err != nil {
		return nil
	}

	return &Candle{
		OpenTime:  openTime,
		CloseTime: closeTIme,
		Open:      open,
		Close:     closeValue,
		High:      high,
		Low:       low,
		Volume:    volume,
	}
}
