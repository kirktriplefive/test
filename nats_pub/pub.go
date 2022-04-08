package main

import (
	//"encoding/json"
	"encoding/json"
	"flag"
	"fmt"
	"log"

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

	order:= getOrdeResponse{
		Order_uid:         "b4b9564ааtfew3t",
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Del:               test.Delivery{
			Delivery_Id: 1,
			Name:        "Test Testov",
			Phone:       "+9720000000",
			Zip:         "2639809",
			City:        "Kiryat Mozkin",
			Address:     "Ploshad Mira 15",
			Region:      "Kraiot",
			Email:       "test@gmail.com",
		},
		Paym:              test.Payment{
			Id:           1,
			Transaction: "b563febауц7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []test.Item{
			{
				ChrtId: 9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price: 453,
				Rid: "ab4219ма40ctest",
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
				Rid: "ab4214aew2test",
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

	b, err := json.Marshal(order)
    if err != nil {
        fmt.Println(err)
        return
    }

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

	msg:=[]byte(string(b))
	err = sc.Publish(subj, msg)
	if err != nil {
		log.Fatalf("Error during publish: %v\n", err)
	}
	//json.Unmarshal(msg, &order)
	//log.Println(order)
	log.Printf("Published [%s] : '%s'\n", subj, msg)
}