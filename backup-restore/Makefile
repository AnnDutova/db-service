
start:
	kubectl create -f deployment/config-map-init.yaml -f deployment/config-map-pas.yaml -f deployment/sts.yaml -f deployment/service.yaml

delete:
	kubectl delete service/db-service
	kubectl delete configmap db-config tables-config
	kubectl delete statefulset.apps/postgres

postgres:
	docker-compose -f postgres.docker-compose.yaml --env-file=config/.env.dev up -d
	docker-compose up -d


