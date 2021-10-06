package services

import (
	"fmt"
	"time"
)

type CandleColor string

type CandleSize int8

const (
	Green CandleColor = "Green"
	Red   CandleColor = "Red"
	Grey  CandleColor = "Grey"

	Zero CandleSize = 0
	S    CandleSize = 1
	M    CandleSize = 2
	L    CandleSize = 3
)

type Candle struct {
	OpenTime  time.Time
	CloseTime time.Time
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64
}

func (c *Candle) GetColor() CandleColor {
	if c.Open-c.Close > 0 {
		return Red
	}

	if c.Open-c.Close < 0 {
		return Green
	}

	return Grey
}

func (c *Candle) GetAvgPrice() float64 {
	return (c.Open + c.Close) / 2
}

func (c *Candle) ToString() string {
	return fmt.Sprintf("%s%d", c.GetColor(), c.GetSize())
}

func CandlesToString(canldes []*Candle) string {
	s := ""

	for _, v := range canldes {
		s += v.ToString()
	}

	return s
}

func (c *Candle) GetSize() CandleSize {
	diff := (c.Close - c.Open) / ((c.Open + c.Close) / 2)

	if c.GetColor() == Red {
		diff *= -1
	}

	if diff < 0.00005 {
		return Zero
	}

	if diff < 0.002 {
		return S
	}

	if diff < 0.004 {
		return M
	}

	return L
}
