package service

import (
	"FinHub/model"
	"FinHub/repository"
	"encoding/json"
	"log"
	"time"
)

type FinancialHubService struct {
	goCardlessApiService   *GoCardlessApiService
	financialHubRepository *repository.FinancialHubRepository
}

func NewFinancialHubService(goCardlessApiService *GoCardlessApiService, financialHubRepository *repository.FinancialHubRepository) *FinancialHubService {
	return &FinancialHubService{goCardlessApiService: goCardlessApiService, financialHubRepository: financialHubRepository}
}

func (f *FinancialHubService) GetTokenByUserId(id int) (*model.Token, error) {
	token, err := f.financialHubRepository.GetToken(id)

	if err != nil {
		log.Println("No token found, requesting a new one.")
		err = f.goCardlessApiService.GetNewToken()
		if err != nil {
			log.Println("Failed to get new token:", err)
			return nil, err
		}
		token, _ = f.financialHubRepository.GetToken(id)
	}

	if token.AccessExpires.Compare(time.Now()) == -1 {
		log.Println("Token expired, requesting a new one.")
		err = f.goCardlessApiService.GetNewToken()
		if err != nil {
			log.Println("Failed to get new token:", err)
			return nil, err
		}
		token, _ = f.financialHubRepository.GetToken(id)
	}

	return token, nil
}

func (f *FinancialHubService) GetBankById(token *model.Token, id string) (*model.Bank, error) {
	bank, err := f.goCardlessApiService.GetBankById(token, id)

	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (f *FinancialHubService) GetAgreementByBankAndUserId(token *model.Token, institutionId string, userId int) (*model.Agreement, error) {
	agreement, err := f.financialHubRepository.GetAgreement(token, institutionId, userId)

	if err != nil {
		log.Println("Agreement not found, requesting a new one")
		err := f.goCardlessApiService.CreateUserAgreement(institutionId, token, userId)
		if err != nil {
			return nil, err
		}
		agreement, err = f.financialHubRepository.GetAgreement(token, institutionId, userId)
		if err != nil {
			return nil, err
		}
	}

	var accessScopeConverted []string
	err = json.Unmarshal([]byte(agreement.AccessScope), &accessScopeConverted)

	agreementConverted := model.Agreement{
		Id:                 agreement.Id,
		UserId:             agreement.UserId,
		AccessToken:        agreement.AccessToken,
		Created:            agreement.Created,
		InstitutionId:      agreement.InstitutionId,
		MaxHistoricalDays:  agreement.MaxHistoricalDays,
		AccessValidForDays: agreement.AccessValidForDays,
		AccessScope:        accessScopeConverted,
		Accepted:           agreement.Accepted,
	}

	return &agreementConverted, nil
}

func (f *FinancialHubService) GetRequisitionsByAgreement(token *model.Token, agreement *model.Agreement) (*model.Requisition, error) {
	requisitions, err := f.financialHubRepository.GetRequisition(agreement.Id)

	if err != nil {
		log.Println("Requisition not found, requesting a new one")
		err := f.goCardlessApiService.CreateRequisition(agreement.InstitutionId, agreement.Id, token)
		if err != nil {
			return nil, err
		}
		requisitions, err = f.financialHubRepository.GetRequisition(agreement.Id)
		if err != nil {
			return nil, err
		}
	}

	return requisitions, nil
}

func (f *FinancialHubService) AuthorizeRequisition(token *model.Token, requisition *model.Requisition) error {
	err := f.goCardlessApiService.AuthorizeRequisition(requisition, token)

	if err != nil {
		return err
	}

	return nil
}

func (f *FinancialHubService) FetchUserAccountsByBank(requisitionId string, token *model.Token) (*model.ListAccountsResponse, error) {
	accounts, err := f.goCardlessApiService.FetchUserAccountsByBank(requisitionId, token)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (f *FinancialHubService) FetchAccountBalance(accountId string, token *model.Token) (*model.Balances, error) {
	balance, err := f.goCardlessApiService.FetchAccontBalance(accountId, token)

	if err != nil {
		return nil, err
	}

	return balance, nil
}
