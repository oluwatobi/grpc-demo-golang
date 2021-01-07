package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"net"
	"os"
	"strings"

	pb "github.com/oluwatobi/grpc-demo-golang/distributions"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	NUMBER_OF_EXPERIMENTS = 2

	port = ":9000"
)

// server is used to implement distributions.DistributionsServer.
type server struct {
	pb.UnimplementedDistributionsServer
}

func calculateDistributionBaseOnAttributes(
	userId string,
	organizationId string,
	latitude float64,
	longitude float64) (pb.DistributionResponse_Variant, bool, error) {
	compositeKey := fmt.Sprintf("%s|%s|%f|%f", userId, organizationId, latitude, longitude)
	h := fnv.New32a()
	h.Write([]byte(compositeKey))
	hashValue := h.Sum32()
	assignmentIndex := hashValue % NUMBER_OF_EXPERIMENTS
	var variant pb.DistributionResponse_Variant
	excludedFromExperiment := false

	if strings.Compare(userId, "EXCLUDED-USER-ID") == 0 {
		variant = pb.DistributionResponse_Variant{
			Name: "CONTROL",
			Type: 0,
		}
		excludedFromExperiment = true
	} else if assignmentIndex == 1 {
		variant = pb.DistributionResponse_Variant{
			Name: "CONTROL",
			Type: 0,
		}
	} else {
		variant = pb.DistributionResponse_Variant{
			Name: "TREATMENT",
			Type: 1,
		}
	}

	return variant, excludedFromExperiment, nil
}

// GetDistribution implements distributions.DistributionsServer
func (s *server) GetVariantDistribution(
	ctx context.Context,
	request *pb.DistributionRequest) (*pb.DistributionResponse, error) {
	userId := request.GetUserId()
	organizationId := request.GetOrganizationId()
	latitude := request.GetLatitude()
	longitude := request.GetLongitude()
	log.WithFields(log.Fields{
		"experimentId":   request.GetExperimentId(),
		"userId":         request.GetUserId(),
		"organizationId": request.GetOrganizationId(),
		"latitude":       request.GetLatitude(),
		"longitude":      request.GetLongitude(),
	}).Info("Request received.")

	assignedVariant, excludedFromExperiment, err := calculateDistributionBaseOnAttributes(
		userId,
		organizationId,
		latitude,
		longitude,
	)

	if err != nil {
		return nil, err
	}

	return &pb.DistributionResponse{
		AssignedVariant:        &assignedVariant,
		ExcludedFromExperiment: excludedFromExperiment,
	}, nil
}

func main() {

	log.SetOutput(os.Stdout)
	log.Info("Server started.")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDistributionsServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
