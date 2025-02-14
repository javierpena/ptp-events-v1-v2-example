// Copyright 2020 The Cloud Native Events Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "bytes"
    "encoding/json"
    "io"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"

    "github.com/redhat-cne/sdk-go/pkg/pubsub"
    "github.com/redhat-cne/cloud-event-proxy/pkg/restclient"
    "github.com/redhat-cne/sdk-go/pkg/types"
    v1pubsub "github.com/redhat-cne/sdk-go/v1/pubsub"
	log "github.com/sirupsen/logrus"
)


var (
	apiAddr            = "ptp-event-publisher-service-NODE_NAME.openshift-ptp.svc.cluster.local:9043"
	apiPath            = "/api/ocloudNotifications/v2/"
	localAPIAddr       = "consumer-events-subscription-service.cloud-events.svc.cluster.local:9043"
	subs               []*pubsub.PubSub
)



func main() {
    var sub1, sub2, sub3 string

    log.Infof("Waiting for 30 seconds before subscribing. In a real deployment, you should check for the publisher health instead.")
	time.Sleep(30 * time.Second)

    // Spin local API
	go server()

    log.Infof("Creating subscriptions.")
    // Subscribe to PTP events using REST API
    nodeName := os.Getenv("NODE_NAME")
    apiAddr = strings.Replace(apiAddr, "NODE_NAME", nodeName, -1)

    s, e := createSubscription(fmt.Sprintf("/cluster/node/%v/sync/sync-status/os-clock-sync-state", nodeName))
    if e != nil {
        fmt.Printf("Failed to create subscription: %v\n", e)
    } else {
        sub1 = s.Resource
    }

    s, e = createSubscription(fmt.Sprintf("/cluster/node/%v/sync/ptp-status/clock-class", nodeName))
    if e != nil {
        fmt.Printf("Failed to create subscription: %v\n", e)
    } else {
        sub2 = s.Resource
    }

    s, e = createSubscription(fmt.Sprintf("/cluster/node/%v/sync/ptp-status/lock-state", nodeName))
    if e != nil {
        fmt.Printf("Failed to create subscription: %v\n", e)
    } else {
        sub3 = s.Resource
    }

    log.Infof("Now getting the current state every 30 seconds")
    for {
        log.Infof("Current state:")
        getCurrentState(sub1)
        getCurrentState(sub2)
        getCurrentState(sub3)
        time.Sleep(30 * time.Second)
    }
}


func server() {
  http.HandleFunc("/event", getEvent)
  http.ListenAndServe(":9043", nil)
}

func getEvent(w http.ResponseWriter, req *http.Request) {
  defer req.Body.Close()
  bodyBytes, err := io.ReadAll(req.Body)
  if err != nil {
    log.Errorf("error reading event %v", err)
  }
  e := string(bodyBytes)
  if e != "" {
    //processEvent(bodyBytes)
    log.Infof("received event %s", string(bodyBytes))
  }

  w.WriteHeader(http.StatusNoContent)
}

// Create PTP event subscriptions POST
func createSubscription(resourceAddress string) (sub pubsub.PubSub, err error) {
  var status int

  subURL := &types.URI{URL: url.URL{Scheme: "http",
    Host: apiAddr,
    Path: fmt.Sprintf("%s%s", apiPath, "subscriptions")}}
  endpointURL := &types.URI{URL: url.URL{Scheme: "http",
    Host: localAPIAddr,
    Path: "event"}}

  sub = v1pubsub.NewPubSub(endpointURL, resourceAddress, "2.0")
  var subB []byte

  if subB, err = json.Marshal(&sub); err == nil {
    rc := restclient.New()
    if status, subB = rc.PostWithReturn(subURL, subB); status != http.StatusCreated {
      err = fmt.Errorf("error in subscription creation api at %s, returned status %d", subURL, status)
    } else {
      err = json.Unmarshal(subB, &sub)
    }
  } else {
    err = fmt.Errorf("failed to marshal subscription for %s", resourceAddress)
  }
  return
}

// Get PTP event state for the resource
func getCurrentState(resource string) {
  //Create publisher
  url := &types.URI{URL: url.URL{Scheme: "http",
    Host: apiAddr,
    Path: fmt.Sprintf("/api/ocloudNotifications/v2/%s/CurrentState", resource)}}
  rc := restclient.New()
  status, event := rc.Get(url)
  if status != http.StatusOK {
    log.Errorf("CurrentState:error %d from url %s, %s", status, url.String(), event)
  } else {
    // Pretty print the received JSON data
    buf := &bytes.Buffer{}
    if err := json.Indent(buf, []byte(event), "", "\t"); err != nil {
        panic(err)
    }
    fmt.Println(buf)    
  }
}
