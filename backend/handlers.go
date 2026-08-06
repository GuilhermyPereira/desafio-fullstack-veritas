package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// validaStatus garante que só existam as três colunas exigidas
func isValidStatus(status string) bool {
	return status == "A Fazer" || status == "Em Progresso" || status == "Concluídas"
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	var taskList []Task
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas (Requisito)
	if strings.TrimSpace(task.Title) == "" {
		http.Error(w, "Título é obrigatório", http.StatusBadRequest)
		return
	}
	if !isValidStatus(task.Status) {
		http.Error(w, "Status inválido", http.StatusBadRequest)
		return
	}

	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	// Gera um ID simples e único usando timestamp
	task.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	tasks[task.ID] = task

	if err := saveTasks(); err != nil {
		http.Error(w, "Falha ao salvar a tarefa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(updatedTask.Title) == "" {
		http.Error(w, "Título é obrigatório", http.StatusBadRequest)
		return
	}
	if !isValidStatus(updatedTask.Status) {
		http.Error(w, "Status inválido", http.StatusBadRequest)
		return
	}

	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		return
	}

	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.Status = updatedTask.Status
	tasks[id] = task

	if err := saveTasks(); err != nil {
		http.Error(w, "Falha ao salvar a tarefa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")

	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	if _, exists := tasks[id]; !exists {
		http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		return
	}

	delete(tasks, id)

	if err := saveTasks(); err != nil {
		http.Error(w, "Falha ao excluir a tarefa", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}