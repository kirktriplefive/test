package main

import (
	//"encoding/json"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/kirktriplefive/test"
	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

// NOTE: Use tls scheme for TLS, e.g. stan-pub -s tls://demo.nats.io:4443 foo hello


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

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
	}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
	bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
	}

func RandomOrder() getOrdeResponse {
	order:= getOrdeResponse{
		Order_uid:         randomString(5),
		TrackNumber:       randomString(7),
		Entry:             randomString(8),
		Del:               test.Delivery{
			Delivery_Id: 1,
			Name:        randomString(5),
			Phone:       randomString(5),
			Zip:         randomString(5),
			City:        randomString(5),
			Address:     randomString(5),
			Region:      randomString(5),
			Email:       randomString(10),
		},
		Paym:              test.Payment{
			Id:           1,
			Transaction: randomString(8),
			RequestId:    "",
			Currency:     randomString(3),
			Provider:     "wbpay",
			Amount:       randomInt(1,10000),
			PaymentDt:    randomInt(1,100000),
			Bank:         "alpha",
			DeliveryCost: randomInt(1,10000),
			GoodsTotal:   randomInt(1,10000),
			CustomFee:    randomInt(1,90),
		},
		Items: []test.Item{
			{
				ChrtId: randomInt(1,10000),
				TrackNumber: randomString(9),
				Price: 453,
				Rid: randomString(3),
				Name:  randomString(10),      
				Sale: randomInt(1,90),      
				Size: "0",       
				TotalPrice:  randomInt(1,10000), 
				NmId:   randomInt(1,10000),     
				Brand: randomString(10),      
				Status: 202,
			},
			{
				ChrtId: 9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price: 453,
				Rid: "ab4214ae0btest",
				Name:  "Mascaras",      
				Sale: 30,      
				Size: "0",       
				TotalPrice:  317, 
				NmId:   2389212,     
				Brand: "Vivienne Sabo",      
				Status: 202,

			},
			{
				ChrtId: 9934930,
				TrackNumber: "WBIeMTESTTRACK",
				Price: 453,
				Rid: randomString(10),
				Name:  "Mascaras",      
				Sale: 30,      
				Size: "0",       
				TotalPrice:  317, 
				NmId:   2389212,     
				Brand: "Vivienne Sabo",      
				Status: 202,

			},

		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmId:              99,
		DateCreated:       "tem",
		OofShard:          "1",
	}
	return order
}

func main() {
	var (
		clusterID string
		clientID  string
		URL       string
		userCreds string
		
	)

	flag.StringVar(&URL, "s", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&URL, "server", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&clusterID, "c", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clusterID, "cluster", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-pub", "The NATS Streaming client ID to connect with")
	flag.StringVar(&clientID, "clientid", "stan-pub", "The NATS Streaming client ID to connect with")

	log.SetFlags(0)
	flag.Parse()


	

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Streaming Example Publisher")}
	// Use UserCredentials
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

	// Connect to NATS
	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
    	log.Fatal(err)
	}
	defer ec.Close()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(ec.Conn))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()
	
	subj:= "foo"
	for i:=1; i<20; i++ {
		if i%3!=0 {
			b, err := json.Marshal(RandomOrder())
			if err != nil {
				fmt.Println(err)
				return
			}
			msg:=[]byte(string(b))
			err = sc.Publish(subj, msg)
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}

		log.Printf("Published [%s] : '%s'\n", subj, msg)
		time.Sleep(time.Second*5)
		} else {
			msg:=[]byte(randomString(10))
			err = sc.Publish(subj, msg)
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}

		log.Printf("Published [%s] : '%s'\n", subj, msg)
		time.Sleep(time.Second*5)
		}
		


		
	}
}