bin_name=sol-rpc-router

build:
	@go build -ldflags="-s -w" -o ./bin/$(bin_name) main.go 
	
run: build
	@./bin/$(bin_name)	

test:
	@go clean -testcache
	@go test ./... -v
