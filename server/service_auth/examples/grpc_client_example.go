package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/youknow2509/cio_verify_face/server/service_auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewAuthServiceClient(conn)

	// Example 1: Validate Token
	fmt.Println("=== Testing Token Validation ===")
	validateReq := &pb.ValidateTokenRequest{
		Token: "your-jwt-token-here",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	validateResp, err := client.ValidateToken(ctx, validateReq)
	if err != nil {
		log.Printf("ValidateToken failed: %v", err)
	} else {
		fmt.Printf("Token validation result: %+v\n", validateResp)
	}

	// Example 2: Get User Info
	fmt.Println("\n=== Testing Get User Info ===")
	userInfoReq := &pb.GetUserInfoRequest{
		UserId: "user-uuid-here",
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	userInfoResp, err := client.GetUserInfo(ctx2, userInfoReq)
	if err != nil {
		log.Printf("GetUserInfo failed: %v", err)
	} else {
		fmt.Printf("User info result: %+v\n", userInfoResp)
	}

	// Example 3: Check User Permission
	fmt.Println("\n=== Testing Check User Permission ===")
	permissionReq := &pb.CheckUserPermissionRequest{
		UserId:    "user-uuid-here",
		CompanyId: "company-uuid-here",
	}

	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel3()

	permissionResp, err := client.CheckUserPermission(ctx3, permissionReq)
	if err != nil {
		log.Printf("CheckUserPermission failed: %v", err)
	} else {
		fmt.Printf("Permission check result: %+v\n", permissionResp)
	}

	// Example 4: Check Device In Company
	fmt.Println("\n=== Testing Check Device In Company ===")
	deviceReq := &pb.CheckDeviceInCompanyRequest{
		DeviceId:  "device-uuid-here",
		CompanyId: "company-uuid-here",
	}

	ctx4, cancel4 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel4()

	deviceResp, err := client.CheckDeviceInCompany(ctx4, deviceReq)
	if err != nil {
		log.Printf("CheckDeviceInCompany failed: %v", err)
	} else {
		fmt.Printf("Device check result: %+v\n", deviceResp)
	}

	// Example 5: Batch Validate Tokens
	fmt.Println("\n=== Testing Batch Validate Tokens ===")
	batchReq := &pb.BatchValidateTokensRequest{
		Tokens: []*pb.TokenValidation{
			{
				Token:     "token-1",
				RequestId: "req-1",
			},
			{
				Token:     "token-2",
				RequestId: "req-2",
			},
		},
	}

	ctx5, cancel5 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel5()

	batchResp, err := client.BatchValidateTokens(ctx5, batchReq)
	if err != nil {
		log.Printf("BatchValidateTokens failed: %v", err)
	} else {
		fmt.Printf("Batch validation result: %+v\n", batchResp)
	}

	fmt.Println("\n=== gRPC Client Tests Completed ===")
}