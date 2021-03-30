default: test

test:
	go test -count=1 -cover ./...
.PHONY: test