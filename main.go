package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math/rand"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func generateEmoji(EmojiText string, DownloadDirectory string) []string {

	//絵文字のURL格納用スライスの定義
	var EmojiURL []string
	//	EmojiURL := []string

	// ログファイルの準備
	NowTime := time.Now()
	LogFileName := NowTime.Format("20060102_150405") + ".txt"
	logfile, err := os.OpenFile("log/"+LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()

	// テキストファイルの読み込み,\rが含まれている場合があるので置換しておく
	f := EmojiText
	f = strings.Replace(f, "\r", "", -1)
	s := strings.Split(f, "\n")

	cnt := 0
	// テキストファイルを行ごとに読み込んでループ
	for _, str := range s {
		Message := str
		OutputFileName := DownloadDirectory + "/" + strconv.Itoa(cnt) + ".png"

		// カンマ区切りなら後半をファイル名にする
		if strings.Contains(Message, ",") {
			slice := strings.Split(Message, ",")
			Message = slice[0]
			OutputFileName = DownloadDirectory + "/" + slice[1] + ".png"
		}

		// 画像サイズを決める
		img := image.NewRGBA(image.Rect(0, 0, 128, 128))

		// 画像背景色を決める
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}

		// フォントファイルの読み込みとパース
		ftBin, err := ioutil.ReadFile("static/fonts/NotoSansCJKjp-Bold.otf")
		if err != nil {
			log.Fatalf("failed to load font: %s", err.Error())
		}
		ft, err := opentype.Parse(ftBin)
		if err != nil {
			log.Fatalf("failed to parse font: %s", err.Error())
		}

		// 文字数を判定する
		if len(Message)/len([]rune(Message)) != 3 {
			fmt.Fprintln(logfile, Message+" ERROR[ごめんなさい！全角文字を使ってね！]")
			continue
		}

		// 文字数による描画設定(フォントサイズ、行数、描画座標)
		var CharSize float64 = 64
		CharSize = 64
		Row1Xvalue, Row1Yvalue, Row2Xvalue, Row2Yvalue := 0, 0, 0, 0
		Row1Text, Row2Text := "", ""
		if len([]rune(Message)) == 2 {
			CharSize = 64
			Row1Xvalue = 0 * 64
			Row1Yvalue = 85 * 64
			Row1Text = Message
		} else if len([]rune(Message)) == 3 {
			CharSize = 42
			Row1Xvalue = 0 * 64
			Row1Yvalue = 75 * 64
			Row1Text = Message
		} else if len([]rune(Message)) == 4 {
			CharSize = 64
			Row1Xvalue = 0 * 64
			Row1Yvalue = 60 * 64
			Row2Xvalue = 0 * 64
			Row2Yvalue = 120 * 64
			Row1Text = Message[0:6]
			Row2Text = Message[6:12]
		} else if len([]rune(Message)) == 6 {
			CharSize = 42
			Row1Xvalue = 0 * 64
			Row1Yvalue = 55 * 64
			Row2Xvalue = 0 * 64
			Row2Yvalue = 105 * 64
			Row1Text = Message[0:9]
			Row2Text = Message[9:18]
		} else if len([]rune(Message)) == 8 {
			CharSize = 32
			Row1Xvalue = 0 * 64
			Row1Yvalue = 50 * 64
			Row2Xvalue = 0 * 64
			Row2Yvalue = 100 * 64
			Row1Text = Message[0:12]
			Row2Text = Message[12:24]
		} else {
			fmt.Fprintln(logfile, Message+" ERROR[ごめんなさい！全角2,3,4,6,8文字以外は対応してません！]")
			continue
		}

		// 文字色の候補
		color := []*color.RGBA{
			{255, 0, 0, 255},
			{255, 219, 0, 255},
			{73, 255, 0, 255},
			{0, 255, 146, 255},
			{0, 146, 255, 255},
			{73, 0, 255, 255},
			{255, 0, 219, 255},
		}

		// 文字色をランダムにあてがう
		rand.Seed(time.Now().UnixNano())
		RandomNumber := rand.Intn(7)

		opt := opentype.FaceOptions{
			Size:    CharSize,
			DPI:     72,
			Hinting: font.HintingNone,
		}
		face, err := opentype.NewFace(ft, &opt)

		// 1行目の描画
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color[RandomNumber]),
			Face: face,
			Dot:  fixed.Point26_6{fixed.Int26_6(Row1Xvalue), fixed.Int26_6(Row1Yvalue)},
		}
		d.DrawString(Row1Text)

		// 2行目の描画
		if len([]rune(Message)) >= 4 {
			d2 := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(color[RandomNumber]),
				Face: face,
				Dot:  fixed.Point26_6{fixed.Int26_6(Row2Xvalue), fixed.Int26_6(Row2Yvalue)},
			}
			d2.DrawString(Row2Text)
		}

		// ファイル出力
		file, err := os.Create(OutputFileName)
		if err != nil {
			panic(err.Error())
		}
		defer file.Close()
		fmt.Fprintln(logfile, Message+" "+OutputFileName)

		if err := png.Encode(file, img); err != nil {
			panic(err.Error())
		}

		//ファイル名をスライスに追加
		EmojiURL = append(EmojiURL, OutputFileName)

		cnt++
	}
	//    if s.Err() != nil {
	//        // non-EOF error.
	//        log.Fatal(s.Err())
	//    }
	return EmojiURL
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		"template/index.gohtml",
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

func handleInputConfirm(w http.ResponseWriter, req *http.Request) {

	EmojiText := req.FormValue("emoji_text")
	EmojiTextColor := req.FormValue("emoji_text_color")

	fmt.Println(EmojiText)
	fmt.Println(EmojiTextColor)

	// 現在時刻でディレクトリを作成
	TimeNow := time.Now().Format("20060102150405")
	DownloadDirectory := "tmp/download/" + TimeNow
	if err := os.Mkdir(DownloadDirectory, 0777); err != nil {
		fmt.Println(err)
	}

	// 絵文字作成関数を実行して、戻り値としてURLを含むスライスを受け取る
	EmojiURL := generateEmoji(EmojiText, DownloadDirectory)

	// downloadディレクトリ内を確認
	files, _ := os.ReadDir("./tmp/download")
	TimeNowYYYYMMDDHH := time.Now().Format("2006010215")

	for _, f := range files {
		fmt.Println(f.Name()[:10] + "      " + TimeNowYYYYMMDDHH)

		i, _ := strconv.Atoi(f.Name()[:10])
		j, _ := strconv.Atoi(TimeNowYYYYMMDDHH)

		// これだと日付かわる瞬間に前日分が全削除されちゃうから危険
		if i <= j-2 {
			fmt.Println(f.Name() + "は削除します")
			if err := os.RemoveAll("tmp/download/" + f.Name()); err != nil {
				fmt.Println(err)
			}
		}

	}

	t := template.Must(template.ParseFiles(
		"template/input_confirm.gohtml",
	))
	if err := t.Execute(w, struct {
		EmojiText      string
		EmojiTextColor string
		EmojiURL       []string
	}{
		EmojiText,
		EmojiTextColor,
		EmojiURL,
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
	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("tmp/"))))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/input_confirm", handleInputConfirm)
	http.HandleFunc("/secret", handleSecret)
	log.Printf("ポート %s で待ち受けを開始します...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("サーバーが異常終了しました: %v", err)
	}
}
