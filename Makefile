.PHONY: build clean

build: clean
	@echo "======================== Building Binary ======================="
	CGO_ENABLED=0 go build -ldflags="-s -w" -v -o dist/ .

clean:
	@echo "======================== Cleaning Project ======================"
	go clean
	rm -f dist/*