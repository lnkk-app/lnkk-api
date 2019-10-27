package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

// PostRequest is used to invoke a Slack Web API method by posting a JSON payload
func PostRequest(ctx context.Context, token, apiMethod string, request interface{}) (*StandardResponse, error) {
	return post(ctx, token, SlackEndpoint+apiMethod, request)
}

// GetRequest is used to query the Slack Web API
func GetRequest(ctx context.Context, token, apiMethod, query string, response interface{}) error {
	url := SlackEndpoint + apiMethod + "?" + query

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// post the request to Slack
	//client := urlfetch.Client(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// unmarshal the response
	return json.NewDecoder(resp.Body).Decode(&response)

}

// post allows to post data to Slack as specified in the Web API documentation.
// See https://api.slack.com/web
func post(ctx context.Context, token, url string, body interface{}) (*StandardResponse, error) {

	m, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// post the request to Slack
	//client := urlfetch.Client(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// unmarshal the response
	var apiResponse StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	return &apiResponse, err
}
