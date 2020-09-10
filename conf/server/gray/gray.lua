
local req = require("request")
local ips = {}
local upstream = ""


function getUpStream()
    return upstream
end



function go2UpStream() 
    local ip = req.getClientIP()
    if ips[ip] ~= nil then
        return true
    end
    return false
end