package model

import "time"

type TokenResponse struct {
	AccessToken    string `json:"access"`
	AccessExpires  int64  `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int64  `json:"refresh_expires"`
}

type Token struct {
	AccessToken    string    `json:"access"`
	AccessExpires  time.Time `json:"access_expires"`
	Refresh        string    `json:"refresh"`
	RefreshExpires time.Time `json:"refresh_expires"`
}

type Bank struct {
	Id                    string   `json:"id"`
	Name                  string   `json:"name"`
	Bic                   string   `json:"bic"`
	TransactionTotalDays  string   `json:"transaction_total_days"`
	Countries             []string `json:"countries"`
	Logo                  string   `json:"logo"`
	MaxAccessValidForDays string   `json:"max_access_valid_for_days"`
}

type AgreementResponse struct {
	Id                 string    `json:"id"`
	UserId             int       `json:"user_id"`
	AccessToken        string    `json:"access"`
	Created            time.Time `json:"created"`
	InstitutionId      string    `json:"institution_id"`
	MaxHistoricalDays  int       `json:"max_historical_days"`
	AccessValidForDays int       `json:"access_valid_for_days"`
	AccessScope        string    `json:"access_scope"`
	Accepted           time.Time `json:"accepted"`
}

type Agreement struct {
	Id                 string    `json:"id"`
	UserId             int       `json:"user_id"`
	AccessToken        string    `json:"access"`
	Created            time.Time `json:"created"`
	InstitutionId      string    `json:"institution_id"`
	MaxHistoricalDays  int       `json:"max_historical_days"`
	AccessValidForDays int       `json:"access_valid_for_days"`
	AccessScope        []string  `json:"access_scope"`
	Accepted           time.Time `json:"accepted"`
}

type Requisition struct {
	ID            string   `json:"id"`
	Redirect      string   `json:"redirect"`
	InstitutionID string   `json:"institution_id"`
	Agreement     string   `json:"agreement"`
	Accounts      []string `json:"accounts"`
	UserLanguage  string   `json:"user_language"`
	Link          string   `json:"link"`
}

type ListAccountsResponse struct {
	RequisitionID string   `json:"id"`
	Status        string   `json:"status"`
	Agreements    string   `json:"agreements"`
	Accounts      []string `json:"accounts"`
	Reference     string   `json:"reference"`
}

type Account struct {
	ID              string
	RequisitionID   string
	Status          string
	Agreements      string
	Reference       string
	BalanceAmount   string
	BalanceCurrency string
	BalanceType     string
	ReferenceDate   string
}

type Balances struct {
	Balances []Balance `json:"balances"`
}

type Balance struct {
	BalanceAmount BalanceAmount `json:"balanceAmount"`
	BalanceType   string        `json:"balanceType"`
	ReferenceDate string        `json:"referenceDate"`
}

type BalanceAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
