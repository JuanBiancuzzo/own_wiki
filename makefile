build: build_proto

build_proto:
	make -C ./core/api build_proto

tests:
	go test ./core/... 

build_tests: build tests