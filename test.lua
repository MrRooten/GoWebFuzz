function dump(o)
    if type(o) == 'table' then
        local s = '{ '
        for k,v in pairs(o) do
            if type(k) ~= 'number' then k = '"'..k..'"' end
            s = s .. '['..k..'] = ' .. dump(v) .. ','
        end
        return s .. '} '
    else
        return tostring(o)
    end
end
local http = require("http")
result = 0

if(res.status_code == 200)
then
    print(dump(res))
    result = 1
    print("It's Vulnable")
end

