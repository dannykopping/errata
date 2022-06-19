.PHONY: eish generate-errata

eish:
	go build -o eish cmd/eish/*.go

generate-errata:
	./eish generate --source=errata.hcl --template golang --package errata > /tmp/x && gofmt /tmp/x > errata.go