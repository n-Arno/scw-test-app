package main

import (
	"crypto/subtle"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxprom "gitlab.com/msvechla/mux-prometheus/pkg/middleware"
)

func root(w http.ResponseWriter, config *Config) {
        e, err := ErrorDbHTML()
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
	db, err := connect(config)
	if err != nil {
                w.WriteHeader(500)
                w.Write(e)
		return
	}

	var news []News
	var count int64
	db.Model(&news).Count(&count)

	if count == 0 {
		db.Create(&News{Title: "Hello", Content: "Welcome to your first news"})
	}

	db.Find(&news)
	result, err := NewsHTML(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(result)
}

func admin(w http.ResponseWriter) {
	result, err := AdminHTML()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(result)
}

func adminNews(w http.ResponseWriter) {
	result, err := AdminNewsHTML()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(result)
}

func addNews(w http.ResponseWriter, r *http.Request, config *Config) {
	result, err := AddNewsHTML()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	db, err := connect(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db.Create(&News{Title: title, Content: content})
	w.WriteHeader(201)
	w.Write(result)
	return
}

func adminConfig(w http.ResponseWriter, config *Config) {
	result, err := AdminConfigHTML(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(result)
}

func updateConfig(w http.ResponseWriter, r *http.Request, config *Config, configPath string) {
	result, err := UpdateConfigHTML()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
        config.Db.Port =  r.PostFormValue("db_port")
        config.Db.Host =  r.PostFormValue("db_host")
        config.Db.Name =  r.PostFormValue("db_name")
        config.Db.User =  r.PostFormValue("db_user")
        config.Db.Pass =  r.PostFormValue("db_pass")
        WriteConfig(configPath, config)
	w.WriteHeader(200)
	w.Write(result)
	return
}

func BasicAuth(w http.ResponseWriter, r *http.Request, username, password, realm string) bool {

	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return false
	}

	return true
}

type Health struct {
	Status string
}

func routers(config *Config, configPath string) *mux.Router {
	username := config.Web.User
	password := config.Web.Pass

	adminPath := "/admin"

	instrumentation := muxprom.NewDefaultInstrumentation()
	topRouter := mux.NewRouter().StrictSlash(true)
	topRouter.Use(instrumentation.Middleware)

	adminRouter := mux.NewRouter().PathPrefix(adminPath).Subrouter().StrictSlash(true)

	topRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		root(w, config)
	})

	topRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		simpleAnswer(w, 200, Health{Status: "OK"})
	})

	topRouter.Path("/metrics").Handler(promhttp.Handler())

	adminRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		admin(w)
	})

	adminRouter.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		adminConfig(w, config)
	}).Methods("GET")

        adminRouter.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
                updateConfig(w, r, config, configPath)
        }).Methods("POST")

	adminRouter.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		adminNews(w)
	}).Methods("GET")

	adminRouter.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		addNews(w, r, config)
	}).Methods("POST")

	topRouter.PathPrefix(adminPath).Handler(negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			if BasicAuth(w, r, username, password, "Provide user name and password") {
				/* Call the next handler iff Basic-Auth succeeded */
				next(w, r)
			}
		}),
		negroni.Wrap(adminRouter),
	))

	return topRouter
}
