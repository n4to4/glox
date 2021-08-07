generated = src/expr.go

.PHONY: all
all: $(generated)
	go run ./src

$(generated):
	go run tool/generate.go ./src

.PHONY: test
test: $(generated)
	go test ./...

.PHONY: t
t:
	make test

.PHONY: clean
clean:
	rm -r $(generated)
