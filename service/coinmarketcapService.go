package service

import (
	"FinHub/model"
	"FinHub/repository"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type CoinmarketcapService struct {
	FinHubRepository *repository.FinancialHubRepository
}

func NewCoinmarketcapService(repo *repository.FinancialHubRepository) *CoinmarketcapService {
	return &CoinmarketcapService{
		FinHubRepository: repo,
	}
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
