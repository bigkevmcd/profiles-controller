package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	profilesv1alpha1 "github.com/bigkevmcd/profiles-controller/api/v1alpha1"
	"github.com/bigkevmcd/profiles-controller/controllers"
)

// APIRouter is an HTTP API for accessing app configurations.
type APIRouter struct {
	*httprouter.Router
	repository controllers.ProfilesRepository
}

// NewRouter creates and returns a new APIRouter.
func NewRouter(r controllers.ProfilesRepository) *APIRouter {
	api := &APIRouter{
		Router:     httprouter.New(),
		repository: r,
	}
	api.HandlerFunc(http.MethodGet, "/profiles", api.SearchProfiles)
	return api
}

// SearchProfile searches profiles for matching fields.
func (a *APIRouter) SearchProfiles(w http.ResponseWriter, r *http.Request) {
	profiles := []profilesv1alpha1.Profile{}
	found, err := a.repository.Search(r.URL.Query().Get("q"))
	if err != nil {
		marshalErrorResponse(w, "failed to search for profiles")
		return
	}
	for i := range found {
		profiles = append(profiles, found[i])
	}
	marshalResponse(w, searchResponse{Profiles: profiles})
}

func marshalResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode response: %s", err)
	}
}

func marshalErrorResponse(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}
