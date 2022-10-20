
PATH=$(shell pwd)

lint:
	protolint lint ./grpc/proto

compile_in_docker: lint
	$(shell mkdir -p ./grpc/gens/utpb)
	$(shell mkdir -p ./grpc/gens/docs)
	$(shell mkdir -p ./grpc/gens/tmp)

	/usr/local/bin/docker run --rm -v $(PATH):/proto -w /proto rvolosatovs/protoc:v4.0.0-rc2 \
		--proto_path=. \
		--proto_path=/usr/include \
		--proto_path=/usr/include/github.com/envoyproxy/protoc-gen-validate \
		--go_out=./grpc/gens/tmp --go_opt=paths=source_relative \
		--go-grpc_out=./grpc/gens/tmp --go-grpc_opt=paths=source_relative \
		--doc_out=./grpc/gens/docs --doc_opt=html,index.html \
		--validate_out="lang=go,paths=source_relative:./grpc/gens/tmp" \
		$(shell find ./grpc/proto -name '*.proto')

clean_tmp:
	$(shell rm -rf ${PATH}/grpc/gens)
	$(shell mkdir -p ${PATH}/grpc/gens)

build_in_docker: clean_tmp compile_in_docker
	$(shell cp -r ${PATH}/grpc/gens/tmp/grpc/proto/* ${PATH}/grpc/gens/utpb/)
	$(shell rm -rf ${PATH}/grpc/gens/tmp)
