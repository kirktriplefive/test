package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kirktriplefive/test"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewItemPostrgres(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) Create(item test.Item) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	var rid string
	createItemQuery := fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING rid", ItemsTable)
	row := tx.QueryRow(createItemQuery, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
	if err :=row.Scan(&rid); err!=nil {
		tx.Rollback()
		return "", err
	}

	return rid, tx.Commit()
}

func (r *ItemPostgres) GetAll() ([]test.Item, error) {
	var item []test.Item
	
	query:= fmt.Sprintf("SELECT * FROM %s", ItemsTable)
	err := r.db.Select(&item, query)
	return item, err
}