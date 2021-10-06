build:
	./scripts/build.bash $(version)

build_mac:
	env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/$(version)-binance_bot-darwin-amd64 -ldflags "-s -w"  main.go