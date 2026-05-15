package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type DnsService interface {
	Add(ctx context.Context, newDns string) (string, error)
	Delete(ctx context.Context, dnsToDelete string) (string, error)
	List(ctx context.Context) ([]string, error)
}

type API struct {
	Service DnsService
}

func NewAPI(service DnsService) *API {
	return &API{
		Service: service,
	}
}

type DNSRequest struct {
	DNS string `json:"dns"`
}

type DNSResponse struct {
	DNS string `json:"dns"`
}

type DNSListResponse struct {
	Items []string `json:"items"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (a *API) handleDNS(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.addDNS(w, r)
	case http.MethodGet:
		a.listDNS(w, r)
	case http.MethodDelete:
		a.removeDNS(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func writeJSON(w http.ResponseWriter, status int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/dns", a.handleDNS)

	return mux
}

func readJSON(w http.ResponseWriter, r *http.Request, request any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	return decode.Decode(request)
}

func caseErrors(w http.ResponseWriter, message string, err error) {
	log.Println(message, " ", err)
	switch {
	case errors.Is(err, ErrIsIncorrect):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, message)
	}
}

func (a *API) addDNS(w http.ResponseWriter, r *http.Request) {
	var req DNSRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "failed to decode request")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	dns, err := a.Service.Add(ctx, req.DNS)
	if err != nil {
		caseErrors(w, "failed to add dns", err)
		return
	}
	writeJSON(w, http.StatusCreated, DNSResponse{
		DNS: dns,
	})
}

func (a *API) listDNS(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	items, err := a.Service.List(ctx)
	if err != nil {
		caseErrors(w, "failed to list DNS", err)
		return
	}
	writeJSON(w, http.StatusOK, DNSListResponse{
		Items: items,
	})
}

func (a *API) removeDNS(w http.ResponseWriter, r *http.Request) {
	var req DNSRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "failed to decode request")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	dns, err := a.Service.Delete(ctx, req.DNS)
	if err != nil {
		caseErrors(w, "failed to delete DNS", err)
		return
	}
	writeJSON(w, http.StatusOK, DNSResponse{
		DNS: dns,
	})
}
