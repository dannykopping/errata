build-eish:
	go build -o eish cmd/eish/*.go

errata:
	./eish generate --eds.file=errata.hcl --template golang --package errata > /tmp/x && gofmt /tmp/x > errors.go