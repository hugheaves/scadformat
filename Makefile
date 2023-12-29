scadformat: test cmd/*.go internal/parser/*.go internal/formatter/*.go internal/logutil/*.go
	go build cmd/scadformat.go

internal/parser/*.go: OpenSCAD.g4
	go generate ./...

test: internal/parser/*.go
	go test ./...

clean:
	rm -f internal/parser/*
	rm -rf internal/parser/.antlr
	rm cmd/version.txt

.PHONY: clean test