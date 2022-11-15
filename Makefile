CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/ninashvl/homework-1/cmd/bot

all: format build test lint

.PHONY: build
build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	go run ${PACKAGE}

.PHONY: generate
generate: install-mockgen
	${MOCKGEN} -source=internal/messages/incoming_msg.go -destination=internal/mocks/messages/messages_mocks.go
	-source=internal/storage/expense_storage/istorage.go -destination=internal/storage/expense_storage/mock/storage.go

.PHONY: lint
lint: install-lint
	${LINTBIN} run

.PHONY: precommit
precommit: format build test lint
	echo "OK"

.PHONY: bindir
bindir:
	mkdir -p ${BINDIR}

.PHONY: format
format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

.PHONY: install-mockgen
install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

.PHONY: install-lint
install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

.PHONY: install-smartimports
install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

.PHONY: migrate-up
migrate-up:
	docker run -v $(CURDIR)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://nina:qwerty@localhost:5003/fin_db?sslmode=disable up

.PHONY: migrate-down
migrate-down:
	docker run -v $(CURDIR)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://nina:qwerty@localhost:5003/fin_db?sslmode=disable down $(MIGRATE_NUM)

.PHONY: migrate-create
migrate-create:
	docker run -v $(CURDIR)/migrations:/migrations --network host migrate/migrate -path=/migrations/  create -ext=sql -dir=/migrations $(MIGRATE_NAME)

.PHONY: docker-run
docker-run:
	sudo docker compose up

ifeq (migrate-create,$(firstword $(MAKECMDGOALS)))

  MIGRATE_NAME:= $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

  # ...and turn it into do-nothing target

  $(eval $(MIGRATE_NAME):;@:)

endif


ifeq (migrate-down,$(firstword $(MAKECMDGOALS)))

  MIGRATE_NUM:= $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

  # ...and turn it into do-nothing target

  $(eval $(MIGRATE_NUM):;@:)

endif