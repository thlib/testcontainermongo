# Mongo testcontainer
A mongo docker testcontainer for go


## Example usage:

```go
initPath, err := filepath.Abs("./initdb")
if err != nil {
    log.Fatalf("%v", err)
}
ctx := context.Background()
container, conn, err := mongotestcontainer.New(ctx, "latest", initPath)
if err != nil {
    log.Fatalf("%v", err)
}
defer Terminate(ctx, container)

fmt.Println(conn)
// Output: mongodb://root:example@localhost:49156/test_db
```

## Run tests

```sh
go test
```
