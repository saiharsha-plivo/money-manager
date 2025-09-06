package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type RecordTypeModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type RecordType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrDuplicateRecordType = errors.New("duplicate currency")
)

func (r *RecordTypeModel) Insert(recordtype *RecordType) error {
	query := `
		INSERT INTO types (name)
		VALUES ($1)
		RETURNING id, version, created_at
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, recordtype.Name).Scan(&recordtype.ID, &recordtype.Version, &recordtype.CreatedAt)
	if err != nil {
		r.ErrorLog.Print(err)
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "type_name_key"`:
			return ErrDuplicateRecordType
		default:
			return err
		}
	}
	return nil
}

func (r *RecordTypeModel) GetAll() ([]*RecordType, error) {
	query := `SELECT id, name, created_at, version FROM types`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recordtypes := make([]*RecordType, 0)
	for rows.Next() {
		var recordtype RecordType
		err := rows.Scan(&recordtype.ID, &recordtype.Name, &recordtype.CreatedAt, &recordtype.Version)
		if err != nil {
			return nil, err
		}
		recordtypes = append(recordtypes, &recordtype)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return recordtypes, nil
}

func (r *RecordTypeModel) GetByID(id int64) (*RecordType, error) {
	query := `
		SELECT id, name, version, created_at 
		FROM types
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var recordtype RecordType

	err := r.DB.QueryRowContext(ctx, query, id).Scan(&recordtype.ID, &recordtype.Name, &recordtype.Version, &recordtype.CreatedAt)

	if err != nil {
		r.ErrorLog.Print(err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &recordtype, err
}

func (r *RecordTypeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM types
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
