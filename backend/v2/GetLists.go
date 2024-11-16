package colorboxd

import "net/http"

func GetLists(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetLists"))
}
