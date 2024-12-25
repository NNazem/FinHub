package api

import (
	"FinHub/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type FinancialHubApi struct {
	FinancialHubService  *service.FinancialHubService
	GoCardlessApiService *service.GoCardlessApiService
	Router               *mux.Router
}

func NewFinancialHubApi(financialHubService *service.FinancialHubService, goCardlessApiService *service.GoCardlessApiService, router *mux.Router) *FinancialHubApi {
	return &FinancialHubApi{FinancialHubService: financialHubService, GoCardlessApiService: goCardlessApiService, Router: router}
}

func (f *FinancialHubApi) InitApi() {
	f.Router.HandleFunc("/token/{id}", f.GetTokenByUserId).Methods("GET")
	f.Router.HandleFunc("/banks/{country}", f.GetAllBanksByCountry).Methods("GET")
	f.Router.HandleFunc("/bank/{id}", f.GetBankById).Methods("GET")
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
