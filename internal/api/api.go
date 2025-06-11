package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
	"todo-service/internal/client"
	"todo-service/internal/config"
	login "todo-service/internal/handlers/Auth/Login"
	register "todo-service/internal/handlers/Auth/Register"
	createtask "todo-service/internal/handlers/DBWorkWithTasks/CreateTask"
	deletetask "todo-service/internal/handlers/DBWorkWithTasks/DeleteTask"
	donetask "todo-service/internal/handlers/DBWorkWithTasks/DoneTask"
	getalltask "todo-service/internal/handlers/DBWorkWithTasks/GetAllTask"
	"todo-service/internal/logger"
	"todo-service/internal/middleware/auth"
	"todo-service/internal/middleware/cors"

	"github.com/gorilla/mux"
)

type API struct {
	log    *slog.Logger
	cfg    *config.Config
	router *mux.Router
	client *client.Client
}

func InitAPI() *API {
	log := logger.InitLogger()
	cfg := config.InitConfig(log)
	ctx := context.Background()
	client, err := client.New(ctx, cfg.GRPCPorts.GRPCApiAuth, cfg.GRPCPorts.GRPCApiDb, 5*time.Second, 4)
	if err != nil {
		log.Info("create client to failed")
		os.Exit(1)
	}
	API := &API{
		log:    log,
		cfg:    cfg,
		router: mux.NewRouter(),
		client: client,
	}
	return API
}
func (api *API) StartServer() {
	api.setupRoutes()
	server := &http.Server{
		Handler:      api.router,
		Addr:         api.cfg.Address,
		ReadTimeout:  api.cfg.TimeOut,
		WriteTimeout: api.cfg.TimeOut,
		IdleTimeout:  api.cfg.IdleTimeout,
	}
	api.log.Debug("server is running on", "addr:", server.Addr)
	api.log.Info("server is running")
	if err := server.ListenAndServe(); err != nil {
		api.log.Debug("server start is failed", "err", err.Error())
		api.log.Info("server is not running")
	}
}
func (api *API) setupRoutes() {
	api.router.Use(cors.New())

	publicRouter := api.router.PathPrefix("").Subrouter()
	publicRouter.Handle("/login", login.New(api.log, api.client)).Methods("POST", "OPTIONS")
	publicRouter.Handle("/register", register.New(api.log, api.client)).Methods("POST", "OPTIONS")

	privateRouter := api.router.PathPrefix("").Subrouter()
	privateRouter.Use(auth.New(api.log, api.cfg.SecretKey))

	privateRouter.Handle("/getalltasks", getalltask.New(api.log, api.client)).Methods("GET", "OPTIONS")
	privateRouter.Handle("/createtask", createtask.New(api.log, api.client)).Methods("POST", "OPTIONS")
	privateRouter.Handle("/deletetask", deletetask.New(api.log, api.client)).Methods("POST", "OPTIONS")
	privateRouter.Handle("/donetask", donetask.New(api.log, api.client)).Methods("POST", "OPTIONS")
}
