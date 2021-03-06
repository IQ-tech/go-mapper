# mapper packages

This package holds helpers for better parses among structs.

#### `New() Mapper`
This method reutrns a new instance of Mapper.

## Mapper types

### `Mapper` interface
This interface contains one method only:

```go
// From creates map from source
type Mapper interface {
    From(src interface{}) Result
}
```

### `Result` interface
This interface contains two methods only:

```go
// From creates map from source
type Result interface {
    Merge(src interface{}) Result
    To(tgr interface{}) error
}
```


#### `From(src interface{}) Result`

This method returns the map `Result` interface.

**Example:**

```go
// mapping some entity
result := mapper.From(entity)

```

#### `Merge(src interface{}) Result`
This method returns the merged map `Result` interface.
```go
// merging some entity
result := result.Merge(entity)

```

#### `To(tgr interface{}) error`
This method receive a pointer as parameter and returns error.
```go
// mapping to some target
err := result.To(&target)
if err != nil {
    fmt.Println(err)
}

```

> **Note:** You can use this methods combined.

**Example**
```go
err := mapper.From(source).To(&target)
if err != nil {
    fmt.Println(err)
}

```

#### TAG `mapper`
The tag `mapper` allow you mapping fields with diferents names.

***Example**
Mapping field Sources from Source to Targes from Target 

```go
type Source struct {
		Name    string
		List    []string
		Sources []Source
	}

	type Target struct {
		Name    string
		List    []string
		Targets []Target `mapper:"Sources"`
	}
```