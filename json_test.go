package json

import (
	"fmt"
	"reflect"
	"testing"
)

type customTypes struct {
	listBuilder      func() any
	listItemPusher   func(any, int, any) any
	objectBuilder    func() any
	objectItemPusher func(any, string, any) any
}

type customMap struct {
	myMap map[string]any
}

func (m *customMap) Push(key string, value any) {
	m.myMap[key] = value
}

type customSlice struct {
	mySlice []any
}

func (m *customSlice) Push(value any) {
	m.mySlice = append(m.mySlice, value)
}

func TestParseFromJson(t *testing.T) {
	data := []struct {
		json          string
		ct            customTypes
		expectedValue any
		expectedError string
	}{
		{"", customTypes{}, nil, syntaxError(0).Error()},
		{"{", customTypes{}, nil, syntaxError(1).Error()},
		{"[", customTypes{}, nil, syntaxError(1).Error()},
		{"}", customTypes{}, nil, syntaxError(0).Error()},
		{"}", customTypes{}, nil, syntaxError(0).Error()},
		{"asd", customTypes{}, nil, syntaxError(0).Error()},
		{`"string"`, customTypes{}, "string", ""},
		{"\r\n\t \"string\"\r\n\t ", customTypes{}, "string", ""},
		{`"\u041f\u0440\u0438\u0432\u0435\u0442"`, customTypes{}, "Привет", ""},
		{`"\u042d\u043c\u043e\u0434\u0437\u0438: \u2705, \ud83e\udd2a"`, customTypes{}, "Эмодзи: ✅, 🤪", ""},
		{"-100", customTypes{}, float64(-100), ""},
		{"-10", customTypes{}, float64(-10), ""},
		{"-1", customTypes{}, float64(-1), ""},
		{"-0.123456789", customTypes{}, -0.123456789, ""},
		{"0.123456789", customTypes{}, 0.123456789, ""},
		{"1", customTypes{}, float64(1), ""},
		{"10", customTypes{}, float64(10), ""},
		{"100", customTypes{}, float64(100), ""},
		{"\r\n\t 100\r\n\t", customTypes{}, float64(100), ""},
		{"\r\n\t -1234.3425525e+3\r\n\t ", customTypes{}, -1234342.5525, ""},
		{"\r\n\t -1234.3425525E+3\r\n\t ", customTypes{}, -1234342.5525, ""},
		{"null", customTypes{}, nil, ""},
		{"true", customTypes{}, true, ""},
		{"false", customTypes{}, false, ""},
		{"[]", customTypes{}, []any{}, ""},
		{"{}", customTypes{}, map[string]any{}, ""},
		{" \r\n\t[]\r\n\t ", customTypes{}, []any{}, ""},
		{" \r\n\t{}\r\n\t ", customTypes{}, map[string]any{}, ""},
		{"[ \r\n\t ]", customTypes{}, []any{}, ""},
		{"{ \r\n\t }", customTypes{}, map[string]any{}, ""},
		{"[1a]", customTypes{}, nil, syntaxError(2).Error()},
		{"[a1]", customTypes{}, nil, syntaxError(1).Error()},
		{"[\"]", customTypes{}, nil, syntaxError(3).Error()},
		{"[\"asd]", customTypes{}, nil, syntaxError(6).Error()},
		{"[truea]", customTypes{}, nil, syntaxError(5).Error()},
		{"[falsea]", customTypes{}, nil, syntaxError(6).Error()},
		{"[nulla]", customTypes{}, nil, syntaxError(5).Error()},
		{"[true,]", customTypes{}, nil, syntaxError(6).Error()},
		{"[1,]", customTypes{}, nil, syntaxError(3).Error()},
		{"[true, false, null]", customTypes{}, []any{true, false, nil}, ""},
		{"[1,2, 3, 4 , \r\n\t5 \r\n\t,6]", customTypes{}, []any{1., 2., 3., 4., 5., 6.}, ""},
		{
			"[-100, -10, -1, -0.1, -0.0123456789, 0, 0.0123456789, 0.1, 1, 10, 100]",
			customTypes{},
			[]any{-100., -10., -1., -0.1, -0.0123456789, 0., 0.0123456789, 0.1, 1., 10., 100.},
			"",
		},
		{
			`["", " ", " \r\n\t", "some text", "\"some name\"", "\u041f\u0440\u0438\u0432\u0435\u0442", "\u042d\u043c\u043e\u0434\u0437\u0438: \u2705, \ud83e\udd2a"]`,
			customTypes{},
			[]any{"", " ", " \r\n\t", "some text", "\"some name\"", "Привет", "Эмодзи: ✅, 🤪"},
			"",
		},
		{"[[],[]]", customTypes{}, []any{[]any{}, []any{}}, ""},
		{
			"[[1,2,3,4],[5,6,7,8]]",
			customTypes{},
			[]any{[]any{1., 2., 3., 4.}, []any{5., 6., 7., 8.}},
			"",
		},
		{
			`[-100, -10, -1, -0.1, -0.0123456789, 0, 0.0123456789, 0.1, 1, 10, 100, "", " ", " \r\n\t", "some text", "\"some name\"", "\u041f\u0440\u0438\u0432\u0435\u0442", "\u042d\u043c\u043e\u0434\u0437\u0438: \u2705, \ud83e\udd2a", [1,2,3], {}]`,
			customTypes{},
			[]any{-100., -10., -1., -0.1, -0.0123456789, 0., 0.0123456789, 0.1, 1., 10., 100., "", " ", " \r\n\t", "some text", "\"some name\"", "Привет", "Эмодзи: ✅, 🤪", []any{1., 2., 3.}, map[string]any{}},
			"",
		},
		{"[{},{}]", customTypes{}, []any{map[string]any{}, map[string]any{}}, ""},
		{`{"key"}`, customTypes{}, nil, syntaxError(6).Error()},
		{`{"key":}`, customTypes{}, nil, syntaxError(7).Error()},
		{`{"key":""}`, customTypes{}, map[string]any{"key": ""}, ""},
		{
			"{\r\n\t \"a\"\r\n\t :\r\n\t \"a\"\r\n\t , \"b\":\"b\", \"c\":\"c\", \"d\":\"d\"}",
			customTypes{},
			map[string]any{"a": "a", "b": "b", "c": "c", "d": "d"},
			"",
		},
		{
			`{"a":"value", "b":1, "c":-1, "d":10, "e":-10, "f": 0.45, "g": -0.45, "h" : false, "i":true, "j":null, "k": "\u041f\u0440\u0438\u0432\u0435\u0442", "l": "\u042d\u043c\u043e\u0434\u0437\u0438: \u2705, \ud83e\udd2a", "m": [], "n": {}}`,
			customTypes{},
			map[string]any{"a": "value", "b": 1., "c": -1., "d": 10., "e": -10., "f": 0.45, "g": -0.45, "h": false, "i": true, "j": nil, "k": "Привет", "l": "Эмодзи: ✅, 🤪", "m": []any{}, "n": map[string]any{}},
			"",
		},
		{
			`[]`,
			customTypes{
				listBuilder:    makeList,
				listItemPusher: pushListItem,
			},
			&customSlice{
				mySlice: []any{},
			},
			"",
		},
		{
			`[1, 2, 3, 4, 5]`,
			customTypes{
				listBuilder:    makeList,
				listItemPusher: pushListItem,
			},
			&customSlice{
				mySlice: []any{1., 2., 3., 4., 5.},
			},
			"",
		},
		{
			`["", " ", " \r\n\t", "some text", "\"some name\"", "\u041f\u0440\u0438\u0432\u0435\u0442", "\u042d\u043c\u043e\u0434\u0437\u0438: \u2705, \ud83e\udd2a"]`,
			customTypes{
				listBuilder:    makeList,
				listItemPusher: pushListItem,
			},
			&customSlice{
				mySlice: []any{"", " ", " \r\n\t", "some text", "\"some name\"", "Привет", "Эмодзи: ✅, 🤪"},
			},
			"",
		},
		{
			`[[1, 2, 3, 4, 5], [1, 2, 3]]`,
			customTypes{
				listBuilder:    makeList,
				listItemPusher: pushListItem,
			},
			&customSlice{
				mySlice: []any{
					&customSlice{
						mySlice: []any{1., 2., 3., 4., 5.},
					},
					&customSlice{
						mySlice: []any{1., 2., 3.},
					},
				},
			},
			"",
		},
		{
			`[[1, 2], [1, 2, 3], [{"key1":"value1","key2":"value2"}]]`,
			customTypes{
				listBuilder:    makeList,
				listItemPusher: pushListItem,
			},
			&customSlice{
				mySlice: []any{
					&customSlice{
						mySlice: []any{1., 2.},
					},
					&customSlice{
						mySlice: []any{1., 2., 3.},
					},
					&customSlice{
						mySlice: []any{map[string]any{"key1": "value1", "key2": "value2"}},
					},
				},
			},
			"",
		},
		{
			`{}`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
			},
			&customMap{
				myMap: map[string]any{},
			},
			"",
		},
		{
			`{"key1":"value1","key2":"value2"}`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
			},
			&customMap{
				myMap: map[string]any{"key1": "value1", "key2": "value2"},
			},
			"",
		},
		{
			`{"key1":{"key1.1":"value1.1","key1.2":"value1.2"},"key2":{"key2.1":"value2.1"}}`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
			},
			&customMap{
				myMap: map[string]any{
					"key1": &customMap{
						myMap: map[string]any{"key1.1": "value1.1", "key1.2": "value1.2"},
					},
					"key2": &customMap{
						myMap: map[string]any{"key2.1": "value2.1"},
					},
				},
			},
			"",
		},
		{
			`{"key1":{"key1.1":"value1.1","key1.2":[1, 2, 3]},"key2":{"key2.1":"value2.1"}}`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
			},
			&customMap{
				myMap: map[string]any{
					"key1": &customMap{
						myMap: map[string]any{"key1.1": "value1.1", "key1.2": []any{1., 2., 3.}},
					},
					"key2": &customMap{
						myMap: map[string]any{"key2.1": "value2.1"},
					},
				},
			},
			"",
		},
		{
			`{}`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
				listBuilder:      makeList,
				listItemPusher:   pushListItem,
			},
			&customMap{
				myMap: map[string]any{},
			},
			"",
		},
		{
			`[]`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
				listBuilder:      makeList,
				listItemPusher:   pushListItem,
			},
			&customSlice{
				mySlice: []any{},
			},
			"",
		},
		{
			`[[1, 2], [1, 2, 3], [{"key1":"value1","key2":"value2"}]]`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
				listBuilder:      makeList,
				listItemPusher:   pushListItem,
			},
			&customSlice{
				mySlice: []any{
					&customSlice{
						mySlice: []any{1., 2.},
					},
					&customSlice{
						mySlice: []any{1., 2., 3.},
					},
					&customSlice{
						mySlice: []any{
							&customMap{
								myMap: map[string]any{"key1": "value1", "key2": "value2"},
							},
						},
					},
				},
			},
			"",
		},
		{
			`{"key1":{"key1.1":"value1.1","key1.2":[1, 2, 3]},"key2":{"key2.1":"value2.1"}}`,
			customTypes{
				objectBuilder:    makeObject,
				objectItemPusher: pushObjectItem,
				listBuilder:      makeList,
				listItemPusher:   pushListItem,
			},
			&customMap{
				myMap: map[string]any{
					"key1": &customMap{
						myMap: map[string]any{
							"key1.1": "value1.1",
							"key1.2": &customSlice{
								mySlice: []any{1., 2., 3.},
							},
						},
					},
					"key2": &customMap{
						myMap: map[string]any{"key2.1": "value2.1"},
					},
				},
			},
			"",
		},
	}

	for _, item := range data {
		j := Bytes([]byte(item.json))
		if item.ct.objectBuilder != nil && item.ct.objectItemPusher != nil {
			j.WithObjectBuilder(item.ct.objectBuilder, item.ct.objectItemPusher)
		}
		if item.ct.listBuilder != nil && item.ct.listItemPusher != nil {
			j.WithListBuilder(item.ct.listBuilder, item.ct.listItemPusher)
		}

		actualValue, err := j.Decode()
		if err != nil {
			if err.Error() != item.expectedError {
				t.Errorf(
					"parse(%q) must return error %q, but %q received.",
					item.json,
					item.expectedError,
					err,
				)
			}
		} else {
			switch v := actualValue.(type) {
			case string, int, float64, bool, nil:
				if v != item.expectedValue {
					t.Errorf(notEquals(item.json, item.expectedValue, actualValue))
				}
			default:
				if reflect.DeepEqual(item.expectedValue, v) == false {
					t.Errorf(
						"parse(%q) must return %q but %q received.",
						item.json,
						reflect.TypeOf(item.expectedValue).String(),
						reflect.TypeOf(actualValue).String(),
					)
				}
			}
		}
	}
}

func notEquals(json string, expectedValue any, actualValue any) string {
	return fmt.Sprintf(
		"parse(%q) must return %q, but %q received.",
		json,
		expectedValue,
		actualValue,
	)
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

func pushListItem(l any, key int, value any) any {
	l.(*customSlice).Push(value)
	return l
}

func makeList() any {
	return &customSlice{
		mySlice: make([]any, 0, 0),
	}
}
