GOOSE_COMMAND=goose -dir ./config/migrations mysql "root:password@/stock_ratings"

build:
	templ generate
	sqlc generate -f ./config/sqlc.yaml
	go build -o ./tmp/stock-management.go cmd/main.go

docker:
	docker build -f dev.Dockerfile --name 'go-stock-management'

mysql: mysql-down
	docker run --name go-stock-management-mysql -v ./mysql_data:/var/lib/mysql -e 'MYSQL_ROOT_PASSWORD=password' -p 3306:3306 --detach 'mysql:latest' 

mysql-down:
	docker rm -f go-stock-management-mysql

goose-create:
	goose -dir ./config/migrations -s mysql "root:password@/stock_ratings" create migration sql

goose-up:
	goose -dir ./config/migrations -s mysql "root:password@/stock_ratings" up

goose-down:
	goose -dir ./config/migrations -s mysql "root:password@/stock_ratings" down

run: build mysql
	./tmp/stock-management