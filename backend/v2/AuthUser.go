package colorboxd

import "net/http"

func AuthUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthUser"))
}
