package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	filename := fmt.Sprintf("os-%s.img", req.VmId)
	path := filepath.Join("/staging", filename)

	// Create a dummy file
	content := fmt.Sprintf("OpenStack VM Export Data for %s\nTimestamp: %s", req.VmId, time.Now().Format(time.RFC3339))
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write export file: %v", err)
	}

	return &pb.ExportVMResponse{ArtifactId: filename, Message: "export successful"}, nil
}

func (s *server) ImportVM(ctx context.Context, req *pb.ImportVMRequest) (*pb.ImportVMResponse, error) {
	// Simulate import delay
	time.Sleep(2 * time.Second)

	filename := req.ArtifactId
	// In a real scenario, we might validate the path or use a configured staging dir
	path := filepath.Join("/staging", filename)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("artifact not found at %s", path)
	} else if err != nil {
		return nil, fmt.Errorf("failed to check artifact: %v", err)
	}

	// Simulate reading/processing
	fmt.Printf("Importing VM from %s (size: %d bytes)\n", path, info.Size())

	newID := fmt.Sprintf("openstack-imported-%d", time.Now().Unix())
	return &pb.ImportVMResponse{NewVmId: newID, Message: "import successful"}, nil
}
