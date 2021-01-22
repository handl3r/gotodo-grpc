package v1

import (
	context "context"
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpchttp/pkg/api/v1"
	"time"
)

const (
	apiVersion = "v1"
)

type todoService struct {
	db *sql.DB
}

func (t *todoService) checkApiVersion(api string) error {
	if len(api) != 0 {
		if api != apiVersion {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implement version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func (t *todoService) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := t.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to connect to database -> "+err.Error())
	}
	return c, nil
}

func (t *todoService) Create(ctx context.Context, request *todo_service.CreateRequest) (*todo_service.CreateResponse, error) {
	if err := t.checkApiVersion(request.Api); err != nil {
		return nil, err
	}
	c, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	reminder, err := ptypes.Timestamp(request.Todo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format ->"+err.Error())
	}
	res, err := t.db.ExecContext(ctx, "INSERT INTO ToDo(`Title`, `Description`, `Reminder`) VALUES(?, ?, ?)",
		request.Todo.Title, request.Todo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert ToDo -> "+err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to get last insert id Todo -> "+err.Error())
	}
	return &todo_service.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

func (t *todoService) Read(ctx context.Context, request *todo_service.ReadRequest) (*todo_service.ReadResponse, error) {
	if err := t.checkApiVersion(request.Api); err != nil {
		return nil, err
	}
	c, err := t.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo WHERE `ID`=?", request.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select item from TODO -> "+err.Error())
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to receive data from ToDo -> "+err.Error())
		}
		return nil, status.Errorf(codes.NotFound, "Todo with ID= %s is not found", request.Id)
	}
	var todo todo_service.Todo
	var reminder time.Time
	if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &reminder ); err != nil {

	}
}

func (t todoService) ReadAll(ctx context.Context, request *todo_service.ReadAllRequest) (*todo_service.ReadAllResponse, error) {
	panic("implement me")
}

func (t todoService) Update(ctx context.Context, request *todo_service.UpdateRequest) (*todo_service.UpdateResponse, error) {
	panic("implement me")
}

func (t todoService) Delete(ctx context.Context, request *todo_service.DeleteRequest) (*todo_service.DeleteResponse, error) {
	panic("implement me")
}

func NewToDoService(db *sql.DB) *todoService {
	return &todoService{
		db: db,
	}
}
