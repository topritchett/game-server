package server

import (
	"log"
	"net/http"

	"github.com/topritchett/game-server/proxmox"
	"github.com/topritchett/game-server/work"
)

// Handler for http requests
type Handler struct {
	mux *http.ServeMux
}

// New http handler
func New(s *http.ServeMux) *Handler {
	h := Handler{s}
	h.registerRoutes()

	return &h
}

// RegisterRoutes for all http endpoints
func (h *Handler) registerRoutes() {
	http.Handle("/static/", http.FileServer(http.Dir("static")))
	h.mux.HandleFunc("/", h.HelloWorld)

	h.mux.HandleFunc("/startvm", h.ServerStartVM)

	h.mux.HandleFunc("/proxurl", h.ServerGetProxUrl)
	h.mux.HandleFunc("/startwork", h.ServerStartWork)
	h.mux.HandleFunc("/pausework", h.ServerPauseWork)
}

func (h *Handler) ServerStartVM(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", r.URL.Path)
	started, err := proxmox.StartVM(proxmox.Auth, proxmox.QemuUrl, "100")
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte(started))
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Hello World!"))
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", r.URL.Path)
	http.ServeFile(w, r, "./static/index.html")
}

func (h *Handler) ServerGetProxUrl(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", r.URL.Path)
	w.Write([]byte(proxmox.GetProxUrl(proxmox.Auth, proxmox.QemuUrl)))
}

func (h *Handler) ServerStartWork(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", r.URL.Path)
	w.Write([]byte(work.StartWorkVMs(proxmox.Auth, proxmox.QemuUrl)))
}

func (h *Handler) ServerPauseWork(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, http.StatusOK, "from", r.RemoteAddr, "to", r.URL.Path)
	w.Write([]byte(work.PauseWorkVMs(proxmox.Auth, proxmox.QemuUrl)))
}
