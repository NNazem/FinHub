package service

import (
	"FinHub/model"
	"FinHub/repository"
	"encoding/json"
	"log"
	"math"
	"sort"
	"strconv"
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

func (f *FinancialHubService) GetUserAmountPerTypologies(userId int) (*model.UserAmountPerCategories, error) {
	amountPerCategories, err := f.financialHubRepository.GetAmountPerTypology(userId)

	if err != nil {
		return nil, err
	}

	var totalSum float64

	for _, category := range amountPerCategories {
		totalSum += category.Amount
	}

	for i, category := range amountPerCategories {
		amountPerCategories[i].Percentage = (category.Amount / totalSum) * 100
	}

	var userAmountPerCategories model.UserAmountPerCategories

	userAmountPerCategories.UserId = userId
	userAmountPerCategories.AmountPerCategory = amountPerCategories

	return &userAmountPerCategories, nil
}

func (f *FinancialHubService) AddCoinToUser(userid string, coin model.AddCryptoRequest) error {
	err := f.financialHubRepository.AddCoinToUser(userid, coin)

	return err
}

func (s *CoinmarketcapService) GetUserCoin(userId int) ([]model.UserCoinsResponse, error) {
	var userCoinsResponse []model.UserCoinsResponse
	userCoins, err := s.FinHubRepository.GetUserCoin(userId)

	for _, coin := range userCoins {
		coinInfo, err := s.FinHubRepository.GetCoin(coin.CoinId)

		if err != nil {
			return nil, err
		}

		var userCoinResponse model.UserCoinsResponse

		userCoinResponse.UserId = coin.UserId
		userCoinResponse.Name = coinInfo.Name
		userCoinResponse.Symbol = coinInfo.Symbol
		userCoinResponse.Slug = coinInfo.Slug
		userCoinResponse.CoinMarketCapId = coinInfo.Id
		userCoinResponse.CoinMarketCapRank = coinInfo.Rank
		userCoinResponse.CoinMarketCapStatus = coinInfo.Status
		userCoinResponse.Amount = coin.Amount
		userCoinResponse.Price = coin.Price

		userCoinsResponse = append(userCoinsResponse, userCoinResponse)
	}

	var userCoinsSlugs []string

	for _, coin := range userCoinsResponse {
		userCoinsSlugs = append(userCoinsSlugs, strconv.Itoa(coin.CoinMarketCapId))
	}

	coinCurrentData, err := s.GetCoinsData(userCoinsSlugs)

	for i, coin := range userCoinsResponse {
		userCoinsResponse[i].CurrentPrice = math.Round((coinCurrentData.Data[strconv.Itoa(coin.CoinMarketCapId)].Quote.USD.Price)*100) / 100
		userCoinsResponse[i].CurrentProfit = math.Round((userCoinsResponse[i].Amount*coinCurrentData.Data[strconv.Itoa(coin.CoinMarketCapId)].Quote.USD.Price-coin.Price)*100) / 100
	}

	if err != nil {
		return nil, err
	}

	return userCoinsResponse, nil
}

func (s *CoinmarketcapService) GetUserCoinsGrouped(userId int) ([]model.UserCoinsResponse, error) {
	var responseCoins []model.UserCoinsResponse

	coins, err := s.FinHubRepository.GetUserCoinsGrouped(userId)

	if err != nil {
		return nil, err
	}

	for _, coin := range coins {
		coinInfo, err := s.FinHubRepository.GetCoin(coin.CoinId)

		if err != nil {
			return nil, err
		}

		var responseCoin model.UserCoinsResponse

		responseCoin.UserId = coin.UserId
		responseCoin.Name = coinInfo.Name
		responseCoin.Symbol = coinInfo.Symbol
		responseCoin.Slug = coinInfo.Slug
		responseCoin.CoinMarketCapId = coinInfo.Id
		responseCoin.CoinMarketCapRank = coinInfo.Rank
		responseCoin.CoinMarketCapStatus = coinInfo.Status
		responseCoin.Amount = coin.Amount

		weightedAveragePrice := coin.Price / coin.Amount

		responseCoin.Price = weightedAveragePrice

		responseCoins = append(responseCoins, responseCoin)
	}

	var userCoinSlug []string

	for _, coin := range responseCoins {
		userCoinSlug = append(userCoinSlug, strconv.Itoa(coin.CoinMarketCapId))
	}

	coinCurrentData, err := s.GetCoinsData(userCoinSlug)

	if err != nil {
		return nil, err
	}

	for i, coin := range responseCoins {
		responseCoins[i].CurrentPrice = math.Round((coinCurrentData.Data[strconv.Itoa(coin.CoinMarketCapId)].Quote.USD.Price)*100) / 100

		responseCoins[i].CurrentProfit = math.Round(coinCurrentData.Data[strconv.Itoa(coin.CoinMarketCapId)].Quote.USD.Price-coin.Price) * coin.Amount

	}

	if err != nil {
		return nil, err
	}

	return responseCoins, nil

}

func (s *CoinmarketcapService) GetCoins() ([]model.AllCoinsResponse, error) {
	coins, err := s.FinHubRepository.GetCoins()

	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Rank < coins[j].Rank
	})

	if err != nil {
		return nil, err
	}

	return coins, nil
}

func (s *CoinmarketcapService) AddUserCoin(coins *model.UserCoins) error {
	err := s.FinHubRepository.AddUserCoin(coins)

	if err != nil {
		return err
	}

	return nil
}
