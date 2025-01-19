package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/akamiya208/go-grpc-tutrial/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"

	pb "github.com/akamiya208/go-grpc-tutrial/internal/pkg/proto"
)

type MockedMySQLClient struct {
	mock.Mock
}

func (m *MockedMySQLClient) GetTask(taskID uint) (models.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockedMySQLClient) GetTasksByName(name string) ([]models.Task, error) {
	args := m.Called(name)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockedMySQLClient) CreateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockedMySQLClient) UpdateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockedMySQLClient) DeleteTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockedMySQLClient) DB() *gorm.DB {
	return nil
}

func TestGetTask(t *testing.T) {
	description := "description"
	testTask := models.Task{
		ID:          1,
		Name:        "task1",
		Description: &description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   gorm.DeletedAt{},
	}

	t.Run("ok", func(t *testing.T) {
		// Mock MySQLClient
		taskID := uint(1)
		mockedMySQLClient := new(MockedMySQLClient)
		mockedMySQLClient.On("GetTask", taskID).Return(testTask, nil)

		// Create a new server
		lis := bufconn.Listen(1024 * 1024)
		s := grpc.NewServer()
		pb.RegisterTaskServiceServer(s, NewTaskServer(mockedMySQLClient))
		go func() {
			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		}()
		defer s.Stop()

		// Create a new client
		conn, err := grpc.NewClient("passthrough://"+lis.Addr().String(),
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		client := pb.NewTaskServiceClient(conn)
		res, _ := client.GetTask(context.Background(), &pb.GetTaskRequest{Id: uint32(taskID)})

		// Assert
		assert.Equal(t, "task1", res.Name)
	})

	t.Run("not found", func(t *testing.T) {
		// Mock MySQLClient
		taskID := uint(1)
		mockedMySQLClient := new(MockedMySQLClient)
		mockedMySQLClient.On("GetTask", taskID).Return(models.Task{}, fmt.Errorf("record not found"))

		// Create a new server
		lis := bufconn.Listen(1024 * 1024)
		s := grpc.NewServer()
		pb.RegisterTaskServiceServer(s, NewTaskServer(mockedMySQLClient))
		go func() {
			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		}()
		defer s.Stop()

		// Create a new client
		conn, err := grpc.NewClient("passthrough://"+lis.Addr().String(),
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		client := pb.NewTaskServiceClient(conn)
		_, err = client.GetTask(context.Background(), &pb.GetTaskRequest{Id: uint32(taskID)})

		// Assert
		assert.Contains(t, err.Error(), "code = NotFound")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Mock MySQLClient
		taskID := uint(1)
		mockedMySQLClient := new(MockedMySQLClient)
		mockedMySQLClient.On("GetTask", taskID).Return(models.Task{}, fmt.Errorf("some error"))

		// Create a new server
		lis := bufconn.Listen(1024 * 1024)
		s := grpc.NewServer()
		pb.RegisterTaskServiceServer(s, NewTaskServer(mockedMySQLClient))
		go func() {
			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		}()
		defer s.Stop()

		// Create a new client
		conn, err := grpc.NewClient("passthrough://"+lis.Addr().String(),
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		client := pb.NewTaskServiceClient(conn)
		_, err = client.GetTask(context.Background(), &pb.GetTaskRequest{Id: uint32(taskID)})

		// Assert
		assert.Contains(t, err.Error(), "code = Internal")
	})
}
