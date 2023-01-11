package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/topritchett/game-server/proxmox"
)

// Handler for http requests
type Handler struct {
	mux *http.ServeMux
}

type vmStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// New http handler
func New(s *http.ServeMux) *Handler {
	h := Handler{s}
	h.registerApiRoutes()

	return &h
}

// RegisterRoutes for all http endpoints
func (h *Handler) registerApiRoutes() {
	h.mux.HandleFunc("/api/", checkGetMethod(h.handleRoot))

	h.mux.HandleFunc("/api/vm/start", checkPostMethod(h.apiStartVM))
	h.mux.HandleFunc("/api/vm/stop", checkPostMethod(h.apiStopVM))
	h.mux.HandleFunc("/api/vm/status/", checkGetMethod(h.apiStatusVM))
}

func (h *Handler) apiStartVM(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", escapeURL(r.URL.Path))
	validateForm(w, r)

	vmid := proxmox.GetVMID(proxmox.Auth, proxmox.QemuUrl, r.FormValue("name"))
	if vmid == "0" {
		w.Write([]byte("No VM found\n"))
		return
	}

	started, err := proxmox.StartVM(proxmox.Auth, proxmox.QemuUrl, vmid)
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte(started))
}

func (h *Handler) apiStopVM(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", escapeURL(r.URL.Path))
	validateForm(w, r)
	vmid := proxmox.GetVMID(proxmox.Auth, proxmox.QemuUrl, r.FormValue("name"))
	if vmid == "0" {
		w.Write([]byte("No VM found\n"))
		return
	}

	started, err := proxmox.PauseVM(proxmox.Auth, proxmox.QemuUrl, vmid)
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte(started))
}

func (h *Handler) apiStatusVM(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", escapeURL(r.URL.Path))
	vmname := strings.TrimPrefix(r.URL.Path, "/api/vm/status/")
	vmid := proxmox.GetVMID(proxmox.Auth, proxmox.QemuUrl, vmname)
	if vmid == "0" {
		w.Write([]byte("No VM found. Please check the name and try again.\n"))
		return
	}

	status := proxmox.VMStatus(proxmox.Auth, proxmox.QemuUrl, vmid)

	vm := &vmStatus{Name: vmname, Status: status}
	json.NewEncoder(w).Encode(vm)
}

func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", escapeURL(r.URL.Path))
	http.ServeFile(w, r, "./static/api.json")
}

// function to escape url
func escapeURL(url string) string {
	escapedURL := strings.Replace(url, "\n", "", -1)
	escapedURL = strings.Replace(escapedURL, "\r", "", -1)
	return escapedURL
}

func checkGetMethod(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func checkPostMethod(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func validateForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) == 0 {
		http.Error(w, "No form data\nPlease provide a name\ne.g. name=vmname\n", http.StatusBadRequest)
		return
	}
	if r.FormValue("name") == "" {
		http.Error(w, "No name provided\nPlease provide a name\ne.g. name=vmname\n", http.StatusBadRequest)
		return
	}
}
