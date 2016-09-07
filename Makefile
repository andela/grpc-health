NUMBER_OF_COMMIT=$(shell git rev-list HEAD --count)
SHA_COMMIT= $(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date +%Y%m%d-%H%M%S)

TAG = $(SHA_COMMIT)-$(NUMBER_OF_COMMIT)
PROJECT = microservices-kube
IMAGE = grpc_health
HEALTH_IMAGE = us.gcr.io/$(PROJECT)/$(IMAGE):$(TAG)

all: run

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o artifact .

image: 
	docker build -t $(HEALTH_IMAGE) -t us.gcr.io/$(PROJECT)/$(IMAGE):latest .

push: image
	gcloud docker push $(HEALTH_IMAGE)

minikube: 
	eval $$(minikube docker-env) && docker build -t us.gcr.io/$(PROJECT)/$(IMAGE):minikube .


run:
	# docker run  us.gcr.io/$(PROJECT)/$(IMAGE):$(TAG)
	PORT=8000 go run main.go

clean: