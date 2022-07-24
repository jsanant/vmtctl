BIN := ./bin/vmtctl
DOCKER_COMPOSE_DOWN = (docker-compose down)
DOCKER_REMOVE_VOLUMES = (docker volume rm `docker volume ls -q`)

# Build only binary
.PHONY: build
build:
	@echo "####### Building the binary #######"
	go build -o ${BIN} .

# Generate endpoints and bring up clustered victoria-metrics
.PHONY: run
run:
	@echo "####### Running the binary #######"
	${BIN}
	@echo "####### Bringing up containers #######"
	docker-compose up -d


# Build a new binary, run it and bring up clustered victoria-metrics
.PHONY: dev
dev: build run

# Clean up after testing
.PHONY: clean
clean:
	@echo "###### Removing containers #######"
	$(DOCKER_COMPOSE_DOWN)

	@echo "###### Removing volumes #######"
	$(DOCKER_REMOVE_VOLUMES)