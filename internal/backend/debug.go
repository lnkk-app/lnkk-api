package backend

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/platform/pkg/platform"
)

// DumpPayload just prints the entire request body to STDOUT
func DumpPayload(c *gin.Context) {

	buf, err := c.Copy().GetRawData()
	if err != nil {
		platform.ReportError(err)
	}

	body, err := url.QueryUnescape(string(buf))
	if err != nil {
		platform.ReportError(err)
	}

	log.Printf("%v", body)

}
