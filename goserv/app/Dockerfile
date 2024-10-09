FROM golang:1.23

RUN apt-get update && apt-get -y install netcat-openbsd openssl ca-certificates python3 python3-pip python3.11-venv && python3 -m venv /venv

RUN /venv/bin/pip install --upgrade pip && \
    /venv/bin/pip install openai

ENV PATH="/venv/bin:$PATH"

# Set the working directory
WORKDIR /app

RUN go env -w GOPROXY=direct

COPY go.mod ./
COPY go.sum ./

RUN mkdir -p certs
RUN openssl req -newkey rsa:4096 -days 365 -nodes -x509 -sha256 \
	-subj "/C=SE/ST=Stockholm/L=Stockholm/O=bolvicstudio/OU=cmaden/CN=localhost" \
	-keyout certs/server.key -out certs/server.pem

RUN mkdir -p /usr/local/share/ca-certificates

RUN cp certs/server.pem /usr/local/share/ca-certificates/server.crt && update-ca-certificates
RUN openssl req -newkey rsa:4096 -nodes \
	-subj "/C=SE/ST=Stockholm/L=Stockholm/O=bolvicstudio/OU=cmaden/CN=localhost" \
	-keyout server.key -out server.csr

RUN echo "[ req ]" > server.ext && \
    echo "distinguished_name = req_distinguished_name" >> server.ext && \
    echo "[ req_distinguished_name ]" >> server.ext && \
    echo "[ v3_req ]" >> server.ext && \
    echo "authorityKeyIdentifier=keyid,issuer" >> server.ext && \
    echo "basicConstraints=CA:FALSE" >> server.ext && \
    echo "keyUsage=digitalSignature,nonRepudiation,keyEncipherment,dataEncipherment" >> server.ext && \
    echo "subjectAltName=@alt_names" >> server.ext && \
    echo "[alt_names]" >> server.ext && \
    echo "DNS.1=localhost" >> server.ext


RUN openssl x509 -req -in server.csr -CA ./certs/server.pem -CAkey ./certs/server.key \
	-CAcreateserial \
	-out server.crt -days 825 -sha256 -extfile server.ext


RUN go mod download && go mod verify

COPY . .

EXPOSE 3333

CMD ["./lauch_serv.sh"]
