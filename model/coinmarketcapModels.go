package model

type CoinResponse struct {
	Status Status              `json:"status"`
	Data   map[string]CoinData `json:"data"`
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
	Notice       string `json:"notice"`
}

type CoinData struct {
	Id                            int       `json:"id"`
	Name                          string    `json:"name"`
	Symbol                        string    `json:"symbol"`
	Slug                          string    `json:"slug"`
	NumMarketPairs                int       `json:"num_market_pairs"`
	DateAdded                     string    `json:"date_added"`
	Tags                          []string  `json:"tags"`
	MaxSupply                     int       `json:"max_supply"`
	CirculatingSupply             float64   `json:"circulating_supply"`
	TotalSupply                   float64   `json:"total_supply"`
	IsActive                      int       `json:"is_active"`
	InfiniteSupply                bool      `json:"infinite_supply"`
	Platform                      *Platform `json:"platform"`
	CmcRank                       int       `json:"cmc_rank"`
	IsFiat                        int       `json:"is_fiat"`
	SelfReportedCirculatingSupply *float64  `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         *float64  `json:"self_reported_market_cap"`
	TvlRatio                      *float64  `json:"tvl_ratio"`
	LastUpdated                   string    `json:"last_updated"`
	Quote                         CoinQuote `json:"quote"`
}

type Platform struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

type CoinQuote struct {
	USD CoinQuoteCurrency `json:"USD"`
}

type CoinQuoteCurrency struct {
	Price                 float64  `json:"price"`
	Volume24h             float64  `json:"volume_24h"`
	VolumeChange24h       float64  `json:"volume_change_24h"`
	PercentChange1h       float64  `json:"percent_change_1h"`
	PercentChange24h      float64  `json:"percent_change_24h"`
	PercentChange7d       float64  `json:"percent_change_7d"`
	PercentChange30d      float64  `json:"percent_change_30d"`
	PercentChange60d      float64  `json:"percent_change_60d"`
	PercentChange90d      float64  `json:"percent_change_90d"`
	MarketCap             float64  `json:"market_cap"`
	MarketCapDominance    float64  `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64  `json:"fully_diluted_market_cap"`
	Tvl                   *float64 `json:"tvl"`
	LastUpdated           string   `json:"last_updated"`
}
