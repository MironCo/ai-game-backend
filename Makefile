# Define binary name and path to main.go
BINARY_NAME=rd-backend
MAIN_PATH=cmd/rd-backend/main.go

# Build the project
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the project
run:
	go run $(MAIN_PATH)

# Clean up generated files
clean:
	rm -f $(BINARY_NAME)
