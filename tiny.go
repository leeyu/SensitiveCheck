package sensitivecheck

import (
    "fmt"
    "strings"
    "net/http"
    "net/url"
    "regexp"
)

type TinyMux struct{
    regex *regexp.Regexp
    params map[int] string
    handle http.HandlerFunc
}

type CMux struct{
    routers []*TinyMux
    staticdir map[string] string
}


func (t *CMux) Listen(port string) {
    print("Listening: %s\n", port[1:])
    http.ListenAndServe(port, t)
}

func (c *CMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    proced := 0
    for _, router := range c.routers {
        if !router.regex.MatchString(r.URL.Path) {
            continue;
        }

        matches := router.regex.FindStringSubmatch(r.URL.Path)
        if len(matches[0]) != len(r.URL.Path) {
            continue;
        }

        params := make(map[string]string)

        if len(router.params) > 0 {
            values := r.URL.Query()
            for i,match := range matches[1:] {
                values.Add(router.params[i], match)
                params[router.params[i]] = match
            }

            r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
            cLog.Notice(r.URL.RawQuery)
            fmt.Printf("%s", params)
        }

        router.handle.ServeHTTP(w, r)
        proced = 1
        break
/*
        if(r.URL.Path == UrlPath) {
            handler.ServeHTTP(w, r)
            proced = 1
            break
        }
        */
    }

    if proced == 0 {
        for path, dir := range c.staticdir {
            if strings.HasPrefix(r.URL.Path, path) {
                file := dir+ "/" + r.URL.Path[len(path):]
                fmt.Println(file)
                http.ServeFile(w, r, file)
                proced = 1
                break
            }
        }
    }


    if(0 == proced) {
        c.FunctionNotFount(w, r)
    }

    cLog.Notice(r.RemoteAddr+"\t"+r.RequestURI)
}

func (t *CMux)Add(patt string, handler http.HandlerFunc) {
//    t.hpatt[patt] = handler

    parts := strings.Split(patt, "/")

    j := 0
    params := make(map[int] string)
    for i,part := range parts {
        if strings.HasPrefix(part, ":") {
            start_index := 1
            expr := "([^/]+)"

            if index := strings.Index(part, "("); index != -1 {
                expr = part[index:]
                part = part[:index]
                }
            params[j] = part[start_index:]
            parts[i] = expr
            j++
        }
    }

    patt = strings.Join(parts, "/")
    regex, regexErr := regexp.Compile(patt)
    if regexErr != nil {
        panic(regexErr)
        return
    }

    route := &TinyMux{}
    route.regex = regex
    route.params = params
    route.handle = handler


    t.routers = append(t.routers, route)
    fmt.Printf("r: %s\n", t.routers)
}

/*//利用默认路由
func (t *TinyMux)AddStatic(patt string, handler http.Handler) {
    t.hpatt[patt] = func(w http.ResponseWriter, r *http.Request) {
        r.URL.Path = "a.html"
        handler.ServeHTTP(w, r)
        }
}
*/

func (t *CMux) AddStatic(path string, dir string) {
    if 0 == len(t.staticdir) {
        t.staticdir = make(map[string]string)
    }
    t.staticdir[path] = dir
}

func Newhttpframe() *CMux {
    cLog.Notice("this is in http framwork")
    return &CMux{}
//    return &TinyMux{ *regexp.Regexp, make(map[int] string), http.HandlerFunc, make(map[string] http.HandlerFunc), make(map[string] string)}
}

func NewRoute() *CMux {
    cLog.Notice("this is in http framwork")
    return &CMux{}
}

func (t *CMux) FunctionNotFount(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("No Action Has Found")
    w.WriteHeader(404)
    err_str := "{\"Errno\":-2,\"Errmsg\":\"error path\",\"Data\":\"0\"}"
    w.Write([]byte(err_str))
}
