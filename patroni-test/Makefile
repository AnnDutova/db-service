patroni-build:
	docker build -t patroni .

patroni-run: patroni-build
	kubectl create configmap init-config --from-file=./db-init-scripts/
	kubectl apply -f service-conf.yaml
	kubectl apply -f sts.yaml
	kubectl apply -f endpoint.yaml
	kubectl apply -f service.yaml
	kubectl apply -f service-repl.yaml
	kubectl apply -f secret.yaml
	kubectl apply -f service-acc.yaml
	kubectl apply -f role.yaml
	kubectl apply -f role-bind.yaml
	kubectl apply -f cluster-role.yaml
	kubectl apply -f cluster-role-bind.yaml

patroni-delete:
	kubectl delete sts patroni-psql
	kubectl delete configmap init-config
	kubectl delete service patroni-psql
	kubectl delete service patroni-psql-config
	kubectl delete service patroni-psql-repl

haproxy-build:
	docker build -t haproxy-cfg ./haproxy

haproxy-run: haproxy-build
	kubectl apply -f haproxy/haproxy-deployment.yaml
	kubectl apply -f haproxy/haproxy-service.yaml
	kubectl apply -f haproxy/haproxy-service-admin.yaml

haproxy-delete:
	kubectl delete service haproxy-admin-svc
	kubectl delete service haproxy-svc
	kubectl delete deployment haproxy

consul-build:
	docker build -t consul ./consul

consul-run: consul-build
	kubectl apply -f consul/service.yaml
	kubectl apply -f consul/sts.yaml

consul-delete:
	kubectl delete sts consul
	kubectl delete service consul

run: consul-run haproxy-run patroni-run
	kubectl get pods -L role

delete: consul-delete haproxy-delete patroni-delete

patroni-test:
	docker build -t patroni -f Dockerfile.patroni .
	kubectl create configmap init-config --from-file=./db-init-scripts/
	kubectl apply -f service-conf.yaml
	kubectl apply -f patroni-conf.yaml
	kubectl apply -f sts.yaml

.PHONY: patroni-run