include .env

migrate_up:
	docker exec -i $$(docker ps | grep app | awk '{{ print $$1 }}') ./migrate -source file://./migrations -database $(DSN) -verbose up

migrate_prev:
	docker exec -i $$(docker ps | grep app | awk '{{ print $$1 }}') ./migrate -source file://./migrations -database $(DSN) -verbose down 1

.PHONY: migrate_up migrate_prev
