package main

import (
    "log"
    "bytes"
    "time"
    "strconv"
    "net/http"
    "io/ioutil"
    "code.google.com/p/go.text/encoding"
    "code.google.com/p/go.text/encoding/japanese"
    "code.google.com/p/go.text/transform"
    "github.com/PuerkitoBio/goquery"
    "appengine"
    "appengine/urlfetch"
)

/**
* インタフェース
**/
type CampingSite interface {
    ScrapeHtml(*http.Request) bool
}

type Facility struct {
    Name string
    Availabilities []*Availability
}

func newFacility(name string) *Facility{
    facility := new(Facility)
    facility.Name = name
    return facility
}

type Availability struct {
    Date time.Time
    Condition string
}

func newAvailability(date time.Time, condition string) *Availability{
    availability := new(Availability)
    availability.Date = date
    availability.Condition = condition
    return availability
}

/**
* 構造体
* イレブンオート
**/
type ElevenAuto struct {
    Name string
    Url string
    Tel string
    Address string
}

func newElevenAuto() *ElevenAuto{
    site := new(ElevenAuto)
    site.Name = "イレブンオート"
    site.Url = "http://www.camp-net.jp/site/status.cgi/site?camp=132"
    site.Tel = "ssss"
    site.Address = "aaaaaa"
    return site
}

func (site *ElevenAuto) ScrapeHtml(req *http.Request) []*Facility {

    context := appengine.NewContext(req)
    client := urlfetch.Client(context)
    resp, _ := client.Get(site.Url)

    defer resp.Body.Close()

    _days := [14]int{}
    _facilities := make(map[int]string)
    _availabilities := make(map[int][]string)

    doc, _ := goquery.NewDocumentFromResponse(resp)
    doc.Find(".table-month").Each(func(_ int, s *goquery.Selection) {
        s.Find("tr").Each(func(tr_idx int, s *goquery.Selection){
            s.Find("td").Each(func(td_idx int, s *goquery.Selection){

                utf8, _ := Decode(japanese.ShiftJIS, s.Text())
                //log.Println(tr_idx, td_idx, utf8)

                // 日付
                if tr_idx == 0 {
                    i, _:= strconv.Atoi(utf8)
                    _days[td_idx] = i
                }
                // ファシリティの名称
                if tr_idx >= 3 && td_idx == 0 {
                    _facilities[tr_idx] = utf8
                }
                // 空き状況
                if tr_idx >= 3 && td_idx > 0 {
                    _availabilities[tr_idx] = append(_availabilities[tr_idx], utf8)
                }
            })
        })
     })

    jst := time.FixedZone("Asia/Tokyo", 9*60*60)
    today := time.Now()
    today = today.UTC().In(jst)

    today = time.Date(today.Year(), today.Month(), _days[0], 0, 0, 0, 0, jst)
    _times := []time.Time{}
    for i := 0; i < 14; i ++ {
        _times = append(_times, today.AddDate(0, 0, i))
    }

    var facilities []*Facility = make([]*Facility, 0)

    for key := range _facilities {
        facility := newFacility(_facilities[key])
        facilities = append(facilities, facility)

        for i := 0; i < len(_availabilities[key]); i++ {
            availability := newAvailability(_times[i], _availabilities[key][i])
            facility.Availabilities = append(facility.Availabilities, availability)
        }
    }
    log.Println(facilities)

    return facilities
}

func Decode(enc encoding.Encoding, s string) (string, error) {
    r := bytes.NewBuffer([]byte(s))
    decoded, err := ioutil.ReadAll(transform.NewReader(r, enc.NewDecoder()))
    return string(decoded), err
}


