package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type service struct {
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	Target       string `yaml:"target"`
	ReqProcessor ReqProcessorInterface
}

type DefaultReqProcessor struct{}

func (d *DefaultReqProcessor) PreFlow(w http.ResponseWriter, r *http.Request) {
	fmt.Println("executing default preflow")
}

func (d *DefaultReqProcessor) PostFlow(q http.ResponseWriter, r *http.Request) {
	fmt.Println("executing default postflow")
}

func (svc service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if svc.ReqProcessor != nil {
		svc.ReqProcessor.PreFlow(w, r)
	}

	targetUrl, err := url.Parse(svc.Target)
	fmt.Println("request", r.URL.Path, r.UserAgent(), r.Header)
	fmt.Println("targetUrl", targetUrl)
	if err != nil {
		http.Error(w, "Invalid target url", http.StatusBadGateway)
		return
	}

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.ServeHTTP(w, r)

	if svc.ReqProcessor != nil {
		svc.ReqProcessor.PostFlow(w, r)
	}

}

func LoadServices() []service {

	entry, err := os.ReadDir("services")
	if err != nil {
		log.Fatal(err)
	}

	services := []service{}
	for _, file := range entry {
		if file.Type().IsRegular() {
			fmt.Println("This is a regular file:", file.Name())
			bytes, err := os.ReadFile("services/" + file.Name())
			if err != nil {
				fmt.Printf("Error reading file %s : %v", file.Name(), err)
				continue
			}
			svc := service{}
			if err = yaml.Unmarshal(bytes, &svc); err != nil {
				fmt.Printf("Error unmarshalling file %s : %v", file.Name(), err)
				continue
			} else {
				services = append(services, svc)
			}

		}
	}

	fmt.Printf("%+v", services)
	return services

}
