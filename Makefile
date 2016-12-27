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

# release
release:
	@echo "[begin] release ==================="
	GOOS=darwin GOARCH=amd64 $(GB) build -ldflags "-X main.build=`git rev-parse HEAD`" $(TARGET_CLOUDCLI)
	GOOS=linux GOARCH=amd64 $(GB) build -ldflags "-X main.build=`git rev-parse HEAD`" $(TARGET_CLOUDCLI)
	GOOS=windows GOARCH=amd64 $(GB) build -ldflags "-X main.build=`git rev-parse HEAD`" $(TARGET_CLOUDCLI)
	tar czvf release.tar.gz ./bin/cloudcli-darwin-amd64 ./bin/cloudcli-windows-amd64.exe ./bin/cloudcli-linux-amd64
	@echo "[end]   release ==================="

cli: clean
	@echo "[begin] build cloudcli ==================="
	$(GB) build -ldflags "-X main.build=`git rev-parse HEAD`" $(TARGET_CLOUDCLI)
	@echo "[end  ] build cloudcli ==================="

.PHONY: clean
clean:
	rm -f $(BIN_CLOLUDCLI)

