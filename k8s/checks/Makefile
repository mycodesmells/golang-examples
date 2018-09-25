build-worker:
	docker build -t mycodesmells/k8s-checks-worker:2 worker/

deploy-worker: build-worker
	kubectl apply -f worker/k8s-worker.yaml

build-master:
	docker build -t mycodesmells/k8s-checks-master:1 master/

deploy-master: build-master
	kubectl apply -f master/k8s-master.yaml

deploy-cache:
	helm install --name cache --set usePassword=false --set cluster.enabled=false stable/redis
