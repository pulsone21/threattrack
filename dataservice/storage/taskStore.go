package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pulsone21/threattrack/lib/entities"
)

type RequestTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerId     string `json:"owner_id"`
	State       string `json:"state"`
	Priority    string `json:"priority"`
	IncidentID  string `json:"incident_id"`
}

type TaskStore struct {
	storage *MySqlStorage
	EntityStore[*entities.Task]
	db *sql.DB
}

func NewTaskStore(storage *MySqlStorage) *TaskStore {
	creatTable, err := LoadRawSQL("tasks/CreateTable.sql")
	if err != nil {
		panic(err)
	}
	storage.Db.Exec(creatTable)
	return &TaskStore{
		storage: storage,
		db:      storage.Db,
	}
}

func (i *TaskStore) Get(ctx context.Context, id string) (*entities.Task, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("tasks/GetById.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	res := i.db.QueryRow(loadedSql, id)
	if res.Err() != nil {
		if res.Err() == sql.ErrNoRows {
			return nil, entities.NotFoundError(fmt.Errorf("no task found"), uri)
		}
		return nil, entities.InternalServerError(res.Err(), uri)
	}
	var task entities.Task
	if err := task.ScanTo(res.Scan); err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	return &task, nil
}

func (i *TaskStore) GetAll(ctx context.Context, qP QueryParameter) (*[]entities.Task, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("tasks/GetAll.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	res, err := i.db.Query(loadedSql, qP.Limit, qP.Offset)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	if res.Err() != nil {
		if res.Err() == sql.ErrNoRows {
			return nil, entities.NotFoundError(fmt.Errorf("no tasks found"), uri)
		}
		return nil, entities.InternalServerError(res.Err(), uri)
	}
	defer res.Close()
	tasks := []entities.Task{}
	for res.Next() {
		var task entities.Task
		if err := task.ScanTo(res.Scan); err != nil {
			return nil, entities.InternalServerError(err, uri)
		} else {
			tasks = append(tasks, task)
		}
	}
	return &tasks, nil
}

func (i *TaskStore) GetQuery(ctx context.Context, qP QueryParameter) (*[]entities.Task, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	whiteList := i.createWhitelist()
	//? maybe a good idea to not give good feedback to make it harder for sql injections?
	for k, v := range qP.Query {
		if k == "incident_id" || k == "owner_id" {
			fmt.Printf("i got %s and it should be an uuid?\n", v)
			if _, err := uuid.Parse(v); err != nil {
				return nil, entities.BadRequestError(fmt.Errorf("whitelist check failed"), uri)
			}
			continue
		}
		if !CheckWhitelist(k, v, whiteList) {
			return nil, entities.BadRequestError(fmt.Errorf("whitelist check failed"), uri)
		}
	}
	rawSql, err := LoadRawSQL("tasks/GetQuery.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	finalSql := FinalizeSQL(rawSql, "tasks", qP)
	rows, err := i.db.Query(finalSql, qP.Limit, qP.Offset)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	if rows.Err() != nil {
		if rows.Err() == sql.ErrNoRows {
			return nil, entities.NotFoundError(fmt.Errorf("no tasks found"), uri)
		}
		return nil, entities.InternalServerError(rows.Err(), uri)
	}
	defer rows.Close()
	var tasks []entities.Task
	for rows.Next() {
		var t entities.Task
		err := t.ScanTo(rows.Scan)
		if err != nil {
			return nil, entities.InternalServerError(rows.Err(), uri)
		}
		tasks = append(tasks, t)
	}
	return &tasks, nil
}

func (i *TaskStore) Create(ctx context.Context, task *entities.Task) (*entities.Task, *entities.ApiError) {
	fmt.Println("creating new task from ", task)
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("tasks/Create.sql")
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	_, err = i.db.Exec(loadedSql)
	if err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	return task, nil
}

func (i *TaskStore) Update(entity entities.Incident) (*entities.Task, *entities.ApiError) {
	panic("not implemented") // TODO: Implement
}

func (i *TaskStore) Delete(ctx context.Context, id string) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	loadedSql, err := LoadRawSQL("tasks/Delete.sql")
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	_, err = i.db.Exec(loadedSql, id)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (i *TaskStore) createWhitelist() Whitelist {
	taskWhitelist := map[string][]string{
		"Priority": {string(entities.Low), string(entities.Medium), string(entities.High), string(entities.Critical)},
		"State":    {string(entities.Backlog), string(entities.Doing), string(entities.Review), string(entities.Done)},
	}
	return taskWhitelist
}
