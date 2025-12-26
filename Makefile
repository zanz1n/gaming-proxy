ifneq ($(wildcard .env),)
include .env
endif

SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

DEBUG ?= 0
PROJECT := proxy
BIN := ./bin
GO := go
TARGETS ?= linux-x86 linux-x86_64 linux-aarch64 linux-riscv64

ARCH := $(shell uname -m | tr '[:upper:]' '[:lower:]')
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

GOOS := $(OS)
GOARCH :=

ifeq ($(ARCH), x86)
	GOARCH = 386
else ifeq ($(ARCH), x86_64)
	GOARCH = amd64
else ifeq ($(ARCH), aarch64)
	GOARCH = arm64
else
	GOARCH = $(ARCH)
endif


LDFLAGS :=

ifneq ($(DEBUG), 1)
	LDFLAGS += -s -w
endif

ifeq ($(GOOS), windows)
	ifeq ($(DEBUG), 1)
	    SUFIX += -debug.exe
	else
	    SUFIX +=.exe
	endif
endif

default: build

all: $(addprefix build-, $(TARGETS))

build-%: $(BIN)
	$(MAKE) DEBUG=$(DEBUG) \
	OS=$(word 1, $(subst -, ,$*)) \
	ARCH=$(word 2, $(subst -, ,$*)) build

BIN_NAME := $(PROJECT)-$(OS)-$(ARCH)$(SUFIX)
OUT := $(BIN)/$(BIN_NAME)

run: build
	$(BIN)/$(BIN_NAME)

build: $(BIN)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build \
	-ldflags "$(LDFLAGS)" -o $(OUT) .

$(BIN):
	mkdir -p $(BIN)

log:
	$(info $(ARCH))
