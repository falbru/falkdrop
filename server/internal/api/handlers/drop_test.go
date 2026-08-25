package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httperror "github.com/falbru/falkdrop/internal/api/errors"
	"github.com/falbru/falkdrop/internal/api/handlers"
	"github.com/falbru/falkdrop/internal/app/auth"
	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/falbru/falkdrop/internal/app/drop/mock"
	"github.com/google/uuid"
)

func TestGetDrop(t *testing.T) {
	t.Run("no dropId is specified", func(t *testing.T) {
		dropService := mock.NewMockDropService()
		dropHandler := handlers.NewDropHandler(dropService)

		req := httptest.NewRequest(http.MethodGet, "/drop/", nil)
		w := httptest.NewRecorder()

		err := dropHandler.Get(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error when no dropId was specified")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror when no dropId was specified")
			return
		}

		if httpError.Code != http.StatusBadRequest {
			t.Errorf("Expected drop handler to throw %v httperror", http.StatusBadRequest)
		}
	})

	t.Run("drop with dropId doesn't exist", func(t *testing.T) {
		dropService := mock.NewMockDropService().WithGetDropWithResourceDownloadUrls(func(ctx context.Context, dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error) {
			return nil, drop.ErrDropNotFound{DropId: dropId}
		})

		dropHandler := handlers.NewDropHandler(dropService)

		req := httptest.NewRequest(http.MethodGet, "/drop/00000", nil)
		req.SetPathValue("dropId", "00000")
		w := httptest.NewRecorder()

		err := dropHandler.Get(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error with non-existent drop")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror with non-existent drop")
			return
		}

		if httpError.Code != http.StatusNotFound {
			t.Errorf("Expected drop handler to throw 404 httperror, but got %v with message: %v", httpError.Code, httpError.Error())
		}
	})

	t.Run("success", func(t *testing.T) {
		dropService := mock.NewMockDropService().WithGetDropWithResourceDownloadUrls(func(ctx context.Context, dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error) {
			expirationDate, _ := time.Parse("2006-Jan-02", "2030-Jan-01")
			return &drop.DropWithResourceDownloadUrls{
				Id:             "12345",
				ExpirationDate: expirationDate,
				Resources: []drop.ResourceWithDownloadUrl{
					drop.ResourceWithDownloadUrl{
						Resource: drop.Resource{
							Id:     drop.ResourceId(uuid.MustParse("11111111-1111-1111-1111-111111111111")),
							Type:   drop.FileResource,
							Name:   nil,
							DropId: nil,
						},
						DownloadUrl: "http://example.org",
					},
				},
			}, nil
		})

		dropHandler := handlers.NewDropHandler(dropService)

		req := httptest.NewRequest(http.MethodGet, "/drop/{dropId}", nil)
		req.SetPathValue("dropId", "12345")
		w := httptest.NewRecorder()

		err := dropHandler.Get(w, req)
		if err != nil {
			t.Errorf("Expected drop not to throw error, but got: %v", err.Error())
			return
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
		}

		var response handlers.DropWithResourceDownloadUrlsDTO
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Id != "12345" {
			t.Errorf("Expected drop id '12345', but got '%v'", response.Id)
		}
		if response.ExpirationDate == "2030-Jan-01" {
			t.Errorf("Expected expiration date to be set")
		}
		if len(response.Resources) != 1 {
			t.Errorf("Expected 1 resource, but got %v", len(response.Resources))
		}
		if len(response.Resources) > 0 {
			if response.Resources[0].Id != "11111111-1111-1111-1111-111111111111" {
				t.Errorf("Expected resource id '11111111-1111-1111-1111-111111111111', but got '%v'", response.Resources[0].Id)
			}
			if response.Resources[0].Type != string(drop.FileResource) {
				t.Errorf("Expected resource type 'file', but got '%v'", response.Resources[0].Type)
			}
			if response.Resources[0].DownloadUrl != "http://example.org" {
				t.Errorf("Expected download url 'http://example.org', but got '%v'", response.Resources[0].DownloadUrl)
			}
		}
	})
}

func TestCreateDrop(t *testing.T) {
	t.Run("bearer token not provided", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.New())
		dropService := mock.NewMockDropService()
		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": ["` + resourceId.String() + `"], "expiry_duration": "P1D"}`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror")
			return
		}

		if httpError.Code != http.StatusUnauthorized {
			t.Errorf("Expected drop handler to throw %v httperror, but got %v", http.StatusUnauthorized, httpError.Code)
		}
	})

	t.Run("unauthorized user", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.New())
		dropService := mock.NewMockDropService().WithCreateDrop(func(ctx context.Context, resourceIds []drop.ResourceId, expiryDuration time.Duration) (*drop.DropWithResourceDownloadUrls, error) {
			return nil, auth.ErrUnauthorized
		})

		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": ["` + resourceId.String() + `"], "expiry_duration": "P1D"}`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror")
			return
		}

		if httpError.Code != http.StatusUnauthorized {
			t.Errorf("Expected drop handler to throw %v httperror, but got %v", http.StatusUnauthorized, httpError.Code)
		}
	})

	t.Run("invalid drop request body", func(t *testing.T) {
		dropService := mock.NewMockDropService()
		dropHandler := handlers.NewDropHandler(dropService)

		req := httptest.NewRequest(http.MethodPost, "/drop", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error with invalid request body")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror with invalid request body")
			return
		}

		if httpError.Code != http.StatusBadRequest {
			t.Errorf("Expected drop handler to throw %v httperror with invalid request body, but got %v", http.StatusBadRequest, httpError.Code)
		}
	})

	t.Run("no resources attached", func(t *testing.T) {
		dropService := mock.NewMockDropService()
		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": [], "expiry_duration": "P1D" }`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error when no resources are attached")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror when no resources are attached")
			return
		}

		if httpError.Code != http.StatusBadRequest {
			t.Errorf("Expected drop handler to throw %v httperror when no resources are attached, but got %v", http.StatusBadRequest, httpError.Code)
		}

		if !strings.Contains(httpError.Error(), "resourceIds is empty") {
			t.Errorf("Expected error message to contain 'resourceIds is empty', but got: %v", httpError.Error())
		}
	})

	t.Run("empty expiry duration", func(t *testing.T) {
		dropService := mock.NewMockDropService()
		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": ["11111111-1111-1111-1111-111111111111"], "expiry_duration": ""}`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror")
			return
		}

		if httpError.Code != http.StatusBadRequest {
			t.Errorf("Expected drop handler to throw %v httperror, but got %v", http.StatusBadRequest, httpError.Code)
		}
	})

	t.Run("invalid expiry duration", func(t *testing.T) {
		dropService := mock.NewMockDropService()
		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": ["11111111-1111-1111-1111-111111111111"], "expiry_duration": "INVALID"}`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror")
			return
		}

		if httpError.Code != http.StatusBadRequest {
			t.Errorf("Expected drop handler to throw %v httperror, but got %v", http.StatusBadRequest, httpError.Code)
		}
	})

	t.Run("one or more resources do not exist", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.New())
		dropService := mock.NewMockDropService().WithCreateDrop(func(ctx context.Context, resourceIds []drop.ResourceId, expiryDuration time.Duration) (*drop.DropWithResourceDownloadUrls, error) {
			return nil, drop.ErrResourcesNotFound{ResourceIds: []drop.ResourceId{resourceId}}
		})
		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": ["` + resourceId.String() + `"], "expiry_duration": "P1D"}`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error when resources do not exist")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror when resources do not exist")
			return
		}

		if httpError.Code != http.StatusNotFound {
			t.Errorf("Expected drop handler to throw %v httperror when resources do not exist, but got %v", http.StatusNotFound, httpError.Code)
		}
	})

	t.Run("one or more resource already belongs to another drop", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.New())
		existingDropId := drop.DropId("existing-drop-id")
		dropService := mock.NewMockDropService().WithCreateDrop(func(ctx context.Context, resourceIds []drop.ResourceId, expiryDuration time.Duration) (*drop.DropWithResourceDownloadUrls, error) {
			return nil, drop.ErrResourceAlreadyBelongsToDrop{ResourceId: resourceId, DropId: existingDropId}
		})
		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": ["` + resourceId.String() + `"]}`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := dropHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected drop handler to throw error when resource already belongs to another drop")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected drop handler to throw httperror when resource already belongs to another drop")
			return
		}

		if httpError.Code != http.StatusBadRequest {
			t.Errorf("Expected drop handler to throw %v httperror when resource already belongs to another drop, but got %v", http.StatusBadRequest, httpError.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
		expirationDate, err := time.Parse("2006-Jan-02", "2030-Jan-01")
		resourceName := "myfile.txt"
		resourceDownloadUrl := "http://example.org/download"

		dropService := mock.NewMockDropService().WithCreateDrop(func(ctx context.Context, resourceIds []drop.ResourceId, expiryDuration time.Duration) (*drop.DropWithResourceDownloadUrls, error) {
			return &drop.DropWithResourceDownloadUrls{
				Id:             "12345",
				ExpirationDate: expirationDate,
				Resources: []drop.ResourceWithDownloadUrl{
					{
						Resource: drop.Resource{
							Id:     resourceId,
							Type:   drop.FileResource,
							Name:   &resourceName,
							DropId: nil,
						},
						DownloadUrl: resourceDownloadUrl,
					},
				},
			}, nil
		})
		dropHandler := handlers.NewDropHandler(dropService)

		reqBody := `{"resource_ids": ["11111111-1111-1111-1111-111111111111"], "expiry_duration": "P1D"}`
		req := httptest.NewRequest(http.MethodPost, "/drop", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err = dropHandler.Create(w, req)
		if err != nil {
			t.Errorf("Expected drop handler not to throw error on success, but got: %v", err.Error())
			return
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
		}

		var response handlers.DropWithResourceDownloadUrlsDTO
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Id != "12345" {
			t.Errorf("Expected drop id '12345', but got '%v'", response.Id)
		}
		if response.ExpirationDate != expirationDate.Format(time.RFC3339) {
			t.Errorf("Expected expiration date '%v', but got '%v'", expirationDate.Format(time.RFC3339), response.ExpirationDate)
		}
		if len(response.Resources) != 1 {
			t.Errorf("Expected 1 resource, but got %v", len(response.Resources))
		}
		if len(response.Resources) > 0 {
			if response.Resources[0].Id != "11111111-1111-1111-1111-111111111111" {
				t.Errorf("Expected resource id '11111111-1111-1111-1111-111111111111', but got '%v'", response.Resources[0].Id)
			}
			if response.Resources[0].Type != string(drop.FileResource) {
				t.Errorf("Expected resource type 'file', but got '%v'", response.Resources[0].Type)
			}
			if response.Resources[0].Name == nil || *response.Resources[0].Name != "myfile.txt" {
				t.Errorf("Expected resource name 'myfile.txt', but got '%v'", response.Resources[0].Name)
			}
			if response.Resources[0].DownloadUrl != "http://example.org/download" {
				t.Errorf("Expected download url 'http://example.org/download', but got '%v'", response.Resources[0].DownloadUrl)
			}
		}
	})
}
