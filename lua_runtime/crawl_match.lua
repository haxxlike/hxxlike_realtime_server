local nk = require("nakama")

local M = {}

function M.match_init(context, setupstate)
    local gamestate = {
        presences = {}
    }

    local tickrate = 1 -- per sec
    local label = ""

    return gamestate, tickrate, label
end

function M.match_join_attempt(context, dispatcher, tick, state, presence, metadata)
    local acceptuser = true
    return state, acceptuser
end

function M.match_join(context, dispatcher, tick, state, presences)
    return state
end

function M.match_leave(context, dispatcher, tick, state, presences)
    return state
end

function M.match_loop(context, dispatcher, tick, state, messages)
    if tick % 60 == 0 then
        return state
    end
    -- state로 받은 url 접속 후 크롤링..
    return state
end

function M.match_terminate(context, dispatcher, tick, state, grace_seconds)
    return nil
end

function M.match_signal(context, dispatcher, tick, state, data)
    return state, "signal received: " .. data
end

return M
