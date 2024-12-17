package main

import (
	"log"
)

func main() {
	token, err := getNewToken()

	if err != nil {
		log.Println(err.Error())
		return
	}

	//banks, err := fetchAllBanksByCountry(token, "IT")

	bank, err := fetchBankById(token, "INTESA_SANPAOLO_BCITITMMXXX")

	agreement, err := CreateUserAgreement(bank.Id, token)

	link, err := CreateLink(bank.Id, agreement.Id, token)

	accounts, err := FetchUserAccountsByBank(link.ID, token)

	balance, err := FetchAccontBalance(agreement.Id, token)
	log.Println(token.AccessToken)

	log.Println(bank.Name)

	log.Println(agreement.Id)

	log.Println(link.Link)

	log.Println(accounts.Accounts)

	log.Println(balance.Balances[0].BalanceAmount.Amount)
}
