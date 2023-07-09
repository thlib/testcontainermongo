# Mongo testcontainer
A mongo docker testcontainer for go


## Example usage:

```go
import (
	"context"
	"path/filepath"
	"regexp"
	"testing"

	testmongo "github.com/thlib/testcontainermongo"
)
initPath, err := filepath.Abs("./initdb")
if err != nil {
    log.Fatalf("%v", err)
}
ctx := context.Background()
c, conn, err := testmongo.New(ctx, "latest",
    testmongo.WithInit(initdb),
    testmongo.WithDb("test_db"),
    testmongo.WithAuth("root", "example"),
)
defer testmongo.Terminate(ctx, container)

fmt.Println(conn)
// Output: mongodb://root:example@localhost:49156/test_db
```

## Run tests

```sh
go test
```
