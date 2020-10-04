run: archai
	./archai $(ARGS)


archai: $(shell find . -type f -regex .*go$)
	go build .

siege:
	bash scripts/siege.sh

docker-compose: docker-compose.yml
	docker-compose up -d

migrate-test: archai
	./archai --migrate --keyspace archai_test

test: docker-compose migrate-test
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress
