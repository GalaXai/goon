BINARY_PATH=./logic/golang/bin/endpoints.out

LOGIC_PATH=./logic/golang/

build:
	cd ${LOGIC_PATH} && go build -o bin/endpoints.out

run: build
	@${BINARY_PATH}