# Mongo testcontainer
A mongo docker testcontainer for go


## Example usage:

```go
initPath, err := filepath.Abs("init.sh")
if err != nil {
    log.Fatalf("%v", err)
}
ctx := context.Background()
mongoC, conn, err := mongotestcontainer.New(ctx, "latest", initPath)
if err != nil {
    log.Fatalf("%v", err)
}
defer Terminate(ctx, mongoC)

fmt.Println(conn)
// Output: mongo://mongo:mongo@localhost:49156/test_db
```

## Run tests

``sh
go test
``
