package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"sync"
)

// ソートを行うための構造体 (行数，元のテキスト，ハッシュ化済みのテキスト)
type sorted struct {
	lineNum int
	line    string
	hashHex string
}

// ソートを行うためのスライス
type ItemList []sorted

func main() {
	// ファイル名を指定
	filename := "test.txt"

	// 文字列と対応する行数（key）でマップを作成
	lineMap := make(map[int]string)

	//作成するゴルーチン（グリーンスレッド）の数（=ファイルの行の総数）
	numWorkers := 0

	// 読み取った行と行数を受け取るチャネル (行数，読み込んだテキスト)
	lines := make(chan struct {
		lineNum int
		line    string
	})

	//ゴルーチンから結果を受け取るチャネル (行数，元のテキスト，ハッシュ化済みのテキスト)
	results := make(chan struct {
		LineNum int
		line    string
		hash    string
	})

	//ソートを行うためのスライス
	sortedLists := []sorted{}

	// ファイルをオープン
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ファイルをオープンできませんでした:", err)
		return
	}
	defer file.Close()

	// ファイルからScannerを作成
	scanner := bufio.NewScanner(file)

	//読み取った文字列と対応する行数をmapに代入，行の総数をカウント
	for scanner.Scan() {
		numWorkers++
		line := scanner.Text()
		lineMap[numWorkers] = line
	}

	// WaitGroupを作成
	var wg sync.WaitGroup

	// 行数分のゴルーチンを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range lines {
				//ハッシュ化してhexダンプする
				hash := sha256.Sum256([]byte(j.line))
				hashHex := hex.EncodeToString(hash[:])
				//resultsチャネルに結果を送信
				results <- struct {
					LineNum int
					line    string
					hash    string
				}{j.lineNum, j.line, hashHex}
			}
		}()
	}

	//チャネルに読み込んだデータのmapを送信
	for lineNum, line := range lineMap {
		lines <- struct {
			lineNum int
			line    string
		}{lineNum, line}
	}

	//チャネルを閉じる
	close(lines)

	// すべてのゴルーチンが完了するのを待つ
	go func() {
		wg.Wait()
		// すべてのキーと値が送信されたらチャネルを閉じる
		close(results)
	}()

	//ソート用スライスにハッシュ化済みのデータを格納
	for i := range results {
		item := sorted{
			lineNum: i.LineNum,
			line:    i.line,
			hashHex: i.hash,
		}
		sortedLists = append(sortedLists, item)
	}

	//keyの順番にソート
	sort.Slice(sortedLists, func(i, j int) bool {
		return sortedLists[i].lineNum < sortedLists[j].lineNum
	})

	//ターミナルに表示
	for i := range sortedLists {
		fmt.Printf("%d行目, 元のテキスト: %s , ハッシュ化: %s \n", sortedLists[i].lineNum, sortedLists[i].line, sortedLists[i].hashHex)
	}
}
