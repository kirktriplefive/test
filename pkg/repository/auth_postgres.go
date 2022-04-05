package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kirktriplefive/test"

)

type AuthPostrgres struct {
	db *sqlx.DB
}

func NewAuthPostrgres(db *sqlx.DB) *AuthPostrgres {
	return &AuthPostrgres{db: db}
}

func (r *AuthPostrgres) CreateUser(user test.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1,$2,$3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err!=nil{
		return 0, err
	}
	return id, nil
}

func (r *AuthPostrgres) GetUser(username, password string) (test.User, error){
	var user test.User
	query:= fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err:=r.db.Get(&user, query, username, password)

	return user, err
}