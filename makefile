BUILD_CMD=go build -o ./tmp/stock-management cmd/main.go
GOOSE_COMMAND=goose -dir ./config/migrations mysql "root:password@/stock_ratings"
MYSQL=root:password@tcp(localhost:3306)/stock_ratings

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

build: build-deps
	$(BUILD_CMD)

build-windows: build-deps
	GOOS=windows GOARCH=amd64 $(BUILD_CMD)

build-mac-intel:build-deps
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD)

build-mac-apple:build-deps
	GOOS=darwin GOARCH=arm64 $(BUILD_CMD)

build-linux:build-deps
	GOOS=linux GOARCH=amd64 $(BUILD_CMD)

build-linux-32:build-deps
	GOOS=linux GOARCH=386 $(BUILD_CMD)

build-deps:
	templ generate
	sqlc generate -f ./config/sqlc.yaml

docker:
	docker build -f dev.Dockerfile --name 'go-stock-management'

mysql: mysql-down
	docker run --name go-stock-management-mysql -v ./mysql_data:/var/lib/mysql -e 'MYSQL_ROOT_PASSWORD=password' -e 'MYSQL_DATABASE=stock_ratings' -p 3306:3306 --detach 'mysql:latest' 

mysql-down:
	docker rm -f go-stock-management-mysql

goose-create:
	goose -dir ./config/migrations -s mysql "$(MYSQL)" create migration sql

goose-up:
	goose -dir ./config/migrations -s mysql "$(MYSQL)" up

goose-down:
	goose -dir ./config/migrations -s mysql "$(MYSQL)" down

run: export MYSQL_CONNECTION_STRING = $(MYSQL)
run: export SIGNING_SECRET = supersecret
run: export ADMIN_USERNAME = admin
run: export ADMIN_PASSWORD = trust
run: build mysql
	./tmp/stock-management
