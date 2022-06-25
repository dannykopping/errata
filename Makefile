.PHONY: eish generate-errata

eish:
	go build -o eish cmd/eish/*.go

generate-errata:
	$(eval TMPFILE := $(shell mktemp))
	./eish generate --source=errata.hcl --template golang --package errata > $(TMPFILE) && gofmt $(TMPFILE) > errata.go
	rm $(TMPFILE)

test:
	go test ./... -count=1