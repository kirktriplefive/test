package repository

import (
	"fmt"
	"log"

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

func (r *OrderPostgres) CreateNewOrder(order test.Order, items []test.Item) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}


	del:=order.Delivery
	var deliveryId int
	createDeliveryQuery := fmt.Sprintf("INSERT INTO %s (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING d_id ", deliveryTable)
	row := tx.QueryRow(createDeliveryQuery, del.Name, del.Phone, del.Zip, del.City, del.Address, del.Region, del.Email )
	if err :=row.Scan(&deliveryId); err!=nil {
		tx.Rollback()
		return " ", err
	}
	var paymentId int
	payment:=order.Payment
	createPaymentQuery := fmt.Sprintf("INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost,goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING RETURNING id ", paymentTable)
	rowp := tx.QueryRow(createPaymentQuery, payment.Transaction,payment.RequestId, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err :=rowp.Scan(&paymentId); err!=nil {
		tx.Rollback()
		return " ", err
	}
	
	var order_id string
	createOrdersQuery := fmt.Sprintf("INSERT INTO %s (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, payment_id, delivery_id ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)  ON CONFLICT  DO NOTHING RETURNING order_uid", ordersTable)
	rowOrder:= tx.QueryRow(createOrdersQuery, order.Order_uid ,order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard, paymentId, deliveryId)
	if err:=rowOrder.Scan(&order_id); err!=nil {
		tx.Rollback()
		return " ", err
	}

	var id int
	for _, value := range items {
		createItemQuery := fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT  DO NOTHING RETURNING rid ", ItemsTable)
		_, err:=tx.Exec(createItemQuery, value.ChrtId, value.TrackNumber, value.Price, value.Rid, value.Name, value.Sale, value.Size, value.TotalPrice, value.NmId, value.Brand, value.Status)
		if err != nil {
			tx.Rollback()
			log.Printf(err.Error())
			return " ", err
		}
		
		createListItemsQuery := fmt.Sprintf("INSERT INTO %s (order_id, item_id) values ($1, $2) RETURNING id", itemsInOrderTable)
		row = tx.QueryRow(createListItemsQuery, order_id, value.Rid)
		if err :=row.Scan(&id); err!=nil {
			tx.Rollback()
			return " ", err
		}


	}
	 return order_id, tx.Commit()
}

func (r *OrderPostgres) Close() error {
	err := r.db.Close()
	return err
}

func (r *OrderPostgres) GetOrdersForCache() ([]test.OrderResponseCache, error) {
	var allOrders []test.OrderResponseCache

	query := fmt.Sprintf("SELECT o.order_uid, o.track_number, o.entry, delivery.name, delivery.phone, delivery.zip, delivery.city, delivery.address, delivery.region, delivery.email, payment.transaction, payment.request_id, payment.currency, payment.provider, payment.amount, payment.payment_dt, payment.bank, payment.delivery_cost, payment.goods_total, payment.custom_fee, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard from %s o LEFT JOIN %s delivery ON o.delivery_id=delivery.d_id LEFT JOIN %s payment ON o.payment_id=payment.id",
						ordersTable, deliveryTable, paymentTable)
	rows, err := r.db.Query(query) 
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
        var ord test.OrderResponseCache
		var allItems []test.Item
        if err := rows.Scan(&ord.Order_uid, &ord.TrackNumber, &ord.Entry, 
							&ord.Del.Name, &ord.Del.Phone, &ord.Del.Zip, 
							&ord.Del.City, &ord.Del.Address, &ord.Del.Region, 
							&ord.Del.Email, &ord.Paym.Transaction, &ord.Paym.RequestId, 
							&ord.Paym.Currency, &ord.Paym.Provider, &ord.Paym.Amount,
							&ord.Paym.PaymentDt, &ord.Paym.Bank, &ord.Paym.DeliveryCost,
							&ord.Paym.GoodsTotal, &ord.Paym.CustomFee,
							&ord.Locale, &ord.InternalSignature, &ord.CustomerId, 
							&ord.DeliveryService, &ord.ShardKey, &ord.SmId,
							&ord.DateCreated, &ord.OofShard); err != nil {
            return nil, err
        }
		id:=ord.Order_uid
		queryItem:=fmt.Sprintf("SELECT item.chrt_id, item.track_number, item.price, item.rid, item.name, item.sale, item.size, item.total_price, item.nm_id, item.brand, item.status FROM %s item LEFT JOIN %s itor ON item.rid=itor.item_id WHERE itor.order_id=$1",
								ItemsTable, itemsInOrderTable)
		itemRows, err := r.db.Query(queryItem, id)
		if err!=nil {
			return nil, err
		}

		defer itemRows.Close()
		for itemRows.Next(){
			var item test.Item
			if err:=itemRows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, 
								&item.Rid, &item.Name, &item.Sale, &item.Size, 
								&item.TotalPrice, &item.NmId, &item.Brand, &item.Status); err != nil {
				return nil, err
								
		}
		allItems = append(allItems,item)
	}
		ord.Items = allItems
        allOrders = append(allOrders, ord)
		allItems = append(allItems[:0], allItems[len(allItems):]...)

    }
    if err = rows.Err(); err != nil {
        return nil, err
    }


	return allOrders, nil
}