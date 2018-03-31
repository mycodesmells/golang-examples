package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/processout/grpc-go-pool"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/mycodesmells/golang-examples/k8s/grpc-pooling/proto/message"
)

func main() {
	addr := os.Getenv("ADDR")
	employeeAddr := os.Getenv("EMPLOYEE_ADDR")

	var factory grpcpool.Factory
	factory = func() (*grpc.ClientConn, error) {
		conn, err := grpc.Dial(employeeAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to start gRPC connection: %v", err)
		}
		log.Infof("Connected to employee at %s", employeeAddr)
		return conn, err
	}

	pool, err := grpcpool.New(factory, 5, 5, time.Second)
	if err != nil {
		log.Fatalf("Failed to create gRPC pool: %v", err)
	}

	http.HandleFunc("/power", func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		query := req.URL.Query()

		baseStr := query.Get("base")
		exponentStr := query.Get("exponent")

		l := log.WithFields(log.Fields{
			"base":     baseStr,
			"exponent": exponentStr,
		})

		base, err := strconv.ParseFloat(baseStr, 32)
		if err != nil {
			msg := "invalid base value"
			l.Warnln(errors.Wrap(err, msg))

			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(msg))
			return
		}

		exponent, err := strconv.ParseFloat(exponentStr, 32)
		if err != nil {
			msg := "invalid exponent value"
			l.Warnln(errors.Wrap(err, msg))

			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(msg))
			return
		}

		ctx := req.Context()
		conn, err := pool.Get(ctx)
		defer conn.Close()
		if err != nil {
			msg := "failed to connect to worker"
			l.Errorln(errors.Wrap(err, msg))

			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(msg))
			return
		}

		client := pb.NewWorkerClient(conn.ClientConn)

		workResp, err := client.Work(ctx, &pb.JobRequest{
			Id:       uuid.NewV4().String(),
			Base:     float32(base),
			Exponent: float32(exponent),
		})
		if err != nil {
			msg := "failed to compute result"
			l.Errorln(errors.Wrap(err, msg))

			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(msg))
			return
		}
		result := workResp.GetResult()

		log.Infof("Job %s, Worker: %s, Result: %f", workResp.GetId(), workResp.GetWorkerId(), result)

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("(%f)^(%f) = %f\n", base, exponent, result)))
		end := time.Now()
		diff := end.Sub(start)
		log.Infof("Finished in time: %v", diff)
	})
	http.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	log.Infof("HTTP server listening on %s", addr)
	http.ListenAndServe(addr, nil)
}
