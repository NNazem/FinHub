package goCardlessApi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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

type Agreement struct {
	Id                 string    `json:"id"`
	Created            time.Time `json:"created"`
	InstitutionId      string    `json:"institution_id"`
	MaxHistoricalDays  int       `json:"max_historical_days"`
	AccessValidForDays int       `json:"access_valid_for_days"`
	AccessScope        []string  `json:"access_scope"`
	Accepted           time.Time `json:"accepted"`
}

type Requisition struct {
	ID            string   `json:"id"`
	Redirect      string   `json:"redirect"`
	InstitutionID string   `json:"institution_id"`
	Agreement     string   `json:"agreement"`
	Accounts      []string `json:"accounts"`
	UserLanguage  string   `json:"user_language"`
	Link          string   `json:"link"`
}

type ListAccountsResponse struct {
	RequisitionID string   `json:"id"`
	Status        string   `json:"status"`
	Agreements    string   `json:"agreements"`
	Accounts      []string `json:"accounts"`
	Reference     string   `json:"reference"`
}

type Balances struct {
	Balances []Balance `json:"balances"`
}

type Balance struct {
	BalanceAmount BalanceAmount `json:"balanceAmount"`
	BalanceType   string        `json:"balanceType"`
	ReferenceDate string        `json:"referenceDate"`
}

type BalanceAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type GoCardlessClient struct {
	Token *Token
}

func NewGoCardlessClient() (*GoCardlessClient, error) {
	token, err := getNewToken()
	if err != nil {
		return nil, err
	}

	return &GoCardlessClient{
		Token: token,
	}, nil
}

func getNewToken() (*Token, error) {
	query := "https://bankaccountdata.gocardless.com/api/v2/token/new/"
	body := struct {
		SecretID  string `json:"secret_id"`
		SecretKey string `json:"secret_key"`
	}{
		SecretID:  os.Getenv("SECRET_ID"),
		SecretKey: os.Getenv("SECRET_KEY"),
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

func FetchBankById(token *Token, id string) (*Bank, error) {
	query := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/institutions/%s/", id)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var bank Bank

	err = json.NewDecoder(res.Body).Decode(&bank)
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func CreateUserAgreement(institutionId string, token *Token) (*Agreement, error) {
	query := "https://bankaccountdata.gocardless.com/api/v2/agreements/enduser/"
	body := struct {
		InstitutionId      string   `json:"institution_id"`
		MaxHistoricalDays  int      `json:"max_historical_days"`
		AccessValidForDays int      `json:"access_valid_for_days"`
		AccessScope        []string `json:"access_scope"`
	}{
		InstitutionId:      institutionId,
		MaxHistoricalDays:  90,
		AccessValidForDays: 90,
		AccessScope:        []string{"balances", "details", "transactions"},
	}

	out, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", query, bytes.NewBuffer(out))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(context.Background())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var agreement Agreement

	err = json.NewDecoder(res.Body).Decode(&agreement)

	if err != nil {
		return nil, err
	}

	return &agreement, nil
}

func CreateLink(institutionid, agreement string, token *Token) (*Requisition, error) {
	query := "https://bankaccountdata.gocardless.com/api/v2/requisitions/"
	body := struct {
		Redirect      string `json:"redirect"`
		InstitutionID string `json:"institution_id"`
		Agreement     string `json:"agreement"`
		//Reference     string `json:"reference"`
		UserLanguage string `json:"user_language"`
	}{
		Redirect:      "http://www.yourwebpage.com",
		InstitutionID: institutionid,
		Agreement:     agreement,
		//Reference:     "123455", //it has to be unique or it will fail
		UserLanguage: "EN",
	}

	out, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", query, bytes.NewReader(out))

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	var requisition Requisition

	err = json.NewDecoder(res.Body).Decode(&requisition)

	if err != nil {
		return nil, err
	}

	return &requisition, nil
}

func FetchUserAccountsByBank(requisitionId string, token *Token) (*ListAccountsResponse, error) {
	query := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/requisitions/%s/", requisitionId)

	req, err := http.NewRequest("GET", query, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	var accounts ListAccountsResponse

	err = json.NewDecoder(res.Body).Decode(&accounts)

	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

func FetchAccountBalance(accountId string, token *Token) (*Balances, error) {
	query := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/accounts/%s/balances/", accountId)

	req, err := http.NewRequest("GET", query, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	var balance Balances

	err = json.NewDecoder(res.Body).Decode(&balance)

	if err != nil {
		return nil, err
	}

	return &balance, nil
}
