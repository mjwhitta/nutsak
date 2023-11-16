-include gomk/main.mk
-include local/Makefile

ifeq ($(unameS),Darwin)
    Seeds = $(shell GOOS=darwin go run ./tools/docextract)
else ifeq ($(unameS),Linux)
    Seeds = $(shell GOOS=linux go run ./tools/docextract)
endif

LDFLAGS += -X 'main.SEEDTYPES=$(Seeds)'

ifneq ($(unameS),Windows)
spellcheck:
	@codespell -f -L hilighter -S ".git,*.pem"
endif

superclean: superclean-default
	@rm -f ./testdata/out*
