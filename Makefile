include .env



.PHONY: protoc
protoc/auth:
	@echo 'Updating proto...'
	protoc --proto_path protos --proto_path vendor.protogen \
	--go_out=services/auth/gen/auth --go_opt=paths=source_relative --go-grpc_out=services/auth/gen/auth --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=${PROTOC_GEN_PATH} \
  	--grpc-gateway_out=services/auth/gen/auth --grpc-gateway_opt=paths=source_relative \
  	--plugin=protoc-gen-grpc-gateway=${PROTOC_GATEWAY_PATH} \
  	--plugin=protoc-gen-openapiv2=${PROTOC_OPENAPIV2_PATH} \
  	--openapiv2_out=services/auth/swagger --openapiv2_opt=allow_merge=true,merge_file_name=swagger.json \
  	--go_out=services/company/gen/auth --go_opt=paths=source_relative \
    --go-grpc_out=services/company/gen/auth --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=services/company/gen/auth --grpc-gateway_opt=paths=source_relative \
  	protos/auth.proto

protoc/comp:
	@echo 'Updating proto...'
	protoc --proto_path protos --proto_path vendor.protogen \
	--go_out=services/company/gen/company --go_opt=paths=source_relative --go-grpc_out=services/company/gen/company --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=${PROTOC_GEN_PATH} \
  	--grpc-gateway_out=services/company/gen/company --grpc-gateway_opt=paths=source_relative \
  	--plugin=protoc-gen-grpc-gateway=${PROTOC_GATEWAY_PATH} \
  	--plugin=protoc-gen-openapiv2=${PROTOC_OPENAPIV2_PATH} \
  	--openapiv2_out=services/company/swagger --openapiv2_opt=allow_merge=true,merge_file_name=swagger.json \
  	protos/company.proto


.PHONY: vendor-proto
vendor-proto:
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
#		@if [ ! -d vendor.protogen/openapiv2 ]; then \
#    		git clone https://github.com/grpc-ecosystem/grpc-gateway.git vendor.protogen/grpc-gateway && \
#    		mkdir -p vendor.protogen/openapiv2/ && \
#    		mv vendor.protogen/grpc-gateway/protoc-gen-openapiv2/* vendor.protogen/openapiv2/ && \
#    		rm -rf vendor.protogen/grpc-gateway ; \
#    	fi





.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=sql -dir=D:/ProgramData/workspacego/diploma/services/auth/migrations ${name}


.PHONY: db/migrations
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path=D:/ProgramData/workspacego/diploma/services/auth/migrations -database=${POSTGRES_URI} up

db/migrations/down:
	@echo 'Running up migrations ...'
	migrate -path=D:/ProgramData/workspacego/diploma/services/auth/migrations -database=${POSTGRES_URI} down 1


.PHONY: all auth company run

run:
	start /b cmd /c "cd services/auth && go run ./cmd/auth 2>&1" & \
	start /b cmd /c "cd services/company && go run ./cmd/company 2>&1"









