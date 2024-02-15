package main

// 导入依赖
import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
)

// 定义全局变量
var img_path = "./goodjob-img/resources"
var img_dict = make(map[string]int)
var name_list []string
var app *gin.Engine

// 定义 ServeHandler 结构体
type ServeHandler struct {
	name string
}

// 定义构造函数
func NewServeHandler(name string) *ServeHandler {
	return &ServeHandler{name: name}
}

// 定义 Call 方法
func (sf *ServeHandler) Call(_gin *gin.Context) {
	path := pickImg(sf.name)
	if path == "404" {
		_gin.AbortWithStatus(http.StatusNotFound)
	} else {
		_gin.File(path)
	}
}

// 随机抓取图片
func pickImg(name string) string {
	if _, exists := img_dict[name]; exists {
	} else {
		fmt.Println(name, "not found")
		return "404"
	}
	randNum := strconv.Itoa(gofakeit.Number(1, img_dict[name]) - 1)
	// 茄子是 gif
	if name == "茄子" {
		return fmt.Sprintf("%s/%s/%s.gif", img_path, name, randNum)
	} else {
		return fmt.Sprintf("%s/%s/%s.png", img_path, name, randNum)
	}
}

// 获取文件夹下文件数量
func lenPath(path string) int {
	count := 0

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	return count
}

// 遍历文件夹注册路由点等
func dealPath(_ string, info os.FileInfo, err error) error {
	if info.IsDir() {
		// 我也不知道为啥，不过这样写就不会崩了
		if info.Name() == "resources" {
			return nil
		}
		// 更新字典, 注册路由
		img_dict[info.Name()] = lenPath(img_path + "/" + info.Name())
		name_list = append(name_list, info.Name())
		handler := NewServeHandler(info.Name())
		app.GET("/"+info.Name(), handler.Call)
	}
	return nil
}

// 加载文件夹
func loadPath() {
	filepath.Walk(img_path, dealPath)
}

// 随机抽取名字
func randName() string {
	return name_list[gofakeit.Number(1, len(name_list))-1]
}

// 直接访问的逻辑
func direct(_gin *gin.Context) {
	_gin.File(pickImg(randName()))
}

// 帮助页面
func help(_gin *gin.Context) {
	_gin.JSON(http.StatusOK, gin.H{
		"用法":   "发送 GET 请求时会从词库随机抽取一个图片，如果需要指定某个人则可以 GET 对应地址",
		"请求方式": "发送 GET 请求获取任意图片，发送 GET 请求到对应地址获取某个人的图片",
		"示例":   "GET https://img-api.justrobot.dev/win11 来获得 win11 的随机图片",
		"更新频率": "API 图库会在每天的 00:05 左右重启进行图库更新，耗时 10 秒以内",
		"速率限制": "图库有访问限制，单 IP 每秒限制 1 次，每 10 秒限制 20 次(包含无效访问), 超过任意限制均返回 429 错误",
	})
}

// 从 git 更新图片
func updateImg() {
	cmd := exec.Command("git", "pull", "-f")
	cmd.Dir = "./goodjob-img"
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}

// 入口点
func main() {
	// 生产模式
	gin.SetMode(gin.ReleaseMode)
	// 初始化 gin
	app = gin.Default()
	// 更新图片
	updateImg()
	// 加载文件夹
	loadPath()
	// 注册固定路由节点
	app.GET("/", direct)
	app.GET("/HELPS", help)
	// Gin, 启动!
	app.Run(":10808")
}
