package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	// 测试配置
	testCount := 100  // 总请求数
	concurrency := 10 // 并发数
	targetURL := "http://localhost:8081/api/product/seckill"
	token := "your-token-here" // 替换为实际的token

	// 统计结果
	successCount := 0
	errorCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 开始时间
	startTime := time.Now()

	// 并发测试
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// 每个并发 worker 处理的请求数
			requestsPerWorker := testCount / concurrency
			if workerID < testCount%concurrency {
				requestsPerWorker++
			}

			for j := 0; j < requestsPerWorker; j++ {
				// 创建请求
				req, err := http.NewRequest("POST", targetURL, strings.NewReader(`{"productId": 1, "captchaId": "test", "captchaStr": "test"}`))
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					continue
				}

				// 设置请求头
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				// 发送请求
				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					continue
				}

				// 检查响应
				if resp.StatusCode == http.StatusOK {
					mu.Lock()
					successCount++
					mu.Unlock()
				} else {
					mu.Lock()
					errorCount++
					mu.Unlock()
				}

				resp.Body.Close()
				time.Sleep(10 * time.Millisecond) // 稍微延迟，避免请求过于密集
			}
		}(i)
	}

	// 等待所有请求完成
	wg.Wait()

	// 结束时间
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// 打印结果
	fmt.Printf("测试完成！\n")
	fmt.Printf("总请求数: %d\n", testCount)
	fmt.Printf("成功数: %d\n", successCount)
	fmt.Printf("失败数: %d\n", errorCount)
	fmt.Printf("耗时: %v\n", duration)
	fmt.Printf("QPS: %.2f\n", float64(testCount)/duration.Seconds())
}
