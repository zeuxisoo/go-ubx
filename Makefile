all:
	@echo "make run"

run:
	@go run *.go

build:
	@bash ./scripts/build.sh
