# init project path
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output

# init SWEEP Region Arg
SWEEP := bj,gz,su
ifdef BAIDUCLOUD_REGION
	SWEEP = $(BAIDUCLOUD_REGION)
endif

# init GO
# init GO in build.sh
# init command params
GO      := go
GOBUILD := $(GO) build
GOTEST  := $(GO) test
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")
# make, make all
all: prepare compile package
# make prepare, download dependencies
prepare:
	# nothing to do
# make compile, go build
compile: build
build:
	$(GOBUILD) -o $(HOMEDIR)/terraform-provider-baiducloud
# make test, test your code
test: test-case
testacc:
	TF_ACC=1 $(GOTEST) -v -cover $(GOPKGS) -timeout 120m

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	$(GOTEST) $(GOPKGS) -v -sweep=$(SWEEP) $(SWEEPARGS)

# make package
package: package-bin
package-bin:
	mkdir -p $(OUTDIR)
	cp terraform-provider-baiducloud  $(OUTDIR)/
# make clean
clean:
	rm -rf $(OUTDIR)
	rm -f $(HOMEDIR)/terraform-provider-baiducloud
# avoid filename conflict and speed up build 
.PHONY: all prepare compile test testacc sweep package clean build
