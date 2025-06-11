package client

import (
	"context"
	"fmt"
	"time"
	"todo-service/internal/middleware/auth"
	"todo-service/internal/models/task"

	"github.com/dinoagera/proto/gen/go/myservice"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	authClient myservice.AuthClient
	dbClient   myservice.DBWorkClient
}

func New(ctx context.Context, addrAuth string, addrDBTasks string, timeout time.Duration, retriesCount int) (*Client, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	}
	authConn, err := grpc.DialContext(ctx, addrAuth, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("while connect to authservice,%w", err)
	}
	dbConn, err := grpc.DialContext(ctx, addrDBTasks, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("while connect to dbservice,%w", err)
	}
	return &Client{
		authClient: myservice.NewAuthClient(authConn),
		dbClient:   myservice.NewDBWorkClient(dbConn),
	}, nil
}
func (c *Client) CreateTask(ctx context.Context, title string, description string) (int64, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return 0, "", fmt.Errorf("Unauthorization")
	}
	resp, err := c.dbClient.CreateTask(ctx, &myservice.CreateRequest{
		Title:       title,
		Description: description,
		Userid:      uid,
	})
	if err != nil {
		return 0, "", fmt.Errorf("create task to failed")
	}
	if resp == nil {
		return 0, "", nil
	}
	return resp.Id, resp.Message, nil
}
func (c *Client) DeleteTask(ctx context.Context, id int64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("Unauthorization")
	}
	resp, err := c.dbClient.DeleteTask(ctx, &myservice.DeleteRequest{
		Id:     id,
		Userid: uid,
	})
	if err != nil {
		return "", fmt.Errorf("delete task to failed")
	}
	if resp == nil {
		return "", nil
	}
	return resp.Message, nil
}
func (c *Client) DoneTask(ctx context.Context, id int64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("Unauthorization")
	}
	resp, err := c.dbClient.DoneTask(ctx, &myservice.DoneRequest{
		Id:     id,
		Userid: uid,
	})
	if err != nil {
		return "", fmt.Errorf("done task to failed")
	}
	if resp == nil {
		return "", nil
	}
	return resp.Message, nil
}
func (c *Client) GetAllTask(ctx context.Context) ([]*task.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Unauthorization")
	}
	resp, err := c.dbClient.GetAllTask(ctx, &myservice.GetAllRequest{
		Userid: uid,
	})
	if err != nil {
		return nil, fmt.Errorf("get all tasks to failed")
	}
	if resp == nil {
		return []*task.Task{}, nil
	}
	tasks := make([]*task.Task, 0, len(resp.Tasks))
	for _, t := range resp.Tasks {
		tasks = append(tasks, &task.Task{
			Id:          t.Id,
			Title:       t.Title,
			Description: t.Description,
			Done:        t.Done,
			Uid:         t.Uid,
		})
	}
	return tasks, nil
}
func (c *Client) Login(ctx context.Context, email string, password string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email is empty")
	}
	if password == "" {
		return "", fmt.Errorf("password is empty")
	}
	if len(password) < 6 {
		return "", fmt.Errorf("password must be at least 6 characters")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.authClient.LoginUser(ctx, &myservice.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}
	if resp == nil {
		return "", nil
	}
	return resp.Token, nil
}
func (c *Client) Registeruser(ctx context.Context, email string, password string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email is empty")
	}
	if password == "" {
		return "", fmt.Errorf("password is empty")
	}
	if len(password) < 6 {
		return "", fmt.Errorf("password must be at least 6 characters")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.authClient.RegisterUser(ctx, &myservice.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to register,err:%w", err)
	}
	return resp.Message, nil
}
