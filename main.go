package main

import (
	"FinHub/repository"
	"FinHub/service"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("properties.env")

	if err != nil {
		log.Println(err.Error())
		return
	}

	db, err := repository.InitDb()

	FinancialHubRepository := &repository.FinancialHubRepository{Db: db}

	GoCardlessApiService, err := service.NewGocardlessApiClient(FinancialHubRepository)

	FinancialHubService := service.NewFinancialHubService(GoCardlessApiService, FinancialHubRepository)

	token, err := FinancialHubService.GetTokenByUserId(1)

	agreement, err := FinancialHubService.GetAgreementByBankAndUserId(token, "INTESA_SANPAOLO_BCITITMMXXX", 1)

	log.Println(agreement.Created)
	log.Println(agreement.AccessScope)
}
