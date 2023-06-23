package client

import (
	pb "github.com/marcusbello/go-crm/proto"
	"google.golang.org/grpc"
)

// Client is a client to the CRM service.
type Client struct {
	client pb.CRMClient
	conn   *grpc.ClientConn
}

// New is the constructor for Client. addr is the server's [host]:[port].
func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		client: pb.NewCRMClient(conn),
		conn:   conn,
	}, nil
}

// Customer is a wrapper around a *pb.Customer that can return Go versions of
// fields and errors if the returned stream has an error.
type Customer struct {
	*pb.Customer
	err error
}

// Proto will give the Customer's proto representation.
func (c Customer) Proto() *pb.Customer {
	return c.Customer
}
