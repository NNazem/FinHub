package service

import (
	"FinHub/model"
	"FinHub/repository"
)

type UserService struct {
	financialHubRepository *repository.FinancialHubRepository
	financialHubService    *FinancialHubService
	coinMarketCapService   *CoinMarketCapService
}

func NewUserService(financialHubRepository *repository.FinancialHubRepository, coinMarketCapService *CoinMarketCapService, financialHubService *FinancialHubService) *UserService {
	return &UserService{financialHubRepository: financialHubRepository, coinMarketCapService: coinMarketCapService, financialHubService: financialHubService}
}

func (u *UserService) GetUserList() ([]int, error) {
	userList, err := u.financialHubRepository.GetUserList()

	if err != nil {
		return nil, err
	}

	return userList, err
}

func (u *UserService) GetUserCoinsGrouped(userId int) ([]model.UserCoinsResponse, error) {
	var userCoinsResponse []model.UserCoinsResponse

	userCoins, err := u.financialHubRepository.GetUserCoinsGrouped(userId)

	if err != nil {
		return nil, err
	}

	return u.financialHubService.calculateProfit(userCoins, userCoinsResponse)

}

func (u *UserService) AddCoinToUser(userid string, coin model.AddCryptoRequest) error {
	err := u.financialHubRepository.AddCoinToUser(userid, coin)

	return err
}

func (u *UserService) AddUserCoin(coins *model.UserCoins) error {
	err := u.financialHubRepository.AddUserCoin(coins)

	if err != nil {
		return err
	}

	return nil
}

func (f *FinancialHubService) GetUserAmountPerCryptos(userId int) (*model.UserAmountPerCrypto, error) {
	amountPerCrypto, err := f.financialHubRepository.GetAmountPerCrypto(userId)

	var cryptos []string

	for _, crypto := range amountPerCrypto {
		cryptos = append(cryptos, crypto.Name)
	}

	coinData, err := f.CoinMarketCapService.GetCoinsData(cryptos)

	if err != nil {
		return nil, err
	}

	for i, crypto := range amountPerCrypto {
		price := coinData.Data[crypto.Name].Quote.USD.Price
		amountPerCrypto[i].CurrentValue = price
	}

	return &model.UserAmountPerCrypto{AmountPerCrypto: amountPerCrypto, UserId: userId}, nil
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

func (f *FinancialHubService) GetUserCoin(userId int) ([]model.UserCoinsResponse, error) {
	var userCoinsResponse []model.UserCoinsResponse
	userCoins, err := f.financialHubRepository.GetUserCoin(userId)

	if err != nil {
		return nil, err
	}

	return f.calculateProfit(userCoins, userCoinsResponse)
}
