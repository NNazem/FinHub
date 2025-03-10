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
	repo := &repository.FinancialHubRepository{Db: db}

	coinService := service.NewCoinMarketCapService(repo)
	finHubService := service.NewFinancialHubService(repo)
	userService := service.NewUserService(repo, coinService, finHubService)

	finHubService.UserService = userService

	go func() {
		for {
			log.Println("Updating crypto data")
			err := coinService.GetCoinsHistoricalData()
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
			err := finHubService.InsertUsersPortfolioTotalValue()
			if err != nil {
				log.Println("Error saving portfolio value: " + err.Error())
			}
			log.Println("Portfolio values saved")
			time.Sleep(5 * time.Minute)
		}
	}()

	r := mux.NewRouter()
	api.NewFinancialHubApi(finHubService, coinService, userService, r).InitApi()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
