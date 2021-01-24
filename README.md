# gokv
KV store written in GO

## Examples

#### List all Keys:

```
curl 127.0.0.1:9000/kv
```

#### Add key and value

```
curl -XPOST 127.0.0.1:9000/kv/Tom -H "Content-Type: application/json" --data '{"age": "34","birthday":"10/23/1986","FavColor": "tan"}'
```

#### Get value at key

```
curl 127.0.0.1:9000/kv/Tom 
```

#### Delete Key

```
curl -XDELETE 127.0.0.1:9000/kv/Tom
```