
run:
	go run main.go

build:
	rm -rf build
	mkdir build
	go build -o build/dashboard

.PHONY: run
.PHONY: build
