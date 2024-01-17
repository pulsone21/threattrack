package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pulsone21/threattrack/lib/entities"
)

type RequestUser struct {
	Firstname string
	Lastname  string
	Email     string
}

type UserStore struct {
	storage *MySqlStorage
	EntityStore[entities.User]
	db *sql.DB
}

func NewUserStore(storage *MySqlStorage) *UserStore {
	creatTable, err := LoadRawSQL("users/CreateTable.sql")
	if err != nil {
		panic(err)
	}
	storage.Db.Exec(creatTable)
	return &UserStore{
		storage: storage,
		db:      storage.Db,
	}
}

func (u *UserStore) Get(ctx context.Context, id string) (*entities.User, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("users/GetById.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	res := u.db.QueryRow(loadedSql, id)
	if res.Err() != nil {
		if res.Err() == sql.ErrNoRows {
			return nil, entities.NotFoundError(fmt.Errorf("no user found"), uri)
		}
		return nil, entities.InternalServerError(res.Err(), uri)
	}
	var usr entities.User
	if err := usr.ScanTo(res.Scan); err != nil {
		return nil, entities.InternalServerError(err, ctx.Value("uri").(string))
	}
	return &usr, nil
}

func (u *UserStore) GetAll(ctx context.Context, qP QueryParameter) (*[]entities.User, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("users/GetAll.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	res, err := u.db.Query(loadedSql, qP.Limit, qP.Offset)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	if res.Err() != nil {
		if res.Err() == sql.ErrNoRows {
			return nil, entities.NotFoundError(fmt.Errorf("no incidents found"), uri)
		}
		return nil, entities.InternalServerError(res.Err(), uri)
	}
	defer res.Close()
	usrs := []entities.User{}
	for res.Next() {
		var usr entities.User
		if err := usr.ScanTo(res.Scan); err != nil {
			return nil, entities.InternalServerError(err, uri)
		} else {
			usrs = append(usrs, usr)
		}
	}
	return &usrs, nil
}

func (u *UserStore) GetQuery(ctx context.Context, qP QueryParameter) (*[]entities.User, *entities.ApiError) {
	// ! This Entity isn't queryable
	return nil, entities.BadRequestError(fmt.Errorf("not applicable"), "/incidenttypes/query")
}

func (u *UserStore) Create(ctx context.Context, user entities.User) (*entities.User, *entities.ApiError) {
	fmt.Println("creating new inc from ", user)
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidenttypes/Create.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	_, err = u.db.Exec(loadedSql, user.Id, user.Firstname, user.Lastname, user.Email, user.Fullname, user.CreatedAt)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	return &user, nil
}

func (u *UserStore) Update(ctx context.Context, entity entities.User) (*entities.User, *entities.ApiError) {
	panic("not implemented") // TODO: Implement
}

func (u *UserStore) Delete(ctx context.Context, id string) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("users/Delete.sql")
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	_, err = u.db.Exec(loadedSql, id)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}
