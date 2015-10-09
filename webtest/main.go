package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	kclient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
)

var (
	local     = flag.Bool("local", false, "set to true if running on local machine not within cluster")
	localPort = flag.Int("localport", 8001, "port that kubectl proxy is running on (local must be true)")

	client *kclient.Client
)

func main() {

	flag.Parse()
	var (
		cfg *kclient.Config
		err error
	)
	if *local {
		cfg = &kclient.Config{Host: fmt.Sprintf("http://localhost:%d", *localPort)}
	} else {
		cfg, err = kclient.InClusterConfig()
		if err != nil {
			glog.Errorf("failed to load config: %v", err)
			os.Exit(1)
		}
	}

	client, err = kclient.New(cfg)
	// fmt.Println("%v", nodes.Items[0].NodeAddress.Address)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRequest))
	http.ListenAndServe(":5000", mux)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	nodes, err := client.Nodes().List(
		labels.Everything(), fields.Everything())

	if err != nil {
		fmt.Printf("failed to get node list %v", err)
	}
	fmt.Fprint(w, "Getting node list:\n")
	for _, n := range nodes.Items {
		fmt.Fprintf(w, "%s\n", getExternalIP(n.Status.Addresses))
	}

}

func getExternalIP(addr []api.NodeAddress) string {
	for _, el := range addr {
		if el.Type == "ExternalIP" {
			return el.Address
		}
	}
	return ""
}
