package grpc

import (
	"context"
	"time"

	"github.com/mrkovshik/memento/internal/model"
	pb "github.com/mrkovshik/memento/proto"
)

func (s *Server) AddUser(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	token, err := s.service.AddUser(ctx, model.User{
		Name:     request.User.Name,
		Email:    request.User.Email,
		Password: request.User.Password,
	})
	if err != nil {
		return &pb.AddUserResponse{Error: err.Error()}, err
	}
	return &pb.AddUserResponse{
		Token: token,
	}, nil
}

func (s *Server) GetToken(ctx context.Context, request *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	token, err := s.service.GetToken(ctx, model.User{
		Name:     request.User.Name,
		Email:    request.User.Email,
		Password: request.User.Password,
	})
	if err != nil {
		return &pb.GetTokenResponse{Error: err.Error()}, err
	}
	return &pb.GetTokenResponse{
		Token: token,
	}, nil
}

func (s *Server) AddCredential(ctx context.Context, request *pb.AddCredentialRequest) (*pb.AddCredentialResponse, error) {
	if err := s.service.AddCredential(ctx, model.Credential{
		Login:    request.Credential.Login,
		Password: request.Credential.Password,
		Meta:     request.Credential.Meta,
	}); err != nil {
		return &pb.AddCredentialResponse{Error: err.Error()}, err
	}
	return &pb.AddCredentialResponse{}, nil
}

func (s *Server) ListCredentials(ctx context.Context, _ *pb.ListCredentialsRequest) (*pb.ListCredentialsResponse, error) {

	credentials, err := s.service.ListCredentials(ctx)
	if err != nil {
		return &pb.ListCredentialsResponse{Error: err.Error()}, err
	}
	response := make([]*pb.Credential, len(credentials))
	for i, credential := range credentials {
		response[i] = &pb.Credential{
			Login:     credential.Login,
			Password:  credential.Password,
			Meta:      credential.Meta,
			Uuid:      credential.UUID.String(),
			CreatedAt: credential.CreatedAt.Format(time.DateTime),
			UpdatedAt: credential.UpdatedAt.Format(time.DateTime),
		}
	}
	return &pb.ListCredentialsResponse{Credentials: response}, nil
}

func (s *Server) AddCard(ctx context.Context, request *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	if err := s.service.AddCard(ctx, model.CardData{
		Number: request.CardData.Number,
		Name:   request.CardData.Name,
		CVV:    uint(request.CardData.Cvv),
		Expiry: request.CardData.Expiry,
		Meta:   request.CardData.Meta,
	}); err != nil {
		return &pb.AddCardResponse{Error: err.Error()}, err
	}
	return &pb.AddCardResponse{}, nil
}

func (s *Server) ListCards(ctx context.Context, _ *pb.ListCardsRequest) (*pb.ListCardsResponse, error) {

	cards, err := s.service.ListCards(ctx)
	if err != nil {
		return &pb.ListCardsResponse{Error: err.Error()}, err
	}
	response := make([]*pb.CardData, len(cards))
	for i, card := range cards {
		response[i] = &pb.CardData{
			Number:    card.Number,
			Name:      card.Name,
			Cvv:       uint32(card.CVV),
			Expiry:    card.Expiry,
			Meta:      card.Meta,
			Uuid:      card.UUID.String(),
			CreatedAt: card.CreatedAt.Format(time.DateTime),
			UpdatedAt: card.UpdatedAt.Format(time.DateTime),
		}
	}
	return &pb.ListCardsResponse{Cards: response}, nil
}

func (s *Server) ListVariousData(ctx context.Context, _ *pb.ListVariousDataRequest) (*pb.ListVariousDataResponse, error) {

	dataList, err := s.service.ListVariousData(ctx)
	if err != nil {
		return &pb.ListVariousDataResponse{Error: err.Error()}, err
	}
	response := make([]*pb.VariousData, len(dataList))
	for i, data := range dataList {
		response[i] = &pb.VariousData{
			DataType:  int32(data.DataType),
			Meta:      data.Meta,
			Uuid:      data.UUID.String(),
			CreatedAt: data.CreatedAt.Format(time.DateTime),
			UpdatedAt: data.UpdatedAt.Format(time.DateTime),
		}
	}
	return &pb.ListVariousDataResponse{Data: response}, nil
}
