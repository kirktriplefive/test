package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kirktriplefive/test"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) Create(order test.Order) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	del:=order.Delivery
	var deliveryId int
	createDeliveryQuery := fmt.Sprintf("INSERT INTO %s (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING d_id", deliveryTable)
	row := tx.QueryRow(createDeliveryQuery, del.Name, del.Phone, del.Zip, del.City, del.Address, del.Region, del.Email )
	if err :=row.Scan(&deliveryId); err!=nil {
		tx.Rollback()
		return 0, err
	}
	var paymentId int
	payment:=order.Payment
	createPaymentQuery := fmt.Sprintf("INSERT INTO %s (request_id, currency, provider, amount, payment_dt, bank, delivery_cost,goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", paymentTable)
	rowp := tx.QueryRow(createPaymentQuery, payment.RequestId, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err :=rowp.Scan(&paymentId); err!=nil {
		tx.Rollback()
		return 0, err
	}
	createOrdersQuery := fmt.Sprintf("INSERT INTO %s (track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, payment_id, delivery_id ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING order_uid", ordersTable)
	_, err = tx.Exec(createOrdersQuery, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard, paymentId, deliveryId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return paymentId, tx.Commit()

}

func (r *OrderPostgres) CreateOrderWithItem(orderId int, itemId string) (int, error){
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	
	var id int
	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (order_id, item_id) values ($1, $2) RETURNING id", itemsInOrderTable)
	row := tx.QueryRow(createListItemsQuery, orderId, itemId)
	if err :=row.Scan(&id); err!=nil {
		tx.Rollback()
		return 0, err
	}
	 return id, tx.Commit()
}

func (r *OrderPostgres) CreateOrderWithNewItem(orderId int, item test.Item) (int, error){
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var rid string
	createItemQuery := fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING rid", ItemsTable)
	row:=tx.QueryRow(createItemQuery, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
	err = row.Scan(&rid)
	if err != nil {
		return 0, err
	}

	var id int
	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (order_id, item_id) values ($1, $2) RETURNING id", itemsInOrderTable)
	row = tx.QueryRow(createListItemsQuery, orderId, rid)
	if err :=row.Scan(&id); err!=nil {
		tx.Rollback()
		return 0, err
	}
	 return id, tx.Commit()
}

type OrderById struct {
	Order_uid int `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry string `json:"entry" db:"entry"`
	Del test.Delivery `json:"delivery"`
	Paym test.Payment `json:"payment"`
	Items []test.Item `json:"items"`
	Locale string `json:"locale" db:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerId string `json:"customer_id"`
	DeliveryService string `json:"delivery_service"`
	ShardKey string `json:"shardkey"`
	SmId int `json:"sm_id"`
	DateCreated string `json:"date_created"`
	OofShard string `json:"oof_shard"`
}

func (r *OrderPostgres) GetById(orderId int) (test.Order, []test.Item, error){
	var order test.Order
	var items []test.Item
	var delivery test.Delivery
	var payment test.Payment
	deliveryQuery := fmt.Sprintf("SELECT td.d_id, td.name, td.phone, td.zip, td.city, td.address, td.region, td.email FROM %s td INNER JOIN %s ot ON ot.delivery_id = td.d_id WHERE ot.order_uid = $1",
								deliveryTable, ordersTable)
	if err:=r.db.Get(&delivery, deliveryQuery, orderId); err != nil {
		return order, nil ,err
	}

	paymentQuery := fmt.Sprintf("SELECT tp.id, tp.request_id, tp.currency, tp.provider, tp.amount, tp.payment_dt, tp.bank, tp.delivery_cost, tp.goods_total, tp.custom_fee FROM %s tp INNER JOIN %s ot ON ot.payment_id = tp.id WHERE ot.order_uid = $1",
								paymentTable, ordersTable)
	if err:=r.db.Get(&payment, paymentQuery, orderId); err != nil {
		return order, nil ,err
	}

	orderQuery := fmt.Sprintf("SELECT ot.track_number, ot.entry, ot.locale, ot.internal_signature, ot.customer_id, ot.delivery_service, ot.shardkey, ot.sm_id, ot.date_created, ot.oof_shard FROM %s ot WHERE ot.order_uid = $1", 
							ordersTable)
	if err:=r.db.Get(&order, orderQuery, orderId); err != nil{
		return order, nil, err
	}
	order.Payment = payment
	order.Delivery = delivery

	itemsQuery:=fmt.Sprintf("SELECT ti.chrt_id, ti.track_number, ti.price, ti.rid, ti.name, ti.sale, ti.size, ti.total_price, ti.nm_id, ti.brand, ti.status FROM %s ti INNER JOIN %s toi ON ti.rid = toi.item_id WHERE toi.order_id = $1",
				ItemsTable, itemsInOrderTable)
	if err:=r.db.Select(&items, itemsQuery, orderId); err != nil {
		return order, nil, err
	}

	return order, items, nil
}