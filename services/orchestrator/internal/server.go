package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "vmmigrator/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedOrchestratorServer
	mu   sync.Mutex
	jobs map[string]*pb.GetMigrationResponse

	providerEndpoints map[string]string
}

func NewServer() *server {
	return &server{
		jobs: make(map[string]*pb.GetMigrationResponse),
		providerEndpoints: map[string]string{
			"vmware":    "localhost:50051",
			"openstack": "localhost:50052",
		},
	}
}

func (s *server) StartMigration(ctx context.Context, req *pb.StartMigrationRequest) (*pb.StartMigrationResponse, error) {
	jobID := time.Now().Format("20060102150405")
	now := time.Now()
	s.mu.Lock()
	s.jobs[jobID] = &pb.GetMigrationResponse{
		JobId:      jobID,
		State:      pb.JobState_QUEUED,
		StartedAt:  nil,
		VmProgress: []*pb.VMProgress{},
	}
	s.mu.Unlock()

	go s.runJob(jobID, req, now)
	return &pb.StartMigrationResponse{JobId: jobID, Message: "started"}, nil
}

func (s *server) runJob(jobID string, req *pb.StartMigrationRequest, startTime time.Time) {
	s.mu.Lock()
	job := s.jobs[jobID]
	job.State = pb.JobState_RUNNING
	job.StartedAt = timestamppb.New(startTime)
	s.mu.Unlock()

	srcName := parseProvider(req.SourceProvider)
	dstName := parseProvider(req.DestProvider)

	srcCli, err := s.newProviderClient(srcName)
	if err != nil {
		s.failJob(jobID, fmt.Sprintf("src provider error: %v", err))
		return
	}
	dstCli, err := s.newProviderClient(dstName)
	if err != nil {
		s.failJob(jobID, fmt.Sprintf("dst provider error: %v", err))
		return
	}

	for _, vmid := range req.VmIds {
		exp, err := srcCli.ExportVM(context.Background(), &pb.ExportVMRequest{
			VmId:          vmid,
			StagingTarget: req.StagingTarget,
		})
		if err != nil {
			s.failJob(jobID, fmt.Sprintf("export failed for %s: %v", vmid, err))
			return
		}
		_, err = dstCli.ImportVM(context.Background(), &pb.ImportVMRequest{
			ArtifactId: exp.ArtifactId,
		})
		if err != nil {
			s.failJob(jobID, fmt.Sprintf("import failed for %s: %v", vmid, err))
			return
		}
		s.mu.Lock()
		s.jobs[jobID].VmProgress = append(s.jobs[jobID].VmProgress, &pb.VMProgress{
			VmId:    vmid,
			State:   pb.JobState_COMPLETED,
			Message: "migrated",
		})
		s.mu.Unlock()
	}

	s.mu.Lock()
	end := time.Now()
	job = s.jobs[jobID]
	job.State = pb.JobState_COMPLETED
	job.FinishedAt = timestamppb.New(end)
	s.mu.Unlock()
}

func (s *server) failJob(jobID, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[jobID]
	if !ok {
		return
	}
	job.State = pb.JobState_FAILED
	job.VmProgress = append(job.VmProgress, &pb.VMProgress{
		VmId:    "job",
		State:   pb.JobState_FAILED,
		Message: msg,
	})
}

func (s *server) newProviderClient(name string) (pb.ProviderClient, error) {
	endpoint, ok := s.providerEndpoints[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// NOTE: in a production system you would manage connection lifecycle properly.
	return pb.NewProviderClient(conn), nil
}

func parseProvider(in string) string {
	for i := 0; i < len(in); i++ {
		if in[i] == ':' {
			return in[:i]
		}
	}
	return in
}

func (s *server) GetMigration(ctx context.Context, req *pb.GetMigrationRequest) (*pb.GetMigrationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[req.JobId]
	if !ok {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}
