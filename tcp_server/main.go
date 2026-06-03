package main

import (
    "fmt"
    "net"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
 
    
    buf := make([]byte, 1024)

    for {

        n, err := conn.Read(buf)

        if(err != nil){
            fmt.Println("Client Disconnected : " , err )
            return 
        }

        msg := string(buf[:n])

        fmt.Println(msg)

        _ , err = conn.Write([]byte("Received " + msg))

         if(err != nil){
            fmt.Println( err )
            return 
        }
    }

}

func main() {
    
    listener, err := net.Listen("tcp", ":8080")
    
    if(err != nil){
        panic(err)
    }
   
    // Close listener when server exits
    defer listener.Close()


    for{
        conn, err := listener.Accept()
        
        if(err != nil){
             panic(err)
        }

        fmt.Println("New Cleint Arrived" ,  conn.RemoteAddr() )

        go handleConnection(conn)

    }

};