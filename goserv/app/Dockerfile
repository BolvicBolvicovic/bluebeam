FROM golang:1.23

RUN apt-get update && apt-get -y install netcat-openbsd certbot openssl ca-certificates vim python3 python3-pip python3.11-venv && python3 -m venv /venv

RUN /venv/bin/pip install --upgrade pip && \
    /venv/bin/pip install openai

RUN /venv/bin/pip install -q -U google-generativeai

ENV PATH="/venv/bin:$PATH"

# Set the working directory
WORKDIR /app

RUN export GOPROXY=https://proxy.golang.org

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o app ./cmd/main.go

ENTRYPOINT ["./launch_serv.sh"]
