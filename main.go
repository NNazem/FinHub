package main

import (
	"FinHub/api"
	"FinHub/repository"
	"FinHub/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

func main() {
	err := godotenv.Load("properties.env")

	if err != nil {
		log.Println(err.Error())
		return
	}
	db, err := repository.InitDb()

	FinancialHubRepository := &repository.FinancialHubRepository{Db: db}

	FinancialHubService := service.NewFinancialHubService(FinancialHubRepository)

	CoinmarketcapService := service.NewCoinmarketcapService(FinancialHubRepository)

	go func() {
		for {
			log.Println("updating crypto data")
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
