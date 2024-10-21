BINARY= bin/endpoints.out
LOGIC_PATH=./logic/
FRONT_PATH=./front

build-backend:
	@ cd ${LOGIC_PATH} && go build -o ${BINARY}
	$(info Finished building backend)


watch-backend:
	@echo "Watching for changes in backend files..."
	@while true; do \
		make run-backend; \
		inotifywait -q -e modify -r ${LOGIC_PATH}; \
		make kill; \
	done

run-backend: build-backend
	@${LOGIC_PATH}${BINARY} &
build-frontend:	
	@cd ${FRONT_PATH} && bun install && bun run build

run-frontend:
	@cd ${FRONT_PATH} && bun run dev

build: build-backend build-frontend

run: kill build-backend run-backend run-frontend

run-no-build:
	@${LOGIC_PATH}${BINARY} &

kill:
	@-pkill -f "${BINARY}" || true
	@echo "Attempted to kill backend process"

run-all:
	@make -j2 build run-backend dev-frontend