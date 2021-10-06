package patterns

import (
	"errors"
	"github.com/PxyUp/binance_bot/services"
)

type CandlePattern struct {
	Color   services.CandleColor
	MinSize services.CandleSize
}

var (
	errSmallLen = errors.New("small length of candles")
	errWrongSP  = errors.New("cant get SP")
)

func GetSR(candles []*services.Candle) (float64, float64, error) {
	minsCount := 0
	maxsCount := 0
	minsSum := 0.0
	maxsSum := 0.0

	if len(candles) < 3 {
		return -1, -1, errSmallLen
	}

	for i := 1; i < len(candles)-1; i++ {
		current := candles[i].GetAvgPrice()

		if current < candles[i+1].GetAvgPrice() && current < candles[i-1].GetAvgPrice() {
			minsCount += 1
			minsSum += current
			continue
		}

		if current > candles[i+1].GetAvgPrice() && current > candles[i-1].GetAvgPrice() {
			maxsCount += 1
			maxsSum += current
			continue
		}
	}

	if minsCount == 0 && maxsCount == 0 {
		return -1, -1, errWrongSP
	}

	return minsSum / float64(minsCount), maxsSum / float64(maxsCount), nil
}

func CandleMultiPatterns(candles []*services.Candle, sellPatterns [][]*CandlePattern, buyPatterns [][]*CandlePattern, useSR bool, livePrice, diffPercent float64) Action {
	countBuy := len(buyPatterns)
	countSell := len(sellPatterns)

	max := countBuy

	if countBuy < countSell {
		max = countSell
	}

	for i := 0; i < max; i++ {
		var left, right []*CandlePattern

		if i < len(buyPatterns) {
			right = buyPatterns[i]
		}

		if i < len(sellPatterns) {
			left = sellPatterns[i]
		}

		action := CandleSinglePattern(candles, left, right, useSR, livePrice, diffPercent)

		if action != Hold {
			return action
		}
	}

	return Hold
}

func CandleSinglePattern(candles []*services.Candle, sellPattern []*CandlePattern, buyPattern []*CandlePattern, useSR bool, livePrice float64, diffPercent float64) Action {
	pricesLength := len(candles)

	if pricesLength < len(sellPattern) || pricesLength < len(buyPattern) {
		return Hold
	}

	patternLength := len(sellPattern)

	if len(buyPattern) > len(sellPattern) {
		patternLength = len(buyPattern)
	}

	isSell := true

	isBuy := true

	slice := candles[pricesLength-patternLength:]

	for i := 0; i < patternLength; i++ {
		current := slice[i]
		if isSell == false && isBuy == false {
			return Hold
		}
		if i < len(buyPattern) && current.GetColor() == buyPattern[i].Color && current.GetSize() >= buyPattern[i].MinSize && isBuy {
			isBuy = true
			continue
		} else {
			isBuy = false
		}
		if i < len(sellPattern) && current.GetColor() == sellPattern[i].Color && current.GetSize() >= sellPattern[i].MinSize && isSell {
			isSell = true
			continue
		} else {
			isSell = false
		}
	}

	if isBuy && !useSR {
		return Buy
	}

	if isSell && !useSR {
		return Sell
	}

	if !useSR {
		return Hold
	}

	// Here lets check with Support Resistance
	s, r, err := GetSR(candles)

	if err != nil {
		return Hold
	}

	// Here need buy until price lower than resistance
	if isBuy && livePrice < r+diffPercent/100*r {
		return Buy
	}

	// Here need sell until price more than support
	if isSell && livePrice > s-diffPercent/100*s {
		return Sell
	}

	return Hold
}
