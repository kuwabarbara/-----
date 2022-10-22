package main

//https://www.jamsystem.com/ancdic/index.html

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

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
	CTime int64  `json:"ctime"`
}

// メインプログラム - サーバーを起動する --- (*3)
func main() {
	/*csvFile, _ := os.Open("kuwa/f.csv")
	reader := csv.NewReader(transform.NewReader(csvFile, japanese.ShiftJIS.NewDecoder()))
	//reader := csv.NewReader(csvFile)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(line[0] + " " + line[1])
	}*/

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
		// fmt.Println("まだファイルないよ")
		fmt.Printf("%+v\n", p)

		fmt.Println(score)
		w.Write([]byte(getFormLogs(p, r.URL.Path[1:])))
	} else {
		fmt.Println("まだファイルないよ", err)
	}

	//fmt.Printf("%+v\n", p)
	//fmt.Println(score)
	//w.Write([]byte(getFormLogs(p, r.URL.Path[1:])))
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

	var search_flag bool
	var result string

	search_flag, result = searchDictionary(r.Form["name"][0])

	//もし検索が発見できなかったら
	if search_flag == false {
		http.Redirect(w, r, "/"+r.Form["logname"][0], 302)
		return
	}

	result = "aaa"
	fmt.Println(result)

	fmt.Println("びびびびい")

	//書き込まれた内容をjsonファイルに書き込む
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
