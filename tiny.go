package sensitivecheck

import (
    "fmt"
    "strings"
    "net/http"
)

type TinyMux struct{
    hpatt map[string] http.HandlerFunc
    spatt map[string] string
}


func (t *TinyMux) Listen(port string) {
    print("Listening: %s\n", port[1:])
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

    for path, dir := range t.spatt {
        if strings.HasPrefix(r.URL.Path, path) {
            file := dir+ "/" + r.URL.Path[len(path):]
            fmt.Println(file)
            http.ServeFile(w, r, file)
            proced = 1
            break
        }
    }


    if(0 == proced) {
        t.FunctionNotFount(w, r)
    }

    cLog.Notice(r.RemoteAddr+"\t"+r.RequestURI)
}

func (t *TinyMux)Add(patt string, handler http.HandlerFunc) {
    t.hpatt[patt] = handler
}

/*//利用默认路由
func (t *TinyMux)AddStatic(patt string, handler http.Handler) {
    t.hpatt[patt] = func(w http.ResponseWriter, r *http.Request) {
        r.URL.Path = "a.html"
        handler.ServeHTTP(w, r)
        }
}
*/

func (t *TinyMux) AddStatic(path string, dir string) {
    t.spatt[path] = dir
}

func Newhttpframe() *TinyMux {
    cLog.Notice("this is in http framwork")
    return &TinyMux{make(map[string] http.HandlerFunc), make(map[string] string)}
}

func (t *TinyMux) FunctionNotFount(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("No Action Has Found")
    w.WriteHeader(404)
    err_str := "{\"Errno\":-2,\"Errmsg\":\"error path\",\"Data\":\"0\"}"
    w.Write([]byte(err_str))
}
