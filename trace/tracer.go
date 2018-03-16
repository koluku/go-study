package trace

import (
	"fmt"
	"io"
)

// Tracer ログの構造体
type Tracer struct {
	out io.Writer
}

// Trace ログを出力
func (t *Tracer) Trace(a ...interface{}) {
	if t == nil || t.out == nil {
		return
	}
	fmt.Fprintln(t.out, a...)
}

// New Tracerに対するポインタを作成
func New(w io.Writer) *Tracer {
	return &Tracer{out: w}
}
