package conf

//CircuitBreaker 熔断配置
type CircuitBreaker struct {
	ForceBreak      bool       `json:"force-break,omitempty"`
	Disable         bool       `json:"disable,omitempty"`
	SwitchWindow    int        `json:"switch-window,omitempty"`
	CircuitBreakers []*Breaker `json:"circuit-breakers,omitempty"`
}

//Breaker URL熔断配置项
type Breaker struct {
	URL              string `json:"url" valid:"ascii,required"`
	RequestPerSecond int    `json:"request-per-second,omitempty"`
	FailedPercent    int    `json:"failed-request,omitempty"`
	RejectPerSecond  int    `json:"reject-per-second,omitempty"`
	Disable          bool   `json:"disable,omitempty"`
}

//NewCircuitBreaker 构建熔断配置
func NewCircuitBreaker(window int) *CircuitBreaker {
	return &CircuitBreaker{
		SwitchWindow: window,
	}
}

//WithDisable 禁用熔断配置
func (c *CircuitBreaker) WithDisable() *CircuitBreaker {
	c.Disable = true
	return c
}

//WithEnable 启用熔断配置
func (c *CircuitBreaker) WithEnable() *CircuitBreaker {
	c.Disable = false
	return c
}

//WithForceBreak 强制熔断
func (c *CircuitBreaker) WithForceBreak(b bool) *CircuitBreaker {
	c.ForceBreak = b
	return c
}

//Append 添加URL的熔断配置
func (c *CircuitBreaker) Append(url string, requestPerSecond int, failedPercent int, rejectPreSecond int) *CircuitBreaker {
	c.CircuitBreakers = append(c.CircuitBreakers, &Breaker{
		URL:              url,
		RequestPerSecond: requestPerSecond,
		FailedPercent:    failedPercent,
		RejectPerSecond:  rejectPreSecond,
	})
	return c
}

//AppendAll 所有URL使用此熔断配置
func (c *CircuitBreaker) AppendAll(requestPerSecond int, failedPercent int, rejectPreSecond int) *CircuitBreaker {
	c.CircuitBreakers = append(c.CircuitBreakers, &Breaker{
		URL:              "*",
		RequestPerSecond: requestPerSecond,
		FailedPercent:    failedPercent,
		RejectPerSecond:  rejectPreSecond,
	})
	return c
}
