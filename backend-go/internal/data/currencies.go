package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type CurrencyModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type Currency struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Rate      float32   `json:"rate"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"version"`
}

var (
	ErrDuplicateCurrency = errors.New("duplicate currency")
)

func (c *CurrencyModel) GetAll() ([]*Currency, error) {
	query := `SELECT id, name, rate, created_at, version FROM currencies`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currencies := []*Currency{}

	for rows.Next() {
		var cur Currency
		err := rows.Scan(&cur.ID, &cur.Name, &cur.Rate, &cur.CreatedAt, &cur.Version)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, &cur)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}

func (c *CurrencyModel) GetByID(id int64) (*Currency, error) {
	query := `
		SELECT id, name, rate, version, created_at 
		FROM currencies 
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var currency Currency

	err := c.DB.QueryRowContext(ctx, query, id).Scan(&currency.ID, &currency.Name, &currency.Rate, &currency.Version, &currency.CreatedAt)

	if err != nil {
		c.ErrorLog.Print(err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &currency, err
}

func (c *CurrencyModel) Insert(currency *Currency) error {
	query := `
		INSERT INTO currencies (name , rate)
		VALUES ($1 , $2)
		RETURNING id, created_at, version 
		`

	args := []interface{}{currency.Name,
		currency.Rate}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&currency.ID, &currency.CreatedAt, &currency.Version)

	if err != nil {
		c.ErrorLog.Print(err)
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "currencies_name_key"`:
			return ErrDuplicateCurrency
		default:
			return err
		}
	}
	return nil
}

func (c *CurrencyModel) Update(currency *Currency) error {
	query := `
		UPDATE currencies 
		SET name=$1, rate=$2, version=version+1
		WHERE id=$3 AND version=$4
		RETURNING version 
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		currency.Name,
		currency.Rate,
		currency.ID,
		currency.Version,
	}

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&currency.Version)

	if err != nil {
		c.ErrorLog.Print(err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (c *CurrencyModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM currencies
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
