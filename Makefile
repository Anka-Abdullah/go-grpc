PROTO_DIR := proto

.PHONY: all generate clean fix regenerate run dev-up

all: generate fix

generate:
	@echo "🔄 Generating gRPC and Go files from proto definitions..."
	@find $(PROTO_DIR) -name "*.proto" | while read file; do \
		protoc --proto_path=$(PROTO_DIR) \
			--go_out=$(PROTO_DIR) \
			--go-grpc_out=$(PROTO_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			$$file; \
	done
	@echo "✅ Generation complete."

fix:
	@echo "🔧 Moving generated files to correct folders..."
	@find $(PROTO_DIR) -type d | while read dir; do \
		if [ -d "$$dir" ] && [ "$$(basename $$dir)" != "$$(basename $(PROTO_DIR))" ]; then \
			subdir="$$dir/$$(basename $$dir)"; \
			if [ -d "$$subdir" ]; then \
				mv "$$subdir"/*.pb.go "$$dir"/ 2>/dev/null || true; \
				rm -rf "$$subdir"; \
			fi \
		fi \
	done
	@echo "✅ Fix complete."

clean:
	@echo "🧹 Cleaning generated files..."
	@find $(PROTO_DIR) -name "*.pb.go" -delete
	@find $(PROTO_DIR) -name "*_grpc.pb.go" -delete
	@find $(PROTO_DIR) -type d -name "*/*.proto" -exec dirname {} \; | xargs -r rm -rf
	@echo "🧽 Cleanup done."

regenerate: clean all

run:
	@echo "🚀 Starting development container (build if needed)..."
	@if ! docker image inspect grpc_server_img:latest > /dev/null 2>&1; then \
		echo "🔧 Image not found. Building first..."; \
		docker compose -f Docker/docker-compose.yml build; \
	else \
		echo "✅ Image 'grpc_server_img' already exists. Skipping build..."; \
	fi
	docker compose -f Docker/docker-compose.yml up


down:
	@echo "🧹 Stopping container and removing volumes..."
	docker compose -f Docker/docker-compose.yml down -v