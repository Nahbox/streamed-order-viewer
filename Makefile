nats-serv:
	cd nats-streaming/ && docker compose up -d
	sleep 5

publish:
	cd nats-streaming/ && go run cmd/main.go
	sleep 5

postgres:
	cd service/ && docker compose up postgres -d
	sleep 5

server:
	cd service/ && go run cmd/main.go

all: nats-serv publish postgres server

clean:
	cd nats-streaming/ && docker compose down
	cd service/ && docker compose down