GOPATH = $(shell go env GOPATH)
OS = $(shell go env GOOS)
all: mod check build

mod:
	go mod tidy

check: mod
	go vet ./...
	golint ./...
	gofumpt -l -w .

build:
	mkdir -p build
	go build -o build/ gping
	@ if [ "$(OS)" = "linux" ]; then sudo setcap cap_net_raw=+ep build/gping; fi

clean:
	go clean
	rm -rf build

install:
	go install github.com/equationzhao/gping
	@ if [ "$(OS)" = "linux" ]; then sudo setcap cap_net_raw=+ep $(GOPATH)/bin/gping; fi

uninstall:
	$(info Uninstalling gping)
	@if [ "$(OS)" = "linux" ]; \
 		then sudo setcap -r $(GOPATH)/bin/gping; \
 	elif [ "$(OS)" = "darwin" ]; \
 	  	then sudo chflags nosuid $(GOPATH)/bin/gping; \
 	elif [ "$(OS)" = "windows" ]; then rm -f "${GOPATH}\\bin\\gping.exe";\
 	fi;
