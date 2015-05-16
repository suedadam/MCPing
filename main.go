package main

import (
        "bytes"
        "io"
        "fmt"
        "net"
        "strings"
        "os"
        "github.com/minero/minero/proto/packet"
        "time"
)

func main() {
        if len(os.Args) < 2 {
                fmt.Printf("Usage: %s <IP>:<port>\n", os.Args[0])
                os.Exit(2)
        }
        addr := os.Args[1]
        _, port, err := net.SplitHostPort(addr)
        if len(port) <= 0 {
                addr += ":25565"
        }

        c, err := net.DialTimeout("tcp", addr, time.Duration(500*time.Millisecond))
        if err != nil {
                fmt.Printf("Error connecting!\n%s\n", err)
                return
        }
        defer c.Close()

        fmt.Println("Connected to:", c.RemoteAddr())

        p := packet.ServerListPing{Magic: 1}
        p.WriteTo(c)

        var buf bytes.Buffer
        io.Copy(&buf, c)

        id, _ := buf.ReadByte()

        if id != packet.PacketDisconnect {
                fmt.Printf("Unexpected packet id! D:\n")
                return
        }

        r := new(packet.Disconnect)
        r.ReadFrom(&buf)

        s := strings.Split(r.Reason, "\x00")

        fmt.Println("Server version:", s[2])
        fmt.Printf("Players: %s/%s\n", s[4], s[5])
        fmt.Println("MOTD:", s[3])
}
