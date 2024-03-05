tidy:
	go mod tidy

build:
	go build -o build/go-bagit main/main.go

install:
	sudo cp build/go-bagit /usr/local/go/bin/

test:
	gotest -v
	
clean:
	rm -rf build/go-bagit

