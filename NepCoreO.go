/*
Package NepCoreO : The LINE Client Protocol of Star Neptune BOT
	===
			The Package is the OpenSource Version of NepCore

			LICENSE: AGPL 3.0

						Copyright(c) 2019 Star Inc. All Rights Reserved.
*/
package NepCoreO

import (
	"context"
	core "" // LINE TalkService Core you own
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// Seq :  int
var Seq int32

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

// SetThread : Set Timeout for goroutine
func SetThread(ms int) (context.Context, context.CancelFunc) {
	var addTime time.Duration = time.Duration(ms) * time.Millisecond
	timeout := time.Now().Add(addTime)
	ctx, deadline := context.WithDeadline(context.Background(), timeout)
	return ctx, deadline
}

// Connect : To connect to LINE Server
func Connect(authToken string, talkPath string, configJSON string) *core.TalkServiceClient {
	var err error
	// Config Headers
	var Configs config
	jsonFile, err := os.Open("config.json")
	deBug("Load JSON config", err)
	defer jsonFile.Close()
	srcJSON, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(srcJSON, &Configs)
	deBug("Login JSON Initialize", err)
	HeaderConfigs := Configs.Header
	// Set Transport
	var transport thrift.TTransport
	TalkURL := fmt.Sprintf("%s%s", Configs.Server, talkPath)
	transport, err = thrift.NewTHttpPostClient(TalkURL)
	deBug("Login Thrift Client Initialize", err)
	// Set Header
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", authToken)
	connect.SetHeader("User-Agent", HeaderConfigs.UserAgent)
	connect.SetHeader("X-Line-Application", HeaderConfigs.Application)
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	// Return Client
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}

// SendText : Send text message to someone
func SendText(client *core.TalkServiceClient, toID string, msgText string) {
	msgObj := core.NewMessage()
	msgObj.ContentType = core.ContentType_NONE
	msgObj.To = toID
	msgObj.Text = msgText
	ctx, _ := SetThread(500)
	_, err := client.SendMessage(ctx, Seq, msgObj)
	deBug("SendMessage - Text", err)
}
