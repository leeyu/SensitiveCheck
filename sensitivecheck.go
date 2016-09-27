package sensitivecheck

import(
"fmt"
"bufio"
"os"
"strings"
"net/http"
"encoding/json"
"strconv"
/*
"bytes"
*/
)

type node struct{
    fail *node
    nc rune
    sons map[rune] *node
    count int
    }

type resp struct{
    Errno int
    Errmsg string
    Data string
    }

type ColorGroup struct {
    ID     int
    Name   string
    Colors []string
}

var root *node

var errres = resp{
    Errno: -1,
    Errmsg: "errno happened",
    Data: "0",
    }

func (root *node) Scheck(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    ct := r.Form["ct"]

    if nil == ct {
        ret := echojson(errres)
        fmt.Fprintf(w, "%s", ret)
        return
        }

    //    blacklist.GetRoot()
    //    black_root := blacklist.GetRoot()
    content := ct[0]
    cLog.Notice("in scheck" + ct[0])

    if "" == content {
        ret := echojson(errres)
        fmt.Fprintf(w, "%s", ret)
        return
        }

    cnt := root.Query(content)

    res := resp{
        Errno: 0,
        Errmsg: "ok",
        Data: strconv.Itoa(cnt),
        }

    ret := echojson(res)
    fmt.Fprintf(w, "%s", ret)
//    fmt.Printf("this is :%s", string(ret), res.errmsg, res.errno)
//    fmt.Fprintf(w, "Catch: %s", cnt)
}


func echojson(res resp) []byte {
    ret,err := json.Marshal(res)
    if err != nil {
        ret,err = json.Marshal(errres)
    }

    return ret
}


//建立tries树
func insert(str string, root *node) {
    cur := root
    str_rune := []rune(str)
    for i:=0; i<len(str_rune); i++ {
        chld, ok := cur.sons[str_rune[i]]
        if ok {
            cur = chld
        } else {
            if nil == cur.sons {
                cur.sons = make(map[rune] *node, 1000)
                }

                n_node := new(node)
                cur.sons[str_rune[i]] = n_node
                n_node.nc = str_rune[i]

                cur = n_node
            }
        }

        cur.count++;
//        fmt.Printf("rune type : %c   %d\n", cur.nc, cur.count)
    }

//建立回朔路径
func buidAcAm(root *node) {
    var queue  [500000]*node
    var fail_n *node

    root.fail = nil
    head := 0
    tail := 0
    queue[head] = root
    head++
    for {
        if head <= tail {
            break
            }
        tmp := queue[tail]
        tail++

        if(nil != tmp && nil != tmp.sons) {
            for ch,t_node := range tmp.sons {

                if tmp == root {
                    t_node.fail = root
//                    fmt.Printf("rune fail root: %c\n", t_node.nc)
                } else {
                    fail_n = tmp.fail //父节点的fail
                    for {
                        if(nil == fail_n) {
                            t_node.fail = root
                            break
                            }
                        //子节点的fail，设置为父节点fail的相同自节点
                        if(nil != fail_n.sons[ch]) {
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
func (root *node) Query(str string) int {
    cur := root
    str_rune := []rune(str)
    var temp *node
    count := 0

    for i:=0; i<len(str_rune); i++ {
        if nil == cur {
            cur = root
            }

        //如果当前字符不匹配节点的子节点,则递归fail指针
        if nil == cur.sons[str_rune[i]] {
            for {
                if cur == nil || cur == root || nil != cur.sons[str_rune[i]] {
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
                if temp == root || nil == temp {
                    break
                    }

                count += temp.count
                temp = temp.fail
                }
            }
        }

        return count
    }

//逐行读取文件
func cat(scanner *bufio.Scanner, root *node) error{

    for scanner.Scan(){
//        fmt.Println(scanner.Text())
        row := scanner.Text()
        insert(strings.TrimSpace(row), root)
        /*
        words := strings.Split(strings.TrimSpace(row), " ")
        for i:=0; i<len(words); i++ {
            if("" != strings.TrimSpace(words[i])) {
                insert(strings.TrimSpace(words[i]), root)
                }
            fmt.Println("aa",words[i])    
        }
        */
    }

    return scanner.Err()
}

func loadFile(fileName string, root *node) {
    f,err := os.OpenFile(fileName,os.O_RDONLY,0660)
    if err != nil{
        fmt.Fprintf(os.Stderr,"%s err read from %s : %s\n",
        fileName,fileName,err)
    }

    cat(bufio.NewScanner(f), root)
    f.Close()
    }

func GetRoot() *node {
    return root
        }

func New() *node {
//    source := []byte("我们俩吃饭是敏感词 他我们 你我们仨 我们仨 我们仨")

    root := new(node)
    file := string("words.txt")
    loadFile(file, root)


     buidAcAm(root)

     return root
     /*
     content := "说不能出现砍死你等恶意词汇"
     cnt := query(content, root)

     fmt.Printf("count :%d\n", cnt)
     */
    }
