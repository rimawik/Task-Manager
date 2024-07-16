GOPATH := $(if $(GOPATH),$(GOPATH),$(shell go env GOPATH))

build:
	go build -o ./bin/main ./cmd/app/main.go

run:
	go run ./cmd/app/main.go

swagger:
	$(GOPATH)/bin/swag init -d cmd/app/,./routes/,./models/,./handlers/

test:
	go test ./handlers/task_test.go