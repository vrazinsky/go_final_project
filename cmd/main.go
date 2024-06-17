package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/vrazinsky/go-final-project/auth"
	"github.com/vrazinsky/go-final-project/handlers"
	"github.com/vrazinsky/go-final-project/store"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	dbfile := "scheduler.db"
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbfile = envFile
	}
	db := store.NewDbService(dbfile, ctx)
	err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	var port int = 7540
	a := auth.NewAuthService(os.Getenv("TODO_PASSWORD"), os.Getenv("AUTH_KEY"))
	r := chi.NewRouter()
	h := handlers.NewHandler(ctx, *db)
	FileServer(r, "/", http.Dir("web"))
	r.Get("/api/nextdate", h.HandleNextTime)
	r.Post("/api/task", a.Auth(h.HandleAddTask))
	r.Get("/api/task", a.Auth(h.HandleGetTask))
	r.Get("/api/tasks", a.Auth(h.HandleGetTasks))
	r.Put("/api/task", a.Auth(h.HandleUpdateTask))
	r.Post("/api/task/done", a.Auth(h.HandleCompleteTask))
	r.Delete("/api/task", a.Auth(h.HandleDeleteTask))

	r.Post("/api/signin", a.HandleSignIn)

	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}
	log.Println("listen on", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatal(err)
	}

}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
