BINARY= bin/endpoints.out
LOGIC_PATH=./logic/golang/
FRONT_PATH=./front

build-backend:
	@ cd ${LOGIC_PATH} && go build -o ${BINARY}

run-backend: build-backend
	@${LOGIC_PATH}${BINARY}

build-frontend:	
	@cd ${FRONT_PATH} && npm install && npm run build

dev-frontend:
	@cd ${FRONT_PATH} && npm run dev

build: build-backend build-frontend

run: run-backend

run-all: run-backend dev-frontend