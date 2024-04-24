package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
)

func main() {
	// ファイル名を指定
	filename := "example.txt"

	// ファイルをオープン
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ファイルをオープンできませんでした:", err)
		return
	}
	defer file.Close()

	// ファイルからScannerを作成
	scanner := bufio.NewScanner(file)

	// ワーカーゴルーチンの数を定義
	numWorkers := 5

	// チャネルを作成してワーカーゴルーチンとメインゴルーチンが情報を共有する
	lines := make(chan string)
	var wg sync.WaitGroup

	// ワーカーゴルーチンを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				hash := sha256.Sum256([]byte(line))
				fmt.Printf("'%s' のSHA256ハッシュ: %x\n", line, hash)
			}
		}()
	}

	// ファイルから一行ずつ読み込んでチャネルに送信
	for scanner.Scan() {
		line := scanner.Text()
		lines <- line
	}

	// チャネルを閉じてゴルーチンが終了するようにする
	close(lines)

	// 全てのワーカーゴルーチンが終了するのを待つ
	wg.Wait()

	// エラーのチェック
	if err := scanner.Err(); err != nil {
		fmt.Println("ファイルの読み込み中にエラーが発生しました:", err)
	}
}
