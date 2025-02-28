package service

import (
	"FinHub/model"
	"FinHub/repository"
	"math"
	"sort"
	"strconv"
)

type FinancialHubService struct {
	financialHubRepository *repository.FinancialHubRepository
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
