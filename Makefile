.PHONY:	build
build:
	@echo "Building..."
	@go build -ldflags "-s -w" -o met main.go

.PHONY:	install
install:
	@echo "Installing..."
	@go install -ldflags "-s -w"
