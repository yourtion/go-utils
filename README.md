# go-utils

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
