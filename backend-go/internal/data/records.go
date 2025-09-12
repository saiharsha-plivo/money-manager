package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/saiharsha/money-manager/pkg/validator"
)

type RecordModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type Record struct {
	ID          int64     `json:"id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	TypeID      int64     `json:"type_id"`
	CurrencyID  int64     `json:"currency_id"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int64     `json:"version"`
}

var (
	ErrDuplicateRecord = errors.New("duplicate record")
)

func (r *RecordModel) Insert(record *Record) error {
	query := `
		INSERT INTO records (amount, description, type_id, currency_id, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{record.Amount, record.Description, record.TypeID, record.CurrencyID, record.UserID}

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&record.ID, &record.CreatedAt, &record.Version)
	if err != nil {
		r.ErrorLog.Print(err.Error())
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "records_pkey"`:
			return ErrDuplicateRecord
		default:
			return err
		}
	}

	return nil
}

func (r *RecordModel) GetByID(id int64) (*Record, error) {
	query := `
		SELECT id, amount, description, type_id, currency_id, created_at, version, user_id
		FROM records
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var record Record

	err := r.DB.QueryRowContext(ctx, query, id).Scan(&record.ID, &record.Amount, &record.Description, &record.TypeID, &record.CurrencyID, &record.CreatedAt, &record.Version, &record.UserID)
	if err != nil {
		r.ErrorLog.Print(err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &record, nil
}

func (r *RecordModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM records
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id)
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

func (r *RecordModel) Update(record *Record) error {
	query := `
		UPDATE records
		SET amount = $1, description = $2, type_id = $3, currency_id = $4, version = version + 1
		WHERE id = $5 AND version = $6 AND user_id = $7
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{&record.Amount, &record.Description, &record.TypeID, &record.CurrencyID, &record.ID, &record.Version, &record.UserID}

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&record.Version)
	if err != nil {
		r.ErrorLog.Print(err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Query Parameters:
// userID: int64
// filters: Filters
// filters.StartDate: time.Time example: 2025-01-01 00:00:00 , can be empty
// filters.EndDate: time.Time example: 2025-01-01 00:00:00
// with pagination
// return: []*Record, error
func (r *RecordModel) GetRecordsForUser(userID int64, filters Filters) ([]*Record, error) {

	var query string
	var args []interface{}

	switch {
	case filters.StartDate.IsZero():
		query = fmt.Sprintf(`
			SELECT id, amount, description, type_id, currency_id, created_at, version
			FROM records
			WHERE user_id = $1
			AND created_at <= $2
			ORDER BY %s %s, id ASC
			LIMIT $3 OFFSET $4`,
			filters.sortColumn(), filters.sortDirection())

		args = []interface{}{userID, filters.EndDate, filters.limit(), filters.offset()}
	default:
		query = fmt.Sprintf(`
			SELECT id, amount, description, type_id, currency_id, created_at, version
			FROM records
			WHERE user_id = $1
			AND created_at >= $2
			AND created_at <= $3
			ORDER BY %s %s, id ASC
			LIMIT $4 OFFSET $5`,
			filters.sortColumn(), filters.sortDirection())
		args = []interface{}{userID, filters.StartDate, filters.EndDate, filters.limit(), filters.offset()}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		r.ErrorLog.Print(err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	records := make([]*Record, 0)
	for rows.Next() {
		var record Record
		err := rows.Scan(&record.ID, &record.Amount, &record.Description, &record.TypeID, &record.CurrencyID, &record.CreatedAt, &record.Version)
		if err != nil {
			r.ErrorLog.Print(err.Error())
			return nil, err
		}
		records = append(records, &record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func ValidateRecord(v *validator.Validator, record *Record) {
	v.Check(record.Amount > 0, "amount", "must be greater than 0")
	v.Check(record.TypeID > 0, "type_id", "must be provided")
	v.Check(record.CurrencyID > 0, "currency_id", "must be provided")
}

func ValidateRecordUpdate(v *validator.Validator, amount *int64, typeID *int64, currencyID *int64) {
	if amount != nil {
		v.Check(*amount > 0, "amount", "must be greater than 0")
	} else {
		v.AddError("amount", "must be provided")
	}
	if typeID != nil {
		v.Check(*typeID > 0, "type_id", "must be provided")
	} else {
		v.AddError("type_id", "must be provided")
	}
	if currencyID != nil {
		v.Check(*currencyID > 0, "currency_id", "must be provided")
	} else {
		v.AddError("currency_id", "must be provided")
	}
}
