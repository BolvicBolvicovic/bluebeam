FROM alpine:latest

RUN apk add mariadb mariadb-client
EXPOSE 3306

ADD tools/init.sh /scripts/init.sh
ADD tools/init-db.sh /scripts/init-db.sh
RUN	chmod +x scripts/init.sh
RUN	chmod +x scripts/init-db.sh

ADD conf/mariadb-server.cnf /etc/my.cnf.d/mariadb-server.cnf

ENTRYPOINT ["/scripts/init.sh"]
