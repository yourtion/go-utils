package metric

// Metric
type Metric struct {
	Name     string         `json:"name"`  // 计数器名称
	Value    float64        `json:"value"` // 计数器指标
	Category string         `json:"type"`  // 计数器类型
	Unit     string         `json:"unit"`  // 计数器单位（只用于显示）
	Function func() float64 `json:"-"`     // 计数器获取参数方法
}

// Get 获取当前值（如果有计算方法则先执行计算方法）
func (metric *Metric) Get() float64 {
	if metric.Function != nil {
		metric.Value = metric.Function()
	}
	return metric.Value
}

// Set 设置当前值
func (metric *Metric) Set(value float64) {
	metric.Value = value
}

// CreateMetric 创建普通 Metric
func CreateMetric(name string, category string, unit string) *Metric {
	return &Metric{
		Name:     name,
		Category: category,
		Unit:     unit,
		Value:    0,
	}
}

// CreateFuncMetric 创建使用 function 计算的 Metric
func CreateFuncMetric(name string, category string, unit string, cb func() float64) *Metric {
	return &Metric{
		Name:     name,
		Category: category,
		Unit:     unit,
		Value:    0,
		Function: cb,
	}
}
