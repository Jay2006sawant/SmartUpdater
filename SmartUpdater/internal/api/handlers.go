package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/smartupdater/internal/models"
	"gorm.io/gorm"
)

// Handler handles HTTP requests
type Handler struct {
	db            *gorm.DB
	githubService *GitHubService
}

// NewHandler creates a new API handler
func NewHandler(db *gorm.DB, githubService *GitHubService) *Handler {
	return &Handler{
		db:            db,
		githubService: githubService,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/repositories", h.ListRepositories).Methods("GET")
	r.HandleFunc("/repositories", h.AddRepository).Methods("POST")
	r.HandleFunc("/repositories/{id}", h.GetRepository).Methods("GET")
	r.HandleFunc("/repositories/{id}", h.UpdateRepository).Methods("PUT")
	r.HandleFunc("/repositories/{id}", h.DeleteRepository).Methods("DELETE")
	r.HandleFunc("/repositories/{id}/updates", h.ListUpdates).Methods("GET")
}

// ListRepositories handles GET /repositories
func (h *Handler) ListRepositories(w http.ResponseWriter, r *http.Request) {
	var repositories []models.Repository
	if err := h.db.Find(&repositories).Error; err != nil {
		http.Error(w, "Failed to fetch repositories", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(repositories)
}

// AddRepository handles POST /repositories
func (h *Handler) AddRepository(w http.ResponseWriter, r *http.Request) {
	var repo models.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch repository info from GitHub
	info, err := h.githubService.GetRepositoryInfo(r.Context(), repo.Owner, repo.Name)
	if err != nil {
		http.Error(w, "Failed to fetch repository info", http.StatusBadRequest)
		return
	}

	repo.LastUpdated = info.LastUpdated
	repo.LastCommit = info.LastCommit
	repo.Active = true

	if err := h.db.Create(&repo).Error; err != nil {
		http.Error(w, "Failed to create repository", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
}

// GetRepository handles GET /repositories/{id}
func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	var repo models.Repository
	if err := h.db.First(&repo, id).Error; err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(repo)
}

// UpdateRepository handles PUT /repositories/{id}
func (h *Handler) UpdateRepository(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	var repo models.Repository
	if err := h.db.First(&repo, id).Error; err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	var updates models.Repository
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update only allowed fields
	repo.UpdateFrequency = updates.UpdateFrequency
	repo.Active = updates.Active

	if err := h.db.Save(&repo).Error; err != nil {
		http.Error(w, "Failed to update repository", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(repo)
}

// DeleteRepository handles DELETE /repositories/{id}
func (h *Handler) DeleteRepository(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	if err := h.db.Delete(&models.Repository{}, id).Error; err != nil {
		http.Error(w, "Failed to delete repository", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUpdates handles GET /repositories/{id}/updates
func (h *Handler) ListUpdates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid repository ID", http.StatusBadRequest)
		return
	}

	var updates []models.UpdateHistory
	if err := h.db.Where("repository_id = ?", id).Find(&updates).Error; err != nil {
		http.Error(w, "Failed to fetch updates", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updates)
} 