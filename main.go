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
			log.Println("Updating crypto data!")
			time.Sleep(10 * time.Hour)
			err := CoinmarketcapService.GetCoinsHistoricalData()
			if err != nil {
				log.Println(err)
			}
			log.Println("crypto data updated")
		}
	}()

	r := mux.NewRouter()
	api.NewFinancialHubApi(FinancialHubService, CoinmarketcapService, r).InitApi()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
