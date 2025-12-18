package metrics

import "fmt"

type TodoMetrics struct{}

func NewTodoMetrics() *TodoMetrics {
	return &TodoMetrics{}
}

func (t TodoMetrics) Inc(key string, value int) {
	fmt.Println("metrics do nothing")
}
