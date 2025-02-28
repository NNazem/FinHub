package model

import "time"

type UserFinanceProduct struct {
	Id        int     `json:"id"`
	UserId    int     `json:"user_id"`
	ProductId int     `json:"product_id"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
}

type UserCoins struct {
	Id     int     `json:"id"`
	UserId int     `json:"user_id"`
	CoinId int     `json:"coin_id"`
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type UserCoinsResponse struct {
	UserId              int     `json:"user_id"`
	Name                string  `json:"name"`
	Symbol              string  `json:"Symbol"`
	Slug                string  `json:"slug"`
	CoinMarketCapId     int     `json:"coin_market_cap_id"`
	CoinMarketCapRank   int     `json:"coin_market_cap_rank"`
	CoinMarketCapStatus int     `json:"coin_market_cap_status"`
	Amount              float64 `json:"amount"`
	Price               float64 `json:"price"`
	CurrentPrice        float64 `json:"current_price"`
	CurrentProfit       float64 `json:"current_profit"`
}

type AllCoinsResponse struct {
	Id   int    `json:"id"`
	Rank int    `json:"rank"`
	Name string `json:"name"`
}

type AmountPerCategory struct {
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
	Category   string  `json:"category"`
}

type UserAmountPerCategories struct {
	UserId            int                 `json:"user_id"`
	AmountPerCategory []AmountPerCategory `json:"amount_per_category"`
}

type AddCryptoRequest struct {
	Coin struct {
		Id   int    `json:"id"`
		Rank int    `json:"rank"`
		Name string `json:"name"`
	} `json:"coin"`
	Amount       string    `json:"amount"`
	PurchaseDate time.Time `json:"purchaseDate"`
	Price        string    `json:"price"`
}
