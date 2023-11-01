-include gomk/main.mk
-include local/Makefile

build: sak;

sak: dir fmt
	@GOARCH=$(GOARCH) GOOS=$(GOOS) GOPATH=$(GOPATH) $(CC) build --ldflags "$(LDFLAGS) -X 'main.SEEDTYPES=$(shell go run ./tools/docextract)'" -o "$(OUT)" $(TRIM) ./cmd/sak

ifneq ($(unameS),Windows)
spellcheck:
	@codespell -f -L hilighter -S ".git,*.pem"
endif

superclean: superclean-default
	@rm -f ./testdata/out*
