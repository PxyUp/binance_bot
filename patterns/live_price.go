package patterns

import (
	"log"
	"math"
	"time"
)

type Action int8

type PriceRow struct {
	Price float64
	Date  time.Time
}

// Buy [0 -1 -1 -1 -1 -1 1 1]
// Sell [0 1 1 1 1 1 -1 -1 -1]

func LivePriceMultiPattern(prices []*PriceRow, sellPatterns [][]int8, buyPatterns [][]int8) Action {
	countBuy := len(buyPatterns)
	countSell := len(sellPatterns)

	max := countBuy

	if countBuy < countSell {
		max = countSell
	}

	for i := 0; i < max; i++ {
		var left, right []int8

		if i < len(buyPatterns) {
			right = buyPatterns[i]
		}

		if i < len(sellPatterns) {
			left = sellPatterns[i]
		}

		action := LivePriceSinglePattern(prices, left, right)

		if action != Hold {
			return action
		}
	}

	return Hold
}

func LivePriceSinglePattern(prices []*PriceRow, sellPattern []int8, buyPattern []int8) Action {
	pricesLength := len(prices)

	if pricesLength < len(sellPattern) || pricesLength < len(buyPattern) {
		return Hold
	}

	patternLength := len(sellPattern)

	if len(buyPattern) > len(sellPattern) {
		patternLength = len(buyPattern)
	}

	if sellPattern[0] != 0 || buyPattern[0] != 0 {
		log.Println("Wrong patter, patter should start from 0")
		return Hold
	}

	isSell := true

	isBuy := true

	slice := prices[pricesLength-patternLength:]

	prev := slice[0]

	for i := 1; i < patternLength; i++ {
		diff := slice[i].Date.Sub(prev.Date)
		if diff > time.Second*2 || diff < 0 {
			log.Println("Low interval")
			return Clean
		}
		if math.MaxFloat64 == prev.Price || slice[i].Price == math.MaxFloat64 {
			return Hold
		}
		if isSell == false && isBuy == false {
			return Hold
		}
		var current int8 = 0
		if slice[i].Price > prev.Price {
			current = 1
		} else {
			if slice[i].Price < prev.Price {
				current = -1
			}
		}
		prev = slice[i]
		if i < len(buyPattern) && current == buyPattern[i] && isBuy {
			isBuy = true
			continue
		} else {
			isBuy = false
		}
		if i < len(sellPattern) && current == sellPattern[i] && isSell {
			isSell = true
			continue
		} else {
			isSell = false
		}
	}

	if isBuy {
		return Buy
	}

	if isSell {
		return Sell
	}

	return Hold
}
