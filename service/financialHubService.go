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

func (f *FinancialHubService) GetUserCoin(userId int) ([]model.UserCoinsResponse, error) {
	var userCoinsResponse []model.UserCoinsResponse
	userCoins, err := f.financialHubRepository.GetUserCoin(userId)

	if err != nil {
		return nil, err
	}

	return f.calculateProfit(userCoins, userCoinsResponse)
}

func (f *FinancialHubService) GetUserCoinsGrouped(userId int) ([]model.UserCoinsResponse, error) {
	var userCoinsResponse []model.UserCoinsResponse

	userCoins, err := f.financialHubRepository.GetUserCoinsGrouped(userId)

	if err != nil {
		return nil, err
	}

	return f.calculateProfit(userCoins, userCoinsResponse)

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

	coinCurrentData, err := f.CoinmarketcapService.GetCoinsData(userCoinsSlugs)

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

func (f *FinancialHubService) AddUserCoin(coins *model.UserCoins) error {
	err := f.financialHubRepository.AddUserCoin(coins)

	if err != nil {
		return err
	}

	return nil
}

func (f *FinancialHubService) GetUserList() ([]int, error) {
	userlist, err := f.financialHubRepository.GetUserList()

	if err != nil {
		return nil, err
	}

	return userlist, err
}

func (f *FinancialHubService) GetUserPortfolioHistoricalValue(userid int) ([]model.UserHistoricalPortfolioValue, error) {
	userHistoricalValues, err := f.financialHubRepository.GetUserPortfolioHistoricalValue(userid)

	if err != nil {
		return nil, err
	}

	return userHistoricalValues, nil
}

func (f *FinancialHubService) InsertUsersPortfolioTotalValue() error {
	ids, err := f.GetUserList()

	if err != nil {
		return err
	}

	for _, id := range ids {
		userCoinsGrouped, err := f.GetUserCoinsGrouped(id)

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
