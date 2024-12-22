package repository

import (
	"FinHub/model"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
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

func (d *FinancialHubRepository) InsertNewRequisition(id, redirect, institutionId, agreement, userLanguage, link string) error {
	sqlStatement := `
    INSERT INTO requisitions (id, redirect, institution_id, agreement, user_language, link)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := d.Db.Exec(sqlStatement, id, redirect, institutionId, agreement, userLanguage, link)

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
