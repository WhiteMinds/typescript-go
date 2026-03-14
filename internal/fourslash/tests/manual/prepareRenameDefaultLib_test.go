package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestPrepareRenameDefaultLib(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
/*setTimeout*/setTimeout(() => {}, 100);
/*console*/console.log("hello");
let arr: /*Array*/Array<number> = [];
let p: /*Promise*/Promise<void>;
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Standard library symbols should not be renameable via prepareRename
	stdlibMarkers := []string{"setTimeout", "console", "Array", "Promise"}
	for _, marker := range stdlibMarkers {
		f.GoToMarker(t, marker)
		f.VerifyPrepareRenameFailed(t, nil /*preferences*/)
	}
}

func TestPrepareRenameUserDefined(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @noLib: true
let /*myVar*/myVar = 42;
function /*myFunc*/myFunc() { return 1; }
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// User-defined symbols should be renameable
	userMarkers := []string{"myVar", "myFunc"}
	for _, marker := range userMarkers {
		f.GoToMarker(t, marker)
		f.VerifyPrepareRenameSucceeded(t, nil /*preferences*/)
	}
}
