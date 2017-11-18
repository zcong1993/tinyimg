generate:
	@go generate ./...

build: generate
	@echo "====> Build tinyimg"
	@sh -c ./build.sh
