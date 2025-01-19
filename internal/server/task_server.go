package server

import (
	"context"

	"github.com/akamiya208/go-grpc-tutrial/internal/pkg/models"
	"github.com/akamiya208/go-grpc-tutrial/internal/pkg/mysql"
	pb "github.com/akamiya208/go-grpc-tutrial/internal/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	mysqlClient mysql.IClient
}

func NewTaskServer(client mysql.IClient) *TaskServer {
	return &TaskServer{mysqlClient: client}
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	task, err := s.mysqlClient.GetTask(uint(req.Id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.TaskResponse{
		Id:          uint32(task.ID),
		Name:        task.Name,
		Description: *task.Description,
	}, nil
}

func (s *TaskServer) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.TaskResponses, error) {
	name := req.Name
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	tasks, err := s.mysqlClient.GetTasksByName(name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responses := make([]*pb.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = &pb.TaskResponse{
			Id:          uint32(task.ID),
			Name:        task.Name,
			Description: *task.Description,
		}
	}

	return &pb.TaskResponses{
		Tasks: responses,
	}, nil
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	task := models.Task{Name: req.Name, Description: &req.Description}
	if err := s.mysqlClient.CreateTask(&task); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TaskResponse{
		Id:          uint32(task.ID),
		Name:        task.Name,
		Description: *task.Description,
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	task, err := s.mysqlClient.GetTask(uint(req.Id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	task.Name = req.Name
	task.Description = &req.Description

	if err := s.mysqlClient.UpdateTask(&task); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TaskResponse{
		Id:          uint32(task.ID),
		Name:        task.Name,
		Description: *task.Description,
	}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.TaskResponse, error) {
	task, err := s.mysqlClient.GetTask(uint(req.Id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if err := s.mysqlClient.DeleteTask(&task); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TaskResponse{
		Id:          uint32(task.ID),
		Name:        task.Name,
		Description: *task.Description,
	}, nil
}
