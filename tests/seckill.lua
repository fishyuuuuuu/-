wrk.method = "POST"

-- 在这里填入实际的商品ID和 JWT token
local productId = 1
local token = "REPLACE_WITH_REAL_JWT_TOKEN"

request = function()
  local body = string.format('{"productId": %d, "captchaId": "wrk", "captchaStr": "1234"}', productId)
  local headers = {
    ["Content-Type"] = "application/json",
    ["Authorization"] = "Bearer " .. token
  }
  return wrk.format(nil, nil, headers, body)
end

