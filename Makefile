
run:
	go run main.go

build:
	rm -rf build
	mkdir build
	go build -o build/dashboard

install:
	go install

.PHONY: run
.PHONY: build
.PHONY: install
