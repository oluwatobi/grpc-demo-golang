package main

import (
	"context"
	"os"
	"time"

	pb "github.com/oluwatobi/grpc-demo-golang/distributions"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:9000"
)

func main() {

	log.SetOutput(os.Stdout)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	latitude := float64(12.3456789)
	longitude := float64(98.7654321)
	distributionRequestControl := pb.DistributionRequest{
		ExperimentId:   "experiment-id",
		UserId:         "EXCLUDED-USER-ID",
		OrganizationId: "organization-id",
		Latitude:       &latitude,
		Longitude:      &longitude,
	}
	distributionRequestRandom := pb.DistributionRequest{
		ExperimentId:   "experiment-id",
		UserId:         "user-id",
		OrganizationId: "organization-id",
		Latitude:       &latitude,
		Longitude:      &longitude,
	}
	client := pb.NewDistributionsClient(conn)
	control, err1 := client.GetVariantDistribution(ctx, &distributionRequestControl)
	random, err2 := client.GetVariantDistribution(ctx, &distributionRequestRandom)
	if err1 != nil {
		log.Fatalf("Could not retrieve variant assignment: %v", err)
	} else {
		log.WithFields(log.Fields{
			"variantName":            control.GetAssignedVariant().GetName(),
			"variantType":            control.GetAssignedVariant().GetType(),
			"excludedFromExperiment": control.GetExcludedFromExperiment(),
		}).Info("Response recieved.")
	}
	if err2 != nil {
		log.Fatalf("Could not retrieve variant assignment: %v", err)
	} else {
		log.WithFields(log.Fields{
			"variantName":            random.GetAssignedVariant().GetName(),
			"variantType":            random.GetAssignedVariant().GetType(),
			"excludedFromExperiment": random.GetExcludedFromExperiment(),
		}).Info("Response recieved.")
	}
}
