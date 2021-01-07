package main

import (
	"context"
	"fmt"

	pb "distributions"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDistributionsClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	distributionRequest := DistributionRequest{
		ExperimentId:   "experiment-id",
		UserId:         "EXCLUDED_USER_ID",
		OrganizationId: "organization-id",
		Latitude:       float64(12.3456789),
		Longitude:      float64(98.7654321),
	}
	r, err := c.GetVariantDistribution(ctx, &distributionRequest)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.WithFields(log.Fields{
		"variantName":            r.GetAssignedVariant().GetName(),
		"variantType":            r.GetAssignedVariant().GetType(),
		"excludedFromExperiment": r.GetExcludedFromExperiment(),
	}).Info("Response recieved.")
}
