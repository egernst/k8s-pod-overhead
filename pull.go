package main

import (
	"encoding/json"
	"fmt"
	"time"
        //"github.com/prometheus/client_golang/prometheus"
        //"github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
    "net/http"
    "net/url"
	"k8s.io/client-go/rest"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)



func makeRequestAndGetValue(client *http.Client, req *http.Request, value interface{}) error {
        // TODO(directxman12): support validating certs by hostname
        response, err := client.Do(req)
        if err != nil {
                return err
        }
        defer response.Body.Close()
        body, err := ioutil.ReadAll(response.Body)
        if err != nil {
                return fmt.Errorf("failed to read response body - %v", err)
        }
        if response.StatusCode == http.StatusNotFound {
                return fmt.Errorf("not found")
        } else if response.StatusCode != http.StatusOK {
                return fmt.Errorf("request failed - %q, response: %q", response.Status, string(body))
        }

        kubeletAddr := "[unknown]"
        if req.URL != nil {
                kubeletAddr = req.URL.Host
        }
        fmt.Printf("Raw response from Kubelet at %s: %s", kubeletAddr, string(body))

        err = json.Unmarshal(body, value)
        if err != nil {
                return fmt.Errorf("failed to parse output. Response: %q. Error: %v", string(body), err)
        }
        return nil
}

func getSummary(client *http.Client) (*stats.Summary, error) {

	scheme := "https"

	url := url.URL{
		Scheme: scheme,
		Host: "kata-k8s-1:10250",
		Path: "/stats/summary",
		RawQuery: "only_cpu_and_memory=true",
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	summary := &stats.Summary{}

	err = makeRequestAndGetValue(client, req, summary)

	return summary, err
}


func recordPodMetrics(client *http.Client) {
        go func() {
                for {

  summary, err := getSummary(client)
  if err == nil {
	fmt.Printf("summary: %v\n", summary)
  } else {
	fmt.Printf(" err: %v\n", err)
  }

                        time.Sleep(2 * time.Second)
                }
        }()
}

func main() {

  // setup a client configuration so we can access.
  // Initially this is 100% totally, terribly insecure:	
  restConfig := &rest.Config{}
  restConfig = rest.AnonymousClientConfig(restConfig)
  restConfig.TLSClientConfig = rest.TLSClientConfig{
	Insecure: true,
	CAData: nil,
	CAFile: "",
   }

  transport, _ := rest.TransportFor(restConfig)

  client := &http.Client{
	  Transport: transport,
  }

  // startup process for pulling pod stats from kubelet
  recordPodMetrics(client)

  // post metrics
  http.Handle("/metrics", promhttp.Handler())
  http.ListenAndServe(":2112", nil)
}
