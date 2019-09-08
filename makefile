dep:
	go get

format:
	go fmt .
	go fmt ./app/...

clean:
	rm shiphand

build:
	go build -o=shiphand *.go
