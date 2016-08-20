package main

import (
    "net"
    "os"
    "fmt"
    "time"
)

type TicToc struct {
    tic time.Time
}

func (ticToc * TicToc) Tic() {
    ticToc.tic = time.Now()
}

func (ticToc * TicToc) TocPrint() {
    fmt.Println(time.Now().Sub(ticToc.tic))
}

func (ticToc * TicToc) Toc() time.Duration {
    return time.Now().Sub(ticToc.tic)
}

func main() {
    var t TicToc

    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s host", os.Args[0])
        os.Exit(1)
    }
    service := os.Args[1]
    fmt.Println("Scanning ", service, ", please wait...")

    d := &net.Dialer{Timeout: 150 * time.Millisecond}
    sem := make(chan bool, 100);
    t.Tic()
    for i := 1; i < 65536; i++ {
        sem <- true
        go func (host string, port int) {
            host_name := fmt.Sprintf("%s:%d", host, port)
            _, err := d.Dial("tcp", host_name)
            if (err != nil) {
                //fmt.Printf("%d: closed\n", port)
            } else {
                fmt.Printf("%d: open\n", port)
            }
            <-sem;
        } (service, i);
    }

    for i := 0; i < cap(sem); i++ {
        sem <- true
    }

    fmt.Println("Scanning took: ", t.Toc())

    os.Exit(0)
}
