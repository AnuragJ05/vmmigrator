package main

import (
	"log"
	"net"

	pb "vmmigrator/proto"
	"vmmigrator/services/orchestrator/internal"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50060")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrchestratorServer(s, internal.NewServer())

	log.Println("orchestrator gRPC listening :50060")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
