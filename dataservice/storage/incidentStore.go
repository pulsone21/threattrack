package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pulsone21/threattrack/lib/entities"
)

type RequestIncident struct {
	Name         string
	Severity     string
	IncidentType int
}

type IncidentStore struct {
	storage *MySqlStorage
	EntityStore[entities.Incident]
	db *sql.DB
}

func NewIncidentStore(storage *MySqlStorage) *IncidentStore {
	creatTable, err := LoadRawSQL("incidents/CreateTable.sql")
	if err != nil {
		panic(err)
	}
	storage.Db.Exec(creatTable)
	return &IncidentStore{
		storage: storage,
		db:      storage.Db,
	}
}

func (i *IncidentStore) Get(ctx context.Context, id string) (*entities.Incident, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidents/GetById.sql")
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
	var inc entities.Incident
	if err := inc.ScanTo(res.Scan); err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	// ? We are not joning the tasks from the sql instead we add them here
	// ? maybe an optimization for later on to do it in the sql, but requieres a rebuild of the scan func
	qP := DefaultParams()
	qP.Query["incident_id"] = id
	tasks, err1 := i.storage.TaskStore.GetQuery(ctx, *qP)
	if err1 != nil {
		return nil, err1
	}
	if *tasks == nil {
		fmt.Println("Tasks are nil, creating empty array")
		tasks = &[]entities.Task{}
	}
	inc.Tasks = *tasks
	return &inc, nil
}

func (i *IncidentStore) GetAll(ctx context.Context, qP QueryParameter) (*[]entities.Incident, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidents/GetAll.sql")
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

	incs := []entities.Incident{}
	for res.Next() {
		var inc entities.Incident
		if err := inc.ScanTo(res.Scan); err != nil {
			return nil, entities.InternalServerError(err, uri)
		} else {
			incs = append(incs, inc)
		}
	}
	return &incs, nil
}

func (i *IncidentStore) GetQuery(ctx context.Context, qP QueryParameter) (*[]entities.Incident, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	rawSql, err := LoadRawSQL("incidents/GetQuery.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	whiteList := i.createWhitelist()
	if whiteList == nil {
		return nil, entities.InternalServerError(fmt.Errorf("couldn't create whitelist for entity"), uri)
	}
	for key, val := range qP.Query {
		if !CheckWhitelist(key, val, whiteList) {
			return nil, entities.BadRequestError(fmt.Errorf("whitelist check failed"), uri)
		}
	}
	finalSql := FinalizeSQL(rawSql, "incidents", qP)
	rows, err := i.db.Query(finalSql, qP.Limit, qP.Offset)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	if rows.Err() != nil {
		if rows.Err() == sql.ErrNoRows {
			return nil, entities.NotFoundError(fmt.Errorf("no incidents found"), uri)
		}
		return nil, entities.InternalServerError(rows.Err(), uri)
	}
	defer rows.Close()
	var incs []entities.Incident
	for rows.Next() {
		var i entities.Incident
		err := i.ScanTo(rows.Scan)
		if err != nil {
			return nil, entities.InternalServerError(rows.Err(), uri)
		}
		incs = append(incs, i)
	}
	return &incs, nil
}

func (i *IncidentStore) Create(ctx context.Context, inc *entities.Incident) (*entities.Incident, *entities.ApiError) {
	fmt.Println("creating new inc from ", inc)
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidents/Create.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	_, err = i.db.Exec(loadedSql, inc.Id, inc.Name, inc.Severity, inc.IncidentType.Id)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	return inc, nil
}

func (i *IncidentStore) Update(entity entities.Incident) (*entities.Incident, *entities.ApiError) {
	panic("not implemented") // TODO: Implement
}

func (i *IncidentStore) Delete(ctx context.Context, id string) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("incidents/Delete.sql")
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	_, err = i.db.Exec(loadedSql, id)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (i *IncidentStore) createWhitelist() Whitelist {
	sql, err := LoadRawSQL("incidenttypes/GetAll.sql")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	res, err := i.db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Close()
	iTs := []string{}

	for res.Next() {
		var iT entities.IncidentType
		err = iT.ScanTo(res.Scan)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		iTs = append(iTs, iT.Name)
	}
	incWhitelist := map[string][]string{
		"Severity": {string(entities.Low), string(entities.Medium), string(entities.High), string(entities.Critical)},
		"Status":   {string(entities.Pending), string(entities.Open), string(entities.Active), string(entities.Closed)},
		"Type":     iTs,
	}
	return incWhitelist
}
