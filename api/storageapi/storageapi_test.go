package storageapi_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/spacelift-io/homework-object-storage/api/storageapi"
	"github.com/spacelift-io/homework-object-storage/storage"
)

func isDataEqual(t *testing.T, expected, got []byte) {
	if len(expected) != len(got) {
		t.Fatalf("slice length is not the same; expected(%d); got(%d)", len(got), len(expected))
	}

	for i := range expected {
		if expected[i] != got[i] {
			t.Fatalf("byte expected: %c, got %c", expected[i], got[i])
		}
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("err")
}
func (e *errorReader) Close() error {
	return nil
}

type ctxKey int

var testCtxKey ctxKey

func TestGet(t *testing.T) {
	log.SetOutput(io.Discard)
	tests := []struct {
		name               string
		id                 string
		storage            storage.Storage
		expectedStatusCode int
		expctedBody        []byte
		validContentType   bool
	}{
		{
			name:               "storage error",
			id:                 "id",
			expectedStatusCode: http.StatusNotFound,
			storage: &storage.StorageMock{
				GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
					return nil, fmt.Errorf("err")
				},
			},
		},
		{
			name:               "reader error",
			id:                 "id",
			expectedStatusCode: http.StatusNotFound,
			storage: &storage.StorageMock{
				GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
					return &errorReader{}, nil
				},
			},
		},
		{
			name:               "object returned /w valid content-type",
			id:                 "id-03",
			expectedStatusCode: http.StatusOK,
			expctedBody:        []byte{'b', 'i', 'e', 'n'},
			validContentType:   true,
			storage: &storage.StorageMock{
				GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
					expectedID := "id-03"
					if id != expectedID {
						t := ctx.Value(testCtxKey).(*testing.T)
						t.Fatalf("storage got wrong id: got(%s); expected(%s)", id, expectedID)
					}

					data := bytes.NewReader([]byte("bien"))
					return io.NopCloser(data), nil
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/object/%s", test.id), nil)

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", test.id)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))

			request = request.WithContext(context.WithValue(request.Context(), testCtxKey, t))

			storageapi.Get(test.storage)(recorder, request)
			if recorder.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("expected status code (%d); got(%d)", test.expectedStatusCode, recorder.Result().StatusCode)
			}

			if test.validContentType {
				if recorder.HeaderMap.Get("Content-Type") != "application/octet-stream" {
					t.Errorf("invalid content-type: got(%s); expected(%s)", recorder.HeaderMap.Get("Content-Type"), "application/octet-stream")
				}
			}

			isDataEqual(t, test.expctedBody, recorder.Body.Bytes())
		})
	}
}

func TestPut(t *testing.T) {
	log.SetOutput(io.Discard)
	tests := []struct {
		name               string
		id                 string
		storage            storage.Storage
		Body               []byte
		expectedStatusCode int
		invalidContentType bool
	}{
		{
			name:               "invalid content type",
			id:                 "id",
			expectedStatusCode: http.StatusBadRequest,
			invalidContentType: true,
		},
		{
			name:               "object placed",
			id:                 "id-03",
			expectedStatusCode: http.StatusCreated,
			Body:               []byte{'b', 'i', 'e', 'n'},
			storage: &storage.StorageMock{
				PutFn: func(ctx context.Context, id string, reader io.Reader, length int64) error {
					t := ctx.Value(testCtxKey).(*testing.T)
					expectedID := "id-03"
					if id != expectedID {
						t.Fatalf("storage got wrong id: got(%s); expected(%s)", id, expectedID)
					}

					if length != 4 {
						t.Fatalf("expected length 4, got %d", length)
					}

					data, _ := io.ReadAll(reader)
					isDataEqual(t, []byte{'b', 'i', 'e', 'n'}, data)

					return nil
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := bytes.NewReader([]byte{})
			if test.Body != nil {
				data = bytes.NewReader(test.Body)
			}
			request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/object/%s", test.id), data)

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", test.id)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))

			request = request.WithContext(context.WithValue(request.Context(), testCtxKey, t))

			if !test.invalidContentType {
				request.Header.Set("Content-Type", "application/octet-stream")
			}

			recorder := httptest.NewRecorder()
			storageapi.Put(test.storage)(recorder, request)
			if recorder.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("expected status code (%d); got(%d)", test.expectedStatusCode, recorder.Result().StatusCode)
			}
		})
	}
}
