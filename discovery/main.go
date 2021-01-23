package main

import (
	"fmt"
	"log"
	"runtime"

	pb "natsapp/order"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/spf13/viper"
)

var orderServiceUri string

func init() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found")
	}
	orderServiceUri = viper.GetString("discovery.orderservice")
	timeoutValue := viper.GetInt32("server.timeoutvalue")
	serverip := viper.GetString("server.serverip")
	issecure := viper.GetBool("server.issecure")

	fmt.Println(timeoutValue)
	fmt.Println(serverip)

	if issecure {
		fmt.Println("Is Secure ")
	}
}

func main() {
	fmt.Println(orderServiceUri)

	// Create server connection
	natsConnection, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Println("Not able to Connect NATS Server " + nats.DefaultURL)
	} else {
		log.Println("Able to Connect NATS Server " + nats.DefaultURL)
	}

	natsConnection.Subscribe("Discovery.OrderService", func(m *nats.Msg) {
		orderServiceDiscovery := pb.ServiceDiscovery{Orderserviceuri: orderServiceUri}
		data, err := proto.Marshal(&orderServiceDiscovery)

		if err == nil {
			natsConnection.Publish(m.Reply, data)
		}
	})

	natsConnection.Subscribe("Discovery.GetSystemInfo", func(m *nats.Msg) {
		viper.SetConfigName("app")
		viper.AddConfigPath("config")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal("Config file not found")
		}

		SystemInfo := pb.GetSystemTime{}
		SystemInfo.Systemtime = viper.GetString("systeminfo.systemtime")
		SystemInfo.Systemdate = viper.GetString("systeminfo.systemdate")
		SystemInfo.Username = viper.GetString("systeminfo.username")
		SystemInfo.Serverip = viper.GetString("systeminfo.serverip")

		log.Println(SystemInfo)

		data1, err := proto.Marshal(&SystemInfo)

		if err == nil {
			natsConnection.Publish(m.Reply, data1)
			log.Println("Sent ", data1)
		}
	})

	// Keep the connection alive
	runtime.Goexit()
}
