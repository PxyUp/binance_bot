package utils

import (
	"github.com/PxyUp/binance_bot/patterns"
	"math"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func Int64InArray(item int64, arr []int64) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func GetAction(actions []patterns.Action) patterns.Action {
	if len(actions) == 0 {
		return patterns.Hold
	}

	first := actions[0]

	for i := 1; i < len(actions); i++ {
		if actions[i] != first {
			return patterns.Hold
		}
	}

	return first
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
