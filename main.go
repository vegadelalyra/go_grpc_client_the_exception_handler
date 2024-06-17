package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vegadelalyra/go_microservices_the_exception_handler/rpc/rpc_auth"
	"google.golang.org/grpc"
)

var loginResponse *rpc_auth.LoginResponse

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		fmt.Println("error in dial", err)
	}
	defer conn.Close()

	// loginClient := rpc_auth.NewLoginServiceClient(conn)
	validateClient := rpc_auth.NewValidateTokenServiceClient(conn)

	// fmt.Println("****** Simple *******")
	// simpleRpc(loginClient)

	// fmt.Println("****** Server Stream *****")
	// serverStream(loginClient)

	// fmt.Println("****** Client Stream ******")
	// clientStream(loginClient)

	fmt.Println("******* Validate ******")
	bidirectionalReceive(validateClient)
}

func simpleRpc(loginClient rpc_auth.LoginServiceClient) {
	var err error
	loginResponse, err = loginClient.LoginSimpleRPC(context.Background(), &rpc_auth.LoginRequest{
		Username:   "steve@yopmail.com",
		Password:   "steve123",
		RememberMe: true,
	})

	if err != nil {
		fmt.Println("Simple RPC login error", err)
	} else {
		fmt.Printf("Simple RPC login\n %+v \n", loginResponse)
	}
}

func serverStream(loginClient rpc_auth.LoginServiceClient) {
	var loginDetails rpc_auth.LoginRequestList = rpc_auth.LoginRequestList{
		Data: make([]*rpc_auth.LoginRequest, 0),
	}
	loginDetails.Data = append(loginDetails.Data, &rpc_auth.LoginRequest{
		Username:   "steve@yopmail.com",
		Password:   "steve123",
		RememberMe: true,
	}, &rpc_auth.LoginRequest{
		Username:   "mark@yopmail.com",
		Password:   "mark@123",
		RememberMe: true,
	})

	stream, err := loginClient.LoginServerStreamRPC(context.Background(), &loginDetails)
	if err != nil {
		fmt.Println("Server stream RPC login error", err)
	}
	var i int = 0
	for {
		loginRes, errRec := stream.Recv()
		if errRec == io.EOF {
			break
		} else if errRec != nil {
			log.Fatalf("receive LoginServerStreamRPC err %v", errRec)
		} else {
			i++
			fmt.Printf("(%d) Server stream RPC login response \n %+v \n", i, loginRes)
		}
	}
}

func clientStream(loginClient rpc_auth.LoginServiceClient) {
	var loginDetails []rpc_auth.LoginRequest
	loginDetails = append(loginDetails, rpc_auth.LoginRequest{
		Username:   "steve@yopmail.com",
		Password:   "steve123",
		RememberMe: true,
	}, rpc_auth.LoginRequest{
		Username:   "mark@yopmail.com",
		Password:   "mark@123",
		RememberMe: true,
	})
	stream, err := loginClient.LoginClientStreamRPC(context.Background())
	if err != nil {
		fmt.Printf("LoginClientStreamRPC connect failed error %+v", err)
	}
	for _, loginCredential := range loginDetails {
		if errSend := stream.Send(&loginCredential); errSend != nil {
			fmt.Printf("credential send error %+v", errSend)
		}
	}

	response, errRecv := stream.CloseAndRecv()
	if errRecv != nil {
		fmt.Printf("CloseAndRecv error %+v", errRecv)
	}
	fmt.Printf("login tokens: \n%+v\n", response)
}

func bidirectionalReceive(validateClient rpc_auth.ValidateTokenServiceClient) {
	stream, err := validateClient.Validate(context.Background())
	if err != nil {
		fmt.Println("bidirectionalReceive RPC error", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a claims : %v", err)
			}
			log.Printf("Got valid - %v - claims : companyId: %s , Username: %s", in.IsValid, in.CompanyId, in.Username)
		}
	}()
	fmt.Println("sending validate request 1st time")
	if err := stream.Send(&rpc_auth.ValidateTokenRequest{
		Token: loginResponse.Token,
	}); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
	time.Sleep(2 * time.Second)
	fmt.Println("sending validate request 2nd time")
	if err := stream.Send(&rpc_auth.ValidateTokenRequest{
		Token: loginResponse.Token + "1",
	}); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
	stream.CloseSend()
	<-waitc
}
