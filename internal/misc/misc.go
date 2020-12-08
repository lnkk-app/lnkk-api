package misc

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// ExtractBodyAsString extracts a requests body, assuming it is a string
func ExtractBodyAsString(c *gin.Context) (string, error) {

	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}
	return "", nil
}
