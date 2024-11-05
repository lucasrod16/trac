.PHONE: build
build:
	go build -ldflags="-s -w" -o ./bin/trac

.PHONY: test
test:
	go test -v -race -failfast -count=1 -cover ./...

.PHONY: clean
clean:
	rm -rf bin/ .trac/
