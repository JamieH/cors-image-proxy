package main
import (
        "io/ioutil"
        "log"
        "net/http"
        "strings"
        "fmt"
        "google.golang.org/appengine"
        "google.golang.org/appengine/urlfetch"
)

const CDNURL = "web.poecdn.com/"

func handle(w http.ResponseWriter, r *http.Request) {
        setDefaultHeaders(w)
        ctx := appengine.NewContext(r)
        client := urlfetch.Client(ctx)
        url := r.URL.Path
        url = strings.TrimLeft(url, "/")
        wecandoit := strings.HasPrefix(url, CDNURL)
        if (wecandoit == false) {
                log.Printf("refusing to proxy: %v", url)
                http.Error(w, "bad url. no proxy4u", http.StatusBadRequest)
                return
        }

        url = fmt.Sprintf("https://%s", url)
        log.Println(url)

        res, err := client.Get(url)
        if err != nil {
                log.Printf("Error getting URL: %v", err)
                http.Error(w, "can't get url soz fam", http.StatusBadRequest)
                return
        }

        if res.Body == nil {
                http.Error(w, "no response fam", http.StatusBadRequest)
                return
        }

        defer res.Body.Close()

        body, err := ioutil.ReadAll(res.Body)
        bodyString := string(body)

        if err != nil {
                log.Printf("Error reading body: %v", err)
                http.Error(w, "can't read body soz fam", http.StatusBadRequest)
                return
        }

        fmt.Fprintln(w, bodyString)
}

func setDefaultHeaders(w http.ResponseWriter) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET")
        w.Header().Set("Vary", "Accept-Encoding")
}

func main() {
        http.HandleFunc("/", handle)
        appengine.Main()
}
