build:
	@echo "[!] BUILDING GO-GATOR"
	go build -o ./bin/go-gator.exe

run: build
	@echo "[!] RUNNING GO-GATOR"
	./bin/go-gator.exe fetch --sources abc

test:
	@echo "[!] RUNNING GO-GATOR TESTS"
	go test ./cmd/parsers/... -v