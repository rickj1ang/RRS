.PHONY: run
run: build
	bin/main

.PHONY: build
build:
	go build -o bin/main cmd/main.go 

.PHONY: push
push:
	git push -u origin main
