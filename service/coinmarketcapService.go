package service

import (
	"FinHub/model"
	"FinHub/repository"
	"encoding/json"
	"io"
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

func (s *CoinmarketcapService) GetCoinLatestData(coin string) (*model.CoinResponse, error) {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=" + coin

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
