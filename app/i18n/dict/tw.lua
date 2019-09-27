local dict = {
    ["hello"] = {
        [1] = "沒問題",
        [2] = "不行~ {{ .Name }}",
        [3] = 4,
    }
}

return require(root .. "/helper")(dict)
