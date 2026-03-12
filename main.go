package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/xmpp"
)

var client *xmpp.Client

// 处理登录请求
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// 使用用户名和密码连接到 XMPP 服务器
		client, err := xmpp.NewClient(fmt.Sprintf("%s@%s", username, "example.com"), password, xmpp.Options{})
		if err != nil {
			http.Error(w, "登录失败", http.StatusInternalServerError)
			return
		}

		// 发送登录成功后的首页
		http.Redirect(w, r, "/chat", http.StatusFound)
		return
	}

	// 渲染登录页面
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "无法加载登录页面", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// 聊天页面
func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20) // 允许最多 10MB 的文件上传
		message := r.FormValue("message")

		// 发送文本消息
		if message != "" {
			err := client.Send(xmpp.Message{
				To:   "example@server.com",
				Body: message,
			})
			if err != nil {
				http.Error(w, "消息发送失败", http.StatusInternalServerError)
				return
			}
		}

		// 处理文件上传
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "文件上传失败", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// 创建一个文件保存目录
		savePath := filepath.Join("public", "upload", "image.png")
		outFile, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "文件保存失败", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// 保存上传的文件
		_, err = outFile.ReadFrom(file)
		if err != nil {
			http.Error(w, "文件保存失败", http.StatusInternalServerError)
			return
		}

		// 返回上传的文件路径
		http.Redirect(w, r, "/chat", http.StatusFound)
		return
	}

	// 渲染聊天页面
	tmpl, err := template.ParseFiles("chat.html")
	if err != nil {
		http.Error(w, "无法加载聊天页面", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	// 路由和处理器
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/chat", chatHandler)

	// 启动 Web 服务
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
