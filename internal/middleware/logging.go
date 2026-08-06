package middleware

import (
	"log"
	"net/http"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {

	// function is returned and logs Method and URL Path
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path) // Method(GET, POST and etc) and r.URL.Path(e.g "/health") are logged
		next(w, r)                                // original handler called
	}

}
