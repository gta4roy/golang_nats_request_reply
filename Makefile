all:
	compile, build

compile:
	protoc -I order/ order/order.proto --go_out=plugins=grpc:order

build:
	go build -o natsapp/client natsapp/client
	go build -o natsapp/discovery natsapp/discovery
	cp -r discovery/config natsapp/

