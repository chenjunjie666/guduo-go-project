package test

import (
	"fmt"
	"github.com/dop251/goja"
	"testing"
)

func TestJsEngine(t *testing.T) {
	script := `
function demo(t, e) {
	if (t) {
		for (var n = t.split(""), i = e.split(""), a = {}, r = [], o = 0; o < n.length / 2; o++)
			a[n[o]] = n[n.length / 2 + o];
		for (var s = 0; s < e.length; s++)
			r.push(a[i[s]]);
		return r.join("")
	}
}
	`
	vm := goja.New()
	_, _ = vm.RunString(script)

	sign, _ := goja.AssertFunction(vm.Get("demo"))


	fmt.Println(sign(goja.Undefined(), vm.ToValue(`Bu2-xmNPndHtpTz.,5+48-%0621973`), vm.ToValue(`xd22uxdzd`)))
}
