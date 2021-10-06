package patterns

import "github.com/PxyUp/binance_bot/services"

var (
	DefaultBuy = []int8{
		0,
		-1,
		-1,
		-1,
		-1,
		1,
		1,
	}
	BuyWithGrow = []int8{
		0,
		-1,
		1,
		-1,
		-1,
		1,
		1,
	}
	SellWithGrow = []int8{
		0,
		1,
		-1,
		1,
		1,
		-1,
		-1,
	}
	DefaultSell = []int8{
		0,
		1,
		1,
		1,
		1,
		-1,
		-1,
	}

	CandleDefaultBuy = []*CandlePattern{
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
	}
	CandleBuyWithSpikeCandle = []*CandlePattern{
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.L,
		},
	}
	CandleBuyWithGrowCandle = []*CandlePattern{
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.L,
		},
		{
			Color:   services.Green,
			MinSize: services.L,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
	}
	CandleBuyWithDoubleGrowCandle = []*CandlePattern{
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.L,
		},
		{
			Color:   services.Green,
			MinSize: services.L,
		},
	}
	CandleBuyWithGrow = []*CandlePattern{
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
	}

	CandleDefaultSell = []*CandlePattern{
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
	}
	CandleSellWithSpike = []*CandlePattern{
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.L,
		},
	}
	CandleSellWithoutGrow = []*CandlePattern{
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
	}
	CandleSellWithGrowCandle = []*CandlePattern{
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.L,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
	}

	CandleSellWithGrow = []*CandlePattern{
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
	}

	AllBuy = [][]int8{
		DefaultBuy,
		BuyWithGrow,
	}

	AllSell = [][]int8{
		SellWithGrow,
		DefaultSell,
	}

	AllCandleBuy = [][]*CandlePattern{
		CandleDefaultBuy,
		CandleBuyWithSpikeCandle,
		CandleBuyWithGrowCandle,
		CandleBuyWithDoubleGrowCandle,
		CandleBuyWithGrow,
	}

	AllCandleSell = [][]*CandlePattern{
		CandleDefaultSell,
		CandleSellWithSpike,
		CandleSellWithoutGrow,
		CandleSellWithGrowCandle,
		CandleSellWithGrow,
	}

	// SHORT patterns

	CandleDefaultShortBuy = []*CandlePattern{
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
	}

	CandleDefaultShortSell = []*CandlePattern{
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Green,
			MinSize: services.S,
		},
		{
			Color:   services.Red,
			MinSize: services.S,
		},
	}
)

func UseCandles(candles ...[]*CandlePattern) [][]*CandlePattern {
	return candles
}

func UsePrice(prices ...[]int8) [][]int8 {
	return prices
}

func UseDefaults(forBuy bool, short bool) [][]*CandlePattern {
	if short {
		if forBuy {
			return UseCandles(CandleDefaultShortBuy)
		}
		return UseCandles(CandleDefaultShortSell)
	}

	if forBuy {
		return AllCandleBuy
	}

	return AllCandleSell
}
