.PHONY: dev

dev:
	@echo "Starting development server..."
	@sass --watch ../../static/style.scss:static/style.css & air

build:
	go build -o bin/endpoints

run: build
	@./bin/endpoints