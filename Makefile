run: archai
	./archai $(ARGS)


archai: $(shell find . -type f -regex .*go$)
	go build .

siege:
	bash scripts/siege.sh

test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress
