package receptor

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry/gunk/urljoiner"
	"github.com/tedsuo/rata"
)

type Client interface {
	CreateTask(CreateTaskRequest) error
}

func NewClient(addr, user, password string) Client {
	return &client{
		user:       user,
		password:   password,
		httpClient: &http.Client{},
		reqGen:     rata.NewRequestGenerator(urljoiner.Join("http://", addr), Routes),
	}
}

type client struct {
	user       string
	password   string
	httpClient *http.Client
	reqGen     *rata.RequestGenerator
}

func (c *client) CreateTask(request CreateTaskRequest) error {
	return c.doRequest(CreateTask, nil, request)
}

func (c *client) doRequest(requestName string, params rata.Params, request interface{}) error {
	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := c.reqGen.CreateRequest(requestName, params, bytes.NewReader(requestJson))
	if err != nil {
		return err
	}

	req.ContentLength = int64(len(requestJson))
	req.SetBasicAuth(c.user, c.password)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		errResponse := Error{}
		json.NewDecoder(res.Body).Decode(&errResponse)
		return errResponse
	}

	return nil
}
