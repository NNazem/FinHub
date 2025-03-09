package main

import (
	"FinHub/api"
	"FinHub/repository"
	"FinHub/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	db, _ := repository.InitDb()

	FinancialHubRepository := &repository.FinancialHubRepository{Db: db}

	FinancialHubService := service.NewFinancialHubService(FinancialHubRepository)

	CoinmarketcapService := service.NewCoinmarketcapService(FinancialHubRepository)

	go func() {
		for {
			log.Println("Updating crypto data")
			err := CoinmarketcapService.GetCoinsHistoricalData()
			time.Sleep(10 * time.Hour)
			if err != nil {
				log.Println(err)
			}
			log.Println("crypto data updated")
		}
	}()

	go func() {
		for {
			log.Println("Saving portfolio value")
			err := FinancialHubService.InsertUsersPortfolioTotalValue()
			if err != nil {
				log.Println("Error saving portfolio value: " + err.Error())
			}
			log.Println("Portfolio values saved")
			time.Sleep(5 * time.Minute)
		}
	}()

	r := mux.NewRouter()
	api.NewFinancialHubApi(FinancialHubService, CoinmarketcapService, r).InitApi()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
