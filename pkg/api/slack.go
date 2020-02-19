package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

// ActionsEndpoint receives callbacks from Slack
func ActionsEndpoint(c *gin.Context) {
	//ctx := appengine.NewContext(c.Request)

	//c.Copy().GetRawData()
	//buf := make([]byte, 1024)
	//num, _ := c.Request.Body.Read(buf)
	//reqBody := string(buf[0:num])

	//log.Printf("Actions, Body: %v", printJSON(reqBody))

	//resp := fmt.Sprintf("%s/r/%s", env.Getenv("BASE_URL", "/"), c.FullPath())
	//c.JSON(http.StatusOK, resp)

	//var action slack.MessageActionPayload
	//err := c.BindJSON(&action)
	//if err != nil {
	//	platform.Report(err)
	//}

	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Printf("Actions, Body: %v", prettyPrintJSON(body))

}

func printJSON(target interface{}) string {
	b, err := json.Marshal(target)
	if err != nil {
		// TODO really ?
		log.Fatal(err)
	}
	return (string)(b)
}

func prettyPrintJSON(target interface{}) string {
	b, err := json.Marshal(target)
	if err != nil {
		// TODO really ?
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, " ", " ")

	return (string)(out.Bytes())
	//out.WriteTo(os.Stdout)
}
