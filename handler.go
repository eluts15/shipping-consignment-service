// consignment-service/handler.go
package main

import (
	"context"
	"log"

	pb "github.com/eluts15/shipping-consignment-service/proto/consignment"
	vesselProto "github.com/eluts15/shipping-vessel-service/proto/vessel"
)

type handler struct {
	session      *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

func (s *handler) GetRepo() Repository {
	return &ConsignmentRepository{s.Session.Clone()}
}

func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()

	log.Println("Creating: ", req)

	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	err = repo.Create(req)
	if err != nil {
		return err
	}

	res.Create = true
	res.Consignment = req
	return nil
}

func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()

	consignments, err := repo.GetAll()
	if err != nil {
		return err
	}

	res.Consignments = consignments
	return nil
}
