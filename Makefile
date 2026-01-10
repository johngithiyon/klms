include internal/api/config/.env

MAIN = cmd/server/main.go

migrate-create:
	migrate create -ext sql -dir internal/api/storage/migrations/ $(name)

migrate-up:
	migrate -path internal/api/storage/migrations/ -database "$(CONN_STR)" up

migrate-down:
	migrate -path internal/api/storage/migrations/ -database "$(CONN_STR)" down 1

migrate-version:
	migrate -path internal/api/storage/migrations/ -database "$(CONN_STR)" version

migrate-cure:
	migrate -path internal/api/storage/migrations/ -database "$(CONN_STR)" force $(verno)

start:
	sudo docker start mypostgres
	sudo docker start myminio
	sudo docker start myredis	
	sudo docker start rabbitmq		
	sudo systemctl start nginx																																																									

stop:
	sudo docker stop mypostgres
	sudo docker stop myminio
	sudo docker stop myredis
	sudo docker stop rabbitmq
	sudo systemctl stop nginx

run:
	go run $(MAIN)