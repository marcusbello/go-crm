package storage

import (
	"context"
	"errors"
	pb "github.com/marcusbello/go-crm/proto/crm"
)

// Data represents our data storage.
type Data interface {
	// AddCustomers adds pet entries into storage.
	AddCustomers(ctx context.Context, customers []*pb.Customer) error
	// UpdateCustomers updates customers entries in storage.
	UpdateCustomers(ctx context.Context, customers []*pb.Customer) error
	// DeleteCustomers deletes customers in storage by their ID. Will not error
	// on IDs not found.
	DeleteCustomers(ctx context.Context, ids []string) error
	// SearchCustomers searches storage for customers entries that match the
	// filter.
	SearchCustomers(ctx context.Context, filter *pb.SearchCustomersReq) chan SearchItem
}

// SearchItem is an item returned by a search.
type SearchItem struct {
	// Customer is the customer that matched the search filters.
	Customer *pb.Customer
	// Error indicates that there was an error. If set the channel
	// will close after this entry.
	Error error
}

// ValidateCustomer validates that *pb.Customer has valid fields.
func ValidateCustomer(ctx context.Context, c *pb.Customer, forUpdate bool) error {
	if forUpdate && c.Id == "" {
		return errors.New("updates must have the Id field set")
	} else {
		if !forUpdate && c.Id != "" {
			return errors.New("cannot set the Id field")
		}
	}
	//c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return errors.New("cannot have a customer without a name")
	}

	return nil

}
