build:
	@echo "[!] BUILDING GO-GATOR"
	go build -o ./bin/go-gator.exe

run: build
	@echo "[!] RUNNING GO-GATOR"
	./bin/go-gator.exe

test:
	@echo "[!] RUNNING GO-GATOR TESTS"
	go test ./cmd/... -v