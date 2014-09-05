package main

import(
    "net/http"
    "time"
    "github.com/martini-contrib/render"
    "appengine"
    "appengine/datastore"
)

type Tweet struct {
    Tweet string
    CreatedAt time.Time
}

type IndexViewModel struct {
    Title string
    Tweets []string
    Now time.Time
}

func IndexRender(r render.Render, req *http.Request) {

    c := appengine.NewContext(req)
    q := datastore.NewQuery("tweet").Order("-CreatedAt")

    temp := []string{}
    var ts []Tweet
    q.GetAll(c, &ts)

    for i := 0; i < len(ts); i++ {
        temp = append(temp, ts[i].Tweet)
    }

    viewModel := IndexViewModel{
        Title : "Martini Datastore Demo",
        Tweets : temp,
        Now : time.Now(),
    }

    r.HTML(200, "index", viewModel)
}

func IndexPostHandler(r render.Render, req *http.Request) {

    err := req.ParseForm()

    _tweet := req.PostFormValue("tweet")

    tweet := Tweet{
        Tweet : _tweet,
        CreatedAt : time.Now(),
    }

    c := appengine.NewContext(req)

    key := datastore.NewIncompleteKey(c, "tweet", nil)

    if key, err = datastore.Put(c, key, &tweet); err != nil {
        return
    }
    _ = err // これがないと defined and not use で error になる

    r.Redirect("/")
}

