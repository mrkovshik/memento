package api

import "context"

type Service interface {
	AddUser(ctx context.Context, username string, password string) error
}
