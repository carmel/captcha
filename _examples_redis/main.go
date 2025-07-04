// example of HTTP server that uses the captcha package.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/carmel/captcha"
	"github.com/go-redis/redis"
)

// customizeRdsStore An object implementing Store interface
type customizeRdsStore struct {
	redisClient *redis.Client
}

// customizeRdsStore implementing Set method of  Store interface
func (s *customizeRdsStore) Set(id string, value string) {
	err := s.redisClient.Set(id, value, time.Minute*10).Err()
	if err != nil {
		panic(err)
	}
}

// customizeRdsStore implementing Get method of  Store interface
func (s *customizeRdsStore) Get(id string, clear bool) (value string) {
	val, err := s.redisClient.Get(id).Result()
	if err != nil {
		panic(err)
	}
	if clear {
		err := s.redisClient.Del(id).Err()
		if err != nil {
			panic(err)
		}
	}
	return val
}

func init() {
	// create redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// init redis store
	customeStore := customizeRdsStore{client}

	captcha.SetCustomStore(&customeStore)

}

// ConfigJsonBody json request body.
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
	// parse request parameters
	// 接收客户端发送来的请求参数
	decoder := json.NewDecoder(r.Body)
	var postParameters ConfigJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	// create base64 encoding captcha
	// 创建base64图像验证码

	var config any
	switch postParameters.CaptchaType {
	case "audio":
		config = postParameters.ConfigAudio
	case "character":
		config = postParameters.ConfigCharacter
	default:
		config = postParameters.ConfigDigit
	}
	captchaId, captcaInterfaceInstance := captcha.Generate(postParameters.Id, config)
	base64blob := captcha.WriteToBase64Encoding(captcaInterfaceInstance)

	// or you can just write the captcha content to the httpResponseWriter.
	// before you put the captchaId into the response COOKIE.
	// captcaInterfaceInstance.WriteTo(w)

	// set json response
	// 设置json响应
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := map[string]any{"code": 1, "data": base64blob, "captchaId": captchaId, "msg": "success"}
	json.NewEncoder(w).Encode(body)
}

// captcha verify http handler
func captchaVerifyHandle(w http.ResponseWriter, r *http.Request) {

	// parse request parameters
	// 接收客户端发送来的请求参数
	decoder := json.NewDecoder(r.Body)
	var postParameters ConfigJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	// verify the captcha
	// 比较图像验证码
	verifyResult := captcha.Verify(postParameters.Id, postParameters.VerifyValue)

	// set json response
	// 设置json响应
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := map[string]any{"code": "error", "data": "验证失败", "msg": "captcha failed"}
	if verifyResult {
		body = map[string]any{"code": "success", "data": "验证通过", "msg": "captcha verified"}
	}
	json.NewEncoder(w).Encode(body)
}

// start a net/http server
// 启动golang net/http 服务器
func main() {

	staticPath := fmt.Sprintf("%s/src/github.com/carmel/captcha/_examples/static", os.Getenv("GOPATH"))

	// serve Vuejs+ElementUI+Axios Web Application
	http.Handle("/", http.FileServer(http.Dir(staticPath)))

	// api for create captcha
	http.HandleFunc("/api/getCaptcha", generateCaptchaHandler)

	// api for verify captcha
	http.HandleFunc("/api/verifyCaptcha", captchaVerifyHandle)

	fmt.Println("Server is at localhost:7777")
	if err := http.ListenAndServe("localhost:7777", nil); err != nil {
		log.Fatal(err)
	}
}
