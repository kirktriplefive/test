package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kirktriplefive/test"
)



func (h *Handler) createOrder(c *gin.Context){
	var input test.Order
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	
	payment_id, err := h.services.Order.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"payment_id": payment_id,
	})

}

func (h *Handler) createOrderWithItem(c *gin.Context){
	orderId,err:=strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order id param")
		return
	}

	var itemId string
	if err := c.BindJSON(&itemId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	

	id,err:=h.services.Order.CreateOrderWithItem(orderId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":id,
	})

}

func (h *Handler) createOrderWithNewItem(c *gin.Context){
	orderId,err:=strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order id param")
		return
	}

	var item test.Item
	if err := c.BindJSON(&item); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	

	id,err:=h.services.Order.CreateOrderWithNewItem(orderId, item)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":id,
	})

}

type getOrderesponse struct {
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

func (h *Handler) getOrderById(c *gin.Context){
	orderId,err:=strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid order id param")
		return
	}

	order, items, err := h.services.Order.GetById(orderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getOrderesponse{
		Order_uid: order.Order_uid,
		TrackNumber: order.TrackNumber,
		Entry: order.Entry,
		Del: order.Delivery,
		Paym: order.Payment,
		Items: items, 
		Locale: order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerId: order.CustomerId,
		DeliveryService: order.DeliveryService,
		ShardKey: order.ShardKey,
		SmId: order.SmId,
		DateCreated: order.DateCreated,
		OofShard: order.OofShard,
	})
}

