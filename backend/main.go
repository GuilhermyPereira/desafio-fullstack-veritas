package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTasksHandler(w, r)
	case http.MethodPost:
		createTaskHandler(w, r)
	case http.MethodPut:
		if r.URL.Path == "/tasks" || r.URL.Path == "/tasks/" {
			http.Error(w, "ID é obrigatório na URL", http.StatusBadRequest)
			return
		}
		updateTaskHandler(w, r)
	case http.MethodDelete:
		if r.URL.Path == "/tasks" || r.URL.Path == "/tasks/" {
			http.Error(w, "ID é obrigatório na URL", http.StatusBadRequest)
			return
		}
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func securityAndLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Loga o método e a rota acessada
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		
		// Injeta cabeçalhos defensivos básicos
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Requisito bônus: carrega tarefas do JSON ao iniciar[cite: 1]
	loadTasks()

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", tasksHandler)
	mux.HandleFunc("/tasks/", tasksHandler)

	// Requisito obrigatório: Configura o CORS para o Frontend[cite: 1]
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
	})

	// Aplica o middleware de segurança e log junto com o CORS
	handler := c.Handler(securityAndLogMiddleware(mux))

	fmt.Println("Backend rodando perfeitamente na porta 8080...")
	fmt.Println("Acesse o Frontend em: http://localhost:3000")
	fmt.Println("API do Backend em: http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", handler))
}