all: client server

re: clean client server

clean:
	rm client.o server.o

client:
	go build -o client.o client.go

server:
	go build -o server.o server.go
