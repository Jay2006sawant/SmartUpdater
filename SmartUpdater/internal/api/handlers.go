package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/smartupdater/internal/models"
	"github.com/smartupdater/internal/services"
)

type Handler struct {
	scheduler *services.SchedulerService
	db        *gorm.DB
}

func NewHandler(scheduler *services.SchedulerService, db *gorm.DB) *Handler {
	return &Handler{
		scheduler: scheduler,
		db:        db,
	}
}

func (h *Handler) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/health", h.healthCheck).Methods("GET")
	r.Handle("/api/v1/metrics", promhttp.Handler())
	r.HandleFunc("/api/v1/repositories", h.listRepositories).Methods("GET")
	r.HandleFunc("/api/v1/repositories", h.addRepository).Methods("POST")
	r.HandleFunc("/api/v1/stats", h.getStats).Methods("GET")
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

type AddRepositoryRequest struct {
	Owner           string `json:"owner"`
	Name            string `json:"name"`
	UpdateFrequency string `json:"update_frequency"`
}

func (h *Handler) addRepository(w http.ResponseWriter, r *http.Request) {
	var req AddRepositoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repo := &models.Repository{
		Owner:           req.Owner,
		Name:            req.Name,
		UpdateFrequency: req.UpdateFrequency,
		Active:          true,
	}

	if err := h.db.Create(repo).Error; err != nil {
		logrus.Errorf("Error creating repository: %v", err)
		http.Error(w, "Error creating repository", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
}

func (h *Handler) listRepositories(w http.ResponseWriter, r *http.Request) {
	var repos []models.Repository
	if err := h.db.Find(&repos).Error; err != nil {
		logrus.Errorf("Error fetching repositories: %v", err)
		http.Error(w, "Error fetching repositories", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(repos)
}

type Stats struct {
	TotalRepositories  int64   `json:"total_repositories"`
	ActiveRepositories int64   `json:"active_repositories"`
	AvgCISuccessRate  float64 `json:"avg_ci_success_rate"`
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	var stats Stats
	h.db.Model(&models.Repository{}).Count(&stats.TotalRepositories)
	h.db.Model(&models.Repository{}).Where("active = ?", true).Count(&stats.ActiveRepositories)
	
	var avgSuccessRate float64
	h.db.Model(&models.Repository{}).Where("active = ?", true).Select("AVG(ci_success_rate)").Row().Scan(&avgSuccessRate)
	stats.AvgCISuccessRate = avgSuccessRate

	json.NewEncoder(w).Encode(stats)
} 