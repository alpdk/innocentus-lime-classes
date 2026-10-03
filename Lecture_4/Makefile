.PHONY: build
build: bin/task1-file-delivery
build: bin/task2-hello-world

./bin/%: ./%/*.go
	@echo "Build $@"
	go build -o $@ -v -trimpath ./$(dir $<)