package errors

import (
	"github.com/prometheus-community/promql-langserver/vendored/go-tools/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b //@complete(" //", Bob_interface)
}
