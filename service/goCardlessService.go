package service

import (
	"FinHub/model"
	"FinHub/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type GoCardlessApiService struct {
	Token            *model.Token
	FinHubRepository *repository.FinancialHubRepository
}

func NewGocardlessApiClient(repo *repository.FinancialHubRepository) (*GoCardlessApiService, error) {

	return &GoCardlessApiService{
		FinHubRepository: repo,
	}, nil

}

func (s *GoCardlessApiService) GetNewToken() error {
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
		return err
	}

	req, err := http.NewRequest("POST", query, bytes.NewBuffer(out))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(context.Background())

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	token := &model.TokenResponse{}

	err = json.NewDecoder(res.Body).Decode(token)

	if err != nil {
		return err
	}

	accessTime := time.Now().Local().Add(time.Second * time.Duration(token.AccessExpires))
	refreshTime := time.Now().Local().Add(time.Second * time.Duration(token.RefreshExpires))

	err = s.FinHubRepository.InsertNewToken(token.AccessToken, token.Refresh, accessTime, refreshTime, 1)

	if err != nil {
		return err
	}

	return nil
}
func (s *GoCardlessApiService) GetAllBanksByCountry(country string) ([]model.Bank, error) {
	query := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/institutions/?country=%s", country)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	token, err := s.FinHubRepository.GetToken(1)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var banks []model.Bank

	err = json.NewDecoder(res.Body).Decode(&banks)
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (s *GoCardlessApiService) GetBankById(id string) (*model.Bank, error) {
	query := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/institutions/%s/", id)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	token, err := s.FinHubRepository.GetToken(1)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var bank model.Bank

	err = json.NewDecoder(res.Body).Decode(&bank)
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (s *GoCardlessApiService) CreateUserAgreement(institutionId string, token *model.Token, userId int) error {
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
		return err
	}

	req, err := http.NewRequest("POST", query, bytes.NewBuffer(out))

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(context.Background())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var agreement model.Agreement

	err = json.NewDecoder(res.Body).Decode(&agreement)

	agreement.UserId = userId
	agreement.AccessToken = token.AccessToken

	s.FinHubRepository.InsertNewAgreement(&agreement)

	if err != nil {
		return err
	}

	return nil
}

func (s *GoCardlessApiService) CreateRequisition(institutionid, agreement string, token *model.Token) error {
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
		return err
	}

	req, err := http.NewRequest("POST", query, bytes.NewReader(out))

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	var requisition model.Requisition

	err = json.NewDecoder(res.Body).Decode(&requisition)

	if err != nil {
		return err
	}

	s.FinHubRepository.InsertNewRequisition(requisition)

	return nil
}

func (s *GoCardlessApiService) FetchUserAccountsByBank(requisitionId string, token *model.Token) (*model.ListAccountsResponse, error) {
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

	var accounts model.ListAccountsResponse

	err = json.NewDecoder(res.Body).Decode(&accounts)

	for i := range accounts.Accounts {
		if s.FinHubRepository.GetAccountById(accounts.Accounts[i]) == 0 {
			s.FinHubRepository.InsertNewAccount(&model.Account{
				ID:            accounts.Accounts[i],
				Reference:     accounts.Reference,
				RequisitionID: accounts.RequisitionID,
				Status:        accounts.Status,
				BalanceAmount: "0",
			})
		}
	}

	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (s *GoCardlessApiService) FetchAccontBalance(accountId string, token *model.Token) (*model.Balances, error) {
	query := fmt.Sprintf("https://bankaccountdata.gocardless.com/api/v2/accounts/%s/balances/", accountId)

	req, err := http.NewRequest("GET", query, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	var balance model.Balances

	err = json.NewDecoder(res.Body).Decode(&balance)

	if err != nil {
		return nil, err
	}

	err = s.FinHubRepository.UpdateAccountBalance(accountId, balance.Balances[1])

	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (s *GoCardlessApiService) AuthorizeRequisition(requisition *model.Requisition, token *model.Token) error {
	query := requisition.Link

	// Open the link in the default browser
	err := exec.Command("open", query).Start()
	if err != nil {
		return fmt.Errorf("failed to open link: %v", err)
	}

	// Wait for the user to complete the authorization
	fmt.Println("Please complete the authorization in the opened browser window and press Enter to continue...")
	fmt.Scanln()

	// Check the status of the requisition after authorization
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("authorization failed with status: %s", res.Status)
	}

	return nil
}

func (s *GoCardlessApiService) DeleteToken(id int) error {
	err := s.FinHubRepository.DeleteToken(id)

	if err != nil {
		return err
	}

	return nil
}
