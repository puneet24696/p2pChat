package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strings"
    "net"
    "encoding/json"
    "io/ioutil"
//    "github.com/pkg/browser"
)

var db = make(map[string]string)

func Welcome(w http.ResponseWriter, r *http.Request) {
    r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
    // attention: If you do not call ParseForm method, the following data can not be obtained form
    fmt.Println(r.Form) // print information on server side.
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello User! \n Go to url/login to validate") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //get request method
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        t.Execute(w, nil)
    } else{
        r.ParseForm()
        if r.Form["username"][0] != ""{
	    fmt.Println(r)
            ip, port, _ := net.SplitHostPort(r.RemoteAddr)
//	    fmt.Printf("%v:%v \n",ip,port)
	    port = r.Form["port"][0]
	    ip_port := ip+":"+port
	    
	    str_usr := r.Form["username"][0]
            for a,b := range(db){
                if str_usr==a{
                    fmt.Fprintf(w,"%v","Username already taken");return
                }else if ip_port == b{
                    fmt.Fprintf(w,"%v","You are already online");return
                }else {continue}
        }
/*        if str_usr != "" {
            client := &http.Client{}
            req,err := http.NewRequest("GET", "/chatroom", nil)
                if err != nil{
                    fmt.Println(err)
                    return
                }
            client.Do(req)
        }else {return}*/

        db[str_usr] = ip_port
	fmt.Fprintf(w,"%v",db)
	/*        
	// logic part of log in
        fmt.Printf("username:%v \n",str_usr)
        //fmt.Println("password:", r.Form["password"])
	fmt.Println("ip_port:", ip_port)        
        */
	}else{fmt.Fprintf(w,"Bad username")}
    }
}

func chatroom(w http.ResponseWriter, r *http.Request){
    r.ParseForm()
    fmt.Println("method:", r.Method)
    if r.Method == "POST"{

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {fmt.Println(err,"Error while reading")}
	//fmt.Printf("%T %v",r.Body,body)
	ek_map := make(map[string]string)
	err = json.Unmarshal(body,&ek_map)
	if err!= nil{fmt.Println(err,"error in unmarshalling")}
	name := ek_map["name"] 
	if ek_map["ip_port"] == db[name]{
	    fmt.Fprintf(w,"%v %v", name,"has been disconnected")
	    delete(db,name)
	}
    }
    for a,b := range(db){
        fmt.Fprintf(w, "%v : %v\n",a,b)
    }
    fmt.Println(r)
}


func main() {
    http.HandleFunc("/", Welcome) // setting router rule
    http.HandleFunc("/login", login)
    http.HandleFunc("/chatroom", chatroom)
    err := http.ListenAndServe(":9090", nil) // setting listening port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    fmt.Println("go worked")
}
