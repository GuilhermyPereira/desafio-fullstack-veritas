package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Estrutura principal baseada nos requisitos
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var (
	tasks      = make(map[string]Task)
	tasksMutex sync.Mutex
	dataFile   = "tasks.json"
)

// loadTasks lê o arquivo JSON ao iniciar o servidor
func loadTasks() {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	file, err := os.ReadFile(dataFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return // Arquivo não existe ainda, começa vazio
		}
		panic(err)
	}

	var taskList []Task
	if err := json.Unmarshal(file, &taskList); err != nil {
		panic(err)
	}

	for _, task := range taskList {
		tasks[task.ID] = task
	}
}

// saveTasks escreve o mapa atual no arquivo JSON
func saveTasks() error {
	var taskList []Task
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	data, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, data, 0644)
}