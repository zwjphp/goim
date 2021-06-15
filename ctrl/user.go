package ctrl

import (
	"goim/util"
	"net/http"
)

func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	util.GenerateCaptchaHandler(w, r)
}
