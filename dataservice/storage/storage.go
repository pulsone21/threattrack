package storage

import (
	"context"
	"net/http"

	"github.com/pulsone21/threattrack/lib/entities"
)

type DBConfig struct {
	Address  string
	Port     string
	User     string
	Password string
	Database string
}

type Storage interface {
	HandleGetAll(context.Context, http.ResponseWriter, *http.Request) (*entities.ApiResponse, *entities.ApiError)
	HandleGetID(context.Context, http.ResponseWriter, *http.Request) (*entities.ApiResponse, *entities.ApiError)
	HandleGetQuery(context.Context, http.ResponseWriter, *http.Request) (*entities.ApiResponse, *entities.ApiError)
	HandleCreate(context.Context, http.ResponseWriter, *http.Request) (*entities.ApiResponse, *entities.ApiError)
	HandleDelete(context.Context, http.ResponseWriter, *http.Request) (*entities.ApiResponse, *entities.ApiError)
	HandleUpdate(context.Context, http.ResponseWriter, *http.Request) (*entities.ApiResponse, *entities.ApiError)
}

type EntityStore[T entities.Entity] interface {
	Get(ctx context.Context, id string) (*T, *entities.ApiError)
	GetAll(ctx context.Context, qP QueryParameter) (*[]T, *entities.ApiError)
	GetQuery(ctx context.Context, qP QueryParameter) (*[]T, *entities.ApiError)
	Create(ctx context.Context, entity T) (*T, *entities.ApiError)
	Update(ctx context.Context, entity T) (*T, *entities.ApiError)
	Delete(ctx context.Context, id string) *entities.ApiError
}

type Whitelist map[string][]string
