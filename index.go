package main

import(
    "net/http"
    "github.com/martini-contrib/render"
)

type IndexViewModel struct {
    Days []string
    Facilities []*Facility
}

func IndexRender(req *http.Request, ren render.Render) {

    var site = newElevenAuto()
    facilities := site.ScrapeHtml(req)

    days := []string{}
    for i := 0; i < len(facilities[0].Availabilities); i++ {
        days = append(days, facilities[0].Availabilities[i].Date.Format("01/02"))
    }

    viewModel := IndexViewModel{
        Days : days,
        Facilities : facilities,
    }

    ren.HTML(200, "index", viewModel)
}

