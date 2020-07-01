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

// ClientInterface : The interface for client
type ClientInterface struct {
	TalkServiceClient *core.TalkServiceClient
	Config            *NepCoreConfig
	talkPath          string
	authToken         string
}

// NepCoreConfig : Configs for NepCore
type NepCoreConfig struct {
	Server string `json:"Server"`
	Header struct {
		UserAgent   string `json:"User-Agent"`
		Application string `json:"X-Line-Application"`
	} `json:"Header"`
}

func deBug(where string, err error) bool {
	if err != nil {
		fmt.Printf("NepCore Error #%s\nReason:\n%s\n\n", where, err)
		return false
	}
	return true
}

// SetRoutine : Set Timeout for goroutine
func SetRoutine(ms int) (context.Context, context.CancelFunc) {
	var addTime time.Duration = time.Duration(ms) * time.Millisecond
	timeout := time.Now().Add(addTime)
	ctx, deadline := context.WithDeadline(context.Background(), timeout)
	return ctx, deadline
}

// NewClientInterface : Set interface for client
func NewClientInterface(talkPath string) *ClientInterface {
	client := new(ClientInterface)
	client.talkPath = talkPath
	client.readConfig()
	client.setProtocol()
	return client
}

func (client *ClientInterface) readConfig() {
	jsonFile, err := os.Open("nepcore.json")
	deBug("Loading JSON config", err)
	defer jsonFile.Close()
	srcJSON, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(srcJSON, &client.Config)
	deBug("JSON config Initialize", err)
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
