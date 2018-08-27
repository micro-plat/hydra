package context

type circuitBreakerParam struct {
	*inputParams
	ext map[string]interface{}
}

func (s *circuitBreakerParam) Clear() {
	s.inputParams.Clear()
	s.ext = nil
}

//IsOpen 熔断开发是否打开
func (s *circuitBreakerParam) IsOpen() bool {
	if v, ok := s.ext["__is_circuit_breaker_"].(bool); ok {
		return v
	}
	return false
}
func (s *circuitBreakerParam) GetDefStatus() int {
	return 503
}
func (s *circuitBreakerParam) GetDefContent() string {
	return ""
}
