export DOCKER_BUILDKIT=1
GPC_PROJECT_ID=my-cloud-collection
SERVICE_NAME=logger-enpoint
CONTAINER_NAME=europe-west1-docker.pkg.dev/$(GPC_PROJECT_ID)/cloud-run/$(SERVICE_NAME)

.PHONY: *

run: build
	docker run -p 8080:8080 $(CONTAINER_NAME)
build:
	docker build -t $(CONTAINER_NAME) .
push: build
	docker push $(CONTAINER_NAME)
deploy: build push
	gcloud run deploy $(SERVICE_NAME)\
		--project $(GPC_PROJECT_ID)\
		--allow-unauthenticated\
		-q\
		--region europe-west1\
		--platform managed\
		--memory 128Mi\
		--image $(CONTAINER_NAME)
test:
	go test ./...
