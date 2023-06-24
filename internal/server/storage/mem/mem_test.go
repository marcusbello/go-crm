package mem

import (
	"context"
	"github.com/kylelemons/godebug/pretty"
	"github.com/marcusbello/go-crm/internal/server/storage"
	pb "github.com/marcusbello/go-crm/proto/crm"
	"google.golang.org/protobuf/proto"
	"sort"
	"strconv"
	"testing"
)

// This tests we implement the interface.
var _ storage.Data = &Data{}

var customers = []*pb.Customer{
	{
		Id:   "0",
		Name: "Adam",
	},
	{
		Id:   "1",
		Name: "Becky",
	},
	{
		Id:   "2",
		Name: "Calvin",
	},
	{
		Id:   "3",
		Name: "David",
	},
	{
		Id:   "4",
		Name: "Elaine",
	},
	{
		Id:   "5",
		Name: "Elaine",
	},
}

// makeCustomers takes the global "customers" var and clones everything in it and puts it into
// a *Data, so we have test data.
func makeCustomers() *Data {
	d := New()

	n := []*pb.Customer{}
	for _, p := range customers {
		n = append(n, proto.Clone(p).(*pb.Customer))
	}

	d.AddCustomers(context.Background(), n)
	return d
}

func TestByNames(t *testing.T) {
	d := makeCustomers()

	got := d.byNames(context.Background(), &pb.SearchCustomersReq{Names: []string{"David", "Elaine"}})
	sort.Strings(got)

	want := []string{"3", "4", "5"}
	if diff := pretty.Compare(want, got); diff != "" {
		t.Errorf("TestByNames: -want/+got:\n%s", diff)
	}
}

func TestDeleteCustomers(t *testing.T) {
	d := makeCustomers()

	deletions := []string{"3", "5", "20"}

	if err := d.DeleteCustomers(context.Background(), deletions); err != nil {
		t.Fatalf("TestDeleteCustomers: got err == %v, want err == nil", err)
	}

	// Don't check the last deletion, it is only there to make sure
	// a non-existent value doesn't do anything.
	for _, id := range deletions[:len(deletions)-1] {
		if _, ok := d.ids[id]; ok {
			t.Errorf("TestDeleteCustomers: found ids[%s]", id)
		}
		i, _ := strconv.Atoi(id)

		if m, ok := d.names[customers[i].Name]; ok {
			if _, ok := m[id]; ok {
				t.Errorf("TestDeleteCustomers: found(%s) in names", id)
			}
		}

	}
}

func TestSearchCustomers(t *testing.T) {
	d := makeCustomers()

	ch := d.SearchCustomers(
		context.Background(),
		&pb.SearchCustomersReq{
			Names: []string{
				"Becky",
				"Calvin",
				"David",
				"Elaine",
			},
		},
	)

	got := []storage.SearchItem{}
	for item := range ch {
		got = append(got, item)
	}

	want := []storage.SearchItem{{Customer: customers[4]}}

	config := pretty.Config{TrackCycles: true}
	if diff := config.Compare(want, got); diff != "" {
		t.Errorf("TestSearchCustomers: -want/+got:\n%s", diff)
	}
}
