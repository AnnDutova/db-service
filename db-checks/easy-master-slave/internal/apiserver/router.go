package apiserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"easy-master-slave/internal/db"
	"easy-master-slave/internal/utils"
	"easy-master-slave/pkg/model"
)

var (
	serialId int64
)

type Router struct {
	Provider *db.DataBaseProvider
	Config   *model.Config
}

func configureRouters(r *Router) {
	http.HandleFunc("/write", r.HandleWrite)
	http.HandleFunc("/read", r.HandleRead)
}

func (rout *Router) HandleWrite(rw http.ResponseWriter, r *http.Request) {
	project := utils.GenerateRandomString()
	result, err := rout.Provider.Write().Exec("INSERT INTO project(id, title) VALUES ($1,$2)",
		serialId, project)
	if err != nil {
		respond(rw, r, http.StatusInternalServerError, result)
	}
	serialId++
	log.Printf("Write to db [ %d, %s ]", serialId, project)
	respond(rw, r, http.StatusOK, nil)
}

func (rout *Router) HandleRead(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows := rout.Provider.Read().QueryRowContext(ctx, "Select title from project WHERE id=$1", serialId)
	if rows.Err() != nil {
		respond(rw, r, http.StatusInternalServerError, rows.Err())
	}

	var project model.Project
	if err := rows.Scan(&project.ID, &project.Title); err != nil {
		respond(rw, r, http.StatusInternalServerError, err)
	}
	log.Printf("Write to db [ %d, %s ]", serialId, project)

	respond(rw, r, http.StatusOK, project)
}

func respond(w http.ResponseWriter, r *http.Request, code int, date interface{}) {
	w.WriteHeader(code)
	if date != nil {
		json.NewEncoder(w).Encode(date)
	}
}
