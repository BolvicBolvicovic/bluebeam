#!/bin/bash

if [ ${MODE} == "dev" ]; then
	mkdir -p /etc/ssl/certs;
	mkdir -p /etc/ssl/private;
	
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/private/selfsigned.key -out /etc/ssl/certs/selfsigned.crt -subj "/CN=bluebeam.dev";
	
	export SSL_CERT="/etc/ssl/certs/selfsigned.crt";
	export SSL_KEY="/etc/ssl/private/selfsigned.key";
elif [ ${MODE} == "prod" ]; then
	certbot certonly --standalone \
		--agree-tos \
        	--non-interactive \
        	--email victor.bolheme@gmai.com \
		-d bluebeam.dev;
	export SSL_CERT="/etc/letsencrypt/live/bluebeam.dev/fullchain.pem";
	export SSL_KEY="/etc/letsencrypt/live/bluebeam.dev/privkey.pem";
fi

while ! nc -z mariadb 3306; do
	sleep 0.1
done

./app
