package main

// 导入依赖
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
)

// 定义全局变量
var img_path = "./goodjob-img/resources"
var res_path = "https://gitee.com/SmallK111407/goodjob-img"
var img_dict = make(map[string][]string)
var aliasMap map[string]string
var name_list []string
var app *gin.Engine

// 定义 ServeHandler 结构体
type ServeHandler struct {
	name string
}

// 定义构造函数
func NewServeHandler(name string) *ServeHandler {
	return &ServeHandler{name: aliasMap[name]}
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

// 获取别名映射
func getAliasMap() error {
	resp, err := http.Get("https://gitee.com/SmallK111407/useless-plugin/raw/main/model/aliasData/alias.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var aliasData map[string][]string
	err = json.Unmarshal(body, &aliasData)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	AliasMap := make(map[string]string)

	for k, v := range aliasData {
		if _, exitsts := img_dict[k]; !exitsts {
			continue
		}
		AliasMap[k] = k
		for _, alias := range v {
			AliasMap[alias] = k
		}
	}

	aliasMap = AliasMap
	return nil
}

// 随机抓取图片
func pickImg(name string) string {
	if _, exists := img_dict[name]; exists {
	} else {
		fmt.Println(name, "not found")
		return "404"
	}
	return fmt.Sprintf("%s/%s/%s", img_path, name, RandItem(img_dict[name]))
}

// 从列表中随机抽取元素
func RandItem(list []string) string {
	return list[gofakeit.Number(1, len(list)-1)]
}

// 获得文件夹下的文件列表
func getFileList(Name string) ([]string, error) {
	dirName := img_path + "/" + Name
	files, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	return fileList, nil
}

// 遍历文件夹更新字典
func dealPath(_ string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		// 我也不知道为啥，不过这样写就不会崩了
		if info.Name() == "resources" {
			return nil
		}
		// 更新字典
		img_dict[info.Name()], err = getFileList(info.Name())
		if err != nil {
			return err
		}
		name_list = append(name_list, info.Name())
	}
	return err
}

// 加载文件夹
func loadPath() {
	filepath.Walk(img_path, dealPath)
}

// 加载路由节点
func LoadPoints() {
	for k := range aliasMap {
		app.GET("/"+k, NewServeHandler(k).Call)
	}
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
		"用法":        "发送 GET 请求时会从词库随机抽取一个图片，如果需要指定某个人则可以 GET 对应地址",
		"请求方式":      "发送 GET 请求获取任意图片，发送 GET 请求到对应地址获取某个人的图片",
		"别名处理":      "以 https://gitee.com/SmallK111407/useless-plugin/blob/main/model/aliasData/alias.json 中的别名为准",
		"快速获得别名映射表": "GET https://img-api.justrobot.dev/AliasMap",
		"示例":        "GET https://img-api.justrobot.dev/oldcitynight 来获得 oldcitynight 的随机图片",
		"更新频率":      "API 图库会在每天的 00:05 左右重启进行图库更新，耗时 10 秒以内",
		"速率限制":      "图库有访问限制，单 IP 每秒限制 1 次，每 10 秒限制 20 次(包含无效访问), 超过任意限制均返回 429 错误",
		"图片类型":      "图片类型可能为 png 或 gif , 如果图库有其他图片类型会原样提供",
	})
}

// 从 git 更新图片
func updateImg() {
	// 确认图片资源存在
	if _, err := os.Stat("./goodjob-img"); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", res_path)
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
	}
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
	// 获取别名映射
	getAliasMap()
	fmt.Println(aliasMap)
	// 注册固定路由节点
	app.GET("/", direct)
	app.GET("/HELPS", help)
	// 加载路由节点
	LoadPoints()
	// Gin, 启动!
	app.Run(":10808")
}
