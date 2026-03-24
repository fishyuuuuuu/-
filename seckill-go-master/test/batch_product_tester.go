package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TestResult struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	ExpectedPrice float64 `json:"expected_price"`
	ActualPrice   float64 `json:"actual_price"`
	Success     bool    `json:"success"`
	Message     string  `json:"message"`
}

type Order struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	UserName    string  `json:"user_name"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	OrderAmount float64 `json:"order_amount"`
	Status      int     `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type OrderListResponse struct {
	Code int      `json:"code"`
	Data []Order  `json:"data"`
}

var products = map[int]struct {
	name  string
	price float64
}{
	1: {"iPhone 15 Pro", 7999.99},
	2: {"MacBook Pro 14", 12999.99},
	3: {"AirPods Pro 2", 1899.99},
	4: {"iPad Pro 12.9", 8999.99},
	5: {"Sony WH-1000XM5", 2999.99},
	6: {"Nintendo Switch OLED", 2199.99},
	7: {"Dyson V15 Detect", 4599.99},
	8: {"Nike Air Max 270", 899.99},
}

func main() {
	baseURL := "http://localhost:8081"
	
	fmt.Println("========================================")
	fmt.Println("    批量商品测试方案 - 开始执行")
	fmt.Println("========================================")
	fmt.Printf("测试时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("测试商品数量: %d\n", len(products))
	fmt.Println("----------------------------------------")
	
	var results []TestResult
	successCount := 0
	failCount := 0
	
	for productID, expected := range products {
		fmt.Printf("\n[测试商品 %d] %s\n", productID, expected.name)
		fmt.Println("----------------------------------------")
		
		userID := 1000 + productID
		
		result := TestResult{
			ProductID:     productID,
			ProductName:   expected.name,
			ExpectedPrice: expected.price,
		}
		
		seckillReq := map[string]interface{}{
			"productId": productID,
			"userId":    userID,
		}
		seckillBody, _ := json.Marshal(seckillReq)
		
		fmt.Printf("  1. 发送秒杀请求 (商品ID: %d, 用户ID: %d)...\n", productID, userID)
		seckillResp, err := http.Post(
			baseURL+"/api/product/seckill",
			"application/json",
			bytes.NewBuffer(seckillBody),
		)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("秒杀请求失败: %v", err)
			results = append(results, result)
			failCount++
			fmt.Printf("  ❌ 秒杀请求失败: %v\n", err)
			continue
		}
		seckillResp.Body.Close()
		
		fmt.Printf("  2. 秒杀请求成功，等待订单创建...\n")
		time.Sleep(500 * time.Millisecond)
		
		fmt.Printf("  3. 查询订单列表验证商品信息...\n")
		orderResp, err := http.Get(baseURL + "/api/order/admin/list")
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("查询订单失败: %v", err)
			results = append(results, result)
			failCount++
			fmt.Printf("  ❌ 查询订单失败: %v\n", err)
			continue
		}
		
		body, _ := io.ReadAll(orderResp.Body)
		orderResp.Body.Close()
		
		var orderList OrderListResponse
		if err := json.Unmarshal(body, &orderList); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("解析订单数据失败: %v", err)
			results = append(results, result)
			failCount++
			fmt.Printf("  ❌ 解析订单数据失败: %v\n", err)
			continue
		}
		
		var targetOrder *Order
		for i := range orderList.Data {
			if orderList.Data[i].ProductID == productID && orderList.Data[i].UserID == userID {
				targetOrder = &orderList.Data[i]
				break
			}
		}
		
		if targetOrder == nil {
			result.Success = false
			result.Message = "未找到对应订单"
			results = append(results, result)
			failCount++
			fmt.Printf("  ❌ 未找到对应订单\n")
			continue
		}
		
		result.ActualPrice = targetOrder.OrderAmount
		
		if targetOrder.ProductName != expected.name {
			result.Success = false
			result.Message = fmt.Sprintf("商品名称不匹配: 期望 '%s', 实际 '%s'", expected.name, targetOrder.ProductName)
			results = append(results, result)
			failCount++
			fmt.Printf("  ❌ 商品名称不匹配: 期望 '%s', 实际 '%s'\n", expected.name, targetOrder.ProductName)
			continue
		}
		
		if targetOrder.OrderAmount != expected.price {
			result.Success = false
			result.Message = fmt.Sprintf("商品价格不匹配: 期望 %.2f, 实际 %.2f", expected.price, targetOrder.OrderAmount)
			results = append(results, result)
			failCount++
			fmt.Printf("  ❌ 商品价格不匹配: 期望 %.2f, 实际 %.2f\n", expected.price, targetOrder.OrderAmount)
			continue
		}
		
		result.Success = true
		result.Message = "测试通过"
		results = append(results, result)
		successCount++
		fmt.Printf("  ✅ 测试通过! 商品名称: %s, 价格: %.2f\n", targetOrder.ProductName, targetOrder.OrderAmount)
	}
	
	fmt.Println("\n========================================")
	fmt.Println("           测试结果汇总")
	fmt.Println("========================================")
	fmt.Printf("测试商品总数: %d\n", len(products))
	fmt.Printf("通过数量: %d\n", successCount)
	fmt.Printf("失败数量: %d\n", failCount)
	fmt.Printf("通过率: %.2f%%\n", float64(successCount)/float64(len(products))*100)
	fmt.Println("----------------------------------------")
	
	if failCount > 0 {
		fmt.Println("\n失败详情:")
		for _, r := range results {
			if !r.Success {
				fmt.Printf("  - 商品ID %d (%s): %s\n", r.ProductID, r.ProductName, r.Message)
			}
		}
	}
	
	fmt.Println("\n========================================")
	fmt.Println("    批量商品测试方案 - 执行完成")
	fmt.Println("========================================")
}
