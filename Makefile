all: client server

client:
	go build -o client.o client.go

server:
	go build -o server.o server.go
