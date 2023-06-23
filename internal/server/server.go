package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/marcusbello/go-crm/internal/server/storage"
	pb "github.com/marcusbello/go-crm/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"sync"
)

// API that implements our grpc server API
type API struct {
	pb.UnimplementedCRMServer

	addr  string
	store storage.Data

	gOpts      []grpc.ServerOption
	grpcServer *grpc.Server

	mu sync.Mutex
}

// Option to customize grpc
type Option func(a *API)

// WithGRPCOpts creates the gRPC server with the options passed.
func WithGRPCOpts(opts ...grpc.ServerOption) Option {
	return func(a *API) {
		a.gOpts = append(a.gOpts, opts...)
	}
}

// New constructor for the API
func New(addr string, store storage.Data, options ...Option) (*API, error) {
	a := &API{addr: addr, store: store}
	for _, o := range options {
		o(a)
	}
	a.grpcServer = grpc.NewServer(a.gOpts...)
	a.grpcServer.RegisterService(&pb.CRM_ServiceDesc, a)
	reflection.Register(a.grpcServer)

	return a, nil
}

// Start starts the server. This blocks until Stop() is called.
func (a *API) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	lis, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(lis)
}

// Stop stops the server.
func (a *API) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.grpcServer.Stop()
}

// AddCustomers adds customers to the CRM
func (a *API) AddCustomers(ctx context.Context, req *pb.AddCustomersReq) (resp *pb.AddCustomersResp, err error) {
	// Actual work.
	ids := make([]string, 0, len(req.Customers))
	for _, p := range req.Customers {
		if err := storage.ValidateCustomer(ctx, p, false); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		p.Id = uuid.New().String()
		ids = append(ids, p.Id)
	}

	if err = a.store.AddCustomers(ctx, req.Customers); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AddCustomersResp{Ids: ids}, nil
}

// UpdateCustomers updates customers on the CRM
func (a *API) UpdateCustomers(ctx context.Context, req *pb.UpdateCustomersReq) (resp *pb.UpdateCustomersResp, err error) {
	for _, p := range req.Customers {
		if err = storage.ValidateCustomer(ctx, p, true); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	if err = a.store.UpdateCustomers(ctx, req.Customers); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateCustomersResp{}, nil
}

// DeleteCustomers deletes customers from the CRM.
func (a *API) DeleteCustomers(ctx context.Context, req *pb.DeleteCustomersReq) (resp *pb.DeleteCustomersResp, err error) {
	if err = a.store.DeleteCustomers(ctx, req.Ids); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteCustomersResp{}, nil
}

// SearchCustomers finds customers in the CRM.
func (a *API) SearchCustomers(ctx context.Context, req *pb.SearchCustomersReq, stream pb.CRM_SearchCustomersServer) (err error) {
	count := 0

	ch := a.store.SearchCustomers(ctx, req)
	for item := range ch {
		count++
		if item.Error != nil {
			return status.Error(codes.Internal, item.Error.Error())
		}
		if err := stream.Send(item.Customer); err != nil {
			return err
		}
	}
	if ctx.Err() != nil {
		return status.Error(codes.DeadlineExceeded, stream.Context().Err().Error())
	}
	return nil
}
