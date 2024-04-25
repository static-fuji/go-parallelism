package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"
)

type Item struct {
	IntValue    int
	StringValue string
}

type ItemList []Item

func main() {
	// ファイル名を指定
	filename := "test.txt"

	// マップを作成
	lineMap := make(map[int]string)

	items := []Item{}

	// ファイルをオープン
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ファイルをオープンできませんでした:", err)
		return
	}
	defer file.Close()

	// ファイルからScannerを作成
	scanner := bufio.NewScanner(file)

	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		lineMap[lineNumber] = line
		lineNumber++
	}

	// チャネルを作成
	ch := make(chan struct {
		key   int
		value string
	}, lineNumber)

	results := make(chan struct {
		worker   int
		newkey   int
		newvalue string
	}, lineNumber)

	// WaitGroupを作成してゴルーチンの完了を待つ
	var wg sync.WaitGroup
	numWorkers := 5

	// マップ内のキーと値をチャネルに送信するゴルーチンを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for pair := range ch {
				results <- struct {
					worker   int
					newkey   int
					newvalue string
				}{i, pair.key, pair.value}
				//time.Sleep(time.Second)
			}
		}(i)
	}

	for key, value := range lineMap {
		ch <- struct {
			key   int
			value string
		}{key, value}
	}

	close(ch)

	// すべてのゴルーチンが完了するのを待つ
	go func() {
		wg.Wait() // すべてのキーと値が送信されたらチャネルを閉じる
		close(results)
	}()

	// チャネルからキーと値を受信して出力
	for pair := range results {
		item := Item{
			IntValue:    pair.newkey,
			StringValue: pair.newvalue,
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].IntValue < items[j].IntValue
	})

	for i := range items {
		fmt.Printf("key: %d, word: %s\n", items[i].IntValue, items[i].StringValue)
	}
}
