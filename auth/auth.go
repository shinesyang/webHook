package auth

import (
	"net/http"
	"webHook/parse"

	"github.com/shinesyang/common"
)

/*认证失败*/

//Basic auth认证

func AuthFailed(w http.ResponseWriter, msg string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="grafana webhook"`)
	http.Error(w, msg, http.StatusUnauthorized)
}

/*中间件密码认证*/

func BaseAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		//fmt.Println(username, password, "----->")
		if !ok {
			AuthFailed(w, "401 Unauthorized!")
			return
		}

		u := parse.CONF.Username
		p := parse.CONF.Password

		if u == "" {
			u = "monitor_webhook"
		}

		if p == "" {
			p = "gx52cGjRVsck"
		}

		if username != u && password != p {
			AuthFailed(w, "401 Password error!")
			return
		}
		common.Logger.Infof("认证成功,用户: %s", username)
		h.ServeHTTP(w, r)
	})

	//http://monitor-msj.akbing.com/d/9CWBz0bik/node-exporter-dashboard?orgId=1&viewPanel=232&var-origin_prometheus=&var-job=jijia_node&var-hostname=All&var-instance=10.15.16.5:6091&var-device=All&var-interval=2m&var-maxmount=%2F&var-show_hostname=mao_jumpserver&var-total=9&var-addr=10.15.16.5&var-project=%E7%8C%AB

}
