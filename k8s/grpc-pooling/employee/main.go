package main

import (
	"context"
	"math"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/mycodesmells/golang-examples/k8s/grpc-pooling/proto/message"
)

func main() {
	addr := os.Getenv("ADDR")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to initialize TCP listen: %v", err)
	}
	defer lis.Close()

	hn, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to read hostname: %v", err)
	}
	worker := employeeServer{WorkerID: hn}

	server := grpc.NewServer()
	pb.RegisterWorkerServer(server, worker)

	log.Infof("Worker initialized, workerID = %s", worker.WorkerID)
	log.Printf("gRPC Listening on %s", lis.Addr().String())
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

type employeeServer struct {
	WorkerID string
}

func (eS employeeServer) Work(ctx context.Context, req *pb.JobRequest) (*pb.JobResponse, error) {
	base := req.GetBase()
	exponent := req.GetExponent()

	result := math.Pow(float64(base), float64(exponent))

	select {
	case <-time.After(time.Second * 10):
		return &pb.JobResponse{
			Id:       req.GetId(),
			WorkerId: eS.WorkerID,
			Result:   float32(result),
		}, nil
	}
}
