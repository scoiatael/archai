run: archai
	./archai


archai: $(shell find . -type f -regex .*go$)
	go build .
