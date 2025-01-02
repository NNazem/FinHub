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
		err = f.goCardlessApiService.DeleteToken(id)
		if err != nil {
			log.Println("Failed to delete token:", err)
			return nil, err
		}
		err = f.goCardlessApiService.GetNewToken()
		if err != nil {
			log.Println("Failed to get new token:", err)
			return nil, err
		}
		token, _ = f.financialHubRepository.GetToken(id)
	}

	return token, nil
}

func (f *FinancialHubService) GetBankById(id string) (*model.Bank, error) {
	bank, err := f.goCardlessApiService.GetBankById(id)

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

func (f *FinancialHubService) GetUserAccountsByBank(requisitionId string, token *model.Token) (*model.ListAccountsResponse, error) {
	accounts, err := f.goCardlessApiService.FetchUserAccountsByBank(requisitionId, token)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (f *FinancialHubService) GetAccountBalance(accountId string, token *model.Token) (*model.Balances, error) {
	balance, err := f.goCardlessApiService.FetchAccontBalance(accountId, token)

	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (f *FinancialHubService) GetUserTotalBalance(userId int) (float32, error) {
	totalBalance, err := f.financialHubRepository.GetBalanceByUserId(userId)

	if err != nil {
		return 0, err
	}

	return totalBalance, nil
}

func (f *FinancialHubService) GetUserAccounts(userId string) ([]model.Account, error) {
	accounts, err := f.financialHubRepository.GetAccountsByUserId(userId)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (f *FinancialHubService) GetAccountTransactions(userId int, accountId string) ([]model.TransactionResponse, error) {

	//transactions, err := f.goCardlessApiService.FetchAccountTransactions(accountId, token)
	transactions, err := f.financialHubRepository.GetAccountTransaction(accountId)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (f *FinancialHubService) GetUserTransactions(userId int) ([]model.TransactionResponse, error) {
	transactions, err := f.financialHubRepository.GetUserTransaction(userId)

	if err != nil {
		return nil, err
	}

	return transactions, nil

}

func (f *FinancialHubService) GetUserTransactionsByMonths(userId int, months int) ([]model.TransactionResponse, error) {
	transactions, err := f.financialHubRepository.GetUserTransactionsByMonths(userId, months)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}
