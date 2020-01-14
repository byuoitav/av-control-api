BUILD_DIR := .build
IGNORE_DIRS := _build
REPO_OWNER := byuoitav
REPO_NAME := av-control-api
DOCKER_REPO := docker.pkg.github.com
DOCKER_IMAGE_BASE := $(DOCKER_REPO)/$(REPO_OWNER)/$(REPO_NAME)

# If we aren't running in Github Actions, just act on all drivers
ifndef GITHUB_ACTIONS
	LIBRARY_AFFECTED := ""
	DRIVER_LIST := $(sort $(filter-out $(IGNORE_DIRS), $(patsubst %/,%,$(wildcard */))))
# If we are in Github Actions then act only on changed drivers
else
	COMMIT_HASH := $(shell git rev-parse --short $(GITHUB_SHA))
	DRIVER_LIST := $(shell _build/github-info get_modified_drivers $(GITHUB_EVENT_PATH))
	ifeq ($(shell _build/github-info is_tagged_release $(GITHUB_EVENT_PATH)), True)
		RELEASE_VERSION := $(shell _build/github-info get_version $(GITHUB_EVENT_PATH))
	endif
	LIBRARY_AFFECTED := $(shell _build/github-info is_library_modified $(GITHUB_EVENT_PATH))
	ifeq ($(LIBRARY_AFFECTED), False)
		undefine LIBRARY_AFFECTED
	endif
endif

# ------------------------------------------
# |                                        |
# |             Info Printout              |
# |                                        |
# ------------------------------------------

.PHONY: info
info:
ifdef RELEASE_VERSION
	@echo "This is a tagged release. Version: $(RELEASE_VERSION)"
else
	@echo "This is not a tagged release."
endif
ifdef COMMIT_HASH
	@echo "Commit hash: $(COMMIT_HASH)"
endif
ifdef LIBRARY_AFFECTED
	@echo "Drivers library affected: True"
else
	@echo "Drivers library affected: False"
endif
	@echo "The following drivers are affected:"
	@for driver in $(DRIVER_LIST); do \
		echo "	$$driver" ; \
	done

# ------------------------------------------
# |                                        |
# |                 Build                  |
# |                                        |
# ------------------------------------------
.PHONY: build
build: build-amd64 build-arm

# ---------- Build AMD 64 Drivers ----------

BUILD_AMD_LIST = $(patsubst %, $(BUILD_DIR)/%-amd64, $(DRIVER_LIST))
.PHONY: build-amd64
build-amd64: $(BUILD_AMD_LIST)

.SECONDEXPANSION:
$(BUILD_AMD_LIST): $(BUILD_DIR)/%-amd64: $$(wildcard %/*.go) %/go.mod %/go.sum %/
	@echo ""
	@echo "-------------------- Building $* for AMD64 --------------------"
	@echo ""
	cd $* && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../$@

# ---------- Build ARM Drivers ----------

BUILD_ARM_LIST = $(patsubst %, $(BUILD_DIR)/%-arm, $(DRIVER_LIST))
.PHONY: build-arm
build-arm: $(BUILD_ARM_LIST)

$(BUILD_ARM_LIST): $(BUILD_DIR)/%-arm: $$(wildcard %/*.go) %/go.mod %/go.sum %/
	@echo ""
	@echo "-------------------- Building $* for ARM --------------------"
	@echo ""
	cd $* && GOOS=linux GOARCH=arm go build -v -o ../$@

# ------------------------------------------
# |                                        |
# |                 Test                   |
# |                                        |
# ------------------------------------------

.PHONY: test
test:
ifdef LIBRARY_AFFECTED
	@echo ""
	@echo "-------------------- Testing Drivers Library --------------------"
	@echo ""
	@go test ./...
endif
	@for driver in $(DRIVER_LIST) ; do \
		echo "" ; \
		echo "-------------------- Testing $$driver Driver --------------------" ; \
		echo "" ; \
		cd $$driver ; \
		go test ./... || exit 2 ; \
		cd .. ; \
	done

# ------------------------------------------
# |                                        |
# |                 Lint                   |
# |                                        |
# ------------------------------------------

.PHONY: lint
lint:
ifdef LIBRARY_AFFECTED
	@echo ""
	@echo "-------------------- Linting Drivers Library --------------------"
	@echo ""
	@golangci-lint run ./...
endif
	@for driver in $(DRIVER_LIST) ; do \
		echo "" ; \
		echo "-------------------- Linting $$driver Driver --------------------" ; \
		echo "" ; \
		cd $$driver ; \
		golangci-lint run ./... || exit 2 ; \
		cd .. ; \
	done

# ------------------------------------------
# |                                        |
# |               Docker                   |
# |                                        |
# ------------------------------------------

DOCKER_LIST = $(patsubst %, %-docker, $(DRIVER_LIST))
.PHONY: docker
docker: $(DOCKER_LIST)

# Build each of the docker containers necessary
.PHONY: $(DOCKER_LIST)
$(DOCKER_LIST): %-docker: $(BUILD_DIR)/%-arm $(BUILD_DIR)/%-amd64
ifdef RELEASE_VERSION
	@echo ""
	@echo "-------------------- Building $*-driver:$(RELEASE_VERSION) Image --------------------"
	@echo ""
	@docker build -f _build/dockerfile --build-arg NAME=$*-amd64 \
	-t $(DOCKER_IMAGE_BASE)/$*-driver:latest $(BUILD_DIR)
	@docker tag $(DOCKER_IMAGE_BASE)/$*-driver:latest $(DOCKER_IMAGE_BASE)/$*-driver:$(RELEASE_VERSION)
	@echo ""
	@echo "-------------------- Building $*-driver-arm:$(RELEASE_VERSION) Image --------------------"
	@echo ""
	@docker build -f _build/dockerfile --build-arg NAME=$*-arm \
	-t $(DOCKER_IMAGE_BASE)/$*-driver-arm:latest $(BUILD_DIR)
	@docker tag $(DOCKER_IMAGE_BASE)/$*-driver-arm:latest $(DOCKER_IMAGE_BASE)/$*-driver-arm:$(RELEASE_VERSION)
else
	@echo ""
	@echo "-------------------- Building $*-driver-dev Image --------------------"
	@echo ""
	@docker build -f _build/dockerfile --build-arg NAME=$*-amd64 \
	-t $(DOCKER_IMAGE_BASE)/$*-driver-dev:latest $(BUILD_DIR)
	@echo ""
	@echo "-------------------- Building $*-driver-dev-arm Image --------------------"
	@echo ""
	@docker build -f _build/dockerfile --build-arg NAME=$*-arm \
	-t $(DOCKER_IMAGE_BASE)/$*-driver-dev-arm:latest $(BUILD_DIR)

ifdef COMMIT_HASH
	@docker tag $(DOCKER_IMAGE_BASE)/$*-driver-dev:latest $(DOCKER_IMAGE_BASE)/$*-driver-dev:$(COMMIT_HASH)
	@docker tag $(DOCKER_IMAGE_BASE)/$*-driver-dev-arm:latest $(DOCKER_IMAGE_BASE)/$*-driver-dev-arm:$(COMMIT_HASH)
endif
endif

.PHONY: docker-login
docker-login:
	@echo ""
	@echo "-------------------- Logging into Docker repository --------------------"
	@echo ""
	docker login $(DOCKER_REPO) -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)

# ------------------------------------------
# |                                        |
# |               Release                  |
# |                                        |
# ------------------------------------------

RELEASE_LIST = $(patsubst %, %-release, $(DRIVER_LIST))
.PHONY: release
release: docker-login $(RELEASE_LIST)

# Release/Push all of the docker containers that have been built
.PHONY: $(RELEASE_LIST)
$(RELEASE_LIST): %-release: %-docker
ifdef RELEASE_VERSION
	@echo ""
	@echo "-------------------- Pushing $*-driver Image --------------------"
	@echo ""
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver:latest
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver-arm:latest
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver:$(RELEASE_VERSION)
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver-arm:$(RELEASE_VERSION)
else
	@echo ""
	@echo "-------------------- Pushing $*-driver-dev Image --------------------"
	@echo ""
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver-dev:latest
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver-dev-arm:latest
ifdef COMMIT_HASH
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver-dev:$(COMMIT_HASH)
	@docker push $(DOCKER_IMAGE_BASE)/$*-driver-dev-arm:$(COMMIT_HASH)
endif
endif

.PHONY: clean
clean:
	rm $(BUILD_DIR) -rf

