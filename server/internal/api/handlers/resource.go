package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/falbru/falkdrop/internal/api/errors"
	"github.com/falbru/falkdrop/internal/app/auth"
	"github.com/falbru/falkdrop/internal/app/drop"
)

type ResourceHandler struct {
	dropService drop.DropService
}

func NewResourceHandler(dropService drop.DropService) *ResourceHandler {
	return &ResourceHandler{
		dropService: dropService,
	}
}

type ResourceDTO struct {
	Id   string  `json:"id"`
	Type string  `json:"type"`
	Name *string `json:"name"`
}

type ResourceWithDownloadUrlDTO struct {
	ResourceDTO
	DownloadUrl string `json:"download_url"`
}

type ResourceWithUploadUrlDTO struct {
	ResourceDTO
	UploadUrl string `json:"upload_url"`
}

type CreateResourceRequest struct {
	ResourceType string  `json:"resource_type"`
	Name         *string `json:"name"`
}

func (handler ResourceHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var req CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperror.New(http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if token == "" {
		return httperror.New(http.StatusUnauthorized, "bearer token not provided")
	}

	ctx := context.WithValue(context.Background(), "token", token)

	resource, err := handler.dropService.CreateResourceWithUploadUrl(ctx, drop.ResourceType(req.ResourceType), req.Name)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			return httperror.New(http.StatusUnauthorized, err.Error())
		}
		return fmt.Errorf("failed to create resource: %w", err)
	}

	response := ResourceWithUploadUrlDTO{
		ResourceDTO: ResourceDTO{
			Id:   resource.Id.String(),
			Type: string(resource.Type),
			Name: resource.Name,
		},
		UploadUrl: resource.UploadUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
