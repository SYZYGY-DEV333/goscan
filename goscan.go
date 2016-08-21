package main

import (
    "net"
    "os"
    "fmt"
    "time"
    "strconv"
)

var addr string = os.Args[1]
var fpst string = os.Args[2]
var lpst string = os.Args[3]

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
    
    fp, err := strconv.Atoi(fpst)
    lp, err := strconv.Atoi(lpst)
    if err != nil {
        fmt.Println("Error: strconv.Atoi")
        os.Exit(1)
    }

    if len(os.Args) != 4 {
        fmt.Fprintf(os.Stderr, "Usage: %s [host] [firstport] [lastport]", os.Args[0])
        os.Exit(1)
    }
    fmt.Println("Scanning", addr, "ports ", fpst, "-", lpst)

    d := &net.Dialer{Timeout: 150 * time.Millisecond}
    sem := make(chan bool, 100);
    t.Tic()
    for i := fp; i < lp; i++ {
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
        } (addr, i);
    }

    for i := 0; i < cap(sem); i++ {
        sem <- true
    }

    fmt.Println("Scanning took: ", t.Toc())

    os.Exit(0)
}
