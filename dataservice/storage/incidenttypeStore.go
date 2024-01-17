package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pulsone21/threattrack/lib/entities"
)

type RequestIncidentType struct {
	Name string
}
type IncidentTypeStore struct {
	storage *MySqlStorage
	EntityStore[entities.IncidentType]
	db *sql.DB
}

func NewIncidentTypeStore(storage *MySqlStorage) *IncidentTypeStore {
	creatTable, err := LoadRawSQL("incidenttypes/CreateTable.sql")
	if err != nil {
		panic(err)
	}
	storage.Db.Exec(creatTable)
	return &IncidentTypeStore{
		storage: storage,
		db:      storage.Db,
	}
}

func (i *IncidentTypeStore) Get(ctx context.Context, id string) (*entities.IncidentType, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidenttypes/GetById.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	res := i.db.QueryRow(loadedSql, id)
	if res.Err() != nil {
		if res.Err() == sql.ErrNoRows {
			return nil, entities.NotFoundError(fmt.Errorf("no incident found"), uri)
		}
		return nil, entities.InternalServerError(res.Err(), uri)
	}
	var iT entities.IncidentType
	if err := iT.ScanTo(res.Scan); err != nil {
		return nil, entities.InternalServerError(err, ctx.Value("uri").(string))
	}
	return &iT, nil
}

func (i *IncidentTypeStore) GetAll(ctx context.Context, qP QueryParameter) (*[]entities.IncidentType, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidenttypes/GetAll.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	res, err := i.db.Query(loadedSql, qP.Limit, qP.Offset)
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
	iTs := []entities.IncidentType{}
	for res.Next() {
		var iT entities.IncidentType
		if err := iT.ScanTo(res.Scan); err != nil {
			return nil, entities.InternalServerError(err, uri)
		} else {
			iTs = append(iTs, iT)
		}
	}
	return &iTs, nil
}

func (i *IncidentTypeStore) GetQuery(ctx context.Context, qP QueryParameter) (*[]entities.IncidentType, *entities.ApiError) {
	// ! This Entity isn't queryable
	return nil, entities.BadRequestError(fmt.Errorf("not applicable"), "/incidenttypes/query")
}

func (i *IncidentTypeStore) Create(ctx context.Context, incidentType entities.IncidentType) (*entities.IncidentType, *entities.ApiError) {
	fmt.Println("creating new inc from ", incidentType)
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidenttypes/Create.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	res, err := i.db.Exec(loadedSql, incidentType.Name)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	//The Id is auto incrementing in the database, so the given id from the parameters are irrelevant,
	//since we now the acutal id only after inserting to db.
	iD, _ := res.LastInsertId()
	iT := entities.IncidentType{
		Name: incidentType.Name,
		Id:   iD,
	}
	return &iT, nil
}

func (i *IncidentTypeStore) Update(ctx context.Context, entity entities.IncidentType) (*entities.IncidentType, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	return nil, entities.NotImplementedError(fmt.Errorf("not implemented"), uri)
}

func (i *IncidentTypeStore) Delete(ctx context.Context, id string) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidenttypes/Delete.sql")
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	_, err = i.db.Exec(loadedSql, id)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}
