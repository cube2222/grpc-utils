package health

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func LaunchHealthCheckHandler() {
	m := mux.NewRouter()
	m.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "I'm alive.")
	})
	log.Fatal(http.ListenAndServe(":6666", m))
}
