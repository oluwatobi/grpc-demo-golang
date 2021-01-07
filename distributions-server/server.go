package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"net"
	"strings"

	pb "distributions"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	NUMBER_OF_EXPERIMENTS = 2

	port = ":50051"
)

// server is used to implement distributions.DistributionsServer.
type server struct {
	pb.UnimplementedDistributionsServer
}

func calculateDistributionBaseOnAttributes(
	userId string,
	organizationId string,
	latitude float64,
	longitude float64) (pb.DistributionResponse_Variant, excludedFromExperiment, error) {
	compositeKey := fmt.Sprintf("%s|%s|%f|%f", userId, brandId, latitude, longitude)
	h := fnv.New32a()
	h.Write([]byte(compositeKey))
	hashValue := h.Sum32()
	assignmentIndex := hashValue % NUMBER_OF_EXPERIMENTS
	var variant pb.Variant
	var excludedFromExperiment bool
	if assignmentIndex == 1 {
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

	if strings.Compare(userId, "EXCLUDED_USERID") {
		excludedFromExperiment = true
	} else {
		excludedFromExperiment = false
	}

	return variant, excludedFromExperiment, nil
}

// GetDistribution implements distributions.DistributionsServer
func (s *server) GetDistribution(
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
	}).Info("Request Recieved.")

	assignedVariant, excludedFromExperiment, error := calculateDistributionBaseOnAttributes(
		userId,
		organizationId,
		latitude,
		longitude,
	)

	if err != nil {
		return nil, err
	}

	return &pb.DistributionResponse{
		AssignedVariant:        assignedVariant,
		ExcludedFromExperiment: excludedFromExperiment,
	}, nil
}

func main() {
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
