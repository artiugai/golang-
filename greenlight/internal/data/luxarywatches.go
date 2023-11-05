package data

import (
	"database/sql"
	"github.com/lib/pq"
	"greenlight.alexedwards.net/internal/validator"
	"time"
)

type Watches struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Price     Price     `json:"Price,omitempty"`
	Type      []string  `json:"WatchesType,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateWatches(v *validator.Validator, watches *Watches) {
	v.Check(watches.Title != "", "title", "must be provided")
	v.Check(len(watches.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(watches.Year != 0, "year", "must be provided")
	v.Check(watches.Year >= 1888, "year", "must be greater than 1888")
	v.Check(watches.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(watches.Price != 0, "price", "must be provided")
	v.Check(watches.Price > 0, "price", "must be a positive integer")
	v.Check(watches.Type != nil, "watchestype", "must be provided")
	v.Check(len(watches.Type) >= 1, "watchestype", "must contain at least 1 genre")
	v.Check(len(watches.Type) <= 5, "watchestype", "must not contain more than 5 watchestype")
	v.Check(validator.Unique(watches.Type), "watchestype", "must not contain duplicate values")

}

type WatchesModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m WatchesModel) Insert(watches *Watches) error {
	// Определите SQL-запрос для вставки новой записи и возвращения сгенерированных системой данных.
	query := `
    INSERT INTO watches (title, year, price, type) VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, version`

	args := []interface{}{watches.Title, watches.Price, watches.Year, watches.Type}
	// Используйте метод QueryRow() для выполнения SQL-запроса на пуле соединений,
	// передавая слайс args в качестве вариативного параметра и сканируя сгенерированные
	// системой значения id, created_at и version в структуру watches.
	return m.DB.QueryRow(query, args...).Scan(&watches.ID, &watches.CreatedAt, &watches.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m WatchesModel) Get(id int64) (*Watches, error) {
	// Тип PostgreSQL bigserial, который мы используем для ID часов, начинает
	// автоинкремент с 1 по умолчанию, так что мы знаем, что у часов не будет ID меньше этого.
	// Чтобы избежать ненужного вызова к базе данных, мы сразу возвращаем ошибку ErrRecordNotFound.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Определите SQL-запрос для получения данных о часах.
	query := `
    SELECT id, created_at, title, year, price, type, version FROM watches
    WHERE id = $1`
	// Объявляем структуру Watches для хранения данных, возвращаемых запросом.
	var watches Watches

	// Выполняем запрос и сканируем результаты непосредственно в структуру watches.
	err := m.DB.QueryRow(query, id).Scan(&watches.ID, &watches.CreatedAt, &watches.Title, &watches.Year, &watches.Price, pq.Array(&watches.Type), &watches.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &watches, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m WatchesModel) Update(movie *Watches) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m WatchesModel) Delete(id int64) error {
	return nil
}

type MockWatchesModel struct{}

func (m MockWatchesModel) Insert(watches *Watches) error {

	return nil //
}

func (m MockWatchesModel) Get(id int64) (*Watches, error) {

	return nil, nil
}

func (m MockWatchesModel) Update(watches *Watches) error {

	return nil
}

func (m MockWatchesModel) Delete(id int64) error {

	return nil
}
