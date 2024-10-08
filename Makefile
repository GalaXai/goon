BINARY= bin/endpoints.out
LOGIC_PATH=./logic/
FRONT_PATH=./front

build-backend:
	@ cd ${LOGIC_PATH} && go build -o ${BINARY}
	$(info Finished building backend)


run-backend: build-backend
	@${LOGIC_PATH}${BINARY} &
build-frontend:	
	@cd ${FRONT_PATH} && npm install && npm run build

dev-frontend:
	@cd ${FRONT_PATH} && npm run dev

build: build-backend build-frontend

run: build-backend run-backend

run-no-build:
	@${LOGIC_PATH}${BINARY} &

kill-backend:
	@pkill -f "${BINARY}" || echo "Backend process not found"

run-all:
	@make -j2 build run-backend dev-frontend