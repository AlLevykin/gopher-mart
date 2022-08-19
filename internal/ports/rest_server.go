package ports

import "context"

type RESTServer interface {
	Start() error
	Stop(ctx context.Context) error
}
