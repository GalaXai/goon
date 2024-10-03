BINARY_PATH=./logic/golang/bin/endpoints

LOGIC_PATH=./logic/golang/

.PHONY: dev

dev:
	@echo "Starting development server..."
	@sass --watch ./static/style.scss:static/style.css & air

build:
	cd ${LOGIC_PATH} && go build -o bin/endpoints

run: build
	@${BINARY_PATH}