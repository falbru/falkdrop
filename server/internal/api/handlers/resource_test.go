package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httperror "github.com/falbru/falkdrop/internal/api/errors"
	"github.com/falbru/falkdrop/internal/api/handlers"
	"github.com/falbru/falkdrop/internal/app/auth"
	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/falbru/falkdrop/internal/app/drop/mock"
	"github.com/google/uuid"
)

func TestCreateResource(t *testing.T) {
	t.Run("bearer token not provided", func(t *testing.T) {
		dropService := mock.NewMockDropService()
		resourceHandler := handlers.NewResourceHandler(dropService)

		reqBody := `{"resource_type": "file", "name": "myfile.txt"}`
		req := httptest.NewRequest(http.MethodPost, "/resource", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		err := resourceHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected resource handler to throw error, but got: %v", err.Error())
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected resource handler to throw httperror, but got %v", err.Error())
			return
		}

		if httpError.Code != http.StatusUnauthorized {
			t.Errorf("Expected resource handler to throw %v httperror, but got %v", http.StatusBadRequest, httpError.Code)
		}
	})

	t.Run("unauthorized user", func(t *testing.T) {
		dropService := mock.NewMockDropService().WithCreateResourceWithUploadUrl(func(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error) {
			return nil, auth.ErrUnauthorized
		})
		resourceHandler := handlers.NewResourceHandler(dropService)

		reqBody := `{"resource_type": "file", "name": "myfile.txt"}`
		req := httptest.NewRequest(http.MethodPost, "/resource", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := resourceHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected resource handler to throw error, but got: %v", err.Error())
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected resource handler to throw httperror, but got %v", err.Error())
			return
		}

		if httpError.Code != http.StatusUnauthorized {
			t.Errorf("Expected resource handler to throw %v httperror, but got %v", http.StatusBadRequest, httpError.Code)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		dropService := mock.NewMockDropService()
		resourceHandler := handlers.NewResourceHandler(dropService)

		req := httptest.NewRequest(http.MethodPost, "/resource", nil)
		w := httptest.NewRecorder()

		err := resourceHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected resource handler to throw error with invalid request body")
			return
		}

		var httpError httperror.HTTPError
		if !errors.As(err, &httpError) {
			t.Errorf("Expected resource handler to throw httperror with invalid request body")
			return
		}

		if httpError.Code != http.StatusBadRequest {
			t.Errorf("Expected resource handler to throw %v httperror with invalid request body, but got %v", http.StatusBadRequest, httpError.Code)
		}
	})

	t.Run("service error on create resource", func(t *testing.T) {
		expectedError := errors.New("failed to create resource in storage")
		dropService := mock.NewMockDropService().WithCreateResourceWithUploadUrl(func(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error) {
			return nil, expectedError
		})
		resourceHandler := handlers.NewResourceHandler(dropService)

		reqBody := `{"resource_type": "file"}`
		req := httptest.NewRequest(http.MethodPost, "/resource", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		err := resourceHandler.Create(w, req)
		if err == nil {
			t.Errorf("Expected resource handler to throw error when service fails")
			return
		}
	})

	t.Run("success with file resource type", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
		resourceName := "myfile.txt"
		uploadUrl := "http://example.org/upload"

		dropService := mock.NewMockDropService().WithCreateResourceWithUploadUrl(func(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error) {
			return &drop.ResourceWithUploadUrl{
				Resource: drop.Resource{
					Id:   resourceId,
					Type: drop.FileResource,
					Name: &resourceName,
				},
				UploadUrl: uploadUrl,
			}, nil
		})
		resourceHandler := handlers.NewResourceHandler(dropService)

		reqBody := `{"resource_type": "file", "name": "myfile.txt"}`
		req := httptest.NewRequest(http.MethodPost, "/resource", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := resourceHandler.Create(w, req)
		if err != nil {
			t.Errorf("Expected resource handler not to throw error on success, but got: %v", err.Error())
			return
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
		}

		var response handlers.ResourceWithUploadUrlDTO
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Id != "11111111-1111-1111-1111-111111111111" {
			t.Errorf("Expected resource id '11111111-1111-1111-1111-111111111111', but got '%v'", response.Id)
		}
		if response.Type != string(drop.FileResource) {
			t.Errorf("Expected resource type '%v', but got '%v'", drop.FileResource, response.Type)
		}
		if response.Name == nil || *response.Name != "myfile.txt" {
			t.Errorf("Expected resource name 'myfile.txt', but got '%v'", response.Name)
		}
		if response.UploadUrl != "http://example.org/upload" {
			t.Errorf("Expected upload url 'http://example.org/upload', but got '%v'", response.UploadUrl)
		}
	})

	t.Run("success with text resource type", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.MustParse("22222222-2222-2222-2222-222222222222"))
		uploadUrl := "http://example.org/upload"

		dropService := mock.NewMockDropService().WithCreateResourceWithUploadUrl(func(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error) {
			return &drop.ResourceWithUploadUrl{
				Resource: drop.Resource{
					Id:   resourceId,
					Type: drop.TextResource,
					Name: nil,
				},
				UploadUrl: uploadUrl,
			}, nil
		})
		resourceHandler := handlers.NewResourceHandler(dropService)

		reqBody := `{"resource_type": "text"}`
		req := httptest.NewRequest(http.MethodPost, "/resource", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := resourceHandler.Create(w, req)
		if err != nil {
			t.Errorf("Expected resource handler not to throw error on success, but got: %v", err.Error())
			return
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
		}

		var response handlers.ResourceWithUploadUrlDTO
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Id != "22222222-2222-2222-2222-222222222222" {
			t.Errorf("Expected resource id '22222222-2222-2222-2222-222222222222', but got '%v'", response.Id)
		}
		if response.Type != string(drop.TextResource) {
			t.Errorf("Expected resource type '%v', but got '%v'", drop.TextResource, response.Type)
		}
		if response.Name != nil {
			t.Errorf("Expected resource name to be nil, but got '%v'", *response.Name)
		}
		if response.UploadUrl != "http://example.org/upload" {
			t.Errorf("Expected upload url 'http://example.org/upload', but got '%v'", response.UploadUrl)
		}
	})

	t.Run("success with optional name field omitted", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.MustParse("33333333-3333-3333-3333-333333333333"))
		uploadUrl := "http://example.org/upload"

		dropService := mock.NewMockDropService().WithCreateResourceWithUploadUrl(func(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error) {
			return &drop.ResourceWithUploadUrl{
				Resource: drop.Resource{
					Id:   resourceId,
					Type: drop.FileResource,
					Name: name,
				},
				UploadUrl: uploadUrl,
			}, nil
		})
		resourceHandler := handlers.NewResourceHandler(dropService)

		reqBody := `{"resource_type": "file"}`
		req := httptest.NewRequest(http.MethodPost, "/resource", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()

		err := resourceHandler.Create(w, req)
		if err != nil {
			t.Errorf("Expected resource handler not to throw error on success, but got: %v", err.Error())
			return
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
		}

		var response handlers.ResourceWithUploadUrlDTO
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Id != "33333333-3333-3333-3333-333333333333" {
			t.Errorf("Expected resource id '33333333-3333-3333-3333-333333333333', but got '%v'", response.Id)
		}
		if response.Type != string(drop.FileResource) {
			t.Errorf("Expected resource type '%v', but got '%v'", drop.FileResource, response.Type)
		}
		if response.Name != nil {
			t.Errorf("Expected resource name to be nil, but got '%v'", *response.Name)
		}
		if response.UploadUrl != "http://example.org/upload" {
			t.Errorf("Expected upload url 'http://example.org/upload', but got '%v'", response.UploadUrl)
		}
	})
}
