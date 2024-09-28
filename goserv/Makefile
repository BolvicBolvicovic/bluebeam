all	:
	mkdir -p ./volumes
	mkdir -p ./volumes/mariadb
	docker-compose -f docker-compose.yml up -d --build

stop	:
	docker-compose -f docker-compose.yml stop

clean	:
	docker-compose -f docker-compose.yml down -v

fclean	: clean
	docker system prune -af
	rm -rf volumes

re	: fclean all

connect_db:
	docker exec -it mariadb mysql -uroot -proot mariadb

.PHONY	: all stop clean fclean re
