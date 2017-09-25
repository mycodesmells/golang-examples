gen_proto:
	protoc -I. --go_out=plugins=grpc:${GOPATH}/src proto/message/message.proto

run/employer:
	ADDR=localhost:8001 EMPLOYEE_ADDR=localhost:8002 go run employer/main.go

run/employee:
	ADDR=localhost:8002 go run employee/main.go

compile/employer:
	echo "Building employer golang binary"
	GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o employer/bin/employer employer/main.go

compile/employee:
	echo "Building employee golang binary"
	GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o employee/bin/employee employee/main.go

build/employer: compile/employer
	eval $(minikube docker-env)
	echo "Building employer Docker image"
	docker build -t mcs/lb-employer employer/

build/employee: compile/employee
	eval $(minikube docker-env)
	echo "Building employee Docker image"
	docker build -t mcs/lb-employee employee/

k8s-docker-env:
	eval $(minikube docker-env)

deploy/employer: build/employer
	kubectl apply -f employer/manifests/employer.yaml
	kubectl apply -f employer/manifests/employer-svc.yaml

deploy/employee: build/employee
	kubectl apply -f employee/manifests/employee.yaml
	kubectl apply -f employee/manifests/employee-svc.yaml
