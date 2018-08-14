package main

import (
    "log"
	"bufio";"os"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
	"net"
    "crypto/tls"
)

var CLIENTS = make(map[string]string)

func get_db(){
    req, err := http.NewRequest("GET","http://192.168.0.84:9090/chatroom",nil)
    if err != nil{fmt.Println(err)}
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil{fmt.Println(err)}
    text, err := ioutil.ReadAll(resp.Body)
    temp:= string(text)
    temp_ar := strings.Split(temp,"\n")
    temp_ar = temp_ar[:len(temp_ar)-1]
    for _,a := range(temp_ar){
        temp_arr := strings.SplitN(a," : ",2)
        CLIENTS[temp_arr[0]]=temp_arr[1]
    }
    for a,b := range(CLIENTS){fmt.Println(a,"is online with ip",b)}
    return
}

func to_write(conn net.Conn){
    for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(conn,"->")
		a,_ := reader.ReadString('n')
        w_mess, err := conn.Write([]byte(a))
        if err!= nil{
            fmt.Println(w_mess,err)
            conn.Close()
            return
       } 
    }
    return
}

func to_read(conn net.Conn){
    for{
        message := make([]byte, 100)
		text, err := conn.Read(message)
        if err != nil {
            fmt.Println(text, err)
            conn.Close()
            return
		}
    }
    return
}


func client_init_chat(addr,s string){
    config:= &tls.Config{
	}
	conn, err := tls.Dial("tcp",addr,config)
    if err != nil {
        log.Println(err)
        return
    }
	go to_read(conn)
	go to_write(conn)
}


func server_init(){}

/*
func main(){
    get_db()
    //fmt.Printf("Online Clients: %v , %T \n",CLIENTS,CLIENTS)
}
*/
