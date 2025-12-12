package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "sport-manager/internal/domain"
    "sport-manager/internal/service"
    "github.com/gorilla/mux"
)

type AthleteHandler struct {
    athleteService *service.AthleteService
}

func NewAthleteHandler(athleteService *service.AthleteService) *AthleteHandler {
    return &AthleteHandler{
        athleteService: athleteService,
    }
}

func (h *AthleteHandler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/api/athletes", h.GetAthletes).Methods("GET")
    router.HandleFunc("/api/athletes/{id}", h.GetAthlete).Methods("GET")
    router.HandleFunc("/api/athletes", h.CreateAthlete).Methods("POST")
    router.HandleFunc("/api/athletes/{id}", h.UpdateAthlete).Methods("PUT")
    router.HandleFunc("/api/athletes/{id}", h.DeleteAthlete).Methods("DELETE")
}

func (h *AthleteHandler) GetAthletes(w http.ResponseWriter, r *http.Request) {
    athletes, err := h.athleteService.GetAll(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(athletes)
}

func (h *AthleteHandler) GetAthlete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    athlete, err := h.athleteService.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    json.NewEncoder(w).Encode(athlete)
}

func (h *AthleteHandler) CreateAthlete(w http.ResponseWriter, r *http.Request) {
    var athlete domain.Athlete
    if err := json.NewDecoder(r.Body).Decode(&athlete); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if err := h.athleteService.Create(r.Context(), &athlete); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(athlete)
}

// ... остальные методы