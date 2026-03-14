package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestPrepareRenameKeyword(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @noLib: true
function f(value: string, /*1*/default: string) {}

const /*2*/default = 1;

function /*3*/default() {}

class /*4*/default {}

const foo = {
    /*5*/[|default|]: 1
}

/*6*/if (true) {}

/*7*/for (let i = 0; i < 10; i++) {}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Keywords used as contextual identifiers in certain positions are not eligible
	keywordMarkers := []string{"1", "2", "3", "4", "6", "7"}
	for _, marker := range keywordMarkers {
		f.GoToMarker(t, marker)
		f.VerifyPrepareRenameFailed(t, nil /*preferences*/)
	}

	// Computed property name using "default" is eligible
	f.GoToMarker(t, "5")
	f.VerifyPrepareRenameSucceeded(t, nil /*preferences*/)
}
