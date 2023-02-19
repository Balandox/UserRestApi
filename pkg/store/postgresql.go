package store

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"my-rest-api/config"
)

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	ErrEmailNotUnique = errors.New("user with this email already exist")
)

type Store struct {
	Db *sql.DB
}

func New(config *config.Config) (*Store, error) {

	storeConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Store.Host, config.Store.Port, config.Store.Username, config.Store.Password, config.Store.Database,
	)

	db, err := sql.Open("postgres", storeConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &Store{db}, nil

}

func (s *Store) Close() error {
	return s.Db.Close()
}
