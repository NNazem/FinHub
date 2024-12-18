package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Token struct {
	AccessToken    string `json:"access"`
	AccessExpires  int64  `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int64  `json:"refresh_expires"`
}

type Bank struct {
	Id                    string   `json:"id"`
	Name                  string   `json:"name"`
	Bic                   string   `json:"bic"`
	TransactionTotalDays  string   `json:"transaction_total_days"`
	Countries             []string `json:"countries"`
	Logo                  string   `json:"logo"`
	MaxAccessValidForDays string   `json:"max_access_valid_for_days"`
}

func getNewToken() (*Token, error) {
	query := "https://bankaccountdata.gocardless.com/api/v2/token/new/"
	body := struct {
		SecretID  string `json:"secret_id"`
		SecretKey string `json:"secret_key"`
	}{
		SecretID:  "d721c529-1f6f-496b-ba8c-d541d5622220",
		SecretKey: "b9a42f189b60ae5ba20f267e411db1b1328963f3f1a9cd280fe6ac36ff63aa6ee8c2d9aa50a338b99dc18166d46aaef66d68a280c7d2cd3fdcfb1c842757cde5",
	}

	out, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", query, bytes.NewBuffer(out))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(context.Background())

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	token := &Token{}

	err = json.NewDecoder(res.Body).Decode(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}
func fetchAllBanksByCountry(token *Token, country string) ([]Bank, error) {
	query := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/institutions/?country=%s", country)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var banks []Bank

	err = json.NewDecoder(res.Body).Decode(&banks)
	if err != nil {
		return nil, err
	}

	return banks, nil
}
