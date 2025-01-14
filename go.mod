module github.com/vegadelalyra/go_grpc_client_the_exception_handler

go 1.22.3

require (
	github.com/vegadelalyra/go_microservices_the_exception_handler/rpc v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.64.0
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/vegadelalyra/go_microservices_the_exception_handler/rpc => ../rpc
