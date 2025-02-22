package freqtrade

import "github.com/kamontat/fthelper/shared/caches"

type CryptoBalance struct {
	Symbol          string  `json:"currency"`
	Balance         float64 `json:"balance"`
	Free            float64 `json:"free"`
	Used            float64 `json:"used"`
	EstStakeBalance float64 `json:"est_stake"`
	StakeSymbol     string  `json:"stake"`
}

type Balance struct {
	Currencies   []*CryptoBalance
	CryptoValue  float64 `json:"total"`
	CryptoSymbol string  `json:"stake"`
	FiatValue    float64 `json:"value"`
	FiatSymbol   string  `json:"symbol"`
}

func EmptyBalance() *Balance {
	return &Balance{
		CryptoValue:  0,
		CryptoSymbol: "UKN",
		FiatValue:    0,
		FiatSymbol:   "UKN",
		Currencies:   make([]*CryptoBalance, 0),
	}
}

// Create new balance
func NewBalance(conn *Connection) *Balance {
	if balance, err := FetchBalance(conn); err == nil {
		return balance
	}
	return EmptyBalance()
}

// Fetch balance to local cache
func FetchBalance(conn *Connection) (*Balance, error) {
	var name = API_BALANCE
	if data, err := conn.Cache(name, conn.ExpireAt(name), func() (interface{}, error) {
		return GetBalance(conn)
	}); err == nil {
		return data.(*Balance), nil
	} else {
		return nil, err
	}
}

// Get balance without cache
func GetBalance(conn *Connection) (*Balance, error) {
	var target = new(Balance)
	var err = GetConnector(conn, API_BALANCE, &target)

	// if fiat value is not exist, use currency rate to est. value
	if target.CryptoValue > 0 && target.FiatValue <= 0 && caches.Global.Has(CACHE_CURRENCY_RATE) {
		var rate = caches.Global.Get(CACHE_CURRENCY_RATE)
		target.FiatValue = target.CryptoValue * rate.Data.(float64)
	}

	return target, err
}
