# Atlas

Atlas is a collection of useful map data structure to make things easier when working with maps in Go.

## Installation

`go get github.com/CaioTeixeira95/atlas`

## Maps

Below you find a list of data structures of this module.

### BiMap

BiMap is a bidirectional map that gives you the ability to search by key or value.

```go
biMap := atlas.NewBiMap[string, string]()
err := biMap.Set("myKey", "myValue")
if err != nil {
    log.Fatal(err)
}

value, ok := biMap.Get("myKey")
fmt.Println(value, ok) // myValue true

key, ok := biMap.GetInverse("myValue")
fmt.Println(key, ok) // myKey true
```

### FrozenMap

FrozenMap gives you a map where the keys are immutable. It returns error when you try to overwrite a key's value.

```go
frozenMap := atlas.NewFrozenMap[string, int64]()
err := frozenMap.Set("key", 10)
if err != nil {
    log.Fatal(err)
}
err = frozenMap.Set("key", 24)
fmt.Println(errors.Is(err, atlas.ErrKeyAlreadySet)) // true
```
