# json 

## json decode examples

### decode json object to golang map
```go
obj, err := json.Bytes([]byte(`{"key1": "value1", "key2": "value2"}`)).Decode()

# obj:
map[string]any{
"key1": "value1",
"key2": "value2",
}
```


### decode json list to golang slice
```go
list, err := json.Bytes([]byte(`["item1", "item2", "item3"]`)).Decode()

# list:
[]any{
    "item1",
    "item2",
    "item3"
}
```

### decode json object to custom structure
```go
package main

import (
	"github.com/URALINNOVATSIYA/json"
)

type customMap struct {
	myMap map[string]any
}

func (m *customMap) Push(key string, value any) {
	m.myMap[key] = value
}

func makeObject() any {
	return &customMap{
		myMap: make(map[string]any),
	}
}

func pushObjectItem(o any, key string, value any) any {
	o.(*customMap).Push(key, value)
	return o
}

func main() {
	cMap, err := json.Bytes([]byte(`{"key1": "value1", "key2": "value2"}`)).WithObjectBuilder(makeObject, pushObjectItem).Decode()

	# cMap:
	&customMap{
		myMap: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}
}
```


### decode json list to custom structure
```go
package main

import (
	"github.com/URALINNOVATSIYA/json"
)

type customSlice struct {
	mySlice []any
}

func (m *customSlice) Push(value any) {
	m.mySlice = append(m.mySlice, value)
}

func pushListItem(l any, key int, value any) any {
	l.(*customSlice).Push(value)
	return l
}

func makeList() any {
	return &customSlice{
		mySlice: make([]int, 0, 0),
	}
}


func main() {
	cSlice, err := json.Bytes([]byte(`[1, 2, 3]`)).WithListBuilder(makeList, pushListItem).Decode()

	# cSlice:
	&customSlice{
		mySlice: []int{1, 2, 3},
	}
}

```