package http

import (
	"Todo"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func (s *Server) mountTodoRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", s.handleTaskIndex)
	r.Post("/", s.handleTaskCreate)
	r.Route("/{todoId}", func(r chi.Router) {
		r.Get("/", s.handleTaskView)
		r.Patch("/", s.handleTaskUpdate)
		r.Delete("/", s.handleTaskDelete)
	})

	return r
}

func (s *Server) handleTaskIndex(w http.ResponseWriter, r *http.Request) {
	var filter Todo.TaskFilter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		Error(w, r, Todo.Errorf(Todo.EINVALID, "Invalid JSON body"))
		return
	}

	tasks, err := s.TaskService.GetTasks(r.Context(), filter)

	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(tasksResponse{
		Tasks: tasks,
	}); err != nil {
		LogError(r, err)
		return
	}
}

func (s *Server) handleTaskView(w http.ResponseWriter, r *http.Request) {

	todoId, err := strconv.Atoi(chi.URLParam(r, "todoId"))

	if err != nil {
		Error(w, r, Todo.Errorf(Todo.EINVALID, "Invalid Id format"))
	}

	task, err := s.TaskService.GetTaskById(r.Context(), todoId)

	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		LogError(r, err)
		return
	}
}

func (s *Server) handleTaskCreate(w http.ResponseWriter, r *http.Request) {

	var task Todo.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		Error(w, r, Todo.Errorf(Todo.EINVALID, "Invalid JSON body"))
		return
	}

	err := s.TaskService.CreateTask(r.Context(), &task)

	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		LogError(r, err)
		return
	}

}

func (s *Server) handleTaskUpdate(w http.ResponseWriter, r *http.Request) {

	var upd Todo.TaskUpdate

	todoId, err := strconv.Atoi(chi.URLParam(r, "todoId"))

	if err != nil {
		Error(w, r, Todo.Errorf(Todo.EINVALID, "Invalid Id format"))
	}

	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, Todo.Errorf(Todo.EINVALID, "Invalid JSON body"))
		return
	}

	task, err := s.TaskService.UpdateTask(r.Context(), todoId, upd)

	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		LogError(r, err)
		return
	}
}

func (s *Server) handleTaskDelete(w http.ResponseWriter, r *http.Request) {

	todoId, err := strconv.Atoi(chi.URLParam(r, "todoId"))

	if err != nil {
		Error(w, r, Todo.Errorf(Todo.EINVALID, "Invalid Id format"))
	}

	if err := s.TaskService.DeleteTask(r.Context(), todoId); err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{}`))

}

type tasksResponse struct {
	Tasks []*Todo.Task `json:tasks`
}
