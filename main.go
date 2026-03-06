package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Middleware para autenticação Admin Token
func requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Admin-Token")
		if token != os.Getenv("ADMIN_TOKEN") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Administrative token required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Handlers
func uploadBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST /blob - upload handler"))
}

func multipartBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST /blob/multipart - multipart handler"))
}

func partBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PUT /blob/:id/part - part handler"))
}

func completeBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST /blob/:id/complete - complete handler"))
}

func listBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /blob - list handler"))
}

func getBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /blob/:id - get handler"))
}

func headBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HEAD /blob/:id - head handler"))
}

func deleteBlobHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DELETE /blob/:id - delete handler"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /health - health handler"))
}

func main() {
	r := mux.NewRouter()

	// Rotas protegidas
	r.Handle("/blob", requireAdmin(http.HandlerFunc(uploadBlobHandler))).Methods("POST")
	r.Handle("/blob/multipart", requireAdmin(http.HandlerFunc(multipartBlobHandler))).Methods("POST")
	r.Handle("/blob/{id}/part", requireAdmin(http.HandlerFunc(partBlobHandler))).Methods("PUT")
	r.Handle("/blob/{id}/complete", requireAdmin(http.HandlerFunc(completeBlobHandler))).Methods("POST")
	r.Handle("/blob", requireAdmin(http.HandlerFunc(listBlobHandler))).Methods("GET")
	r.Handle("/blob/{id}", requireAdmin(http.HandlerFunc(deleteBlobHandler))).Methods("DELETE")

	// Rotas públicas
	r.HandleFunc("/blob/{id}", getBlobHandler).Methods("GET")
	r.HandleFunc("/blob/{id}", headBlobHandler).Methods("HEAD")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	log.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
