package plugin

import (
	"fmt"
	"net/http"
)

type PluginInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
	AddNext(PluginInterface)
}

type Plugin2 struct {
	NextPlugin PluginInterface
}

func (p Plugin2) Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("plugin2 executed")
	if p.NextPlugin != nil {
		p.NextPlugin.Handle(w, r)
	}

}

func (p Plugin2) AddNext(pl PluginInterface) {
	p.NextPlugin = pl
}
