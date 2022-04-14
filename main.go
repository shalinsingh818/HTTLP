package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"io"
	"io/ioutil"
)


type Gateway struct {
    NetIp         string
    MacAddr       string
    Id            int
    Channel       chan string
    ClientRequest *http.Request
}


var gateways = make(map[string]Gateway)
var gatewayCount = 0
var host = "http://127.0.0.1:8080"


//checks to see if request is active in gateway
//gateway struct holds request that was sent, you can monitor the status
func (g *Gateway) checkRequest() bool{
    var status bool = true
    select{
       case <-g.ClientRequest.Context().Done():
        var test = <-g.ClientRequest.Context().Done()
        fmt.Println(test)
        status = false
       default:
        status = true
    }
    return status
}

//Gets an IP/Remote addr from a HTTP request
func GetIP(r *http.Request) string {
    forwarded := r.Header.Get("X-FORWARDED-FOR")
    if forwarded != "" {
        return forwarded
    }
    return r.RemoteAddr
}


//push response to a gateway, input mac address as parameter
func pushHandler(w http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)
    body, err := ioutil.ReadAll(req.Body)

    if err != nil {
        w.WriteHeader(400)
    }
    mac := params["macAddr"]
    gateways[mac].Channel <- string(body)
    fmt.Println("Pushed to: ", gateways[mac])
}



//Poll a response --> send response to server, save response for server push
func pollRequest(w http.ResponseWriter, r *http.Request) {
    //gateway channel, {Holds the responses}
    gChan := make(chan string)
    requestAddr := GetIP(r)
    requestMacAddr := r.FormValue("mac-address")
    //check if there is an active gateway with this mac address
    if gateways[requestMacAddr].MacAddr == requestMacAddr {
        gateways[requestMacAddr].Channel <- "Disconnecting existing connection"
        delete(gateways, requestMacAddr)
    }
    //create a client instance for the gateway
    gateways[requestMacAddr] = Gateway{requestAddr, requestMacAddr, gatewayCount, gChan, r}
    connectionMessage := fmt.Sprintf("Received connection from %v : %v", requestAddr, requestMacAddr)
    fmt.Println(connectionMessage)
    io.WriteString(w, <-gateways[requestMacAddr].Channel)
}


//routes for the server, look at postman documentation for more info
func serverHandler(Port string) {
    fmt.Println("Control server listening on port", Port)
    myPort := fmt.Sprintf(":%s", Port)
    serverRouter := mux.NewRouter().StrictSlash(true)
    serverRouter.HandleFunc("/poll", pollRequest).Methods("POST")
    serverRouter.HandleFunc("/push/{macAddr}", pushHandler).Methods("POST")
    log.Fatal(http.ListenAndServe(myPort, serverRouter))
}



func main(){
	serverHandler("8080")
	fmt.Println("Test")
}
