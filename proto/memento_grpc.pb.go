// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/memento.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Memento_AddUser_FullMethodName                 = "/api.grpc.Memento/AddUser"
	Memento_GetToken_FullMethodName                = "/api.grpc.Memento/GetToken"
	Memento_AddCredential_FullMethodName           = "/api.grpc.Memento/AddCredential"
	Memento_ListCredentials_FullMethodName         = "/api.grpc.Memento/ListCredentials"
	Memento_AddCard_FullMethodName                 = "/api.grpc.Memento/AddCard"
	Memento_ListCards_FullMethodName               = "/api.grpc.Memento/ListCards"
	Memento_AddVariousData_FullMethodName          = "/api.grpc.Memento/AddVariousData"
	Memento_DownloadVariousDataFile_FullMethodName = "/api.grpc.Memento/DownloadVariousDataFile"
	Memento_ListVariousData_FullMethodName         = "/api.grpc.Memento/ListVariousData"
)

// MementoClient is the client API for Memento service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MementoClient interface {
	AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserResponse, error)
	GetToken(ctx context.Context, in *GetTokenRequest, opts ...grpc.CallOption) (*GetTokenResponse, error)
	AddCredential(ctx context.Context, in *AddCredentialRequest, opts ...grpc.CallOption) (*AddCredentialResponse, error)
	ListCredentials(ctx context.Context, in *ListCredentialsRequest, opts ...grpc.CallOption) (*ListCredentialsResponse, error)
	AddCard(ctx context.Context, in *AddCardRequest, opts ...grpc.CallOption) (*AddCardResponse, error)
	ListCards(ctx context.Context, in *ListCardsRequest, opts ...grpc.CallOption) (*ListCardsResponse, error)
	AddVariousData(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[AddVariousDataRequest, AddVariousDataResponse], error)
	DownloadVariousDataFile(ctx context.Context, in *DownloadVariousDataFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadVariousDataFileResponse], error)
	ListVariousData(ctx context.Context, in *ListVariousDataRequest, opts ...grpc.CallOption) (*ListVariousDataResponse, error)
}

type mementoClient struct {
	cc grpc.ClientConnInterface
}

func NewMementoClient(cc grpc.ClientConnInterface) MementoClient {
	return &mementoClient{cc}
}

func (c *mementoClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddUserResponse)
	err := c.cc.Invoke(ctx, Memento_AddUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mementoClient) GetToken(ctx context.Context, in *GetTokenRequest, opts ...grpc.CallOption) (*GetTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTokenResponse)
	err := c.cc.Invoke(ctx, Memento_GetToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mementoClient) AddCredential(ctx context.Context, in *AddCredentialRequest, opts ...grpc.CallOption) (*AddCredentialResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddCredentialResponse)
	err := c.cc.Invoke(ctx, Memento_AddCredential_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mementoClient) ListCredentials(ctx context.Context, in *ListCredentialsRequest, opts ...grpc.CallOption) (*ListCredentialsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCredentialsResponse)
	err := c.cc.Invoke(ctx, Memento_ListCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mementoClient) AddCard(ctx context.Context, in *AddCardRequest, opts ...grpc.CallOption) (*AddCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddCardResponse)
	err := c.cc.Invoke(ctx, Memento_AddCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mementoClient) ListCards(ctx context.Context, in *ListCardsRequest, opts ...grpc.CallOption) (*ListCardsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCardsResponse)
	err := c.cc.Invoke(ctx, Memento_ListCards_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mementoClient) AddVariousData(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[AddVariousDataRequest, AddVariousDataResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Memento_ServiceDesc.Streams[0], Memento_AddVariousData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[AddVariousDataRequest, AddVariousDataResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Memento_AddVariousDataClient = grpc.ClientStreamingClient[AddVariousDataRequest, AddVariousDataResponse]

func (c *mementoClient) DownloadVariousDataFile(ctx context.Context, in *DownloadVariousDataFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadVariousDataFileResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Memento_ServiceDesc.Streams[1], Memento_DownloadVariousDataFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadVariousDataFileRequest, DownloadVariousDataFileResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Memento_DownloadVariousDataFileClient = grpc.ServerStreamingClient[DownloadVariousDataFileResponse]

func (c *mementoClient) ListVariousData(ctx context.Context, in *ListVariousDataRequest, opts ...grpc.CallOption) (*ListVariousDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListVariousDataResponse)
	err := c.cc.Invoke(ctx, Memento_ListVariousData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MementoServer is the server API for Memento service.
// All implementations must embed UnimplementedMementoServer
// for forward compatibility.
type MementoServer interface {
	AddUser(context.Context, *AddUserRequest) (*AddUserResponse, error)
	GetToken(context.Context, *GetTokenRequest) (*GetTokenResponse, error)
	AddCredential(context.Context, *AddCredentialRequest) (*AddCredentialResponse, error)
	ListCredentials(context.Context, *ListCredentialsRequest) (*ListCredentialsResponse, error)
	AddCard(context.Context, *AddCardRequest) (*AddCardResponse, error)
	ListCards(context.Context, *ListCardsRequest) (*ListCardsResponse, error)
	AddVariousData(grpc.ClientStreamingServer[AddVariousDataRequest, AddVariousDataResponse]) error
	DownloadVariousDataFile(*DownloadVariousDataFileRequest, grpc.ServerStreamingServer[DownloadVariousDataFileResponse]) error
	ListVariousData(context.Context, *ListVariousDataRequest) (*ListVariousDataResponse, error)
	mustEmbedUnimplementedMementoServer()
}

// UnimplementedMementoServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMementoServer struct{}

func (UnimplementedMementoServer) AddUser(context.Context, *AddUserRequest) (*AddUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedMementoServer) GetToken(context.Context, *GetTokenRequest) (*GetTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetToken not implemented")
}
func (UnimplementedMementoServer) AddCredential(context.Context, *AddCredentialRequest) (*AddCredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCredential not implemented")
}
func (UnimplementedMementoServer) ListCredentials(context.Context, *ListCredentialsRequest) (*ListCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCredentials not implemented")
}
func (UnimplementedMementoServer) AddCard(context.Context, *AddCardRequest) (*AddCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCard not implemented")
}
func (UnimplementedMementoServer) ListCards(context.Context, *ListCardsRequest) (*ListCardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCards not implemented")
}
func (UnimplementedMementoServer) AddVariousData(grpc.ClientStreamingServer[AddVariousDataRequest, AddVariousDataResponse]) error {
	return status.Errorf(codes.Unimplemented, "method AddVariousData not implemented")
}
func (UnimplementedMementoServer) DownloadVariousDataFile(*DownloadVariousDataFileRequest, grpc.ServerStreamingServer[DownloadVariousDataFileResponse]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadVariousDataFile not implemented")
}
func (UnimplementedMementoServer) ListVariousData(context.Context, *ListVariousDataRequest) (*ListVariousDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVariousData not implemented")
}
func (UnimplementedMementoServer) mustEmbedUnimplementedMementoServer() {}
func (UnimplementedMementoServer) testEmbeddedByValue()                 {}

// UnsafeMementoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MementoServer will
// result in compilation errors.
type UnsafeMementoServer interface {
	mustEmbedUnimplementedMementoServer()
}

func RegisterMementoServer(s grpc.ServiceRegistrar, srv MementoServer) {
	// If the following call pancis, it indicates UnimplementedMementoServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Memento_ServiceDesc, srv)
}

func _Memento_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MementoServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Memento_AddUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MementoServer).AddUser(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Memento_GetToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MementoServer).GetToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Memento_GetToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MementoServer).GetToken(ctx, req.(*GetTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Memento_AddCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MementoServer).AddCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Memento_AddCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MementoServer).AddCredential(ctx, req.(*AddCredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Memento_ListCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MementoServer).ListCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Memento_ListCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MementoServer).ListCredentials(ctx, req.(*ListCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Memento_AddCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MementoServer).AddCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Memento_AddCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MementoServer).AddCard(ctx, req.(*AddCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Memento_ListCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MementoServer).ListCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Memento_ListCards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MementoServer).ListCards(ctx, req.(*ListCardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Memento_AddVariousData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MementoServer).AddVariousData(&grpc.GenericServerStream[AddVariousDataRequest, AddVariousDataResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Memento_AddVariousDataServer = grpc.ClientStreamingServer[AddVariousDataRequest, AddVariousDataResponse]

func _Memento_DownloadVariousDataFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadVariousDataFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MementoServer).DownloadVariousDataFile(m, &grpc.GenericServerStream[DownloadVariousDataFileRequest, DownloadVariousDataFileResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Memento_DownloadVariousDataFileServer = grpc.ServerStreamingServer[DownloadVariousDataFileResponse]

func _Memento_ListVariousData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVariousDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MementoServer).ListVariousData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Memento_ListVariousData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MementoServer).ListVariousData(ctx, req.(*ListVariousDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Memento_ServiceDesc is the grpc.ServiceDesc for Memento service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Memento_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.grpc.Memento",
	HandlerType: (*MementoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddUser",
			Handler:    _Memento_AddUser_Handler,
		},
		{
			MethodName: "GetToken",
			Handler:    _Memento_GetToken_Handler,
		},
		{
			MethodName: "AddCredential",
			Handler:    _Memento_AddCredential_Handler,
		},
		{
			MethodName: "ListCredentials",
			Handler:    _Memento_ListCredentials_Handler,
		},
		{
			MethodName: "AddCard",
			Handler:    _Memento_AddCard_Handler,
		},
		{
			MethodName: "ListCards",
			Handler:    _Memento_ListCards_Handler,
		},
		{
			MethodName: "ListVariousData",
			Handler:    _Memento_ListVariousData_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AddVariousData",
			Handler:       _Memento_AddVariousData_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadVariousDataFile",
			Handler:       _Memento_DownloadVariousDataFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/memento.proto",
}
