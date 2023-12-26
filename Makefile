scadformat: cmd/*.go internal/parser/*.go internal/formatter/*.go internal/logutil/*.go
	go build cmd/scadformat.go

internal/parser/*.go: OpenSCAD.g4
	go generate ./...

.PHONY: clean

clean:
	rm -f internal/parser/*
	rm -rf internal/parser/.antlr
	rm cmd/version.txt
