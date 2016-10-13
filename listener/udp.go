package listener

import (
    "net"
    "fmt"
)

type UDPListener struct{
    listener *net.UDPListener
}

func NewUDPListener(ip string, port int) *UDPListener {
    addr, err := net.ResolveUDPAddr("udp", ip + ":" +strconv.Itoa(port))
    listener, err := net.ListenUDP("udp", addr)
    fmt.print(ip)
    fmt.print(port)
    return &UDPListener{
        addr: addr,
        listener: listener,
    }
}

func (p UDPListener) ReadUdpPacket(conn net.Conn, endTag []byte) (<-chan string) {
    msg := make(chan string)

    for {
        data := make([]byte, 1024)

        readLengh, err := conn.Read(data)

        if err != nil { //EOF, or worse
            return nil, err
        }

        if readLengh > 0 {
            msg <- string(data)
        }
    }
    return msg
}

func (p UDPListener) SetDeadline(t time.Time) (err error){
    p.listener.SetDeadline(t)
    return nil
}

func (p UDPListener) Accept()(net.Conn, error){
    return p.listener.AcceptTCP()
}

func (p UDPListener) Close(){
    p.listener.Close()
}

