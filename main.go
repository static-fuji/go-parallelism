package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// ファイル名を指定
	filename := "test.txt"

	// ファイルをオープン
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ファイルをオープンできませんでした:", err)
		return
	}
	defer file.Close()

	// ファイルからScannerを作成
	scanner := bufio.NewScanner(file)

	// マップを作成して行を格納
	ch := make(chan map[int]string)
	lineMap := make(map[int]string)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		lineMap[lineNumber] = line
		lineNumber++
	}

	// エラーのチェック
	if err := scanner.Err(); err != nil {
		fmt.Println("ファイルの読み込み中にエラーが発生しました:", err)
		return
	}

	// マップの内容を出力
	for lineNumber, line := range lineMap {
		go func() {
			fmt.Println(lineNumber, line, lineMap)
			ch <- lineMap
		}()
		<-ch
	}
	close(ch)
}
