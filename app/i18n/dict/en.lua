dict = {
    ["hello"] = {
        [1] = "ok",
        [2] = "no",
        [3] = 4,
    }
}

return require(root .. "/helper")(dict)
