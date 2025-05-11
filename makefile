APP_NAME = osu-ha
BUILD_DIR = build

.PHONY: all build run clean docker docker-run

all: build

build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run:
	@echo "🚀 Running $(APP_NAME)..."
	go run ./cmd/server

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(BUILD_DIR)

docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-run:
	@echo "🐳 Running Docker container..."
	docker run -it --rm \
		-p 8081:8081 \
		--env-file .env \
		--name $(APP_NAME) $(APP_NAME):latest
