.PHONY: console
console:
	hasura console --admin-secret=secret --no-browser --skip-update-check

.PHONY: db.create
db.create:
	docker-compose up -d db
	docker-compose exec db sh -c "PGPASSWORD=secret createdb postgres -h db -U hasura" || exit $$(($$? - 1))

.PHONY: db.create.test
	docker-compose up -d db
	docker-compose exec db sh -c "PGPASSWORD=secret createdb test -h db -U hasura" || exit $$(($$? - 1))

.PHONY: db.migrate
db.migrate:
	docker-compose up -d hasura
	until curl -s -o /dev/null http://127.0.0.1:8080/healthz; do sleep 1; done
	hasura migrate apply --database-name default --admin-secret=secret --skip-update-check

.PHONY: db.migrate.test
db.migrate.test:
	HASURA_GRAPHQL_DATABASE_URL=postgres://hasura:secret@db:5432/test docker-compose up -d hasura
	until curl -s -o /dev/null http://127.0.0.1:8080/healthz; do sleep 1; done
	hasura migrate apply --database-name default --admin-secret=secret --skip-update-check

.PHONY: db.drop
db.drop:
	docker-compose up -d db
	docker-compose exec db sh -c "PGPASSWORD=secret dropdb postgres -h db -U hasura" || exit $$(($$? - 1))

.PHONY: db.drop.test
db.drop.test:
	docker-compose up -d db
	docker-compose exec db sh -c "PGPASSWORD=secret dropdb test -h db -U hasura" || exit $$(($$? - 1))

.PHONY: db.status
db.status:
	hasura migrate status --database-name=default --admin-secret=secret --skip-update-check

.PHONY: db.status.test
db.status.test:
	HASURA_GRAPHQL_DATABASE_URL=postgres://hasura:secret@db:5432/test docker-compose up -d hasura
	hasura migrate status --database-name=default --admin-secret=secret --skip-update-check

.PHONY: db.reset
db.reset: db.drop db.create	db.migrate

.PHONY: db.reset.test
db.reset.test: db.drop.test db.create.test db.migrate.test

.PHONY: setup_es
setup_es:
	docker-compose up -d elasticsearch
	curl -XPUT 'http://localhost:9200/go-elastic-search_item' -H 'kbn-xsrf: reporting' -H 'Content-Type: application/json'
