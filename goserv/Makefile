all	:
	if grep -q "MODE=dev" ./app/.env; then \
		if grep -q "127.0.0.1 bluebeam.dev" /etc/hosts; then \
			echo "Local DNS already set"; \
		else \
			echo "127.0.0.1 bluebeam.dev" >> /etc/hosts; \
			echo "Local DNS set"; \
		fi \
	fi
	mkdir -p ./volumes
	mkdir -p ./volumes/mariadb
	docker-compose -f docker-compose.yml up -d --build

stop	:
	docker-compose -f docker-compose.yml stop

clean	:
	docker-compose -f docker-compose.yml down -v

fclean	: clean
	docker system prune -af
	sed -i '/127.0.0.1 bluebeam.dev/d' /etc/hosts 
	rm -rf volumes

full_re	: fclean all

re: clean all

connect_db:
	docker exec -it mariadb mysql -uroot -proot mariadb

.PHONY	: all stop clean fclean re full_re
