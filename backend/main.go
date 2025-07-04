package colorboxd

import (
	// "net/http"
	// _ "net/http/pprof"

	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	if err := LoadEnv(); err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
	}

	functions.HTTP("AuthUser", HTTPAuthUser)
	functions.HTTP("GetLists", HTTPGetLists)
	functions.HTTP("SortList", HTTPSortListById)
	functions.HTTP("WriteList", HTTPWriteList)

	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
}
