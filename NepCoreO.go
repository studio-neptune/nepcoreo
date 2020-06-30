/*
Package NepCoreO : The LINE Client Protocol of Star Neptune BOT
	===
			The Package is the OpenSource Version of NepCore

			LICENSE: Apache License 2.0

						Copyright(c) 2020 Star Inc. All Rights Reserved.
*/
package NepCoreO

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	core "github.com/star-inc/olsb_cores/libs/NepCoreO" // Replace LINE TalkService Core you own

	"github.com/star-inc/thrift_go/thrift"
)

// Seq :  int
var Seq int32

// ClientInterface : Set Timeout for goroutine
type ClientInterface struct {
	TalkServiceClient *core.TalkServiceClient
	Config            *config
	talkPath          string
	authToken         string
}

type headerConfig struct {
	UserAgent   string `json:"User-Agent"`
	Application string `json:"X-Line-Application"`
}

type config struct {
	Server string       `json:"Server"`
	Header headerConfig `json:"Header"`
}

func deBug(where string, err error) bool {
	if err != nil {
		fmt.Printf("NepCore Error #%s\nReason:\n%s\n\n", where, err)
		return false
	}
	return true
}

func readConfig(configObj *config) {
	jsonFile, err := os.Open("config.json")
	deBug("Loading JSON config", err)
	defer jsonFile.Close()
	srcJSON, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(srcJSON, &configObj)
	deBug("JSON config Initialize", err)
}

// SetRoutine : Set Timeout for goroutine
func SetRoutine(ms int) (context.Context, context.CancelFunc) {
	var addTime time.Duration = time.Duration(ms) * time.Millisecond
	timeout := time.Now().Add(addTime)
	ctx, deadline := context.WithDeadline(context.Background(), timeout)
	return ctx, deadline
}

// NewClientInterface : To connect to LINE Server
func NewClientInterface(talkPath string) *ClientInterface {
	client := new(ClientInterface)
	client.talkPath = talkPath
	readConfig(client.Config)
	client.setProtocol()
	return client
}

func (client *ClientInterface) setProtocol() {
	// Set Transport
	apiURL := fmt.Sprintf("%s%s", client.Config.Server, client.talkPath)
	transport, err := thrift.NewTHttpPostClient(apiURL)
	deBug("Login Thrift Client Initialize", err)

	// Set Header
	connect := transport.(*thrift.THttpClient)
	connect.SetHeader("User-Agent", client.Config.Header.UserAgent)
	connect.SetHeader("X-Line-Application", client.Config.Header.Application)
	if client.authToken != "" {
		connect.SetHeader("X-Line-Access", client.authToken)

	}
	protocol := thrift.NewTCompactProtocolFactory().GetProtocol(connect)

	// Configure Client
	client.TalkServiceClient = core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}

// Authorize : To connect to LINE Server
func (client *ClientInterface) Authorize(authToken string) {
	client.authToken = authToken
	client.setProtocol()
}

// SendText : Send text message to someone
func (client *ClientInterface) SendText(targetID string, contentText string) {
	msgObj := core.NewMessage()
	msgObj.ContentType = core.ContentType_NONE
	msgObj.To = targetID
	msgObj.Text = contentText
	ctx, _ := SetRoutine(500)
	_, err := client.TalkServiceClient.SendMessage(ctx, Seq, msgObj)
	deBug("SendMessage - Text", err)
}
