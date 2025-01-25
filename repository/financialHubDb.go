package repository

import (
	"FinHub/model"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 8080
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

type FinancialHubRepository struct {
	Db *sql.DB
}

func InitDb() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db, nil
}

func (d *FinancialHubRepository) GetToken(id int) (*model.Token, error) {
	sqlStatament := `
	SELECT AccessToken, AccessExpires, Refresh, RefreshExpires
	FROM tokens
	WHERE userId = $1
	`

	var token model.Token

	err := d.Db.QueryRow(sqlStatament, id).Scan(&token.AccessToken, &token.AccessExpires, &token.Refresh, &token.RefreshExpires)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (d *FinancialHubRepository) GetAgreement(token *model.Token, institutionId string, userId int) (*model.AgreementResponse, error) {
	sqlStatement := `
	SELECT id, userId, accessToken, created, institution_id, max_historical_days, access_valid_for_days, access_scope, accepted
	FROM agreements
	WHERE userId = $1 AND accessToken = $2 AND institution_id = $3
	`

	var agreement model.AgreementResponse
	err := d.Db.QueryRow(sqlStatement, userId, token.AccessToken, institutionId).Scan(&agreement.Id, &agreement.UserId, &agreement.AccessToken, &agreement.Created, &agreement.InstitutionId, &agreement.MaxHistoricalDays, &agreement.AccessValidForDays, &agreement.AccessScope, &agreement.Accepted)

	if err != nil {
		return nil, err
	}

	return &agreement, nil
}

func (d *FinancialHubRepository) GetRequisition(agreementId string) (*model.Requisition, error) {
	sqlStatement := `
	SELECT id, redirect, institution_id, agreement, user_language, link
	FROM requisitions
	WHERE agreement = $1
	`

	var requisition model.Requisition
	err := d.Db.QueryRow(sqlStatement, agreementId).Scan(&requisition.ID, &requisition.Redirect, &requisition.InstitutionID, &requisition.Agreement, &requisition.UserLanguage, &requisition.Link)

	if err != nil {
		return nil, err
	}

	return &requisition, nil
}

func (d *FinancialHubRepository) InsertNewToken(accessToken, refresh string, accessExpires, refreshExpires time.Time, id int) error {
	sqlStatement := `
	INSERT INTO tokens (userId, AccessToken, AccessExpires, Refresh, RefreshExpires)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := d.Db.Exec(sqlStatement, id, accessToken, accessExpires, refresh, refreshExpires)

	return err
}

func (d *FinancialHubRepository) InsertNewBank(id, name, bic, country, logo string, transactionTotalDays, maxAccessValidForDays int) error {
	sqlStatement := `
	INSERT INTO banks (id, name, bic, transactionTotalDays, country, logo, maxAccessValidForDays)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := d.Db.Exec(sqlStatement, id, name, bic, transactionTotalDays, country, logo, maxAccessValidForDays)

	return err
}

func (d *FinancialHubRepository) InsertNewAgreement(agreement *model.Agreement) error {
	sqlStatement := `
    INSERT INTO agreements (id, userId, accessToken, created, institution_id, max_historical_days, access_valid_for_days, access_scope, accepted)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	accessScopeJSON, err := json.Marshal(agreement.AccessScope)
	if err != nil {
		return fmt.Errorf("failed to serialize access_scope: %w", err)
	}

	_, err = d.Db.Exec(sqlStatement, agreement.Id, agreement.UserId, agreement.AccessToken, agreement.Created, agreement.InstitutionId, agreement.MaxHistoricalDays, agreement.AccessValidForDays, accessScopeJSON, agreement.Accepted)

	return err
}

func (d *FinancialHubRepository) InsertNewRequisition(requisition model.Requisition) error {
	sqlStatement := `
    INSERT INTO requisitions (id, redirect, institution_id, agreement, user_language, link)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := d.Db.Exec(sqlStatement, requisition.ID, requisition.Redirect, requisition.InstitutionID, requisition.Agreement, requisition.UserLanguage, requisition.Link)

	return err
}

func (d *FinancialHubRepository) InsertNewAccount(account *model.Account) error {
	query := `
    INSERT INTO accounts (
        id, requisition_id, status, agreements, reference, 
        balance_amount, balance_currency, balance_type, reference_date
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := d.Db.Exec(query,
		account.ID, account.RequisitionID, account.Status, account.Agreements,
		account.Reference, account.BalanceAmount, account.BalanceCurrency,
		account.BalanceType, account.ReferenceDate)
	if err != nil {
		return fmt.Errorf("failed to insert new account: %v", err)
	}
	return nil
}

func (d *FinancialHubRepository) GetAccountById(accountId string) int {
	sqlStatement := `
	SELECT count(*)
	FROM accounts
	WHERE id = $1
	`

	var count int

	err := d.Db.QueryRow(sqlStatement, accountId).Scan(&count)

	if err != nil {
		return 0
	}

	return count
}

func (d *FinancialHubRepository) UpdateAccountBalance(accountId string, balance model.Balance) error {
	sqlStatement := `
	UPDATE accounts
	SET balance_amount = $1, balance_currency = $2, balance_type = $3, reference_date = $4
	WHERE id = $5
	`

	_, err := d.Db.Exec(sqlStatement, balance.BalanceAmount.Amount, balance.BalanceAmount.Currency, balance.BalanceType, balance.ReferenceDate, accountId)

	if err != nil {
		return err
	}

	return nil
}

func (d *FinancialHubRepository) GetBalanceByUserId(userId int) (float32, error) {
	sqlStatement := `
	SELECT SUM(CAST(balance_amount AS FLOAT))
	FROM users 
	join agreements on users.id = agreements.userId
	join requisitions on agreements.id = requisitions.agreement
	join accounts on requisitions.id = accounts.requisition_id
	WHERE users.id = $1
	`

	var balance float32

	err := d.Db.QueryRow(sqlStatement, userId).Scan(&balance)

	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (d *FinancialHubRepository) DeleteToken(id int) error {
	sqlStatament := `
	DELETE FROM tokens
	WHERE userId = $1
	`

	_, err := d.Db.Exec(sqlStatament, id)

	if err != nil {
		return err
	}

	return nil
}

func (d *FinancialHubRepository) GetAccountTransaction(accountId string) ([]model.TransactionResponse, error) {
	sqlStatement := `
	SELECT *
	FROM transactions
	WHERE account_id = $1
	`

	var transactions []model.TransactionResponse

	rows, err := d.Db.Query(sqlStatement, accountId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var transaction model.TransactionResponse
		err = rows.Scan(&transaction.TransactionID, &transaction.Booked, &transaction.Pending, &transaction.AccountID, &transaction.InstitutionID, &transaction.BookingDate, &transaction.ValueDate, &transaction.Amount, &transaction.Currency, &transaction.RemittanceInformationUnstructured, &transaction.InternalTransactionID, &transaction.InsertTime, &transaction.UpdateTime)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	log.Println("transactions")
	return transactions, nil
}

func (d *FinancialHubRepository) GetAccountsByUserId(id string) ([]model.Account, error) {
	sqlStatement := `
	SELECT id, requisition_id, status, agreements, reference, balance_amount, balance_currency, balance_type, reference_date, userId
	FROM accounts
	WHERE userId = $1
	`

	var accounts []model.Account

	rows, err := d.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var account model.Account
		err = rows.Scan(&account.ID, &account.RequisitionID, &account.Status, &account.Agreements, &account.Reference, &account.BalanceAmount, &account.BalanceCurrency, &account.BalanceType, &account.ReferenceDate, &account.UserId)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (d *FinancialHubRepository) GetUserTransaction(id int) ([]model.TransactionResponse, error) {
	sqlStatament := `
	SELECT b.transaction_id, b.booked, b.pending, b.account_id, b.institution_id, b.booking_date, b.value_date, b.amount, b.currency, b.remittance_information_unstructured, b.internal_transactionid
	FROM transactions b
	JOIN accounts a ON a.id = b.account_id
	where a.userId = $1
	`

	var transactions []model.TransactionResponse

	rows, err := d.Db.Query(sqlStatament, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var transaction model.TransactionResponse
		err = rows.Scan(&transaction.TransactionID, &transaction.Booked, &transaction.Pending, &transaction.AccountID, &transaction.InstitutionID, &transaction.BookingDate, &transaction.ValueDate, &transaction.Amount, &transaction.Currency, &transaction.RemittanceInformationUnstructured, &transaction.InternalTransactionID)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (d *FinancialHubRepository) GetUserTransactionsByMonths(id int, months int) ([]model.TransactionResponse, error) {
	sqlStatament := `
	SELECT b.transaction_id, b.booked, b.pending, b.account_id, b.institution_id, b.booking_date, b.value_date, b.amount, b.currency, b.remittance_information_unstructured, b.internal_transactionid
	FROM transactions b
	JOIN accounts a ON a.id = b.account_id
	where a.userId = $1 AND b.booking_date > current_date - interval '1 month' * $2
	`

	var transactions []model.TransactionResponse

	rows, err := d.Db.Query(sqlStatament, id, months)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var transaction model.TransactionResponse
		err = rows.Scan(&transaction.TransactionID, &transaction.Booked, &transaction.Pending, &transaction.AccountID, &transaction.InstitutionID, &transaction.BookingDate, &transaction.ValueDate, &transaction.Amount, &transaction.Currency, &transaction.RemittanceInformationUnstructured, &transaction.InternalTransactionID)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (d *FinancialHubRepository) TruncateTable(tableName string) error {
	sqlStatement := fmt.Sprintf("TRUNCATE TABLE %s", tableName)

	_, err := d.Db.Exec(sqlStatement)

	if err != nil {
		return err
	}

	return nil
}

func (d *FinancialHubRepository) UpdateCoin(product *model.CoinHistoricalData) error {
	sqlStatement := `
	UPDATE coins
	SET rank = $2, name = $3, symbol = $4, slug = $5, is_active = $6, status = $7, first_historical_data = $8, last_historical_data = $9
	WHERE id = $1
	`

	_, err := d.Db.Exec(sqlStatement, product.Id, product.Rank, product.Name, product.Symbol, product.Slug, product.IsActive, product.Status, product.FirstHistoricalData, product.LastHistoricalData)

	if err != nil {
		return err
	}

	return nil
}

func (d *FinancialHubRepository) AddCoin(product *model.CoinHistoricalData) error {
	sqlStatement := `
	INSERT INTO coins (id, rank, name, symbol, slug, is_active, status, first_historical_data, last_historical_data)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := d.Db.Exec(sqlStatement, product.Id, product.Rank, product.Name, product.Symbol, product.Slug, product.IsActive, product.Status, product.FirstHistoricalData, product.LastHistoricalData)

	if err != nil {
		return err
	}

	return nil
}

func (d *FinancialHubRepository) GetCoin(id int) (model.CoinHistoricalData, error) {
	sqlStatement := `
	SELECT id, rank, name, symbol, slug, is_active, status, first_historical_data, last_historical_data
	FROM coins
	WHERE id = $1
	`

	var coin model.CoinHistoricalData

	err := d.Db.QueryRow(sqlStatement, id).Scan(&coin.Id, &coin.Rank, &coin.Name, &coin.Symbol, &coin.Slug, &coin.IsActive, &coin.Status, &coin.FirstHistoricalData, &coin.LastHistoricalData)

	if err != nil {
		return coin, err
	}

	return coin, nil
}

func (d *FinancialHubRepository) IsCoinPresent(id int) bool {
	sqlStatement := `
	SELECT count(*)
	FROM coins
	WHERE id = $1
	`

	var count int

	err := d.Db.QueryRow(sqlStatement, id).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
}

func (d *FinancialHubRepository) AddUserCoin(coin *model.UserCoins) error {
	sqlStatement := `
	INSERT INTO user_coins (user_id, coin_id, amount, price)
	VALUES ($1, $2, $3, $4)
	`

	_, err := d.Db.Exec(sqlStatement, coin.UserId, coin.CoinId, coin.Amount, coin.Price)

	if err != nil {
		return err
	}

	return nil
}

func (d *FinancialHubRepository) GetUserCoin(userId int) ([]model.UserCoins, error) {
	sqlStatement := `
	SELECT user_id, coin_id, amount, price
	FROM user_coins
	WHERE user_id = $1
	`

	var coins []model.UserCoins

	rows, err := d.Db.Query(sqlStatement, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var coin model.UserCoins
		err = rows.Scan(&coin.UserId, &coin.CoinId, &coin.Amount, &coin.Price)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, nil
}
