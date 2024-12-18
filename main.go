package main

import (
	"FinHub/goCardlessApi"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()

	cardlessClient, err := goCardlessApi.NewGoCardlessClient()

	if err != nil {
		log.Println(err.Error())
		return
	}

	//banks, err := fetchAllBanksByCountry(token, "IT")

	bank, err := goCardlessApi.FetchBankById(cardlessClient.Token, "INTESA_SANPAOLO_BCITITMMXXX")

	agreement, err := goCardlessApi.CreateUserAgreement(bank.Id, cardlessClient.Token)

	link, err := goCardlessApi.CreateLink(bank.Id, agreement.Id, cardlessClient.Token)

	accounts, err := goCardlessApi.FetchUserAccountsByBank(link.ID, cardlessClient.Token)

	balance, err := goCardlessApi.FetchAccountBalance(agreement.Id, cardlessClient.Token)
	log.Println(cardlessClient.Token.AccessToken)

	log.Println(bank.Name)

	log.Println(agreement.Id)

	log.Println(link.Link)

	log.Println(accounts.Accounts)

	log.Println(balance.Balances[0].BalanceAmount.Amount)
}
