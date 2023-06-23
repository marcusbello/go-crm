package mem

import (
	"context"
	"errors"
	"fmt"
	"github.com/marcusbello/go-crm/internal/server/storage"
	pb "github.com/marcusbello/go-crm/proto"
	"sync"
)

// Data implements storage.Data.
type Data struct {
	mu    sync.RWMutex // protects the items in this block
	names map[string]map[string]*pb.Customer
	ids   map[string]*pb.Customer

	// searches contains all the search calls that must be done
	// when we do a search. This is populated in New().
	searches []func(context.Context, *pb.SearchCustomersReq) []string
}

// New is the constructor for Data.
func New() *Data {
	d := Data{
		names: map[string]map[string]*pb.Customer{},
		ids:   map[string]*pb.Customer{},
	}
	d.searches = []func(context.Context, *pb.SearchCustomersReq) []string{
		d.byNames,
	}
	return &d
}

// AddCustomers implements storage.Data.AddCustomers().
func (d *Data) AddCustomers(ctx context.Context, customers []*pb.Customer) error {

	d.mu.RLock()
	// Make sure that none of these IDs somehow exist already.
	for _, p := range customers {
		if _, ok := d.ids[p.Id]; ok {
			return errors.New(fmt.Sprintf("customer with ID(%s) is already present", p.Id))
		}
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.populate(ctx, customers)
	return nil
}

// UpdateCustomers implements storage.Data.AddCustomers().
func (d *Data) UpdateCustomers(ctx context.Context, customers []*pb.Customer) error {
	d.mu.RLock()
	// Make sure that ALL of these IDs somehow exist already.
	for _, p := range customers {
		if _, ok := d.ids[p.Id]; !ok {
			return errors.New(fmt.Sprintf("customer with ID(%s) doesn't exist", p.Id))
		}
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.populate(ctx, customers)
	return nil
}

func (d *Data) populate(ctx context.Context, customers []*pb.Customer) {

	for _, c := range customers {
		d.ids[c.Id] = c
		if v, ok := d.names[c.Name]; ok {
			v[c.Id] = c
		} else {
			d.names[c.Name] = map[string]*pb.Customer{
				c.Id: c,
			}
		}
	}
}

// DeleteCustomers implements storage.Data.DeleteCustomers().
func (d *Data) DeleteCustomers(ctx context.Context, ids []string) error {

	d.mu.Lock()
	defer d.mu.Unlock()

	for _, id := range ids {
		p, ok := d.ids[id]
		if !ok {
			continue
		}
		delete(d.ids, id)
		if v, ok := d.names[p.Name]; ok {
			if len(v) == 1 {
				delete(d.names, p.Name)
			} else {
				delete(v, id)
			}
		}

	}
	return nil
}

// SearchCustomers implements storage.Data.SearchCustomers).
func (d *Data) SearchCustomers(ctx context.Context, filter *pb.SearchCustomersReq) chan storage.SearchItem {
	customersCh := make(chan storage.SearchItem, 1)

	go func() {
		defer close(customersCh)
		d.searchCustomers(ctx, filter, customersCh)
	}()

	return customersCh
}

func (d *Data) searchCustomers(ctx context.Context, filter *pb.SearchCustomersReq, out chan storage.SearchItem) {

	d.mu.RLock()
	defer d.mu.RUnlock()

	filters := 0
	if len(filter.Names) > 0 {
		filters++
	}

	// They didn't provide filters, so just return everything.
	if filters == 0 {
		d.returnAll(ctx, out)
		return
	}

	searchCh := make(chan []string, len(d.searches))
	wg := sync.WaitGroup{}
	wg.Add(len(d.searches))

	goCount := 0
	// Spin off our searches.
	for _, search := range d.searches {
		goCount++
		search := search
		go func() {
			defer wg.Done()
			r := search(ctx, filter)
			select {
			case <-ctx.Done():
			case searchCh <- r:
			}
		}()
	}

	// Wait for our searches to complete then close our searchCh.
	go func() { wg.Wait(); close(searchCh) }()

	// Collect all IDs from searches and count them. When one hits
	// the total number of filters send the matching customer to the caller.
	m := map[string]int{}
	matchCh := make(chan string, 1)
	go func() {
		defer close(matchCh)
		for ids := range searchCh {
			for _, id := range ids {
				count := m[id]
				count++
				m[id] = count
				if count == filters {
					matchCh <- id
				}
			}
		}
	}()

	// This handles all our matches getting returned.
	valCount := 0
	defer func() {
		if valCount > 0 {
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case id, ok := <-matchCh:
			if !ok {
				return
			}

			out <- storage.SearchItem{Customer: d.ids[id]}
			valCount++
		}
	}
}

// returnAll streams all the customers that we have.
func (d *Data) returnAll(ctx context.Context, out chan storage.SearchItem) {

	count := 0
	for _, c := range d.ids {
		count++
		select {
		case <-ctx.Done():
			return
		case out <- storage.SearchItem{Customer: c}:
		}
	}
}

// byNames returns IDs of customers that have the names matched in the filter.
func (d *Data) byNames(ctx context.Context, filter *pb.SearchCustomersReq) []string {
	if len(filter.Names) == 0 {
		return nil
	}

	count := 0
	var ids []string
	for _, n := range filter.Names {
		count++
		if ctx.Err() != nil {
			return nil
		}
		c, ok := d.names[n]
		if !ok {
			continue
		}
		for id := range c {
			ids = append(ids, id)
		}
	}
	return ids
}
