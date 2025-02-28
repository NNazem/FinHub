package api

import (
	"FinHub/model"
	"FinHub/service"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type FinancialHubApi struct {
	FinancialHubService  *service.FinancialHubService
	CoinmarketcapService *service.CoinmarketcapService
	Router               *mux.Router
}

func NewFinancialHubApi(financialHubService *service.FinancialHubService, coinMarketcapApiService *service.CoinmarketcapService, router *mux.Router) *FinancialHubApi {
	return &FinancialHubApi{FinancialHubService: financialHubService, CoinmarketcapService: coinMarketcapApiService, Router: router}
}

func (f *FinancialHubApi) InitApi() {

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Your frontend origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Apply the middleware to your router
	f.Router.Use(corsMiddleware)

	// PortFolio tracker part
	f.Router.HandleFunc("/coin/{coin}", f.GetCoinLatestData).Methods("GET")
	f.Router.HandleFunc("/coinInfo/{coin}", f.GetCoinInfo).Methods("GET")
	f.Router.HandleFunc("/coinsHistoricalData", f.GetCoinsHistoricalData).Methods("GET")
	f.Router.HandleFunc("/userCoins", f.AddUserCoins).Methods("POST")
	f.Router.HandleFunc("/userCoins/{userId}", f.GetUserCoin).Methods("GET")
	f.Router.HandleFunc("/userCoinsGrouped/{userId}", f.GetUserCoinGrouped).Methods("GET")
	f.Router.HandleFunc("/coins", f.GetCoins).Methods("GET")
	f.Router.HandleFunc("/userAmountPerTypologies/{userId}", f.GetUserAmountPerTypologies).Methods("GET")
	f.Router.HandleFunc("/addCrypto/{userId}", f.AddCrypto).Methods("POST")

	// Utils
	f.Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})
}

func (f *FinancialHubApi) GetCoinLatestData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coin := params["coin"]

	coinResponse, err := f.CoinmarketcapService.GetCoinsData([]string{coin})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coinResponse)

	return
}

func (f *FinancialHubApi) GetCoinInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coin := params["coin"]

	coinInfoResponse, err := f.CoinmarketcapService.GetCoinInfo(coin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coinInfoResponse)

	return
}

func (f *FinancialHubApi) GetCoinsHistoricalData(w http.ResponseWriter, r *http.Request) {
	err := f.CoinmarketcapService.GetCoinsHistoricalData()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (f *FinancialHubApi) AddUserCoins(w http.ResponseWriter, r *http.Request) {
	coin := &model.UserCoins{}

	err := json.NewDecoder(r.Body).Decode(coin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = f.CoinmarketcapService.AddUserCoin(coin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (f *FinancialHubApi) GetUserCoin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]

	atoi, err := strconv.Atoi(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coins, err := f.CoinmarketcapService.GetUserCoin(atoi)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coins)

	return
}

func (f *FinancialHubApi) GetUserCoinGrouped(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]

	atoi, err := strconv.Atoi(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coins, err := f.CoinmarketcapService.GetUserCoinsGrouped(atoi)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coins)

	return
}

func (f *FinancialHubApi) GetCoins(w http.ResponseWriter, r *http.Request) {
	coins, err := f.CoinmarketcapService.GetCoins()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coins)

	return
}

func (f *FinancialHubApi) GetUserAmountPerTypologies(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]

	userIdConverted, err := strconv.Atoi(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAmountPerCategories, err := f.FinancialHubService.GetUserAmountPerTypologies(userIdConverted)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userAmountPerCategories)
}

func (f *FinancialHubApi) AddCrypto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := params["userId"]

	coin := model.AddCryptoRequest{}

	err := json.NewDecoder(r.Body).Decode(&coin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = f.FinancialHubService.AddCoinToUser(userId, coin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coin)
	return
}
