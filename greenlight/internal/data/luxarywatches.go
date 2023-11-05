package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"greenlight.alexedwards.net/internal/validator"
	"time"
)

// Define a custom ErrEditConflict error.
var ErrEditConflict = errors.New("edit conflict: record has been modified")

type Watches struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Brand     []string  `json:"watchesBrand,omitempty"`
	Material  []string  `json:"watchesMaterial,omitempty"`
	Version   int32     `json:"version"`
}

type WatchesModel struct {
	DB *sql.DB
}

func (w *WatchesModel) Insert(watches *Watches) error {
	query := `
		INSERT INTO watches (title, year, price, brand, material)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`
	args := []interface{}{watches.Title, watches.Year, watches.Price, pq.Array(watches.Brand), pq.Array(watches.Material)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return w.DB.QueryRowContext(ctx, query, args...).Scan(&watches.ID, &watches.CreatedAt, &watches.Version)
}

func (w WatchesModel) Get(id int64) (*Watches, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, title, year, price, brand, material, version
FROM watches
WHERE id = $1`

	var watch Watches
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := w.DB.QueryRowContext(ctx, query, id).Scan(
		&watch.ID,
		&watch.CreatedAt,
		&watch.Title,
		&watch.Year,
		&watch.Price,
		pq.Array(&watch.Brand),
		pq.Array(&watch.Material),
		&watch.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &watch, nil
}

func (w WatchesModel) Update(watch *Watches) error {
	query := `
UPDATE watches
SET title = $1, year = $2, price = $3, brand = $4, material = $5, version = version + 1
WHERE id = $6 AND version = $7
RETURNING version`

	args := []interface{}{
		watch.Title,
		watch.Year,
		watch.Price,
		pq.Array(watch.Brand),
		pq.Array(watch.Material),
		watch.ID,
		watch.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := w.DB.QueryRowContext(ctx, query, args...).Scan(&watch.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (w WatchesModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
DELETE FROM watches
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := w.DB.ExecContext(ctx, query, id)
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

func (w WatchesModel) GetAll(title string, filters Filters) ([]*Watches, Metadata, error) {
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, year, price, brand, material, version
FROM watches
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, filters.limit(), filters.offset()}
	rows, err := w.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	watches := []*Watches{}
	for rows.Next() {
		var watch Watches
		err := rows.Scan(
			&totalRecords,
			&watch.ID,
			&watch.CreatedAt,
			&watch.Title,
			&watch.Year,
			&watch.Price,
			pq.Array(&watch.Brand),
			pq.Array(&watch.Material),
			&watch.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		watches = append(watches, &watch)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return watches, metadata, nil
}

func (w Watches) MarshalJSON() ([]byte, error) {
	type WatchAlias Watches
	aux := struct {
		WatchAlias
		Brand    []string `json:"watchesBrand,omitempty"`
		Material []string `json:"watchesMaterial,omitempty"`
	}{
		WatchAlias: WatchAlias(w),
		Brand:      w.Brand,
		Material:   w.Material,
	}
	return json.Marshal(aux)
}

func ValidateWatches(v *validator.Validator, watches *Watches) {
	v.Check(watches.Title != "", "title", "must be provided")
	v.Check(len(watches.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(watches.Year != 0, "year", "must be provided")
	v.Check(watches.Year >= 1888, "year", "must be greater than 1888")
	v.Check(watches.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(watches.Price >= 0, "price", "must be a positive number")
	v.Check(len(watches.Brand) > 0, "watchesBrand", "must be provided")
	v.Check(len(watches.Material) > 0, "watchesMaterial", "must be provided")
}
