FROM golang:1.23

RUN apt-get update && apt-get -y install netcat-openbsd openssl

# Set the working directory
WORKDIR /app

RUN go env -w GOPROXY=direct

COPY go.mod ./
COPY go.sum ./

RUN openssl req -newkey rsa:4096 -days 365 -nodes -x509 -sha256 \
	-subj "/C=SE/ST=Stockholm/L=Stockholm/O=bolvicstudio/OU=cmaden/CN=localhost" \
	-keyout server.key -out server.crt

RUN go mod download && go mod verify

COPY . .

EXPOSE 3333

CMD ["./lauch_serv.sh"]
