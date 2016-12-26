# makefile for cloudcli
#  build with git revision number
#  example:  gb build -ldflags "-X main.build=`git rev-parse HEAD`" cmd/cloudcli

# varible
GB = gb
CURRENT_REVISION = `git rev-parse HEAD`
LD_FLAGS := "-X main.build=$(CURRENT_REVISION)"

TARGET_CLOUDCLI = cmd/cloudcli
BIN_CLOLUDCLI = ./bin/cloudcli
RELEASE_DIR = ./release

# target
all: cli

release: cli
	@echo "[begin] release ==================="
	-mkdir $(RELEASE_DIR)
	cp $(BIN_CLOLUDCLI) $(RELEASE_DIR)
	@echo "[end]   release ==================="

cli: clean
	@echo "[begin] build cloudcli ==================="
	$(GB) build -ldflags "-X main.build=`git rev-parse HEAD`" $(TARGET_CLOUDCLI)
	@echo "[end  ] build cloudcli ==================="

.PHONY: clean
clean:
	rm -f $(BIN_CLOLUDCLI)

clean_release:
	rm -rf $(RELEASE_DIR)
