package main

//https://www.jamsystem.com/ancdic/index.html

import (
	"fmt"
	"strconv"

	//"html"
	"net/http"
	//"time"
	json "github.com/takoyaki-3/go-json"
)

const logFile = "logs" // データの保存先 --- (*1)

var score int

// Log 掲示板に保存するデータを構造体で定義 --- (*2)
type Log struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Body  string `json:"body"`
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

	fmt.Println("logs" + r.URL.Path[1:] + ".json")

	var p []Log
	if err := json.LoadFromPath("logs" + r.URL.Path[1:] + ".json", &p);err!=nil{
		fmt.Println("ファイル読み込み時のerror：",err)
		return
	}

	fmt.Printf("%+v\n", p)
	fmt.Println(score)
	w.Write([]byte(getFormLogs(p, r.URL.Path[1:])))
}

// gateで書き込まれた内容を処理する
func writegateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // フォームを解析 --- (*10)
	var log Log
	log.Name = r.Form["name"][0]
	if log.Name == "" {
		log.Name = "名無し"
	}
	fmt.Print("aaaaaa")
	fmt.Print(log.Name)
	addLog(log, log.Name)
	http.Redirect(w, r, "/", 302) // リダイレクト --- (*13)
}

// 書き込まれた内容を処理する
func writelogHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // フォームを解析 --- (*10)

	//lognameがlogのファイルの名前
	//nameが入力内容

	fmt.Print("いｄｓ" + r.Form["logname"][0] + "あｊｓｊ")
	fmt.Print("いｄｓ" + r.Form["name"][0] + "あｊｓｊ")
	//fmt.Print("saasasa" + r.URL.Path[1:] + "fefsf")

	//書き込まれた内容をjsonファイルに書き込む

	fmt.Println("びびびびい")

	/*var log Log
	log.Name = r.Form["logname"][0]
	log.Body = r.Form["name"][0]

	fmt.Println("ぎゃぎゃがｙ")

	if log.Name == "" {
		log.Name = "名無し"
	}
	logs := loadLogs2(r.Form["logname"][0]) // 既存のデータを読み出し --- (*11)

	fmt.Println("じゃじゃじゃっじゃ")

	log.ID = len(logs) + 1
	log.CTime = time.Now().Unix()
	logs = append(logs, log)             // 追記 --- (*12)
	saveLogs(logs, r.Form["logname"][0]) // 保存
	*/

	var log Log
	log.Name = r.Form["name"][0]
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
	return "<div>" + namae + "aaaa" + log.Name + "xxx" + strconv.Itoa(score) + "aaa </div>" +
		"<div><form action='/writelog' method='POST'>" +
		"<input type='hidden' name='logname' value='" + namae + "'>" +
		"名前: <input type='text' name='name'><br>" +
		"<input type='submit' value='書込'>" +
		"</form></div><hr>"
}

// func saveLogs(logs []Log, namae string) {
// 	json.DumpToFile(&logs, logFile+namae+".json")
// }

func addLog(log Log, namae string) {
	var logs []Log
	json.LoadFromPath(logFile+namae+".json",&logs)
	logs = append(logs, log)
	json.DumpToFile(logs,logFile+namae+".json")
}

func loadLogs2(namae string) []Log {
	var logs []Log
	json.LoadFromPath(logFile + namae + ".json", &logs)
	return logs
}
