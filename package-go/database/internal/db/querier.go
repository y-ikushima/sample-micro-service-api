// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	CreateSystem(ctx context.Context, arg CreateSystemParams) (System, error)
	DeleteSystem(ctx context.Context, id uuid.UUID) error
	GetSystem(ctx context.Context, id uuid.UUID) (System, error)
	GetSystemByName(ctx context.Context, systemname string) (System, error)
	GetSystems(ctx context.Context) ([]System, error)
	GetSystemsByEmail(ctx context.Context, mailaddress string) ([]System, error)
	GetSystemsByLocalGovernment(ctx context.Context, localgovernmentid sql.NullString) ([]System, error)
	SearchSystems(ctx context.Context, arg SearchSystemsParams) ([]System, error)
	UpdateSystem(ctx context.Context, arg UpdateSystemParams) (System, error)
	UpdateSystemContact(ctx context.Context, arg UpdateSystemContactParams) (System, error)
}

var _ Querier = (*Queries)(nil)
