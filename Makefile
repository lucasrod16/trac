.PHONY: test
test:
	go test -v -race -failfast -count=1 -cover ./...