package freqtrade

import (
	"fmt"
	"time"
)

type Profit struct {
	RealizedCryptoProfit float64 `json:"profit_closed_coin"`
	RealizedFiatProfit   float64 `json:"profit_closed_fiat"`
	// Percent is number from 0 to 1 represent percentage of profit from start balance
	RealizedPercentProfit float64 `json:"profit_closed_ratio"`

	UnrealizedCryptoProfit float64 `json:"profit_all_coin"`
	UnrealizedFiatProfit   float64 `json:"profit_all_fiat"`
	// Percent is number from 0 to 1 represent percentage of profit from start balance
	UnrealizedPercentProfit float64 `json:"profit_all_ratio"`

	TotalTrades  int `json:"trade_count"`
	ClosedTrades int `json:"closed_trade_count"`
	WinTrades    int `json:"winning_trades"`
	LossTrades   int `json:"losing_trades"`

	// first trade timestamp (millisecond)
	FirstTradeTimestamp int64 `json:"first_trade_timestamp"`

	// latest trade timestamp (millisecond)
	LastTradeTimestamp int64 `json:"latest_trade_timestamp"`

	// format 00:00:00
	AverageDuration string `json:"avg_duration"`

	BestPair string  `json:"best_pair"`
	BestRate float64 `json:"best_rate"`
}

func (p *Profit) GetAverageDuration() time.Duration {
	var h, m, s int
	n, err := fmt.Sscanf(p.AverageDuration, "%d:%d:%d", &h, &m, &s)
	if err != nil || n != 3 {
		return -1
	}

	var second = (h * 3600) + (m * 60) + s
	return time.Duration(second) * time.Second
}

func EmptyProfit() *Profit {
	return &Profit{
		UnrealizedCryptoProfit:  0,
		UnrealizedFiatProfit:    0,
		UnrealizedPercentProfit: 0,
		RealizedCryptoProfit:    0,
		RealizedFiatProfit:      0,
		RealizedPercentProfit:   0,

		TotalTrades:  0,
		ClosedTrades: 0,
		WinTrades:    0,
		LossTrades:   0,

		FirstTradeTimestamp: 0,
		LastTradeTimestamp:  0,
		AverageDuration:     "00:00:00",

		BestPair: "",
		BestRate: 0,
	}
}

func NewProfit(conn *Connection) *Profit {
	if profit, err := FetchProfit(conn); err == nil {
		return profit
	}
	return EmptyProfit()
}

func FetchProfit(conn *Connection) (*Profit, error) {
	var name = API_PROFIT
	if data, err := conn.Cache(name, conn.ExpireAt(name), func() (interface{}, error) {
		return GetProfit(conn)
	}); err == nil {
		return data.(*Profit), nil
	} else {
		return nil, err
	}
}

func GetProfit(conn *Connection) (*Profit, error) {
	var target = new(Profit)
	var err = GetConnector(conn, API_PROFIT, &target)
	return target, err
}
