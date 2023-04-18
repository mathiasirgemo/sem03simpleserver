package main

import (
	"io"
	"log"
	"net"
	"sync"
	"github.com/uia-worker/is105sem03/mycrypt"
)

func main() {
    var wg sync.WaitGroup
    server, err := net.Listen("tcp", "172.17.0.4:8788")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("bundet til %s", server.Addr().String())
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            log.Println("før server.Accept() kallet")
            conn, err := server.Accept()
            if err != nil {
                return
            }
            go func(c net.Conn) {
                defer c.Close()
                for {
                    buf := make([]byte, 1024)
                    n, err := c.Read(buf)
                    if err != nil {
                        if err != io.EOF {
                            log.Println(err)
                        }
                        return // fra for løkke
                    }
                    switch msg := string(buf[:n]); msg {
                    case "ping":
                        _, err = c.Write([]byte("pong"))
                    default:
                        dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
                        log.Println("Dekryptert melding: ", string(dekryptertMelding))
                        _, err = c.Write([]byte(string(dekryptertMelding)))
                    }
                    if err != nil {
                        if err != io.EOF {
                            log.Println(err)
                        }
                        return // fra for løkke
                    }
                }
            }(conn)
        }
    }()
    wg.Wait()
}

