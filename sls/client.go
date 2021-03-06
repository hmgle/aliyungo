package sls

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"net/http"
	"time"
	"encoding/json"
)

type Client struct {
	accessKeyId     string //Access Key Id
	accessKeySecret string //Access Key Secret
	debug           bool
	httpClient      *http.Client
	version         string
	internal        bool
	region          common.Region
	endpoint        string
}

type Project struct {
	client      *Client
	Name        string `json:"projectName,omitempty"`
	Description string `json:"description,omitempty"`
}

type LogItem struct {
	Time    time.Time
	Content map[string]string
}

type LogGroupItem struct {
	Logs   []*LogItem
	Topic  string
	Source string
}

const (
	SLSDefaultEndpoint = "sls.aliyuncs.com"
	SLSAPIVersion = "0.6.0"
	METHOD_GET = "GET"
	METHOD_POST = "POST"
	METHOD_PUT = "PUT"
	METHOD_DELETE = "DELETE"
)

// NewClient creates a new instance of ECS client
func NewClient(region common.Region, internal bool, accessKeyId, accessKeySecret string) *Client {
	return &Client{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		internal:        internal,
		region:          region,
		version:         SLSAPIVersion,
		endpoint:        SLSDefaultEndpoint,
		httpClient:      &http.Client{},
	}
}

func (client *Client) Project(name string) (*Project, error) {

	newClient := client.forProject(name)

	req := &request{
		method: METHOD_GET,
		path: "/",
	}

	project := &Project{}

	if err := newClient.requestWithJsonResponse(req, project); err != nil {
		return nil, err
	}
	project.client = newClient
	return project, nil
}

func (client *Client) forProject(name string) *Client {
	newclient := *client

	region := string(client.region)
	if client.internal {
		region = fmt.Sprintf("%s-intranet", region)
	}
	newclient.endpoint = fmt.Sprintf("%s.%s.%s", name, region, SLSDefaultEndpoint)
	return &newclient
}

func (proj *Project) Delete() error {
	req := &request{
		method: METHOD_DELETE,
		path: "/",
	}

	return proj.client.requestWithClose(req)
}

func (client *Client) CreateProject(name string, description string) error {
	project := &Project{
		Name: name,
		Description: description,
	}
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	req := &request{
		method: METHOD_POST,
		path: "/",
		payload:     data,
		contentType: "application/json",
	}

	newClient := client.forProject(name)
	return newClient.requestWithClose(req)
}

//
//func marshal() ([]byte, error) {
//
//	logGroups := []*LogGroup{}
//	tmp := []*LogGroupItem
//	for _, logGroupItem := range tmp {
//
//		logs := []*Log{}
//		for _, logItem := range logGroupItem.Logs {
//			contents := []*Log_Content{}
//			for key, value := range logItem.Content {
//				contents = append(contents, &Log_Content{
//					Key: proto.String(key),
//					Value: proto.String(value),
//				})
//			}
//
//			logs = append(logs, &Log{
//				Time: proto.Uint32(uint32(LogItem.Time.Unix())),
//				Contents: contents,
//			})
//		}
//
//		logGroup := &LogGroup{
//			Topic: proto.String(LogGroupItem.Topic),
//			Source: proto.String(LogGroupItem.Source),
//			Logs: logs,
//		}
//		logGroups = append(logGroups, logGroup)
//	}
//
//	return proto.Marshal(&LogGroupList{
//		LogGroupList: logGroups,
//	})
//}
