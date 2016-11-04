package sensitivecheck

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type TinyMux struct {
	regex      *regexp.Regexp
	params     map[int]string
	handle     Handle
	controller reflect.Type
	//	handle http.HandlerFunc
}

type CMux struct {
	routers   []*TinyMux
	staticdir map[string]string
}

type Handle func(ctx *Context)

func (t *CMux) Listen(port string) {
	print("Listening: %s\n", port[1:])
	http.ListenAndServe(port, t)
}

func (c *CMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proced := 0
	path := "/"
	pathB := "/"
	for _, router := range c.routers {

		if "/" != r.URL.Path {
			path = strings.ToLower(strings.Trim(r.URL.Path, "/"))
			pathB = path + "/"
		}

		fmt.Printf("\n %s %s %s\n", path, pathB, router.regex)
		if !router.regex.MatchString(path) {
			if !router.regex.MatchString(pathB) {
				continue
			}

			path = pathB
		}

		matches := router.regex.FindStringSubmatch(path)
		if len(matches[0]) != len(path) {
			continue
		}

		params := make(map[string]string)

		if len(router.params) > 0 {
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(router.params[i], match)
				params[router.params[i]] = match
			}

			r.URL.RawQuery = url.Values(values).Encode()
			cLog.Notice(r.URL.RawQuery)
		}

		ctx := NewContext(params, w, r)

		Action := "Index"
		if len(ctx.Request.Form["a"]) > 0 {
			Action = strings.Title(ctx.Request.Form["a"][0])
		}

		rc := reflect.New(router.controller)

		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(ctx)
		in[1] = reflect.ValueOf(router.controller.Name())
		InitFun := rc.MethodByName("Init")
		InitFun.Call(in)

		ActionFun := rc.MethodByName(Action)
		if !ActionFun.IsValid() {
			fmt.Printf("error, %s \n", ActionFun)
			break
		}
		inp := make([]reflect.Value, 0)
		ActionFun.Call(inp)

		proced = 1
		break
	}

	if proced == 0 {
		for path, dir := range c.staticdir {
			if strings.HasPrefix(r.URL.Path, path) {
				file := dir + "/" + r.URL.Path[len(path):]
				fmt.Println(file)
				http.ServeFile(w, r, file)
				proced = 1
				break
			}
		}
	}

	if 0 == proced {
		c.FunctionNotFount(w, r)
	}

	cLog.Notice(r.RemoteAddr + "\t" + r.RequestURI)
}

//func (t *CMux) Add(patt string, handler http.HandlerFunc) {
func (t *CMux) Add(patt string, ci ControllerInterface) {
	//    t.hpatt[patt] = handler

	parts := strings.Split(patt, "/")

	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		//action的位置不能有非正则表达
		if i == 2 && part != "" && (!strings.HasPrefix(part, ":")) {
			fmt.Println("注册%s有错误!Action位置不能有非正则表达式哦!", patt)
			os.Exit(1)
		}

		if strings.HasPrefix(part, ":") {
			start_index := 1
			expr := "([^/]+)"
			expr = "([^/].)"

			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part[start_index:]
			parts[i] = expr
			j++

		}
	}

	patt = strings.Trim(strings.Join(parts, "/"), "/")
	if "" == patt {
		patt = "/"
	}

	fmt.Printf("patt: %s\n", patt)
	regex, regexErr := regexp.Compile(strings.ToLower(patt))
	if regexErr != nil {
		panic(regexErr)
		return
	}

	cv := reflect.Indirect(reflect.ValueOf(ci)).Type()
	route := &TinyMux{}
	route.regex = regex
	route.params = params
	route.controller = cv

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
