package main

import (
	"flag"
	"log"
	"net"
)

var (
	TYPE     string // udp/tcp
	ENDPOINT string // ip:port
	FLAG     string // s/c
	HELP     bool
)

func init() {
	flag.StringVar(&TYPE, "protocal", "tcp", "Choose (tcp/udp) protocal to testing.")
	flag.StringVar(&ENDPOINT, "address", "", "address for bind or connect.")
	flag.StringVar(&FLAG, "role", "server", "client or server.")
	flag.BoolVar(&HELP, "help", false, "this help.")
}

type Network interface {
	Client()
	Server()
}

type UdpProtocal struct {
	addr string
}

type TcpProtocal struct {
	addr string
}

func NewNetWork(protcal, addr string) Network {
	if protcal == "tcp" {
		return &TcpProtocal{addr: addr}
	}
	return &UdpProtocal{addr: addr}
}

func (tcp *TcpProtocal) Client() {
	conn, err := net.Dial("tcp", tcp.addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	buf := make([]byte, 2048)
	_, err = conn.Read(buf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("remote addr : %s -> local addr %s\n",
		conn.RemoteAddr().String(), conn.LocalAddr().String())
}

func tcpServerProc(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		cnt, err := conn.Read(buf)
		if err != nil {
			return
		}
		cnt, err = conn.Write(buf[0:cnt])
		if err != nil {
			break
		}
	}
}

func (tcp *TcpProtocal) Server() {
	listen, err := net.Listen("tcp", tcp.addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for {
		conn, err2 := listen.Accept()
		if err2 != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("remote addr : %s -> local addr %s\n",
			conn.RemoteAddr().String(), conn.LocalAddr().String())
		go tcpServerProc(conn)
	}
}

func (udp *UdpProtocal) Client() {
	conn, err := net.Dial("udp", udp.addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	buf := make([]byte, 2048)
	_, err = conn.Read(buf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("[remote addr : %s -> local addr %s]\n",
		conn.RemoteAddr().String(), conn.LocalAddr().String())
}

func (udp *UdpProtocal) Server() {
	addr, err := net.ResolveUDPAddr("udp", udp.addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	buff := make([]byte, 2048)

	for {
		cnt, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			continue
		}
		_, err = conn.WriteToUDP(buff[:cnt], addr)
		if err != nil {
			continue
		}
		log.Printf("[remote addr : %s -> local addr %s]\n",
			addr.String(), conn.LocalAddr().String())
	}
}

func main() {
	flag.Parse()

	if HELP || ENDPOINT == "" {
		flag.Usage()
		return
	}

	if FLAG != "server" && FLAG != "client" {
		flag.Usage()
		return
	}

	if TYPE != "tcp" && TYPE != "udp" {
		flag.Usage()
		return
	}

	log.Printf("[%s %s://%s]\n", FLAG, TYPE, ENDPOINT)

	network := NewNetWork(TYPE, ENDPOINT)

	if FLAG == "server" {
		network.Server()
	} else {
		network.Client()
	}
}
