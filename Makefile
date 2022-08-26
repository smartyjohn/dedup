all: clean
	go build -o workspace/dedup github.com/smartyjohn/dedup

clean:
	mkdir -p workspace
	rm -rf workspace/*

test: clean
	go test -o workspace/dedup.test github.com/smartyjohn/dedup

.PHONY: all clean test
