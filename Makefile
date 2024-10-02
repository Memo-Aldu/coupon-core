build:
	@go build -o bin/couponcore

run: build 
	@./bin/couponcore

test:
	@go test -v ./...
	

