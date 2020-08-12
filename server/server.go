package server

import (
	"log"

	"bytes"
	"github.com/CloudyKit/jet"
	"github.com/savsgio/atreugo"
)

var server *atreugo.Atreugo

var markupSet *jet.Set

type Page struct {
	template *jet.Template
	vars     jet.VarMap
}

func get_markup(name string) Page {
	t, err := markupSet.GetTemplate(name)
	if err != nil {
		panic(err)
	}

	v := make(jet.VarMap)
	v.Set("skin", "skin")

	return Page{t, v}
}

func InitTemplates() {
	markupSet = jet.NewHTMLSet("./views")
	markupSet.SetDevelopmentMode(true)
}

func InitRoutes() {

	server = atreugo.New(atreugo.Config{
		Addr: "0.0.0.0:8086",
	})

	server.GET("/", _Get_Index)
	server.GET("/users/{id}", _Get_User)

	server.GET("/qrc/{id}.png", _Get_QR)

	server.POST("/users", _Post_User)
	server.POST("/validate", _Post_Validate)

}

func Run() {
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func render_markup(ctx *atreugo.RequestCtx, page Page) error {

	w := bytes.Buffer{}

	if err := page.template.Execute(&w, page.vars, nil); err != nil {
		panic(err)
	}
	return ctx.HTTPResponse(w.String(), 200)
}

func he(err error) {
	if err != nil {
		panic(err)
	}
}
