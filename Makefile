
.PHONY dep
dep:
	go mod tidy
	go mod download

.PHONY build
build: dep
	go build .cmd/main -o bin/descriptinator
	chmod +x bin/descriptinator