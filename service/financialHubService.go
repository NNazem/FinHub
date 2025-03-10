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
	CoinMarketCapService   *CoinMarketCapService
	UserService            *UserService
}

func NewFinancialHubService(financialHubRepository *repository.FinancialHubRepository) *FinancialHubService {
	return &FinancialHubService{financialHubRepository: financialHubRepository}
}

func (f *FinancialHubService) calculateProfit(userCoins []model.UserCoins, userCoinsResponse []model.UserCoinsResponse) ([]model.UserCoinsResponse, error) {
	for _, coin := range userCoins {
		coinInfo, err := f.financialHubRepository.GetCoin(coin.CoinId)

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

	coinCurrentData, err := f.CoinMarketCapService.GetCoinsData(userCoinsSlugs)

	for i, coin := range userCoinsResponse {
		totalCost := coin.Price * coin.Amount

		currentValue := coinCurrentData.Data[strconv.Itoa(coin.CoinMarketCapId)].Quote.USD.Price * coin.Amount

		currentProfit := currentValue - totalCost

		userCoinsResponse[i].CurrentPrice = math.Round(currentValue*100) / 100

		userCoinsResponse[i].CurrentProfit = math.Round(currentProfit*100) / 100
	}

	if err != nil {
		return nil, err
	}

	return userCoinsResponse, nil
}

func (f *FinancialHubService) GetCoins() ([]model.AllCoinsResponse, error) {
	coins, err := f.financialHubRepository.GetCoins()

	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Rank < coins[j].Rank
	})

	if err != nil {
		return nil, err
	}

	return coins, nil
}

func (f *FinancialHubService) GetUserPortfolioHistoricalValue(userid int) ([]model.UserHistoricalPortfolioValue, error) {
	userHistoricalValues, err := f.financialHubRepository.GetUserPortfolioHistoricalValue(userid)

	if err != nil {
		return nil, err
	}

	return userHistoricalValues, nil
}

func (f *FinancialHubService) InsertUsersPortfolioTotalValue() error {
	ids, err := f.UserService.GetUserList()

	if err != nil {
		return err
	}

	for _, id := range ids {
		userCoinsGrouped, err := f.UserService.GetUserCoinsGrouped(id)

		if err != nil {
			return err
		}

		var profit float64

		for _, coin := range userCoinsGrouped {
			profit += coin.CurrentProfit
		}

		err = f.financialHubRepository.InsertUserPortfolioTotalValue(id, profit)

		if err != nil {
			return err
		}
	}
	return nil
}
