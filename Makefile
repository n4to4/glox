generated = src/expr.go src/stmt.go

.PHONY: all
all: $(generated)
	go build -o glox ./src

.PHONY: repl
repl:
	go run ./src

$(generated): tool/generate.go
	go run tool/generate.go ./src

.PHONY: test
test: $(generated)
	go test ./...

.PHONY: test2
test2: all test
	./glox examples/fibonacci.lox

.PHONY: clean
clean:
	rm -r $(generated)
