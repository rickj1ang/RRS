.PHONY: run
run: build
	bin/rrs

.PHONY: build
build:
	go build -o bin/rrs ./cmd 

.PHONY: push
push:
	git push -u origin main
