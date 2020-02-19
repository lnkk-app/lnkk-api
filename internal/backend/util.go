package backend

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/lnkk-ai/lnkk/pkg/platform"
)

// DumpPayload just prints the entire request body to STDOUT
func DumpPayload(c *gin.Context) {

	buf, err := c.Copy().GetRawData()
	if err != nil {
		platform.Report(err)
	}

	body, err := url.QueryUnescape(string(buf))
	if err != nil {
		platform.Report(err)
	}

	log.Printf("%v", body)

}
