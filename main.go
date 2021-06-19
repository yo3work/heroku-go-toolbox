package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		"template/index.gohtml",
		"template/index02.gohtml",
		"template/_about.gohtml",
		"template/_copyright.gohtml",
		"template/_footer.gohtml",
		"template/_masthead.gohtml",
		"template/_navigation.gohtml",
	))
	if err := t.Execute(w, struct {
		UserName  string
		Time      time.Time
		RecordsId []int
	}{
		"ゲスト",
		time.Now(),
		[]int{1, 2, 3, 4, 5, 6},
	}); err != nil {
		log.Printf("テンプレート %s の実行に失敗！: %v", t.Name(), err)
		http.Error(w, "内部エラーです", http.StatusInternalServerError)
	}

}

func handleSecret(w http.ResponseWriter, r *http.Request) {
	user, password, _ := r.BasicAuth()
	HerokuEnvUser := os.Getenv("USER")
	HerokuEnvPassword := os.Getenv("PASSWORD")
	if user != HerokuEnvUser || password != HerokuEnvPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "認証に失敗しました", http.StatusUnauthorized)
		return
	}
	log.Printf("%s %s", r.Method, r.RequestURI)
	w.Write([]byte("秘密のページです！"))
}

func main() {
	port := os.Getenv("PORT") //実行時に Heroku が指定するポート番号を取得
	if len(port) == 0 {
		port = "8080"
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/secret", handleSecret)
	log.Printf("ポート %s で待ち受けを開始します...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("サーバーが異常終了しました: %v", err)
	}
}
