package main

import(
    "net/http"  // appengineでは必要
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
)

func init() { // appengineではmainではなくinit

    m := martini.Classic()

    m.Use(render.Renderer(render.Options{
        Directory: "templates",
        Layout: "layout",
        Extensions: []string{".tmpl"},
        Charset: "utf-8",
    }))

    m.NotFound(func (r render.Render){
        r.Redirect("/")
    })

    m.Get("/", IndexRender)
    m.Post("/", IndexPostHandler)

    //m.Run() // appengineでは↓http.Handle("/", m)
    http.Handle("/", m)

}

