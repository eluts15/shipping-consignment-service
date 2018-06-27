/*Code is responsible for interacting with our Mongodb database.*/
package main

import (
	pb "github.com/eluts15/microservices/shipping/consignment-service/proto/consignment"
	"gopkg.in.v2"
)

const (
	dbName               = "shipping"
	cosignmentCollection = "consignments"
)

type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

// Each database request has its on db session/connection.
type ConsignmentRepository struct {
	session *mgo.Session
}

func (repo *ConsignmentRepository) Create(*pb.Consignment) error {
	return repo.collection().Insert(consignment)
}

func (repo *ConsignmentRepository) GetAll([]*pb.Consignment, error) {
	var consignments []*pb.Consignment

	err := repo.collection().Find(nil).All(&consignments)
	return consignments, err
}

// Close closes the database session after each query is ran.
// This is safer and more efficient.
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}
