// Copyright 2021 Star Inc.(https://starinc.xyz)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package NepCoreO

import (
	"fmt"
	core "github.com/star-inc/olsb_cores/libs/NepCoreO" // Replace LINE TalkService Core you own
	"github.com/star-inc/thrift_go/thrift"
)

// NewClient : Set up a talk service client
func NewClient(path string, config *Config) *core.TalkServiceClient {
	// Set Transport
	apiEndpoint := fmt.Sprintf("%s%s", config.Server, path)
	transport, err := thrift.NewTHttpClient(apiEndpoint)
	if err != nil {
		panic(err)
	}

	// Set Header
	connect := transport.(*thrift.THttpClient)
	connect.SetHeader("User-Agent", config.Header.UserAgent)
	connect.SetHeader("X-Line-Application", config.Header.Application)
	if config.Header.Access != "" {
		connect.SetHeader("X-Line-Access", config.Header.Access)
	}

	// Configure Client
	protocol := thrift.NewTCompactProtocolFactory().GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
