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

	banks, err := fetchAllBanksByCountry(token, "IT")

	log.Println(token.AccessToken)

	for i := range banks {
		log.Println(banks[i])
	}
}
