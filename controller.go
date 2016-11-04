package sensitivecheck

import (
	"encoding/json"
	"fmt"
)

type ControllerInterface interface {
	Init(ct *Context, cn string)
	Prepare()
	Preaccess()
	Postaccess()
	Finish()
	Render()
}

type Resp struct {
	Errno  int
	Errmsg string
	Data   string
}

type Controller struct {
	//    Tpl *template.Template
	Context  *Context
	TplNames string
	Errsys   Resp
}

func (c *Controller) Init(ct *Context, cn string) {
	fmt.Println("\nN:this is in base controller!\n")
	c.TplNames = ""
	c.Context = ct
	c.Errsys = Resp{
		Errno:  -1,
		Errmsg: "errno happened",
		Data:   "0",
	}
}

func (c *Controller) Prepare() {
}

func (c *Controller) Preaccess() {
}

func (c *Controller) Postaccess() {
}

func (c *Controller) Finish() {
}

func (c *Controller) Render() {
}

func (c *Controller) Echojson(res interface{}) []byte {
	ret, err := json.Marshal(res)
	if err != nil {
		ret, err = json.Marshal(errres)
	}

	return ret
}
