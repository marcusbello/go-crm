package client

import (
	"context"
	"fmt"
	"github.com/marcusbello/go-crm/internal/server/storage"
	pb "github.com/marcusbello/go-crm/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
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

// Error indicates if there was an error in the Customer output stream.
func (c Customer) Error() error {
	return c.err
}

// CallOption are optional options for an RPC call.
type CallOption func(co *callOptions)

type callOptions struct {
	trace *string
}

// TraceID will cause the RPC call to execute a trace on the service and return "s" to the ID.
// If s == nil, this will ignore the option. If "s" is not set after the call finishes, then
// no trace was made.
func TraceID(s *string) CallOption {
	return func(co *callOptions) {
		if s == nil {
			return
		}
		co.trace = s
	}
}

// AddCustomers adds customers to the service and returns their unique identities in the
// same order as being added.
func (c *Client) AddCustomers(ctx context.Context, customers []*pb.Customer, options ...CallOption) ([]string, error) {
	if len(customers) == 0 {
		return nil, nil
	}

	for _, p := range customers {
		if err := storage.ValidateCustomer(ctx, p, false); err != nil {
			return nil, err
		}
	}
	var header metadata.MD
	ctx, gOpts, f := handleCallOptions(ctx, &header, options)
	defer f()

	resp, err := c.client.AddCustomers(ctx, &pb.AddCustomersReq{Customers: customers}, gOpts...)
	if err != nil {
		return nil, err
	}
	return resp.Ids, nil
}

// UpdateCustomers updates customers that already exist in the system.
func (c *Client) UpdateCustomers(ctx context.Context, customers []*pb.Customer, options ...CallOption) error {
	if len(customers) == 0 {
		return nil
	}

	for _, p := range customers {
		if err := storage.ValidateCustomer(ctx, p, true); err != nil {
			return err
		}
	}

	var header metadata.MD
	ctx, gOpts, f := handleCallOptions(ctx, &header, options)
	defer f()

	_, err := c.client.UpdateCustomers(ctx, &pb.UpdateCustomersReq{Customers: customers}, gOpts...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCustomers deletes pets with the IDs passed. If the ID doesn't exist, the
// system ignores it.
func (c *Client) DeleteCustomers(ctx context.Context, ids []string, options ...CallOption) error {
	if len(ids) == 0 {
		return nil
	}

	var header metadata.MD
	ctx, gOpts, f := handleCallOptions(ctx, &header, options)
	defer f()

	_, err := c.client.DeleteCustomers(ctx, &pb.DeleteCustomersReq{Ids: ids}, gOpts...)
	if err != nil {
		return err
	}
	return nil
}

// SearchCustomers searches the pet store for pets matching the filter. If the filter contains
// no entries, then all pets will be returned.
func (c *Client) SearchCustomers(ctx context.Context, filter *pb.SearchCustomersReq, options ...CallOption) (chan Customer, error) {
	if filter == nil {
		return nil, fmt.Errorf("the filter cannot be nil")
	}

	var header metadata.MD
	ctx, gOpts, f := handleCallOptions(ctx, &header, options)

	stream, err := c.client.SearchCustomers(ctx, filter, gOpts...)
	if err != nil {
		return nil, err
	}
	ch := make(chan Customer, 1)
	go func() {
		defer close(ch)
		defer f()

		for {
			// u for customer as c cant be used as the variable name again
			u, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				ch <- Customer{err: err}
				return
			}
			ch <- Customer{Customer: u}
		}
	}()
	return ch, nil
}

// options wrapper
func handleCallOptions(ctx context.Context, header *metadata.MD, options []CallOption) (context.Context, []grpc.CallOption, func()) {
	opts := callOptions{}
	for _, o := range options {
		o(&opts)
	}
	var gOpts []grpc.CallOption

	if opts.trace != nil {
		(*header)["trace"] = nil
		gOpts = append(gOpts, grpc.Header(header))
	}

	f := func() {
		if opts.trace != nil {
			if len((*header)["otel.traceID"]) != 0 {
				*opts.trace = (*header)["otel.traceID"][0]
			}
		}
	}

	return ctx, gOpts, f
}
