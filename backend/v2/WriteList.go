package colorboxd

import "net/http"

func WriteList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WriteList"))
}
