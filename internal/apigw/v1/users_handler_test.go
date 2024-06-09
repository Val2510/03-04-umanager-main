package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockUsersClient struct {
	mock.Mock
}

func (m *mockUsersClient) ListUsers(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.ListUsersResponse, error) {
	panic("unimplemented")
}

func (m *mockUsersClient) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	panic("unimplemented")
}

func (m *mockUsersClient) CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.Empty), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUsersClient) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.Empty), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUsersClient) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.User, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestPostUsers(t *testing.T) {
	client := new(mockUsersClient)
	handler := newUsersHandler(client)

	t.Run("success", func(t *testing.T) {
		user := apiv1.UserCreate{
			Id:       "123",
			Username: "testuser",
			Password: "testpass",
		}

		userJSON, err := json.Marshal(user)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		client.On("CreateUser", mock.Anything, &pb.CreateUserRequest{
			Id:       user.Id,
			Username: user.Username,
			Password: user.Password,
		}).Return(&pb.Empty{}, nil)

		handler.PostUsers(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)
		client.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		user := apiv1.UserCreate{
			Id:       "123",
			Username: "testuser",
			Password: "testpass",
		}

		userJSON, err := json.Marshal(user)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		grpcErr := status.Error(codes.Internal, "create user error")
		client.On("CreateUser", mock.Anything, &pb.CreateUserRequest{
			Id:       user.Id,
			Username: user.Username,
			Password: user.Password,
		}).Return(nil, grpcErr)

		handler.PostUsers(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		client.AssertExpectations(t)
	})
}

func TestDeleteUsersId(t *testing.T) {
	client := new(mockUsersClient)
	handler := newUsersHandler(client)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/123", nil)
		w := httptest.NewRecorder()

		client.On("DeleteUser", mock.Anything, &pb.DeleteUserRequest{Id: "123"}).Return(&pb.Empty{}, nil)

		handler.DeleteUsersId(w, req, "123")

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		client.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/123", nil)
		w := httptest.NewRecorder()

		grpcErr := status.Error(codes.Internal, "delete user error")
		client.On("DeleteUser", mock.Anything, &pb.DeleteUserRequest{Id: "123"}).Return(nil, grpcErr)

		handler.DeleteUsersId(w, req, "123")

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		client.AssertExpectations(t)
	})
}
