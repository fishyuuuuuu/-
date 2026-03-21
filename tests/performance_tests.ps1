Param(
  [string]$BaseUrl = "http://localhost:8081"
)

Write-Host "== 功能探活 =="
Invoke-WebRequest -UseBasicParsing "$BaseUrl/api/product/list" | Select-Object StatusCode

Write-Host "`n== wrk 并发性能测试（需要本机安装 wrk） =="
Write-Host "示例：对商品列表接口进行 12 线程、400 连接、30 秒压测"
Write-Host "    wrk -t12 -c400 -d30s $BaseUrl/api/product/list"

Write-Host "`n== wrk 秒杀性能测试脚本示例 =="
Write-Host "确保 tests\seckill.lua 中填好 token 和商品ID，然后执行："
Write-Host "    wrk -t12 -c1000 -d30s -s tests/seckill.lua $BaseUrl/api/product/seckill"

