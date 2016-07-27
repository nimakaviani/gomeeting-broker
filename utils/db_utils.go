package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"database/sql"

	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
)

type DataStore interface {
	Init() error
	Close() error
	SaveToken(oauth2.Token) error
	LoadToken() (oauth2.Token, error)
}

type datastore struct {
	db *sql.DB

	User     string
	Password string
	Name     string
	Host     string
}

func NewDBDataStore(user, password, name, host string) DataStore {
	return &datastore{
		User:     user,
		Name:     name,
		Password: password,
		Host:     host,
	}
}

func (d *datastore) Init() error {
	var err error

	dbInfo := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		d.User,
		d.Password,
		d.Name,
		d.Host,
	)

	d.db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}

	return nil
}

func (d *datastore) Close() error {
	return d.db.Close()
}

func (d *datastore) SaveToken(token oauth2.Token) error {
	var err error

	_, err = d.db.Exec("DELETE FROM credentials;")
	if err != nil {
		return err
	}

	tokenJson, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}

	_, err = d.db.Exec(fmt.Sprintf(
		"INSERT INTO credentials (name, credential) VALUES ('%s', '%s');",
		"oauth",
		string(tokenJson),
	))
	if err != nil {
		return err
	}

	return nil
}

func (d *datastore) LoadToken() (oauth2.Token, error) {
	rows, err := d.db.Query("SELECT name, credential FROM credentials;")
	if err != nil {
		return oauth2.Token{}, err
	}
	defer rows.Close()

	var name, credential string
	for rows.Next() {
		err = rows.Scan(&name, &credential)
		if err != nil {
			return oauth2.Token{}, err
		}

		token := oauth2.Token{}
		err = json.Unmarshal([]byte(credential), &token)
		if err == nil {
			return token, nil
		}
	}

	err = rows.Err()
	if err != nil {
		return oauth2.Token{}, err
	}

	return oauth2.Token{}, errors.New("Could not find token")
}
