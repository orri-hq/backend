# Default rule/target
.DEFAULT_GOAL := run

COMPONENT := orri-backend


# Sign the code to avoid stupid OSX filewall alerts every run
SIGN_CERT := local-dev-cert

# =============================================================================
# Provide image specific parameters
#
include Makefile_inc


.PHONY: print-env
print-env:
	@echo VERSION: $(VERSION)
	@echo IMAGE_REPOSITORY: $(IMAGE_REPOSITORY)
	@echo IMAGE_NAME: $(IMAGE_NAME)
	@echo


## Build multi-architecture docker images
.PHONY: run
run: print-env
	@echo "Running Orri Backend on Local"
	@rm -f ./$(COMPONENT)
	@go build -o $(COMPONENT) .
	@codesign -s "$(SIGN_CERT)" ./$(COMPONENT)
	@./$(COMPONENT)   


.PHONY: build-cloud
build-cloud: print-env
	@echo "Building Orri Backend for Cloud (x86_64)"
	@docker buildx build --platform linux/amd64 --rm -t $(IMAGE_NAME) --progress=plain -f ./Dockerfile .    
	@docker push $(IMAGE_NAME):cloud


## Deploy the application
.PHONY: deploy
deploy: build-cloud
	@echo "Deploying Orri Backend"
	@kubectl apply  -f ./yaml/service.yaml
	@cat ./yaml/virtual-service.yaml | sed "s~<domain>~$(DOMAIN)~g" | kubectl apply -f -
	@cat ./yaml/deployment.yaml | sed "s~<image-name>~$(IMAGE_NAME):cloud~g" | kubectl apply -f -

