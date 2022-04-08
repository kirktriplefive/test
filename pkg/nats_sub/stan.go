package nats_sub

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/kirktriplefive/test"
	"github.com/kirktriplefive/test/pkg/service"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

type getOrdeResponse struct {
	Order_uid string `json:"order_uid"`
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

type Subscriber struct {
    URL string
	ClusterID   string
    ClientID    string
	Subject     string
}

type Client struct {
	M         *sync.Mutex
	Host      string
	ClusterID string
	ClientID  string
	Subject   string
	Client    stan.Conn
	Service   service.Order
}

func NewSubscriber(sub Client) *Client {
	return &Client{
		M:         sub.M,
		Host:      sub.Host,
		ClusterID: sub.ClusterID,
		ClientID:  sub.ClientID,
		Subject:   sub.Subject,
		Service:   sub.Service,
	}
}


func (c *Client) ConnectToStan() error {
	nc, err := nats.Connect(
		c.Host,
	)
	if err != nil {
		return err
	}

	sc, err := stan.Connect(c.ClusterID, c.ClientID, stan.NatsConn(nc), stan.SetConnectionLostHandler(c.reconnectToStan))
	if err != nil {
		nc.Close()
		return errors.Wrap(err, "err connect to stan "+c.Host)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", c.Host, c.ClusterID, c.ClientID)
	c.Client = sc
	_, err = c.Client.Subscribe(c.Subject,
		c.Mcb,
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.AckWait(time.Second*50))
	if err != nil {
		logrus.Info(err, "err subscribe to stan "+c.Host)
	} else {
		logrus.Info("subscribe successful")
	}

	return nil
}

func (c *Client) Close() {
	c.M.Lock()
	defer c.M.Unlock()

	time.Sleep(3 * time.Second)

	if c.Client == nil {
		return
	}

	err := c.Client.Close()
	if err != nil {
		logrus.Error("Can't close STAN?")
		return
	}

	c.Client.NatsConn().Close()
}

func (c *Client) reconnectToStan(_ stan.Conn, _ error) {
	c.M.Lock()
	defer c.M.Unlock()

	err := c.Client.Close()
	if err != nil {
		logrus.Error("err close stan conn", err)
		return
	}

	c.Client.NatsConn().Close()
	c.Client = nil

	err = c.ConnectToStan()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("reconnect to stan")
}

func (c *Client) Mcb(msg *stan.Msg) {
	var newOrder test.Order
	var items []test.Item
	logrus.Printf("deeesxedex")
	c.M.Lock()
	defer c.M.Unlock()
	var order getOrdeResponse
	if err := json.Unmarshal(msg.Data, &order); err == nil {
		items = order.Items
		newOrder.Payment = order.Paym
		newOrder.Delivery = order.Del
		newOrder.Order_uid = order.Order_uid
		newOrder.TrackNumber = order.TrackNumber
		newOrder.Entry = order.Entry
		newOrder.Locale = order.Locale
		newOrder.InternalSignature = order.InternalSignature
		newOrder.CustomerId = order.CustomerId
		newOrder.DeliveryService = order.DeliveryService
		newOrder.ShardKey = order.ShardKey
		newOrder.SmId = order.SmId
		newOrder.DateCreated = order.DateCreated
		newOrder.OofShard = order.OofShard
		logrus.Printf("Order")
		if id,err:=c.Service.CreateNewOrder(newOrder, items); err != nil {
			logrus.Error(err, ": err save order")
			err = msg.Ack()
			if err != nil {
				logrus.Error(err)
			}
			return
		} else {
			logrus.Println("id: %s",id)
		}
	} else {logrus.Println("inshaala: %s",err)}
	if err := msg.Ack(); err != nil {
		logrus.Error(err)
		return
	}

}