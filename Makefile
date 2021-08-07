generated = src/expr.go

all: $(generated)
	go run ./src

$(generated):
	go run tool/generate.go ./src

clean:
	rm -r $(generated)
