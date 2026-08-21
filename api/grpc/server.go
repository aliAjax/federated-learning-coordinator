package grpcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/model"
	app "github.com/example/federated-learning-coordinator/internal/round"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"net"
)

type jsonCodec struct{}

func (jsonCodec) Name() string                    { return "json" }
func (jsonCodec) Marshal(v any) ([]byte, error)   { return json.Marshal(v) }
func (jsonCodec) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }

type API struct{ App *app.Service }
type CoordinatorServer interface{}
type Empty struct{}
type RoundRequest struct {
	ID string `json:"id"`
}
type UpdateRequest struct {
	RoundID string              `json:"round_id"`
	Update  model.UpdateRequest `json:"update"`
}

func (a *API) GetRound(ctx context.Context, in *RoundRequest) (*model.Round, error) {
	return a.App.GetRound(ctx, in.ID)
}
func (a *API) SubmitUpdate(ctx context.Context, in *UpdateRequest) (*model.Update, error) {
	return a.App.Submit(ctx, in.RoundID, in.Update)
}
func (a *API) Health(context.Context, *Empty) (*map[string]string, error) {
	v := map[string]string{"status": "ok"}
	return &v, nil
}
func Register(s *grpc.Server, a *API) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "federated.Coordinator", HandlerType: (*CoordinatorServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "GetRound", Handler: roundHandler}, {MethodName: "SubmitUpdate", Handler: updateHandler}, {MethodName: "Health", Handler: healthHandler}}}, a)
}
func roundHandler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(RoundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor != nil {
		info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/federated.Coordinator/GetRound"}
		return interceptor(ctx, in, info, func(ctx context.Context, req any) (any, error) { return srv.(*API).GetRound(ctx, req.(*RoundRequest)) })
	}
	return srv.(*API).GetRound(ctx, in)
}
func updateHandler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor != nil {
		info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/federated.Coordinator/SubmitUpdate"}
		return interceptor(ctx, in, info, func(ctx context.Context, req any) (any, error) {
			return srv.(*API).SubmitUpdate(ctx, req.(*UpdateRequest))
		})
	}
	return srv.(*API).SubmitUpdate(ctx, in)
}
func healthHandler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(*API).Health(ctx, in)
}
func Serve(ctx context.Context, addr string, a *API) (func(), error) {
	encoding.RegisterCodec(jsonCodec{})
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("grpc listen: %w", err)
	}
	s := grpc.NewServer()
	Register(s, a)
	go func() { <-ctx.Done(); s.GracefulStop(); _ = lis.Close() }()
	go func() { _ = s.Serve(lis) }()
	return func() { s.GracefulStop(); _ = lis.Close() }, nil
}
