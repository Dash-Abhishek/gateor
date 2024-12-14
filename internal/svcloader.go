package internal

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type service struct {
	Name          string `yaml:"name"`
	Path          string `yaml:"basepath"`
	StripBasePath bool   `yaml:"stripBasepath"`
	Target        target `yaml:"target"`
}

type target struct {
	Host    string `yaml:"host"`
	Timeout int    `yaml:"timeout"`
}

func LoadSvc() {

	entry, err := os.ReadDir("services")
	if err != nil {
		log.Default().Println("Error reading services directory", err)
	}

	services := []service{}
	for _, file := range entry {
		if file.Type().IsRegular() {
			bytes, err := os.ReadFile("services/" + file.Name())
			if err != nil {
				log.Default().Printf("Error reading file %s : %v\n", file.Name(), err)
				continue
			}
			svc := service{}
			if err = yaml.Unmarshal(bytes, &svc); err != nil {
				log.Default().Printf("Error unmarshalling file %s : %v\n", file.Name(), err)
				continue
			} else {
				services = append(services, svc)
			}

		}
	}

	fmt.Printf("services loaded: %v", len(services))

	mux := InitializeMux()
	for _, svc := range services {
		mux.Handle(svc.Path+"/", svc)
	}

}
