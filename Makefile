NUMBER_OF_COMMIT=$(shell git rev-list HEAD --count)
SHA_COMMIT= $(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date +%Y%m%d-%H%M%S)

TAG = $(SHA_COMMIT)-$(NUMBER_OF_COMMIT)
PROJECT = andela-kube
IMAGE = health_check
HEALTH_IMAGE = us.gcr.io/$(PROJECT)/$(IMAGE)

all: run

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o artifact .

image: build
	docker build -t $(HEALTH_IMAGE) -t us.gcr.io/$(PROJECT)/$(IMAGE):latest .

push: image
	gcloud docker push $(HEALTH_IMAGE)

minikube:
	eval $$(minikube docker-env) && docker build -t us.gcr.io/$(PROJECT)/$(IMAGE) .


run:
	# docker run  us.gcr.io/$(PROJECT)/$(IMAGE):$(TAG)
	PORT=8000 go run main.go

clean:
