
dep:
	go mod tidy
	go mod download

build: dep bin/descriptinator

bin/descriptinator:
	go build -o bin/descriptinator cmd/main.go
	chmod +x bin/descriptinator