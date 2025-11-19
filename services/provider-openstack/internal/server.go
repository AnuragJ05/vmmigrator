package internal

import (
	"context"
	"fmt"

	pb "vmmigrator/proto/proto"
)

type server struct {
	pb.UnimplementedProviderServer
}

func NewServer() *server { return &server{} }

func (s *server) ListVMs(ctx context.Context, req *pb.ListVMsRequest) (*pb.ListVMsResponse, error) {
	v := &pb.VMInfo{Id: "os-201", Name: "os-vm-201", Cpu: 2, MemoryMb: 4096}
	return &pb.ListVMsResponse{Vms: []*pb.VMInfo{v}}, nil
}

func (s *server) ExportVM(ctx context.Context, req *pb.ExportVMRequest) (*pb.ExportVMResponse, error) {
	artifact := fmt.Sprintf("staging/os-%s.img", req.VmId)
	return &pb.ExportVMResponse{ArtifactId: artifact, Message: "stubbed export from openstack"}, nil
}

func (s *server) ImportVM(ctx context.Context, req *pb.ImportVMRequest) (*pb.ImportVMResponse, error) {
	return &pb.ImportVMResponse{NewVmId: "openstack-imported-001", Message: "stubbed import to openstack"}, nil
}
