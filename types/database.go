package types

import "context"

type DB interface {
	All(context.Context) ([]Link, error)
	Get(context.Context, string) (Link, error)
	Create(context.Context, Link) error
	Delete(context.Context, string) error
}
