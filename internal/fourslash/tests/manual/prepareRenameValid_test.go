package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestPrepareRenameValidTargets(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @noLib: true
let /*variable*/myVar = 42;
function /*function*/myFunc() { return 1; }
class /*class*/MyClass {}
interface /*interface*/MyInterface { x: number; }
type /*typeAlias*/MyType = string;
enum /*enumName*/MyEnum { A, B }
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// All user-defined symbols should be renameable
	markers := []string{"variable", "function", "class", "interface", "typeAlias", "enumName"}
	for _, marker := range markers {
		f.GoToMarker(t, marker)
		f.VerifyPrepareRenameSucceeded(t, nil /*preferences*/)
	}
}

func TestPrepareRenameLabel(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @noLib: true
/*label*/myLabel: for (let i = 0; i < 10; i++) {
    break myLabel;
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Labels should be renameable
	f.GoToMarker(t, "label")
	f.VerifyPrepareRenameSucceeded(t, nil /*preferences*/)
}

func TestPrepareRenameImportPath(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
export const x = 0;
// @Filename: /b.ts
import * as a from "/*importPath*/./a";
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	prefsFalse := &lsutil.UserPreferences{
		AllowRenameOfImportPath: core.TSFalse,
	}

	f.GoToMarker(t, "importPath")

	// Import path strings are not currently eligible for rename via prepareRename
	// (isNodeEligibleForRename only allows Identifier/PrivateIdentifier)
	f.Configure(t, prefsFalse)
	f.VerifyPrepareRenameFailed(t, prefsFalse)
}

func TestPrepareRenameDefaultExportKeyword(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
export default function foo() {}
// @Filename: /b.ts
import { /*defaultKeyword*/default as foo } from "./a";
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// "default" in import { default as foo } should not be renameable
	f.GoToMarker(t, "defaultKeyword")
	f.VerifyPrepareRenameFailed(t, nil /*preferences*/)
}
