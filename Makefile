-include gomk/main.mk
-include local/Makefile

ifeq ($(unameS),windows)
    Seeds = $(shell set-item env:GOOS "$(unameS)"; go run ./tools/docextract)
else
    Seeds = $(shell GOOS=$(unameS) go run ./tools/docextract)
endif

LDFLAGS += -X 'main.SEEDTYPES=$(Seeds)'

superclean: superclean-default
ifeq ($(unameS),windows)
	@remove-item -force ./testdata/out*
else
	@rm -f ./testdata/out*
endif
