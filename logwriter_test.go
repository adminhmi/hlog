package hlog_test

import (
	"github.com/adminhmi/hlog"
	"log"
	"net/http"
)

func ExampleLogger_Writer_httpServer() {
	logger := hlog.New()
	w := logger.Writer()
	defer w.Close()

	srv := http.Server{
		// create a stdlib log.Logger that writes to
		// Hlog.Logger.
		ErrorLog: log.New(w, "", 0),
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func ExampleLogger_Writer_stdlib() {
	logger := hlog.New()
	logger.Formatter = &hlog.JSONFormatter{}

	// Use Hlog for standard log output
	// Note that `log` here references stdlib's log
	// Not Hlog imported under the name `log`.
	log.SetOutput(logger.Writer())
}
