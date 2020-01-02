package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/ganeshrahul23/snippetbox/pkg/models"
	"github.com/ganeshrahul23/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type snippetInterface interface {
	Insert(string, string, string) (int, error)
	Get(int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

type userInterface interface {
	Insert(string, string, string) error
	Authenticate(string, string) (int, error)
	Get(int) (*models.User, error)
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      snippetInterface
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         userInterface
}

type contextKey string

var contextKeyUser = contextKey("user")

func main() {
	addr := flag.String("addr", ":8080", "HTTP Network Address")
	dsn := flag.String("dsn", "web:ganesh23@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache(".\\ui\\html\\")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	//err := http.ListenAndServe(*addr, mux)
	err = srv.ListenAndServeTLS(".\\tls\\cert.pem", ".\\tls\\key.pem")
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
