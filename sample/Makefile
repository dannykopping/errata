.PHONY: generate-errata

generate-errata:
	$(eval TMPFILE := $(shell mktemp))
	../eish generate --source=errata.hcl --template golang --package errata > $(TMPFILE) && gofmt $(TMPFILE) > errata/errors.go
	rm $(TMPFILE)