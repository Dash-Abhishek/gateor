package internal

import (
	"fmt"
	"gateor/pkg"
	"gateor/plugin"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type service struct {
	Name          string `yaml:"name"`
	Path          string `yaml:"basepath"`
	StripBasePath bool   `yaml:"stripBasepath"`
	Target        target `yaml:"target"`
	RateLimit     int    `yaml:"rateLimit"`
	PluginChain   plugin.PluginInterface
}

type target struct {
	Host    string `yaml:"host"`
	Timeout int    `yaml:"timeout"`
}

func LoadSvc() {

	pkg.Log.Info("Loading services")
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
				// add plugin chain based on service configs
				pl := plugin.NewLeakyBucketRateLimit(svc.RateLimit, 10)
				pl.AddNext(plugin.Plugin2{})
				svc.PluginChain = pl
				services = append(services, svc)
			}

		}
	}

	fmt.Printf("services loaded: %v\n", len(services))

	mux := InitializeMux()
	for _, svc := range services {
		mux.Handle(svc.Path+"/", svc)
	}

}
