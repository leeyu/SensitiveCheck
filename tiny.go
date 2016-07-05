package sensitivecheck

import (
    "fmt"
    "net/http"
)

type TinyMux struct{
    hpatt map[string] http.HandlerFunc
}

func (t *TinyMux) Listen(port string) {
    fmt.Printf("Listening: %s\n", port[1:])
    http.ListenAndServe(port, t)
}

func (t *TinyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    proced := 0
    for UrlPath, handler := range t.hpatt {
        if(r.URL.Path == UrlPath) {
            handler.ServeHTTP(w, r)
            proced = 1
            break
        }
    }
    if(0 == proced) {
        t.FunctionNotFount(w, r)
    }
}

func (t *TinyMux)Add(patt string, handler http.HandlerFunc) {
    t.hpatt[patt] = handler
}

func Newhttpframe() *TinyMux {
    return &TinyMux{make(map[string] http.HandlerFunc)}
}

func (t *TinyMux) FunctionNotFount(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("No Action Has Found")
    w.WriteHeader(404)
    w.Write([]byte("Nothing to see here"))
}
