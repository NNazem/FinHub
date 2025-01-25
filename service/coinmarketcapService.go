package service

import (
	"FinHub/model"
	"FinHub/repository"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type CoinmarketcapService struct {
	FinHubRepository *repository.FinancialHubRepository
}

func NewCoinmarketcapService(repo *repository.FinancialHubRepository) *CoinmarketcapService {
	return &CoinmarketcapService{
		FinHubRepository: repo,
	}
}

func (s *CoinmarketcapService) GetCoinData(coin string) (*model.CoinResponse, error) {
	url := "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?id=" + coin

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CMC_PRO_API_KEY", os.Getenv("COINMARKETCAP_API_KEY"))

	res, err := http.DefaultClient.Do(req)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	coinResponse := &model.CoinResponse{}

	err = json.Unmarshal(body, coinResponse)

	if err != nil {
		return nil, err
	}

	return coinResponse, nil
}

func (s *CoinmarketcapService) GetCoinInfo(coin string) (*model.CoinInfoResponse, error) {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/info?symbol=" + coin

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CMC_PRO_API_KEY", os.Getenv("COINMARKETCAP_API_KEY"))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	coinInfoResponse := &model.CoinInfoResponse{}

	err = json.Unmarshal(body, coinInfoResponse)

	if err != nil {
		return nil, err
	}

	return coinInfoResponse, nil
}

func (s *CoinmarketcapService) GetCoinsData(coins []string) (*model.CoinResponse, error) {
	baseUrl := "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?id="

	for i, coin := range coins {
		if i == 0 {
			baseUrl += coin
		} else {
			baseUrl += "," + coin
		}
	}

	req, err := http.NewRequest("GET", baseUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CMC_PRO_API_KEY", os.Getenv("COINMARKETCAP_API_KEY"))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	var coinInfoResponse model.CoinResponse

	err = json.NewDecoder(res.Body).Decode(&coinInfoResponse)

	if err != nil {
		return nil, err
	}

	return &coinInfoResponse, nil
}

func (s *CoinmarketcapService) GetCoinsHistoricalData() error {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("X-CMC_PRO_API_KEY", os.Getenv("COINMARKETCAP_API_KEY"))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	coinHistoricalData := &model.CoinHistoricalResponse{}

	err = json.Unmarshal(body, coinHistoricalData)

	for _, coin := range coinHistoricalData.Data {

		isPresent := s.FinHubRepository.IsCoinPresent(coin.Id)

		if isPresent {
			err = s.FinHubRepository.UpdateCoin(&coin)

			if err != nil {
				return err
			}
		} else {
			log.Println("Adding coin: ", coin.Name)
			err = s.FinHubRepository.AddCoin(&coin)

			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *CoinmarketcapService) AddUserCoin(coins *model.UserCoins) error {
	err := s.FinHubRepository.AddUserCoin(coins)

	if err != nil {
		return err
	}

	return nil
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
