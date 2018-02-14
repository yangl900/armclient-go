all: 
	go test -v
	GOOS=linux go build -o armclient .

clean:
	rm -rf dist/
	rm armclient