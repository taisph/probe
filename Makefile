export BUILD_BRANCH_NAME := $(or $(BUILD_BRANCH_NAME),$(shell git symbolic-ref --short HEAD 2>/dev/null))
export BUILD_ID
export BUILD_REPO_NAME
export BUILD_REVISION_ID
export BUILD_SHORT_SHA := $(or $(BUILD_SHORT_SHA),$(shell git describe --match="" --always --abbrev --dirty 2>/dev/null))
export BUILD_TAG_NAME := $(or $(BUILD_TAG_NAME),$(shell git describe --tags --exact 2>/dev/null),$(shell git describe --tags 2>/dev/null | cut -d '-' -f-2))
override BUILD_VARIANT := $(or $(BUILD_VARIANT),debug)
export BUILD_VARIANT
override BUILD_VERSION := $(or $(BUILD_VERSION),$(or $(BUILD_TAG_NAME),0.0.0)+$(subst $(subst ,, ),.,$(strip $(BUILD_SHORT_SHA) $(shell echo "$(BUILD_BRANCH_NAME)" | sed -e 's/[^a-zA-Z0-9-]/-/g;s/^-\+//g;s/-\+$$//g') $(BUILD_VARIANT))))
export BUILD_VERSION

export SKAFFOLD_UPDATE_CHECK := false

CMDS ?= $(patsubst cmd/%,%,$(wildcard cmd/*))

.PHONY: all

all:

#
# Development targets.
#

.PHONY: up down dev prune check

up:
	scripts/dev.sh up

down:
	scripts/dev.sh down

dev:
	scripts/dev.sh dev

prune:
	scripts/dev.sh clean

check:
	scripts/dev.sh check

#
# Build targets.
#

.PHONY: build $(CMDS) clean test

ifeq ($(strip $(CMDS)),)
build:
	scripts/build.sh
else
build: $(CMDS)
endif

$(CMDS):
	scripts/build.sh $(patsubst %,./cmd/%,$@)

clean:
	rm -vrf out

test:
	scripts/test.sh

#
# Ops targets.
#

.PHONY: env image images render

env:
	mkdir -p out
	@set | grep ^BUILD_ | tee out/.env.build

image images:
	scripts/ops.sh build

render:
	scripts/ops.sh render

release:
	scripts/ops.sh release
