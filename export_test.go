package astquery

import "testing"

func DebugON(t *testing.T) {
	org := isDebug
	isDebug = true
	t.Cleanup(func() {
		isDebug = org
	})
}
