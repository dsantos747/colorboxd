package colorboxd

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("AuthUser", HTTPAuthUser)
	functions.HTTP("GetLists", HTTPGetLists)
	functions.HTTP("SortList", HTTPSortListById)
	functions.HTTP("WriteList", HTTPWriteList)
}
