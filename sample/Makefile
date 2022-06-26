.PHONY: db generate-errata

db:
	which sqlite3 >/dev/null || { echo >&2 "sqlite3 not installed"; exit 1; }
	rm users.sqlite3
	cat database.sql | sqlite3 users.sqlite3

generate-errata:
	$(eval TMPFILE := $(shell mktemp))
	../eish generate --source=errata.hcl --template golang --package errata > $(TMPFILE) && gofmt $(TMPFILE) > errata/errors.go
	rm $(TMPFILE)