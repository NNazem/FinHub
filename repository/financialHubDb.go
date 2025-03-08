package repository

import (
	"FinHub/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type FinancialHubRepository struct {
	Db *sql.DB
}

func InitDb() (*sql.DB, error) {
	user, err1 := os.LookupEnv("POSTGRE_USER")

	log.Println(user, err1)

	pass, err2 := os.LookupEnv("POSTGRE_PASSWORD")

	log.Println(pass, err2)

	dbName, err3 := os.LookupEnv("POSTGRE_DB_NAME")

	log.Println(dbName, err3)

	host, err4 := os.LookupEnv("POSTGRE_HOST")

	log.Println(host, err4)

	log.Println("User : " + user)
	log.Println("Pass : " + pass)
	log.Println("DbName : " + dbName)
	log.Println("Host : " + host)

	connStr := fmt.Sprintf("user='%s' password='%s' host='%s' dbname='%s'", user, pass, host, dbName)
	db, err := sql.Open("postgres", connStr)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected.")

	return db, nil
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

func (d *FinancialHubRepository) GetUserCoinsGrouped(userId int) ([]model.UserCoins, error) {
	sqlStatament := `
	SELECT user_id, coin_id, SUM(AMOUNT), SUM(PRICE)
	FROM user_coins
	WHERE user_id = $1
	GROUP BY user_id, coin_id`

	var coins []model.UserCoins

	rows, err := d.Db.Query(sqlStatament, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var coin model.UserCoins
		err := rows.Scan(&coin.UserId, &coin.CoinId, &coin.Amount, &coin.Price)

		if err != nil {
			return nil, err
		}

		coins = append(coins, coin)
	}

	return coins, nil

}

func (d *FinancialHubRepository) GetCoins() ([]model.AllCoinsResponse, error) {
	sqlStatement := `
	SELECT id, rank, name
	FROM coins
	`

	var coins []model.AllCoinsResponse

	rows, err := d.Db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var coin model.AllCoinsResponse
		err := rows.Scan(&coin.Id, &coin.Rank, &coin.Name)

		if err != nil {
			return nil, err
		}

		coins = append(coins, coin)
	}

	return coins, nil
}

func (d *FinancialHubRepository) GetAmountPerTypology(userId int) ([]model.AmountPerCategory, error) {
	sqlStatement := `
	SELECT sum(amount), typology
	FROM user_finance_products_integrated
	WHERE user_id = $1
	group by typology
	`

	var amountPerCategories []model.AmountPerCategory

	rows, err := d.Db.Query(sqlStatement, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var amountPerCategory model.AmountPerCategory
		err := rows.Scan(&amountPerCategory.Amount, &amountPerCategory.Category)

		if err != nil {
			return nil, err
		}

		amountPerCategories = append(amountPerCategories, amountPerCategory)
	}

	return amountPerCategories, nil
}

func (d *FinancialHubRepository) GetAmountPerCrypto(userId int) ([]model.AmountPerCrypto, error) {
	sqlStatement := `
	SELECT sum(amount), c.id
	FROM user_coins uc
	JOIN coins c ON uc.coin_id = c.id
	WHERE uc.user_id = $1
	group by c.id
	`

	var amountPerCryptos []model.AmountPerCrypto

	rows, err := d.Db.Query(sqlStatement, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var amountPerCrypto model.AmountPerCrypto

		err := rows.Scan(&amountPerCrypto.Amount, &amountPerCrypto.Name)

		if err != nil {
			return nil, err
		}

		amountPerCryptos = append(amountPerCryptos, amountPerCrypto)
	}

	return amountPerCryptos, nil
}

func (d *FinancialHubRepository) AddCoinToUser(userid string, coin model.AddCryptoRequest) error {
	sqlstatament :=
		`INSERT INTO user_coins (user_id, coin_id, amount, price) VALUES ($1, $2, $3, $4)`

	_, err := d.Db.Exec(sqlstatament, userid, coin.Coin.Id, coin.Amount, coin.Price)

	if err != nil {
		return err
	}

	return nil
}
