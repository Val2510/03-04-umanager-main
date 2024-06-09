package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "gitlab.com/robotomize/gb-golang/homework/03-04-umanager/internal/apigw/v1"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockLinksClient struct {
	mock.Mock
}

func (m *mockLinksClient) CreateLink(ctx context.Context, in *pb.CreateLinkRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.Empty), args.Error(1)
	}
	return &pb.Empty{}, args.Error(1)
}

func (m *mockLinksClient) DeleteLink(ctx context.Context, in *pb.DeleteLinkRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.Empty), args.Error(1)
	}
	return &pb.Empty{}, args.Error(1)
}

func (m *mockLinksClient) GetLink(ctx context.Context, in *pb.GetLinkRequest, opts ...grpc.CallOption) (*pb.Link, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.Link), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockLinksClient) GetLinkByUserID(ctx context.Context, in *pb.GetLinksByUserId, opts ...grpc.CallOption) (*pb.ListLinkResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.ListLinkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockLinksClient) UpdateLink(ctx context.Context, in *pb.UpdateLinkRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.Empty), args.Error(1)
	}
	return &pb.Empty{}, args.Error(1)
}

func (m *mockLinksClient) ListLinks(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.ListLinkResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.ListLinkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestPostLinks(t *testing.T) {
	client := new(mockLinksClient)
	handler := v1.NewLinksHandler(client)

	t.Run("success", func(t *testing.T) {
		link := apiv1.LinkCreate{
			Id:     "123",
			Title:  "Test Link",
			Url:    "http://example.com",
			Images: []string{"http://example.com/image1.jpg"},
			Tags:   []string{"test"},
			UserId: "user1",
		}

		linkJSON, err := json.Marshal(link)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewReader(linkJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		client.On("CreateLink", mock.Anything, &pb.CreateLinkRequest{
			Id:     link.Id,
			Title:  link.Title,
			Url:    link.Url,
			Images: link.Images,
			Tags:   link.Tags,
			UserId: link.UserId,
		}).Return(&pb.Empty{}, nil)

		handler.PostLinks(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)
		client.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		link := apiv1.LinkCreate{
			Id:     "123",
			Title:  "Test Link",
			Url:    "http://example.com",
			Images: []string{"http://example.com/image1.jpg"},
			Tags:   []string{"test"},
			UserId: "user1",
		}

		linkJSON, err := json.Marshal(link)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewReader(linkJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		grpcErr := status.Error(codes.Internal, "create link error")
		client.On("CreateLink", mock.Anything, &pb.CreateLinkRequest{
			Id:     link.Id,
			Title:  link.Title,
			Url:    link.Url,
			Images: link.Images,
			Tags:   link.Tags,
			UserId: link.UserId,
		}).Return(nil, grpcErr)

		handler.PostLinks(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		client.AssertExpectations(t)
	})
}

func TestDeleteLinksId(t *testing.T) {
	client := new(mockLinksClient)
	handler := v1.NewLinksHandler(client)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/links/123", nil)
		w := httptest.NewRecorder()

		client.On("DeleteLink", mock.Anything, &pb.DeleteLinkRequest{Id: "123"}).Return(&pb.Empty{}, nil)

		handler.DeleteLinksId(w, req, "123")

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		client.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/links/123", nil)
		w := httptest.NewRecorder()

		grpcErr := status.Error(codes.Internal, "delete link error")
		client.On("DeleteLink", mock.Anything, &pb.DeleteLinkRequest{Id: "123"}).Return(nil, grpcErr)

		handler.DeleteLinksId(w, req, "123")

		resp := w.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		client.AssertExpectations(t)
	})
}
