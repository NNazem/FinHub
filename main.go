package main

import (
	"FinHub/api"
	"FinHub/repository"
	"FinHub/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	/*
			token, err := FinancialHubService.GetTokenByUserId(1)

			agreement, err := FinancialHubService.GetAgreementByBankAndUserId(token, "INTESA_SANPAOLO_BCITITMMXXX", 1)

			requisition, err := FinancialHubService.GetRequisitionsByAgreement(token, agreement)

			//err = FinancialHubService.AuthorizeRequisition(token, requisition)

			accounts, err := FinancialHubService.GetUserAccountsByBank(requisition.ID, token)

			log.Println(agreement.Created)
			log.Println(agreement.AccessScope)
			log.Println(requisition.ID)
			log.Println(accounts)

			for _, _ = range accounts.Accounts {
				balance, err := FinancialHubService.GetAccountBalance(accounts.Accounts[1], token)

				if err != nil {
					log.Println(err)
				}

				log.Println(balance)
			}

		balance, err := FinancialHubService.GetUserTotalBalance(1)

		if err != nil {
			log.Println(err)
		}

		log.Println(balance)
	*/

	r := mux.NewRouter()
	api.NewFinancialHubApi(FinancialHubService, GoCardlessApiService, r).InitApi()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
