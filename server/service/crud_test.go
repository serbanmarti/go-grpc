package service

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"

	crudv1 "github.com/serbanmarti/go-grpc/proto_gen/crud/v1"
	"github.com/serbanmarti/go-grpc/proto_gen/crud/v1/crudv1connect"
)

func TestCrudService_Create(t *testing.T) {
	tests := []struct {
		name    string
		reqData crudv1.CreateRequest
		resData crudv1.CreateResponse
	}{
		{
			name: "Test create record",
			reqData: crudv1.CreateRequest{
				Name: "Test Record",
			},
			resData: crudv1.CreateResponse{},
		},
	}

	client := crudv1connect.NewCrudServiceClient(
		newInsecureClient(),
		"http://127.0.0.1:8080",
		connect.WithGRPC(),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			res, err := client.Create(ctx, connect.NewRequest(&tt.reqData))
			assert.NoError(t, err)

			if !cmp.Equal(
				tt.resData, *res.Msg,
				cmpopts.IgnoreUnexported(crudv1.CreateResponse{}),
				cmpopts.IgnoreFields(crudv1.CreateResponse{}, "Id"),
			) {
				t.Errorf("want[-], got[+]\n%v", cmp.Diff(
					tt.resData, *res.Msg,
					cmpopts.IgnoreUnexported(crudv1.CreateResponse{}),
					cmpopts.IgnoreFields(crudv1.CreateResponse{}, "Id"),
				))
			}

			if _, err := ksuid.Parse(res.Msg.Id); err != nil {
				t.Errorf("ID is not a valid KSUID: %v", err)
			}
		})
	}
}

func TestCrudService_Read(t *testing.T) {
	tests := []struct {
		name        string
		reqData     crudv1.ReadRequest
		resData     crudv1.ReadResponse
		expectedErr error
	}{
		{
			name: "Test read record one",
			reqData: crudv1.ReadRequest{
				Id: "2imgNBCejbjXehOazVerssNsgcz",
			},
			resData: crudv1.ReadResponse{
				Id:   "2imgNBCejbjXehOazVerssNsgcz",
				Name: "Test Record 1",
			},
		},
		{
			name: "Test read record two",
			reqData: crudv1.ReadRequest{
				Id: "2imgN7lkpYjE16akMMn52Uvkgln",
			},
			resData: crudv1.ReadResponse{
				Id:   "2imgN7lkpYjE16akMMn52Uvkgln",
				Name: "Test Record 2",
			},
		},
		{
			name: "Test read non-existent record",
			reqData: crudv1.ReadRequest{
				Id: "2imgqwcM6MabAQBULm8VtXvfF86",
			},
			expectedErr: connect.NewError(connect.CodeNotFound, fmt.Errorf("record not found")),
		},
	}

	client := crudv1connect.NewCrudServiceClient(
		newInsecureClient(),
		"http://127.0.0.1:8080",
		connect.WithGRPC(),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			res, err := client.Read(ctx, connect.NewRequest(&tt.reqData))

			if tt.expectedErr == nil {
				assert.NoError(t, err)

				if !cmp.Equal(
					tt.resData, *res.Msg,
					cmpopts.IgnoreUnexported(crudv1.ReadResponse{}),
				) {
					t.Errorf("want[-], got[+]\n%v", cmp.Diff(
						tt.resData, *res.Msg,
						cmpopts.IgnoreUnexported(crudv1.ReadResponse{}),
					))
				}

				if _, err := ksuid.Parse(res.Msg.Id); err != nil {
					t.Errorf("ID is not a valid KSUID: %v", err)
				}
			} else {
				if !cmp.Equal(
					tt.expectedErr.Error(), err.Error(),
				) {
					t.Errorf("want[-], got[+]\n%v", cmp.Diff(
						tt.expectedErr.Error(), err.Error(),
					))
				}
			}
		})
	}
}

func TestCrudService_Update(t *testing.T) {
	tests := []struct {
		name        string
		reqData     crudv1.UpdateRequest
		resData     crudv1.UpdateResponse
		expectedErr error
	}{
		{
			name: "Test update record one",
			reqData: crudv1.UpdateRequest{
				Id:          "2imgNBCejbjXehOazVerssNsgcz",
				UpdatedName: "Test Record 1 - updated",
			},
			resData: crudv1.UpdateResponse{
				Id:   "2imgNBCejbjXehOazVerssNsgcz",
				Name: "Test Record 1 - updated",
			},
		},
		{
			name: "Test update non-existent record",
			reqData: crudv1.UpdateRequest{
				Id: "2imgqwcM6MabAQBULm8VtXvfF86",
			},
			expectedErr: connect.NewError(connect.CodeNotFound, fmt.Errorf("record not found")),
		},
	}

	client := crudv1connect.NewCrudServiceClient(
		newInsecureClient(),
		"http://127.0.0.1:8080",
		connect.WithGRPC(),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			res, err := client.Update(ctx, connect.NewRequest(&tt.reqData))

			if tt.expectedErr == nil {
				assert.NoError(t, err)

				if !cmp.Equal(
					tt.resData, *res.Msg,
					cmpopts.IgnoreUnexported(crudv1.UpdateResponse{}),
				) {
					t.Errorf("want[-], got[+]\n%v", cmp.Diff(
						tt.resData, *res.Msg,
						cmpopts.IgnoreUnexported(crudv1.UpdateResponse{}),
					))
				}

				if _, err := ksuid.Parse(res.Msg.Id); err != nil {
					t.Errorf("ID is not a valid KSUID: %v", err)
				}
			} else {
				if !cmp.Equal(
					tt.expectedErr.Error(), err.Error(),
				) {
					t.Errorf("want[-], got[+]\n%v", cmp.Diff(
						tt.expectedErr.Error(), err.Error(),
					))
				}
			}
		})
	}
}

func TestCrudService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		reqData     crudv1.DeleteRequest
		resData     crudv1.DeleteResponse
		expectedErr error
	}{
		{
			name: "Test delete record one",
			reqData: crudv1.DeleteRequest{
				Id: "2imgNBCejbjXehOazVerssNsgcz",
			},
			resData: crudv1.DeleteResponse{
				Id: "2imgNBCejbjXehOazVerssNsgcz",
			},
		},
		{
			name: "Test delete non-existent record",
			reqData: crudv1.DeleteRequest{
				Id: "2imgqwcM6MabAQBULm8VtXvfF86",
			},
			expectedErr: connect.NewError(connect.CodeNotFound, fmt.Errorf("record not found")),
		},
	}

	client := crudv1connect.NewCrudServiceClient(
		newInsecureClient(),
		"http://127.0.0.1:8080",
		connect.WithGRPC(),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			res, err := client.Delete(ctx, connect.NewRequest(&tt.reqData))

			if tt.expectedErr == nil {
				assert.NoError(t, err)

				if !cmp.Equal(
					tt.resData, *res.Msg,
					cmpopts.IgnoreUnexported(crudv1.DeleteResponse{}),
				) {
					t.Errorf("want[-], got[+]\n%v", cmp.Diff(
						tt.resData, *res.Msg,
						cmpopts.IgnoreUnexported(crudv1.DeleteResponse{}),
					))
				}

				if _, err := ksuid.Parse(res.Msg.Id); err != nil {
					t.Errorf("ID is not a valid KSUID: %v", err)
				}
			} else {
				if !cmp.Equal(
					tt.expectedErr.Error(), err.Error(),
				) {
					t.Errorf("want[-], got[+]\n%v", cmp.Diff(
						tt.expectedErr.Error(), err.Error(),
					))
				}
			}
		})
	}
}
