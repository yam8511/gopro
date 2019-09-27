function TurnString( t )
    if (type(t) == 'table') then
        for k, v in pairs(t) do
            if (type(v) == 'table') then
                if type(k) == 'string' then
                    t[k] = TurnString(v)
                else
                    t[k] = nil
                    t[k .. ''] = TurnString(v)
                end
            else
                if type(k) == 'string' then
                    t[k] = v .. ''
                else
                    t[k] = nil
                    t[k .. ''] = v .. ''
                end
            end
        end

        return t
    end

    return t .. ''
end

return TurnString
