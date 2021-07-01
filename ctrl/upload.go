package ctrl

import (
	"fmt"
	"goim/server"
	"goim/util"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	filesImageMax = 1024 * 1024 * 2 // 2MB
	filesVideoMax = 1024 * 1024 * 10 // 2MB
)

func init() {
	// 运行程序时创建文件夹
	os.MkdirAll("./mnt", os.ModePerm)
}

func UploadLocal(w http.ResponseWriter, r *http.Request) {
	// 获取上传的资源
	srcfile, head, err := r.FormFile("file")
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	defer srcfile.Close()
	// TODO 创建一个新文件
	suffix := ".png"
	// 如果前端文件名称包含后缀
	ofilename := head.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	// 如果前端指定filetype
	filetype := r.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}

	// 文件格式
	files := map[string]string{
		"image/jpeg":               "jpg",
		"image/gif":                "git",
		"image/png":                "png",
		"image/svg+xml":            "svg",
		"image/vnd.microsoft.icon": "ico",
		"audio/mpeg":               "mp3",
		"audio/wav":                "wav",
		"audio/webm":               "weba",
		"video/mpeg":               "mpeg",
		"video/webm":               "webm",
		"video/mp4":                "mp4",
	}
	video := map[string]string{
		"mp4": "map4",
		"webm": "webm",
		"mpeg": "mpeg",
	}
	value, ok := files[head.Header.Get("Conent-Type")]
	if !ok {
		util.RespFail(w, "文件类型错误，请重新上传")
		return
	}
	if _, ok := video[value]; ok {
		// 文件大小
		if head.Size > filesVideoMax {
			util.RespFail(w, "video文件大小不能超过 "+strconv.Itoa(filesVideoMax/1024/1024)+" Mb")
			return
		}
	} else {
		// 文件大小
		if head.Size > filesImageMax {
			util.RespFail(w, "文件大小不能超过 "+strconv.Itoa(filesImageMax/1024/1024)+" Mb")
			return
		}
	}
	filename := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	// TODO 将文件路径转换成url地址
	urlFileName := "mnt/" + filename
	dstfile, err := os.Create(urlFileName)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	// TODO将源文件内容copy到新文件
	_, err = io.Copy(dstfile, srcfile)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	dstfile.Close()
	var qn server.UploadTokenService
	qiniuUrl, err := qn.UploadQiNiuYun(urlFileName, urlFileName)
	if err != nil {
		util.RespFail(w, "OSS错误")
		return
	}
	err = os.Remove(urlFileName)
	if err != nil {

	}
	// 相应到前端
	util.RespOk(w, qiniuUrl, "")
}