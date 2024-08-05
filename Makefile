.PHONY:
run:
	cd cmd && go build -o server.bin && ./server.bin

.PHONY:
test:
	go test -count=1 ./tests