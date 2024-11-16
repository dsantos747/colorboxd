package colorboxd

import "net/http"

func SortList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SortList"))
}
