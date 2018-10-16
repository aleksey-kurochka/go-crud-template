package rest

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/kgoralski/go-crud-template/dao"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

const (
	defaultConfigFilePath  = "./_conf"
	configFilePathUsage    = "config file directory. Config file must be named 'conf_{env}.yml'."
	configFilePathFlagName = "configFilePath"
	envUsage               = "environment for app, prod, dev, test"
	envDefault             = "dev"
	envFlagname            = "env"
)

var configFilePath string
var env string

func init() {
	flag.StringVar(&configFilePath, configFilePathFlagName, defaultConfigFilePath, configFilePathUsage)
	flag.StringVar(&env, envFlagname, envDefault, envUsage)
	flag.Parse()
	configuration(configFilePath, env)
}

// Server Instance which contains router and dao
type Server struct {
	*http.Server
	r  *chi.Mux
	db *dao.BankAPI
}

// NewServer creates new Server with db connection pool
func NewServer() *Server {
	router := chi.NewRouter()
	server := &Server{db: setupDB(viper.GetString("database.URL")), r: router}
	server.routes()
	return server
}

func (s *Server) routes() {
	s.r.Get("/rest/banks/", commonHeaders(s.getBanksHandler))
	s.r.Get("/rest/banks/{id:[0-9]+}", commonHeaders(s.getBankByIDHandler))
	s.r.Post("/rest/banks/", commonHeaders(s.createBankHanlder))
	s.r.Delete("/rest/banks/{id:[0-9]+}", commonHeaders(s.deleteBankByIDHandler))
	s.r.Put("/rest/banks/{id:[0-9]+}", commonHeaders(s.updateBankHanlder))
	s.r.Delete("/rest/banks/", commonHeaders(s.deleteAllBanksHandler))
}

// Start launching the server
func (s *Server) Start() {
	log.Fatal(http.ListenAndServe(viper.GetString("Server.port"), s.r))
}

func configuration(path string, env string) {
	if flag.Lookup("test.v") != nil {
		env = "test"
		path = "../_conf"
	}
	log.Println("Environment is: " + env + " configFilePath is: " + path)
	viper.SetConfigName("conf_" + env)
	viper.AddConfigPath(path) // working directory
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal: %+v", err))
	}
}

func setupDB(dbURL string) *dao.BankAPI {
	var db, err = dao.NewBankAPI(dbURL)
	if err != nil {
		log.Fatal(fmt.Errorf("fatal: %+v", err))
	}
	return db
}
