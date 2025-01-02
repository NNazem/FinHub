package api

import (
	"FinHub/service"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type FinancialHubApi struct {
	FinancialHubService  *service.FinancialHubService
	GoCardlessApiService *service.GoCardlessApiService
	CoinmarketcapService *service.CoinmarketcapService
	Router               *mux.Router
}

func NewFinancialHubApi(financialHubService *service.FinancialHubService, goCardlessApiService *service.GoCardlessApiService, router *mux.Router) *FinancialHubApi {
	return &FinancialHubApi{FinancialHubService: financialHubService, GoCardlessApiService: goCardlessApiService, Router: router}
}

func (f *FinancialHubApi) InitApi() {

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // Your frontend origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Apply the middleware to your router
	f.Router.Use(corsMiddleware)

	f.Router.HandleFunc("/token/{id}", f.GetTokenByUserId).Methods("GET")
	f.Router.HandleFunc("/banks/{country}", f.GetAllBanksByCountry).Methods("GET")
	f.Router.HandleFunc("/bank/{id}", f.GetBankById).Methods("GET")
	f.Router.HandleFunc("/transactions/{userId}/{accountId}", f.GetAccountTransaction).Methods("GET")
	f.Router.HandleFunc("/transactions/{userId}", f.GetUserTransactions).Methods("GET")
	f.Router.HandleFunc("/transactions/{userId}/months/{months}", f.GetUserTransactionsByMonths).Methods("GET")
	f.Router.HandleFunc("/accounts/{userId}", f.GetUserAccounts).Methods("GET")
	f.Router.HandleFunc("/coin/{coin}", f.GetCoinLatestData).Methods("GET")
}

func (f *FinancialHubApi) GetTokenByUserId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	atoi, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := f.FinancialHubService.GetTokenByUserId(atoi)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)

	return
}

func (f *FinancialHubApi) GetAllBanksByCountry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	country := params["country"]

	banks, err := f.GoCardlessApiService.GetAllBanksByCountry(country)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(banks)

	return
}

func (f *FinancialHubApi) GetBankById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	bank, err := f.GoCardlessApiService.GetBankById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bank)

	return
}

func (f *FinancialHubApi) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]

	accounts, err := f.FinancialHubService.GetUserAccounts(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)

	return
}

func (f *FinancialHubApi) GetAccountTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	accountId := params["accountId"]

	atoi, err := strconv.Atoi(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//transactions, err := f.FinancialHubService.GetAccountTransactions(atoi, accountId)

	transactions, err := f.FinancialHubService.GetAccountTransactions(atoi, accountId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)

	return
}

func (f *FinancialHubApi) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]

	atoi, err := strconv.Atoi(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions, err := f.FinancialHubService.GetUserTransactions(atoi)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)

	return
}

func (f *FinancialHubApi) GetUserTransactionsByMonths(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	months := params["months"]

	atoi, err := strconv.Atoi(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	atoiMonths, err := strconv.Atoi(months)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions, err := f.FinancialHubService.GetUserTransactionsByMonths(atoi, atoiMonths)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)

	return
}

func (f *FinancialHubApi) GetCoinLatestData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coin := params["coin"]

	coinResponse, err := f.CoinmarketcapService.GetCoinLatestData(coin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coinResponse)

	return
}
