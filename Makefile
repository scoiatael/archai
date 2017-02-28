run: archai
	./archai


archai: $(shell find . -type f -regex .*go$)
	go build .

siege:
	bash scripts/siege.sh

test:
	go test
	cd persistence && go test
