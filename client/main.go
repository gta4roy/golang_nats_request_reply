package main

import (
	"log"

	pb "natsapp/order"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
)

func main() {
	// Create NATS server connection
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	defer natsConnection.Close()
	log.Println("Connected to " + nats.DefaultURL)
	msg, err := natsConnection.Request("Discovery.OrderService", nil, 1000*time.Millisecond)
	if err == nil && msg != nil {
		var orderServiceDiscovery pb.ServiceDiscovery
		err := proto.Unmarshal(msg.Data, &orderServiceDiscovery)
		if err != nil {
			log.Fatalf("Error on unmarshal: %v", err)
		}
		address := orderServiceDiscovery.Orderserviceuri
		log.Println("OrderService endpoint found at:", address)
	}

	log.Println("Client Requesting to fetch the system info from the NATS Server ")
	msgSecondRequest, err1 := natsConnection.Request("Discovery.GetSystemInfo", nil, 1000*time.Millisecond)
	if err1 == nil && msgSecondRequest != nil {
		var systemInfo pb.GetSystemTime
		err := proto.Unmarshal(msgSecondRequest.Data, &systemInfo)
		if err != nil {
			log.Fatalf("Error on unmarshal: %v", err)
		}

		log.Println(" System Info Details :")
		log.Println("System Time ", systemInfo.Systemtime)
		log.Println("System Date ", systemInfo.Systemdate)
		log.Println("System Username ", systemInfo.Username)
		log.Println("System Server Ip ", systemInfo.Serverip)
	}

}
