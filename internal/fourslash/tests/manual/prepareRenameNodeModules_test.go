package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestPrepareRenameNodeModules(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /index.ts
import { Foo } from "foo";
declare const f: Foo;
f./*notOkProperty*/bar;
let /*okLocal*/myLocal = 1;
// @Filename: /tsconfig.json
 { }
// @Filename: /node_modules/foo/package.json
 { "types": "index.d.ts" }
// @Filename: /node_modules/foo/index.d.ts
export interface Foo {
    bar: string;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Accessing a property from a node_modules type should not be renameable
	f.GoToMarker(t, "notOkProperty")
	f.VerifyPrepareRenameFailed(t, nil /*preferences*/)

	// Local variables should be renameable
	f.GoToMarker(t, "okLocal")
	f.VerifyPrepareRenameSucceeded(t, nil /*preferences*/)
}
