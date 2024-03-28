package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("ファイル読み取り処理を開始します")
	// fileを開く
	f, err := os.Open("test.txt")
	// 読み取り時の例外処理
	if err != nil {
		fmt.Println("error")
	}
	// close
	defer f.Close()

	// byte型スライスの作成
	buf := make([]byte, 1024)
	for {
		// nはバイト数を示す
		n, err := f.Read(buf)
		// バイト数が0になることは、読み取り終了を示す
		if n == 0 {
			break
		}
		if err != nil {
			break
		}
		// バイト型スライスを文字列型に変換してファイルの内容を出力
		fmt.Println(string(buf[:n]))
	}
}
