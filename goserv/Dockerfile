FROM	golang:1.23

COPY	go.mod go.sum ./
RUN		go mod download && go mod verify

COPY	. .
RUN		go build -v -o app ./...

CMD		["app"]
