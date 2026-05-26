package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/falbru/falkdrop/internal/app/drop"
)

type ResourceHandler struct {
	dropService *drop.DropService
}

func NewResourceHandler(dropService *drop.DropService) *ResourceHandler {
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

func (handler ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	var requestBody CreateResourceRequest
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	resource, err := handler.dropService.CreateResourceWithUploadUrl(drop.ResourceType(requestBody.ResourceType), requestBody.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
