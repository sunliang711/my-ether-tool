.PHONY:	build
build:
	@echo "Building..."
	@go build -o met main.go

.PHONY:	install
install:
	@echo "Installing..."
	@go install
