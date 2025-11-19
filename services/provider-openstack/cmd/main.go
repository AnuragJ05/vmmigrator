package main

import (
	"log"
	"net"

	pb "vmmigrator/proto"
	"vmmigrator/services/provider-openstack/internal"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProviderServer(s, internal.NewServer())

	log.Println("provider-openstack gRPC listening :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
