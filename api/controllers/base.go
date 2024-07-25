package controllers

import (
	"api-center/cmd/database"
	"api-center/configs"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(config *configs.Config) {

	server.DB = database.ConnectDatabase(config)

	// err := database.MigrageTable(server.DB)

	// if err != nil {
	// 	log.Fatal("failed to migrate table: ", err)
	// }

	server.Router = mux.NewRouter()
	server.initializeRoutes(config)
	server.InitializeStaticAsset()
}

func (server *Server) Run(addr string) {
	srv := &http.Server{
		Addr:         fmt.Sprint(":", addr),
		Handler:      server.Router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on port " + addr)
	log.Fatal(srv.ListenAndServe())
}
