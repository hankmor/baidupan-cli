GO111MODULE=on
CGO_ENABLED=0

.PHONY: test
test:
	go test ./...

.PHONY: build_linux
build_linux: clean copy_res
	@echo 'building...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/baidupan-cli .
	@echo 'success!'

.PHONY: build_windows
build_windows: clean copy_res
	@echo 'building...'
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/baidupan-cli.exe .
	@echo 'success!'

.PHONY: build_mac
build_mac: clean copy_res
	@echo 'building...'
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/baidupan-cli .
	@echo 'success!'

.PHONY: clean
clean:
	@rm -f bin/*

.PHONY: copy_res
copy_res:
	@cp ./config.yaml ./bin/
	@echo 'copying resource...'