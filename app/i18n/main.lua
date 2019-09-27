dict = {}

if (lang == 'en')
then
    dict = require(root .. "/dict/en")
else
    dict = require(root .. "/dict/tw")
end

-- for k, v in pairs(dict) do
--     print(k, v)
--     for a, b in pairs(v) do
--         print(a,b, type(a), type(b))
--     end
-- end

-- print(key)
keys = split(key)

output = dict

for k, v in pairs(keys) do
    if (output == nil) then
        return ''
    end
    output = output[v]
end

output = mapping(output)

return output
