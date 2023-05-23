package api

import (
	"context"
)

type Dvd struct {
	Path string
	Ip   string
}

type HypervDvdClient interface {
	CreateDvd(ctx context.Context, path string, ip string) (err error)
	DeleteDvd(ctx context.Context, path string) (err error)
	GetDvd(ctx context.Context, path string, ip string) (result Dvd, err error)
}
