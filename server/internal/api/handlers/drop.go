package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	httperror "github.com/falbru/falkdrop/internal/api/errors"
	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/google/uuid"
)

type DropHandler struct {
	dropService drop.DropService
}

func NewDropHandler(dropService drop.DropService) *DropHandler {
	return &DropHandler{
		dropService: dropService,
	}
}

type DropWithResourceDownloadUrlsDTO struct {
	Id             string                       `json:"id"`
	ExpirationDate string                       `json:"expiration_date"`
	Resources      []ResourceWithDownloadUrlDTO `json:"resources"`
}

type CreateDropRequest struct {
	ResourceIds []uuid.UUID `json:"resource_ids"`
}

func (handler DropHandler) Get(w http.ResponseWriter, r *http.Request) error {
	dropId := r.PathValue("dropId")

	if dropId == "" {
		return httperror.New(http.StatusBadRequest, "dropId is required")
	}

	d, err := handler.dropService.GetDropWithResourceDownloadUrls(context.Background(), drop.DropId(dropId))
	if err != nil {
		var dropNotFoundError drop.ErrDropNotFound
		if errors.As(err, &dropNotFoundError) {
			return httperror.NotFound(dropNotFoundError.Error())
		}

		return fmt.Errorf("failed to get drop: %w", err)
	}

	if d == nil {
		return httperror.NotFound(fmt.Sprintf("drop with dropId %v not found", dropId))
	}

	resourcesWithDownloadUrls := make([]ResourceWithDownloadUrlDTO, len(d.Resources))
	for i, resource := range d.Resources {
		resourcesWithDownloadUrls[i] = ResourceWithDownloadUrlDTO{
			ResourceDTO: ResourceDTO{
				Id:   resource.Id.String(),
				Type: string(resource.Type),
				Name: resource.Name,
			},
			DownloadUrl: resource.DownloadUrl,
		}
	}

	response := DropWithResourceDownloadUrlsDTO{
		Id:             string(d.Id),
		ExpirationDate: d.ExpirationDate.String(),
		Resources:      resourcesWithDownloadUrls,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("failed to encode response: %s", err.Error())
	}

	return nil
}

func (handler DropHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var req CreateDropRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperror.New(http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	if len(req.ResourceIds) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return httperror.New(http.StatusBadRequest, "resourceIds is empty")
	}

	resourceIds := make([]drop.ResourceId, len(req.ResourceIds))
	for i, resourceId := range req.ResourceIds {
		resourceIds[i] = drop.ResourceId(resourceId)
	}

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(context.Background(), "token", token)

	drp, err := handler.dropService.CreateDrop(ctx, resourceIds)
	if err != nil {
		var errResourcesNotFound drop.ErrResourcesNotFound
		var errResourceAlreadyBelongsToDrop drop.ErrResourceAlreadyBelongsToDrop
		if errors.As(err, &errResourcesNotFound) {
			return httperror.NotFound(err.Error())
		} else if errors.As(err, &errResourceAlreadyBelongsToDrop) {
			return httperror.New(http.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("failed to create drop: %w", err)
	}

	resourcesWithDownloadUrls := make([]ResourceWithDownloadUrlDTO, len(drp.Resources))
	for i, resource := range drp.Resources {
		resourcesWithDownloadUrls[i] = ResourceWithDownloadUrlDTO{
			ResourceDTO: ResourceDTO{
				Id:   resource.Id.String(),
				Type: string(resource.Type),
				Name: resource.Name,
			},
			DownloadUrl: resource.DownloadUrl,
		}
	}

	response :=
		DropWithResourceDownloadUrlsDTO{
			Id:             string(drp.Id),
			ExpirationDate: drp.ExpirationDate.String(),
			Resources:      resourcesWithDownloadUrls,
		}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
