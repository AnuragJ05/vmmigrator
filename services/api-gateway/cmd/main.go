package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "vmmigrator/proto/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to orchestrator
	conn, err := grpc.Dial("orchestrator:50060", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to orchestrator: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorClient(conn)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/migrations", func(c *gin.Context) {
		var body struct {
			SourceProvider string   `json:"source_provider"`
			DestProvider   string   `json:"dest_provider"`
			VMIDs          []string `json:"vm_ids"`
			StagingTarget  string   `json:"staging_target"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := client.StartMigration(ctx, &pb.StartMigrationRequest{
			SourceProvider: body.SourceProvider,
			DestProvider:   body.DestProvider,
			VmIds:          body.VMIDs,
			StagingTarget:  body.StagingTarget,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"job_id": res.JobId, "message": res.Message})
	})

	r.GET("/migrations/:id", func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := client.GetMigration(ctx, &pb.GetMigrationRequest{JobId: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	log.Println("API Gateway listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run api-gateway: %v", err)
	}
}
