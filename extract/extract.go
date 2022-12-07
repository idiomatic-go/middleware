package extract

import (
	"github.com/idiomatic-go/middleware/accesslog"
)

func Test() {
	accesslog.SetExtract(nil)
}
