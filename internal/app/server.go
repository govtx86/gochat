// Package app provides the main application.
package app

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type User struct {
	conn     net.Conn
	username string
}

var Users = make(map[string]User)

func runListener(address string, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		fmt.Println("err:", err)
	}
	defer listener.Close()
	fmt.Println("Listening on", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		usernameByte, err := bufio.NewReader(conn).ReadString('\n')
		username := strings.Trim(usernameByte, "\r\n")
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		for _, user := range Users {
			if user.username == username {
				conn.Write([]byte("409\n"))
				conn.Close()
				continue
			}
		}
		Users[username] = User{
			conn:     conn,
			username: username,
		}
		conn.Write([]byte("200\n"))
		go handleClient(Users[username])
	}
}

func handleClient(user User) {
	defer func() {
		user.conn.Close()
		for _, u := range Users {
			u.conn.Write([]byte(user.username + " left!\n"))
		}
		delete(Users, user.username)
		broadcastUserList()
	}()
	broadcastUserList()
	for {
		msg, err := bufio.NewReader(user.conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")
		fmt.Println("Received: " + user.username + " : " + msg)
		broadcast(msg, user.username)
	}
}

func broadcast(msg string, username string) {
	for _, u := range Users {
		if u.username == username {
			u.conn.Write([]byte("You: " + msg + "\n"))
			continue
		}
		u.conn.Write([]byte(username + ": " + msg + "\n"))
	}
}

func broadcastUserList() {
	var users string
	for _, u := range Users {
		users = users + u.username + "#$"
	}
	for _, u := range Users {
		u.conn.Write([]byte("#srvc:" + users + "\n"))
	}
}

func RunServer(address string) {
	port := 8080
	runListener(address, port)
}
