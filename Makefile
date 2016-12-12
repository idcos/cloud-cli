# makefile for cloudcli
#  build with git revision number
#  example:  gb build -ldflags "-X main.build=`git rev-parse HEAD`" cmd/cloudcli

# varible
GB = gb
CURRENT_REVISION = `git rev-parse HEAD`
LD_FLAGS := "-X main.build=$(CURRENT_REVISION)"

TARGET_CLOUDCLI = cmd/cloudcli
BIN_CLOLUDCLI = ./bin/cloudcli

# target
all:
	@echo "[begin] build all ==================="
	rm -f $(BIN_CLOLUDCLI)
	$(GB) build -ldflags $(LD_FLAGS) $(TARGET_CLOUDCLI)
	@echo "[end  ] build all ==================="

cli:
	@echo "[begin] build cloudcli ==================="
	rm -f $(BIN_CLOLUDCLI)
	$(GB) build -ldflags "-X main.build=`git rev-parse HEAD`" $(TARGET_CLOUDCLI)
	@echo "[end  ] build cloudcli ==================="

.PHONY: clean
clean:
	rm -f $(BIN_CLOLUDCLI)

