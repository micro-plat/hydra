
local filter=[]
local ips={}
local upstream=""

function main() 
    local ip=req.get_client_ip()
    if ips[ip]~=nil then
        return ""
    end
    return upstream
end