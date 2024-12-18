package main

import (
	"FinHub/goCardlessApi"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err.Error())
		return
	}

	GoCardlessApiService, err := goCardlessApi.NewGocardlessApiClient()

	//banks, err := fetchAllBanksByCountry(token, "IT")

	bank, err := goCardlessApi.FetchBankById(GoCardlessApiService.Token, "INTESA_SANPAOLO_BCITITMMXXX")

	agreement, err := goCardlessApi.CreateUserAgreement(bank.Id, GoCardlessApiService.Token)

	link, err := goCardlessApi.CreateLink(bank.Id, agreement.Id, GoCardlessApiService.Token)

	accounts, err := goCardlessApi.FetchUserAccountsByBank(link.ID, GoCardlessApiService.Token)

	balance, err := goCardlessApi.FetchAccontBalance(agreement.Id, GoCardlessApiService.Token)
	log.Println(GoCardlessApiService.Token.AccessToken)

	log.Println(bank.Name)

	log.Println(agreement.Id)

	log.Println(link.Link)

	log.Println(accounts.Accounts)

	log.Println(balance.Balances[0].BalanceAmount.Amount)
}
