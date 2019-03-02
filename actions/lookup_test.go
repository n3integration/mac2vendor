package actions

import (
	"testing"
)

func TestLookup(t *testing.T) {
	mac = "84:38:35:77:aa:52"

	t.Run("Default", func(t *testing.T) {
		if err := lookupAction(nil); err != nil {
			t.Error("failed to lookup mac: ", err)
		}
	})

	t.Run("Quiet", func(t *testing.T) {
		quiet = true
		if err := lookupAction(nil); err != nil {
			t.Error("failed to lookup mac: ", err)
		}
	})
}
