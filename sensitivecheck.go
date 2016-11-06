package sensitivecheck

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	//	"strconv"
	"strings"
	/*
	   "bytes"
	*/)

type node struct {
	fail  *node
	nc    rune
	sons  map[rune]*node
	count int
	word  []rune
}

type resp struct {
	Errno  int
	Errmsg string
	Data   string
}

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
}

var Root *node

var errres = resp{
	Errno:  -1,
	Errmsg: "errno happened",
	Data:   "0",
}

func (Root *node) Scheck(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ct := r.Form["ct"]

	if nil == ct {
		ret := echojson(errres)
		fmt.Fprintf(w, "Catch: %s", ret)
		return
	}

	//    blacklist.GetRoot()
	//    black_Root := blacklist.GetRoot()
	content := ct[0]
	cLog.Notice("in scheck" + ct[0])

	if "" == content {
		ret := echojson(errres)
		fmt.Fprintf(w, "Catch: %s", ret)
		return
	}

	cnt := Root.Query(content)

	res := resp{
		Errno:  0,
		Errmsg: "ok",
		Data:   cnt,
	}

	ret := echojson(res)
	//    fmt.Printf("this is :%s", string(ret), res.errmsg, res.errno)
	fmt.Fprintf(w, "Catch: %s", ret)
}

func echojson(res resp) []byte {
	ret, err := json.Marshal(res)
	if err != nil {
		ret, err = json.Marshal(errres)
	}

	return ret
}

//建立tries树
func insert(str string, Root *node) {
	cur := Root
	str_rune := []rune(str)
	for i := 0; i < len(str_rune); i++ {
		chld, ok := cur.sons[str_rune[i]]
		if ok {
			cur = chld
		} else {
			if nil == cur.sons {
				//cur.sons = make(map[rune]*node, 1000)
				cur.sons = make(map[rune]*node)
			}

			n_node := new(node)
			cur.sons[str_rune[i]] = n_node
			n_node.nc = str_rune[i]

			cur = n_node
		}
	}

	cur.count++
	cur.word = str_rune
	//	fmt.Printf("rune type : %c   %s\n", cur.nc, str)
	//        fmt.Printf("rune type : %c   %d\n", cur.nc, cur.count)
}

//建立回朔路径
func buidAcAm(Root *node) {
	var queue [500000]*node
	var fail_n *node

	Root.fail = nil
	head := 0
	tail := 0
	queue[head] = Root
	head++
	for {
		if head <= tail {
			break
		}
		tmp := queue[tail]
		tail++

		if nil != tmp && nil != tmp.sons {
			for ch, t_node := range tmp.sons {

				if tmp == Root {
					t_node.fail = Root
					//                    fmt.Printf("rune fail Root: %c\n", t_node.nc)
				} else {
					fail_n = tmp.fail //父节点的fail
					for {
						if nil == fail_n {
							t_node.fail = Root
							break
						}
						//子节点的fail，设置为父节点fail的相同自节点
						if nil != fail_n.sons[ch] {
							t_node.fail = fail_n.sons[ch]
							////////////////                            fmt.Printf("rune fail dup: %c\n", ch)
							break
						}
						fail_n = fail_n.fail
					}
				}

				//把子节点加入到处理队列中
				queue[head] = t_node
				head++
			}
		}
	}

}

//匹配tries树，是否命中敏感词
func (Root *node) Query(str string) string {
	cur := Root
	str_rune := []rune(str)
	var temp *node
	count := 0
	var res [][]rune
	had := 0

	for i := 0; i < len(str_rune); i++ {
		if nil == cur {
			cur = Root
		}

		//如果当前字符不匹配节点的子节点,则递归fail指针
		if nil == cur.sons[str_rune[i]] {
			for {
				if cur == nil || cur == Root || nil != cur.sons[str_rune[i]] {
					break
				}

				cur = cur.fail
			}
		}

		//     fmt.Printf("debug2 :%c %c\n", str_rune[i], cur)

		//计算count
		if nil != cur.sons[str_rune[i]] {
			cur = cur.sons[str_rune[i]]
		}

		//     fmt.Printf("debug1 :%c\n", cur.nc)

		temp = cur
		if nil != temp {
			for {
				if temp == Root || nil == temp {
					break
				}

				count += temp.count
				if len(temp.word) > 0 {
					had = 0

					if len(res) > 0 {
						for _, stword := range res {
							if string(stword) == string(temp.word) {
								had = 1
							}
						}
					}

					if 0 == had {
						res = append(res, temp.word)
					}
				}

				temp = temp.fail
			}
		}
	}

	restr := ""
	var strl []string
	if len(res) > 0 {
		for _, word := range res {
			strl = append(strl, string(word))
		}
		restp, err := json.Marshal(strl)
		if err != nil {
			restp, err = json.Marshal(errres)
		}
		restr = string(restp)

	}

	return restr
}

//逐行读取文件
func cat(scanner *bufio.Scanner, Root *node) error {

	for scanner.Scan() {
		//        fmt.Println(scanner.Text())
		row := scanner.Text()
		insert(strings.TrimSpace(row), Root)
		/*
		   words := strings.Split(strings.TrimSpace(row), " ")
		   for i:=0; i<len(words); i++ {
		       if("" != strings.TrimSpace(words[i])) {
		           insert(strings.TrimSpace(words[i]), Root)
		           }
		       fmt.Println("aa",words[i])
		   }
		*/
	}

	return scanner.Err()
}

func loadFile(fileName string, Root *node) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0660)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s err read from %s : %s\n",
			fileName, fileName, err)
	}

	cat(bufio.NewScanner(f), Root)
	f.Close()
}

func GetRoot() *node {
	return Root
}

func New() *node {
	//    source := []byte("我们俩吃饭是敏感词 他我们 你我们仨 我们仨 我们仨")

	Root = new(node)
	file := string("words.txt")
	loadFile(file, Root)

	buidAcAm(Root)

	return Root
	/*
		   content := "说不能出现砍死你等恶意词汇"
		   cnt := query(content, Root)

		   fmt.Printf("count :%d\n", cnt)
	i/e	*/
}
