# 快速生成 base64 编码图片验证码字符串.base64 图形验证码(captcha)为 golang 而设计.

支持多种样式,算术,数字,字母,混合模式,语音模式.

Base64 是网络上最常见的用于传输 8Bit 字节代码的编码方式之一。Base64 编码可用于在 HTTP 环境下传递较长的标识信息, 直接把 base64 当成是字符串方式的数据就好了
减少了 http 请求；数据就是图片；
为 APIs 微服务而设计

#### 为什么 base64 图片 for RESTful 服务

      Data URIs 支持大部分浏览器,IE8之后也支持.
      小图片使用base64响应对于RESTful服务来说更便捷

#### [godoc 文档](https://godoc.org/github.com/carmel/captcha)

#### 在线 Demo [Playground Powered by Vuejs+elementUI+Axios](http://captcha.mojotv.cn)

[![Playground](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/captcha.png "Playground")](http://captcha.mojotv.cn/ "Playground")
[![28+58.png](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/28%2B58%3D%3F.png)](http://captcha.mojotv.cn/ "Playground")
[![ACNRfd.png](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/ACNRfd.png)](http://captcha.mojotv.cn/ "Playground")
[![rW4npZ.png](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/rW4npZ.png)](http://captcha.mojotv.cn/ "Playground")
[![99+73.png](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/99%2B73%3D%3F.png)](http://captcha.mojotv.cn/ "Playground")
[![ctOv6N.png](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/ctOv6N.png)](http://captcha.mojotv.cn/ "Playground")
[![gGncJC.png](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/gGncJC.png)](http://captcha.mojotv.cn/ "Playground")
[![108360.png](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/108360.png)](http://captcha.mojotv.cn/ "Playground")
[wav file](https://raw.githubusercontent.com/mojocn/captcha/master/examples/static/1lNMVxfysfSQJXvjR1LX.wav)

### 安装 golang 包

```sh
go get -u github.com/carmel/captcha
```

### 创建图像验证码

```go
import "github.com/carmel/captcha"
func demoCodeCaptchaCreate() {
	//config struct for digits
	//数字验证码配置
	var configD = captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	//config struct for audio
	//声音验证码配置
	var configA = captcha.ConfigAudio{
		CaptchaLen: 6,
		Language:   "zh",
	}
	//config struct for Character
	//字符,公式,验证码配置
	var configC = captcha.ConfigCharacter{
		Height:             60,
		Width:              240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               captcha.CaptchaModeNumber,
		ComplexOfNoiseText: captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}
	//创建声音验证码
	//Generate 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyA, capA := captcha.Generate("", configA)
	//以base64编码
	base64stringA := captcha.CaptchaWriteToBase64Encoding(capA)
	//创建字符公式验证码.
	//Generate 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyC, capC := captcha.Generate("", configC)
	//以base64编码
	base64stringC := captcha.CaptchaWriteToBase64Encoding(capC)
	//创建数字验证码.
	//Generate 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyD, capD := captcha.Generate("", configD)
	//以base64编码
	base64stringD := captcha.CaptchaWriteToBase64Encoding(capD)

	fmt.Println(idKeyA, base64stringA, "\n")
	fmt.Println(idKeyC, base64stringC, "\n")
	fmt.Println(idKeyD, base64stringD, "\n")
}

```

### 验证图像验证码

```go
import "github.com/carmel/captcha"
func verfiyCaptcha(idkey,verifyValue string){
    verifyResult := captcha.Verify(idkey, verifyValue)
    if verifyResult {
        //success
    } else {
        //fail
    }
}
```

#### 使用 golang 搭建 API 服务

```go
// example of HTTP server that uses the captcha package.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/carmel/captcha"
	"log"
	"net/http"
)

//ConfigJsonBody json request body.
type ConfigJsonBody struct {
	Id              string
	CaptchaType     string
	VerifyValue     string
	ConfigAudio     captcha.ConfigAudio
	ConfigCharacter captcha.ConfigCharacter
	ConfigDigit     captcha.ConfigDigit
}

// captcha create http handler
func generateCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	//parse request parameters
	//接收客户端发送来的请求参数
	decoder := json.NewDecoder(r.Body)
	var postParameters ConfigJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	//create base64 encoding captcha
	//创建base64图像验证码

	var config interface{}
	switch postParameters.CaptchaType {
	case "audio":
		config = postParameters.ConfigAudio
	case "character":
		config = postParameters.ConfigCharacter
	default:
		config = postParameters.ConfigDigit
	}
	//Generate 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	captchaId, digitCap := captcha.Generate(postParameters.Id, config)
	base64Png := captcha.CaptchaWriteToBase64Encoding(digitCap)

	//or you can do this
	//你也可以是用默认参数 生成图像验证码
	//base64Png := captcha.GenerateCaptchaPngBase64StringDefault(captchaId)

	//set json response
	//设置json响应

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := map[string]interface{}{"code": 1, "data": base64Png, "captchaId": captchaId, "msg": "success"}
	json.NewEncoder(w).Encode(body)
}
// captcha verify http handler
func captchaVerifyHandle(w http.ResponseWriter, r *http.Request) {

	//parse request parameters
	//接收客户端发送来的请求参数
	decoder := json.NewDecoder(r.Body)
	var postParameters ConfigJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	//verify the captcha
	//比较图像验证码
	verifyResult := captcha.Verify(postParameters.Id, postParameters.VerifyValue)

	//set json response
	//设置json响应
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := map[string]interface{}{"code": "error", "data": "验证失败", "msg": "captcha failed"}
	if verifyResult {
		body = map[string]interface{}{"code": "success", "data": "验证通过", "msg": "captcha verified"}
	}
	json.NewEncoder(w).Encode(body)
}

//start a net/http server
//启动golang net/http 服务器
func main() {

	//serve Vuejs+ElementUI+Axios Web Application
	http.Handle("/", http.FileServer(http.Dir("./static")))

	//api for create captcha
	//创建图像验证码api
	http.HandleFunc("/api/getCaptcha", generateCaptchaHandler)

	//api for verify captcha
	http.HandleFunc("/api/verifyCaptcha", captchaVerifyHandle)

	fmt.Println("Server is at localhost:3333")
	if err := http.ListenAndServe("localhost:3333", nil); err != nil {
		log.Fatal(err)
	}
}
```

#### [使用 redis 做储存](examples_redis/main.go)

#### 运行 demo 代码

    cd $GOPATH/src/github.com/carmel/captcha/_examples
    go run main.go

#### 访问 `http://localhost:777`
