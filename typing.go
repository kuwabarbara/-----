package main

//https://www.jamsystem.com/ancdic/index.html

import (
	//"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"

	//"log"
	"os"
	//"text/template"

	"strconv"

	//"html"
	"net/http"
	//"time"
	json "github.com/takoyaki-3/go-json"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const logFile = "logs" // データの保存先 --- (*1)

var score int

// Log 掲示板に保存するデータを構造体で定義 --- (*2)
type Log struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Body  string `json:"body"`
	Kukki string `json:"kukki"`
	CTime int64  `json:"ctime"`
}

// メインプログラム - サーバーを起動する --- (*3)
func main() {

	score = 0

	println("server - http://localhost:8888")
	// URIに対応するハンドラを登録 --- (*4)
	http.HandleFunc("/", gateHandler)
	http.HandleFunc("/writegate", writegateHandler)
	http.HandleFunc("/writelog", writelogHandler)
	//http.HandleFunc("/show", showHandler)
	//http.HandleFunc("/write", writeHandler)
	// サーバーを起動 --- (*5)
	http.ListenAndServe(":8888", nil)
}

// 辞書に乗っているか調べる
func searchDictionary(name string) (bool, string) {
	small_alphabet := [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	big_alphabet := [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	if name == "" {
		return false, "名前の部分が空白になっているね！"
	}

	for i := 0; i < len(small_alphabet); i++ {
		if string(name[0]) == small_alphabet[i] || string(name[0]) == big_alphabet[i] {
			csvFile, _ := os.Open("kuwa/" + small_alphabet[i] + ".csv")
			reader := csv.NewReader(transform.NewReader(csvFile, japanese.ShiftJIS.NewDecoder()))
			//reader := csv.NewReader(csvFile)

			for {
				line, err := reader.Read()
				if err == io.EOF {
					break
				}
				if name == line[0] {
					fmt.Println("発見")
					fmt.Println(line[1])
					return true, line[1]
				}

				//fmt.Println(line[0] + " " + line[1])

			}
		}
	}

	fmt.Println("見つからなかった")
	return false, "入力した単語は検索できなかったな！"

}

// 最初の画面
func gateHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)

	//urlに余分なものがついてない場合
	if r.URL.Path == "/" {
		w.Write([]byte(getFormGate()))
		return
	}

	// favicon.icoだったら読み込まれたらなにもしない
	if r.URL.Path[1:] == "favicon.ico" {
		return
	}

	var p []Log
	if err := json.LoadFromPath("logs"+r.URL.Path[1:]+".json", &p); err == nil {
		//fmt.Printf("%+v\n", p)
		fmt.Println(score)

		w.Write([]byte(getFormLogs(p, r.URL.Path[1:])))
	} else {
		fmt.Println("まだファイルないよ", err)
	}

}

// gateで書き込まれた内容を処理する
func writegateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // フォームを解析 --- (*10)
	var log Log
	//log.Name = r.Form["name"][0]
	log.Name = "一番最初のlog"
	if log.Name == "" {
		log.Name = "名無し"
	}
	fmt.Print("aaaaaa")
	fmt.Print(log.Name)
	addLog(log, r.Form["name"][0])
	http.Redirect(w, r, "/", 302) // リダイレクト --- (*13)
}

// 書き込まれた内容を処理する
func writelogHandler(w http.ResponseWriter, r *http.Request) {
	//クッキーを設定
	//そのためにまずクッキーを取得
	cookiecookie, err := r.Cookie("hoge")
	if err != nil {
		//log.Fatal("Cookie: ", err)
		//クッキーが存在しなかった場合作成
		cookie := &http.Cookie{
			Name:  "hoge",
			Value: "ransu", //ここで本当は乱数を入れたい
		}
		http.SetCookie(w, cookie)
	} else {
		//クッキーは存在している
		fmt.Println("くわー！")
		fmt.Println(cookiecookie.Name)
		//http.Redirect(w, r, "/"+r.Form["logname"][0], 302)
		//return
	}

	r.ParseForm() // フォームを解析 --- (*10)

	//lognameがlogのファイルの名前
	//nameが入力内容
	fmt.Print("いｄｓ" + r.Form["logname"][0] + "あｊｓｊ")
	fmt.Print("いｄｓ" + r.Form["name"][0] + "あｊｓｊ")
	//fmt.Print("saasasa" + r.URL.Path[1:] + "fefsf")

	//現在のlogを取得する
	var last string //最後に入力された内容
	var p []Log
	if err := json.LoadFromPath("logs"+r.Form["logname"][0]+".json", &p); err == nil {
		log := p[len(p)-1]
		last = log.Name
	} else {
		fmt.Println("このjsonファイル開けないよ", err)
		http.Redirect(w, r, "/"+r.Form["logname"][0], 302)
		return
	}

	//打ち込んだ文字と前回打ち込んだ文字とでしりとりになっているか
	if last != "一番最初のlog" && last[len(last)-1] != r.Form["name"][0][0] {
		http.Redirect(w, r, "/"+r.Form["logname"][0], 302)
		return
	}

	var search_flag bool
	var result string

	search_flag, result = searchDictionary(r.Form["name"][0])

	//もし検索が発見できなかったら
	if search_flag == false {
		http.Redirect(w, r, "/"+r.Form["logname"][0], 302)
		return
	}

	fmt.Println(result)

	fmt.Println("びびびびい")

	//書き込まれた内容をjsonファイルに書き込む

	//クッキーの取得を行う
	cookiecookie2, err := r.Cookie("hoge")
	if err != nil {
		//log.Fatal("Cookie: ", err)
		fmt.Println("クッキーが取得できない")
		http.Redirect(w, r, "/"+r.Form["logname"][0], 302)
		return
	}
	v := cookiecookie2.Value

	var log Log
	log.Name = r.Form["name"][0]
	log.Body = result
	log.Kukki = v
	if log.Name == "" {
		log.Name = "名無し"
	}

	addLog(log, r.Form["logname"][0]) // 保存

	score += 10

	http.Redirect(w, r, "/"+r.Form["logname"][0], 302) // リダイレクト --- (*13)
}

// gate用の書き込みフォーム
func getFormGate() string {
	return "<div><form action='/writegate' method='POST'>" +
		"名前: <input type='text' name='name'><br>" +
		"<input type='submit' value='書込'>" +
		"</form></div><hr>"
}

// logの内容を読み込んで表示する
func getFormLogs(logs []Log, namae string) string {
	log := logs[len(logs)-1]
	return "<div>" + namae + "    " + log.Name + "   " + log.Body + "  " + strconv.Itoa(score) + "    </div>" +
		"<div><form action='/writelog' method='POST'>" +
		"<input type='hidden' name='logname' value='" + namae + "'>" +
		"名前: <input type='text' name='name'><br>" +
		"<input type='submit' value='書込'>" +
		"</form></div><hr>"
}

func addLog(log Log, namae string) {
	var logs []Log
	json.LoadFromPath(logFile+namae+".json", &logs)
	logs = append(logs, log)
	json.DumpToFile(logs, logFile+namae+".json")
}

func loadLogs2(namae string) []Log {
	var logs []Log
	json.LoadFromPath(logFile+namae+".json", &logs)
	return logs
}
