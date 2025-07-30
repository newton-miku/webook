wrk.method="POST"
wrk.headers["Content-Type"]="application/json"

local radom = math.random
local function uuid()
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and radom(0, 0xf) or radom(8, 0xb)
        return string.format('%x', v)
    end)
end

-- 初始化
function init(args)
    cnt = 0
    prefix = uuid()
end

function request()
    cnt = cnt + 1
    local url = "/users/signup"
    local body = '{"email":"'..prefix..'-'..cnt..'@qq.com","password":"theTestPwd@123"}'
    return wrk.format("POST", url, {}, body)
end