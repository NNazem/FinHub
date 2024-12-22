package service

import (
	"FinHub/model"
	"FinHub/repository"
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

	return agreement, nil
}
