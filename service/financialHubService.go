package service

import (
	"FinHub/model"
	"FinHub/repository"
	"sort"
	"strconv"
)

type FinancialHubService struct {
	financialHubRepository *repository.FinancialHubRepository
	CoinmarketcapService   *CoinmarketcapService
}

func NewFinancialHubService(financialHubRepository *repository.FinancialHubRepository) *FinancialHubService {
	return &FinancialHubService{financialHubRepository: financialHubRepository}
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

func (f *FinancialHubService) GetUserAmountPerCryptos(userId int) (*model.UserAmountPerCrypto, error) {
	amountPerCrypto, err := f.financialHubRepository.GetAmountPerCrypto(userId)

	var cryptos []string

	for _, crypto := range amountPerCrypto {
		cryptos = append(cryptos, crypto.Name)
	}

	coinData, err := f.CoinmarketcapService.GetCoinsData(cryptos)

	if err != nil {
		return nil, err
	}

	for i, crypto := range amountPerCrypto {
		price := coinData.Data[crypto.Name].Quote.USD.Price
		amountPerCrypto[i].CurrentValue = price
	}

	return &model.UserAmountPerCrypto{AmountPerCrypto: amountPerCrypto, UserId: userId}, nil
}

func (f *FinancialHubService) AddCoinToUser(userid string, coin model.AddCryptoRequest) error {
	err := f.financialHubRepository.AddCoinToUser(userid, coin)

	return err
}

func (s *CoinmarketcapService) GetUserCoin(userId int) ([]model.UserCoinsResponse, error) {
	var userCoinsResponse []model.UserCoinsResponse
	userCoins, err := s.FinHubRepository.GetUserCoin(userId)

	if err != nil {
		return nil, err
	}

	return s.calculateProfit(userCoins, userCoinsResponse)
}

func (s *CoinmarketcapService) GetUserCoinsGrouped(userId int) ([]model.UserCoinsResponse, error) {
	var userCoinsResponse []model.UserCoinsResponse

	userCoins, err := s.FinHubRepository.GetUserCoinsGrouped(userId)

	if err != nil {
		return nil, err
	}

	return s.calculateProfit(userCoins, userCoinsResponse)

}

func (s *CoinmarketcapService) calculateProfit(userCoins []model.UserCoins, userCoinsResponse []model.UserCoinsResponse) ([]model.UserCoinsResponse, error) {
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
		totalCost := coin.Price * coin.Amount

		currentValue := coinCurrentData.Data[strconv.Itoa(coin.CoinMarketCapId)].Quote.USD.Price * coin.Amount

		currentProfit := currentValue - totalCost

		userCoinsResponse[i].CurrentPrice = currentValue

		userCoinsResponse[i].CurrentProfit = currentProfit
	}

	if err != nil {
		return nil, err
	}

	return userCoinsResponse, nil
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
