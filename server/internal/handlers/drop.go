package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/google/uuid"
)

type DropHandler struct {
	dropService *drop.DropService
}

func NewDropHandler(dropService *drop.DropService) *DropHandler {
	return &DropHandler{
		dropService: dropService,
	}
}

type DropWithResourceDownloadUrls struct {
	Id             string                       `json:"id"`
	ExpirationDate string                       `json:"expiration_date"`
	Resources      []ResourceWithDownloadUrlDTO `json:"resources"`
}

type CreateDropRequest struct {
	ResourceIds []uuid.UUID `json:"resource_ids"`
}

func (handler DropHandler) Get(w http.ResponseWriter, r *http.Request) {
	dropId := r.PathValue("dropId")

	if dropId == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "dropId is required")
		return
	}

	drop, err := handler.dropService.GetDropWithResourceDownloadUrls(drop.DropId(dropId))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	resourcesWithDownloadUrls := make([]ResourceWithDownloadUrlDTO, len(drop.Resources))
	for i, resource := range drop.Resources {
		resourcesWithDownloadUrls[i] = ResourceWithDownloadUrlDTO{
			ResourceDTO: ResourceDTO{
				Id:   resource.Id.String(),
				Type: string(resource.Type),
			},
			DownloadUrl: resource.DownloadUrl,
		}
	}

	response := DropWithResourceDownloadUrls{
		Id:             string(drop.Id),
		ExpirationDate: drop.ExpirationDate.String(),
		Resources:      resourcesWithDownloadUrls,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (handler DropHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	var requestBody CreateDropRequest
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	resourceIds := make([]drop.ResourceId, len(requestBody.ResourceIds))
	for i, resourceId := range requestBody.ResourceIds {
		resourceIds[i] = drop.ResourceId(resourceId)
	}

	drop, err := handler.dropService.CreateDrop(resourceIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	resourcesWithDownloadUrls := make([]ResourceWithDownloadUrlDTO, len(drop.Resources))
	for i, resource := range drop.Resources {
		resourcesWithDownloadUrls[i] = ResourceWithDownloadUrlDTO{
			ResourceDTO: ResourceDTO{
				Id:   resource.Id.String(),
				Type: string(resource.Type),
			},
			DownloadUrl: resource.DownloadUrl,
		}
	}

	response :=
		DropWithResourceDownloadUrls{
			Id:             string(drop.Id),
			ExpirationDate: drop.ExpirationDate.String(),
			Resources:      resourcesWithDownloadUrls,
		}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
