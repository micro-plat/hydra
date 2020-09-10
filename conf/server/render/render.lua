

local rep = require("response")

function getStatus()
    return 302
end


function getContent()
    return rep.getContent()
end