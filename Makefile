include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the ./cmd/app application
.PHONY: run
run:
	go run ./cmd/app

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	/home/sergey/go/bin/staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build: build the cmd/api application
.PHONY: build
build:
	@echo 'Building cmd/app...'
	go build -ldflags=${linker_flags} -o=./bin/app ./cmd/app
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/app ./cmd/app

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = "95.143.0.106"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh developer@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api developer@${production_host_ip}:~/notification
	rsync -rP --delete ./migrations developer@${production_host_ip}:~
	rsync -P ./remote/production/notification.service developer@${production_host_ip}:~
	ssh -t developer@${production_host_ip} '\
		cd ~/notification \
		&& sudo mv ~/notification.service /etc/systemd/system/ \
		&& sudo systemctl enable notification \
		&& sudo systemctl restart notification \
		&& sudo service apache2 restart \
	'
