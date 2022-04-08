package test


type Delivery struct {
	Delivery_Id int    `json:"d_id" db:"d_id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Phone       string `json:"phone" db:"phone" binding:"required"`
	Zip         string `json:"zip" db:"zip" binding:"required"`
	City        string `json:"city" db:"city" binding:"required"`
	Address     string `json:"address" db:"address" binding:"required"`
	Region      string `json:"region" db:"region"`
	Email       string `json:"email" db:"email"`
}

type Payment struct {
	Id           int    `json:"id" db:"id"`
	Transaction  string `json:"transaction" db:"transaction" binding:"required"`
	RequestId    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency" binding:"required"`
	Provider     string `json:"provider" db:"provider" binding:"required"`
	Amount       int    `json:"amount" db:"amount" binding:"required"`
	PaymentDt    int    `json:"payment_dt" db:"payment_dt" binding:"required"`
	Bank         string `json:"bank" db:"bank" binding:"required"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost" binding:"required"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total" binding:"required"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" db:"chrt_id" binding:"required"`
	TrackNumber string `json:"track_number" db:"track_number" binding:"required"`
	Price       int    `json:"price" db:"price" binding:"required"`
	Rid         string `json:"rid" db:"rid" binding:"required"`
	Name        string `json:"name" db:"name" binding:"required"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size" binding:"required"`
	TotalPrice  int    `json:"total_price" db:"total_price" binding:"required"`
	NmId        int    `json:"nm_id" db:"nm_id" binding:"required"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status" binding:"required"`
}

//type Order struct {
//	Order_uid int `json:"order_uid" db:"order_uid"`
//	TrackNumber string `json:"track_number" db:"track_number"`
//	Entry string `json:"entry" db:"entry"`
//	Locale string `json:"locale" db:"locale"`
//	InternalSignature string `json:"internal_signature" db:"internal_signature"`
//	CustomerId string `json:"customer_id" db:"customer_id"`
//	DeliveryService string `json:"delivery_service" db:"delivery_service"`
//	ShardKey string `json:"shardkey" db:"shardkey"`
//	SmId int `json:"sm_id" db:"sm_id"`
//	DateCreated timestamp.Timestamp `json:"date_created" db:"date_created"`
//	OofShard string `json:"oof_shard" db:"oof_shard"`
//	PaymentId int `json:"payment_id" db:"payment_id"`
//	DeliveryId int `json:"delivery_id" db:"delivery_id"`
//}

type Order struct {
	Order_uid         string    `json:"order_uid" db:"order_uid"`
	TrackNumber       string `json:"track_number" db:"track_number"`
	Entry             string `json:"entry" db:"entry" binding:"required"`
	Locale            string `json:"locale" db:"locale" binding:"required"`
	InternalSignature string `json:"internal_signature" db:"internal_signature"`
	CustomerId        string `json:"customer_id" db:"customer_id" binding:"required"`
	DeliveryService   string `json:"delivery_service" db:"delivery_service" binding:"required"`
	ShardKey          string `json:"shardkey" db:"shardkey" binding:"required"`
	SmId              int    `json:"sm_id" db:"sm_id" binding:"required"`
	DateCreated       string `json:"date_created" db:"date_created" binding:"required"`
	OofShard          string `json:"oof_shard" db:"oof_shard" binding:"required"`
	Delivery `json:"delivery"`
	Payment `json:"payment"`
}

type OrderItem struct {
	Id      int
	OrderId int
	ItemId  int
}

type OrderResponseCache struct {
	Order_uid string `json:"order_uid"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Entry string `json:"entry" db:"entry"`
	Del Delivery `json:"delivery" db:"delivery"`
	Paym Payment `json:"payment"`
	Items []Item `json:"items"`
	Locale string `json:"locale" db:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerId string `json:"customer_id"`
	DeliveryService string `json:"delivery_service"`
	ShardKey string `json:"shardkey"`
	SmId int `json:"sm_id"`
	DateCreated string `json:"date_created"`
	OofShard string `json:"oof_shard"`
}

