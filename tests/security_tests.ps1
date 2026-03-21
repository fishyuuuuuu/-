Param(
  [string]$BaseUrl = "http://localhost:8081"
)

Write-Host "== JWT 认证测试 =="
Write-Host "1) 使用错误密码尝试登录（预期 401 或 错误提示）"
try {
  Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/api/user/login" -Method POST -Body (@{phone="13300000000"; password="wrong"} | ConvertTo-Json) -ContentType "application/json"
} catch {
  Write-Host ("Response: " + $_.Exception.Message)
}

Write-Host "`n== 限流测试（接口频繁调用，预期部分 429） =="
for ($i=1; $i -le 20; $i++) {
  try {
    Invoke-WebRequest -UseBasicParsing "$BaseUrl/api/product/list" | Out-Null
    Write-Host "[$i] ok"
  } catch {
    if ($_.Exception.Response.StatusCode.value__ -eq 429) {
      Write-Host "[$i] hit rate limit (429)"
    } else {
      Write-Host "[$i] error: " $_.Exception.Message
    }
  }
}

