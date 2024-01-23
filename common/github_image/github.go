package github_image

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

type GithubServe struct {
	RepoInterface
}

func (serve *GithubServe) Push(PATH, filename, content string) (string, string, string) {
	return Push(PATH, filename, content)
}

func (serve *GithubServe) GetFiles() []map[string]interface{} {
	return GetFiles()
}

func (serve *GithubServe) Del(filepath, sha string) string {
	return DelFile(filepath, sha)
}

func Push(PATH, filename, content string) (string, string, string) {

	url := "https://api.github.com/repos/" +
		OWNER + "/" +
		REPO +
		"/contents/" +
		PATH +
		"/" + filename

	fmt.Println(23333333)
	fmt.Println(url)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("PUT")
	req.Header.SetBytesKV([]byte("Content-Type"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+ TOKEN))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]string)
	args["content"] = content
	args["message"] = "upload pic for repo-image-hosting"
	args["branch"] = BRANCH

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	// 发起请求
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	var mapResult map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	d := ""
	p := ""
	s := ""

	_, ok := mapResult["content"]

	if ok {
		if mapResult["content"] != nil {
			path := mapResult["content"].(map[string]interface{})["path"].(string)
			d = "https://cdn.jsdelivr.net/gh/" + OWNER + "/" + REPO + "@" + BRANCH + "/" + path
			p = path
			s = mapResult["content"].(map[string]interface{})["sha"].(string)
		}
	}

	return d, p, s
}

func GetFiles() []map[string]interface{} {
	// 初始化请求与响应
	url := "https://api.github.com/repos/" +
		OWNER + "/" +
		REPO +
		"/contents/" +
		PATH +
		"?ref=" + BRANCH

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("GET")
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+TOKEN))
	// 设置请求的目标网址
	req.SetRequestURI(url)

	// 发起请求
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	//log.Println(string(body),url)

	var mapResult []map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	for _, v := range mapResult {
		v["download_url"] = "https://cdn.jsdelivr.net/gh/" + OWNER + "/" + REPO + "@" + BRANCH + "/" + v["path"].(string)
	}
	return mapResult
}

func DelFile(filepath, sha string) string {
	// 初始化请求与响应
	url := "https://api.github.com/repos/" +
		OWNER + "/" +
		REPO +
		"/contents/" + filepath

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("DELETE")
	req.Header.SetBytesKV([]byte("Content-Type"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+TOKEN))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]string)
	args["sha"] = sha
	args["message"] = "delete pic for repo-image-hosting"
	args["branch"] = BRANCH

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	// 发起请求
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	return string(body)
}
