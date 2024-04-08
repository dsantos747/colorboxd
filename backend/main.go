package colorboxd

import (
	// "net/http"
	// _ "net/http/pprof"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("AuthUser", HTTPAuthUser)
	functions.HTTP("GetLists", HTTPGetLists)
	functions.HTTP("SortList", HTTPSortListById)
	functions.HTTP("WriteList", HTTPWriteList)

	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
}
