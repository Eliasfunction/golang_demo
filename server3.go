package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type client struct {
	C    chan string
	Name string
	Addr string
}

//創建MAP 儲存在線用戶名稱
var onlineMap map[string]client

//創建主頻道傳遞用戶資訊
var message = make(chan string)

func writeMsgToclient(clnt client, conn net.Conn) {
	for msg := range clnt.C {
		conn.Write([]byte(msg + "\n"))
	}
}

func makemsg(clnt client, msg string) (buf string) {
	buf = "[" + clnt.Addr + "]" + clnt.Name + ": " + msg
	return
}

//handlerconnect is handlerconnect
func handlerconnect(conn net.Conn) {
	defer conn.Close()

	//活躍度 避免掛機
	hasdate := make(chan bool)

	//查詢用戶地址
	netAddr := conn.RemoteAddr().String()

	//產生新上線用戶的資料物件 默認IP+PORT
	clnt := client{make(chan string), netAddr, netAddr}

	//將新用戶加到MAP中 key: ip+port value client
	onlineMap[netAddr] = clnt

	//給當前用戶發送訊息
	go writeMsgToclient(clnt, conn)

	//通知使用者上限消息到Channel中
	message <- makemsg(clnt, "login")

	isquit := make(chan bool) //連線狀態

	//匿名程序 處理客戶端發送的訊息 //goroutine
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				isquit <- true
				fmt.Println("偵測到客戶端退出: \n", clnt.Name, netAddr)
				return
			}
			if err != nil {
				fmt.Println("conn  read err:", err)
				return
			}

			//將讀到的用戶訊息 儲存到msg string
			msg := string(buf[:n-1])

			//獲取在線上用戶的列表
			if msg == "who" && len(msg) == 3 {
				conn.Write([]byte("online user list:\n"))
				//訪問MAP 獲取在線上用戶

				for _, user := range onlineMap {
					userinfo := user.Addr + ":" + user.Name + "\n"
					conn.Write([]byte(userinfo))
				}
				//判斷用戶改名
			} else if len(msg) >= 8 && msg[:6] == "rename" {
				newname := strings.Split(msg, " ")[1]
				clnt.Name = newname       //修改struct name
				onlineMap[netAddr] = clnt //更新再現用戶列表
				conn.Write([]byte("rename successful\n"))
			} else {
				//將讀到的用戶訊息 寫入到主頻道中
				message <- makemsg(clnt, msg)
			}
			hasdate <- true //活躍度偵測
		}
	}()
	//防中斷
	for {
		//監聽頻道上的資料
		select {
		case <-isquit:
			delete(onlineMap, clnt.Addr)
			message <- makemsg(clnt, "logout")
			return
		case <-hasdate:
			//偵測到活躍 重設下方計時器
		case <-time.After(time.Second * 20):
			delete(onlineMap, clnt.Addr)
			message <- makemsg(clnt, "timeout")
			return
		}
	}
}

// manager is manager
func manager() {
	//初始化map
	onlineMap = make(map[string]client)

	//循環更新MASSAGE
	for {
		//監聽channel是否有資料 有則存到MSG
		msg := <-message

		//循環發送給所有客戶端
		for _, clnt := range onlineMap {
			clnt.C <- msg
		}
	}
}

//main is main
func main() {
	//創建監聽Socket
	Listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("Listen err", err)
		return
	}
	defer Listener.Close()

	//創建管理MAP與channel
	go manager()

	//循環監聽賀戶端連接請求
	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Accept err", err)
			return
		}
		//啟動程序處理客戶端請求資料
		go handlerconnect(conn)
	}

}
