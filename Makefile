protos:
	protoc contract/contract.proto --go_out=plugins=grpc:.

SERVER_CN=localhost

generate-ssl:
	mkdir ssl
	openssl genrsa -passout pass:1111 -des3 -out ssl/ca.key 4096
	openssl req -passin pass:1111 -new -x509 -days 3650 -key ssl/ca.key -out ssl/ca.crt -subj "/CN=${SERVER_CN}"
	openssl genrsa -passout pass:1111 -des3 -out ssl/server.key 4096
	openssl req -passin pass:1111 -new -key ssl/server.key -out ssl/server.csr -subj "/CN=${SERVER_CN}"
	openssl x509 -req -passin pass:1111 -days 365 -in ssl/server.csr -CA ssl/ca.crt -CAkey ssl/ca.key -set_serial 01 -out ssl/server.crt
	openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in ssl/server.key -out ssl/server.pem