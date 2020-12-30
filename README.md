# go-utils

![.github/workflows/go.yml](https://github.com/yourtion/go-utils/workflows/.github/workflows/go.yml/badge.svg?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/yourtion/go-utils.svg)](https://pkg.go.dev/github.com/yourtion/go-utils)

## Node.js + PM2

```go
// Init and get instance
pm := pm2.GetInstance()
// Add action
pm.AddAction("gc", func(opt interface{}) string {
    runtime.GC()
    return "ok"
})
// Add Metric
pm.AddMetric()
```
