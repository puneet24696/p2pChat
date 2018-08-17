package main

import (
	"encoding/json"
    "log"
	"bufio"
    "os"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
	"net"
    "crypto/tls"
	"bytes"
)

var CLI_conn = make(map[string]net.Conn)
var CLIENTS = make(map[string]string)
var ip_port string
var my_name string
var db_url = "http://192.168.0.84:9090/chatroom"

func get_db(){
    req, err := http.NewRequest("GET",db_url,nil)
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
    fmt.Println("started writing")
    scanner := bufio.NewScanner(os.Stdin)
    for {
	fmt.Printf("%v",conn)
	scanner.Scan()
	a := scanner.Text()
        w_mess, err := conn.Write([]byte(a))
        if err!= nil || a == "bye"{
            fmt.Println(w_mess,err)
            conn_closer(conn)
            return
       } 
    }
    return
}
/*
func to_write_beta(){
    fmt.Println("started writing")
    for {
	var name string
	fmt.Println("to <space> message")
	fmt.Scanln(&name)
	friend, a := strings.SplitN(name," ",2)[0], strings.SplitN(name," ",2)[1]	
	fmt.Println("me","->",friend)
        w_mess, err := CLI_conn[CLIENTS[friend]].Write([]byte(a))
        if err!= nil || a == "bye"{
            fmt.Println(w_mess,err)
            conn_closer(conn)
            return
       } 
    }
    return

}
*/
func to_read(conn net.Conn){

    fmt.Println("started reading")
    for{
        message := make([]byte, 1024)
	text, err := conn.Read(message)
        if err != nil {
            fmt.Println(text, err)
            conn_closer(conn)
            return
	}
	fmt.Println("him->",string(message))
    }
    return

}


func client_init_chat(ip_port_friend string){
    fmt.Println("entered client func")
    cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }
    config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
    fmt.Println("made config")
    conn, err := tls.Dial("tcp", ip_port_friend, config)
    fmt.Println("entered client func -> post dial")
    if err != nil {
        log.Println(err)
        return
    }
    _, _ = conn.Write([]byte("hello\n"))
    fmt.Println("entered client func -> about to read and write")
	go to_read(conn)
	go to_write(conn)
	fmt.Println("Chat has been initialized")
	return
}


func server_init(ip_port string){
        cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
        if err != nil {
            log.Fatalf("server: loadkeys: %s", err)
    }
        config := &tls.Config{Certificates: []tls.Certificate{cert}}
        fmt.Println("server func -> about to listen")
	listener, err := tls.Listen("tcp",ip_port,config)
	if err != nil{
		log.Println(err, "error in listener")
		return
	}
	defer listener.Close()
	for {
                fmt.Println("server func -> about to accept")
		conn,err :=listener.Accept()
		CLI_conn[conn.RemoteAddr().String()] = conn
		if err!=nil{
			fmt.Println(err,"err in listener/accept")
			continue
		}
		fmt.Println("connected to ",conn)
		go to_read(conn)
		go to_write(conn)
	}
}

func conn_closer(conn net.Conn){
	temp_map := make(map[string]string)	
	temp_map["name"] = my_name
	temp_map["ip_port"]= ip_port
	temp_json,_ := json.Marshal(temp_map)
	req, err := http.NewRequest("POST", db_url,bytes.NewBuffer(temp_json))
	req.Header.Set("Content-Type", "application/json")
	if err != nil{
		fmt.Println(req,err,"error occured in connection closer,New request,retry")
		return
	}
	client := &http.Client{}
	resp,err := client.Do(req)
	if err != nil {
		fmt.Println(resp,err,"error occured in conection closer,client do,retry")
		return
	}else if resp.StatusCode == 200 {
		conn.Close()
		fmt.Println("You have been disconnected")
		return
	}
	
}

func main(){
	my_name = os.Args[2]
	ip_port = os.Args[1]
	get_db()
	fmt.Println("Test 1 executed")
	go server_init(ip_port)
	fmt.Println("Test 1 executed")
	fmt.Println("Which friend do you want to connect","->")
	//for{
            var friend string    
            fmt.Printf("%v","now add friend")
	    fmt.Scanln(&friend)
	    if friend != "none" {
	        fmt.Printf("%v %T","Test 1 executed",friend)
	    
	        go client_init_chat(friend)
	    
	        fmt.Println("Test 1 executed")
            }
	//}
	fmt.Printf("%v,%v","no friends",friend)
	for{}
    
}

