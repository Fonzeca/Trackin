package domain

import "context"

type ISender interface {
	SendMessage(context context.Context, destination string, message []byte) error
}
