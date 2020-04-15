package mapper

import (
	"reflect"
	"testing"
	"time"
)

func TestFromTo(t *testing.T) {
	m := New()

	timeNow := time.Now()

	type args struct {
		e        interface{}
		v        interface{}
		expected interface{}
	}

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

	type ModelLevel2 struct {
		Name string
	}

	type ModelLevel1 struct {
		ModelLevel2 ModelLevel2
		Name        string
	}

	type ModelSliceField struct {
		Name string
		Dt   time.Time
		Subs []ModelLevel1
	}

	type Model struct {
		Name        string
		ModelLevel1 ModelLevel1
	}

	type MapModel struct {
		Name      string
		MapInt    map[int]int
		MapString map[string]string
		Map       map[string]ModelLevel1
	}

	type ViewLevel1 struct {
		NameLevel2 string `mapper:"ModelLevel2.Name"`
		NameLevel1 string `mapper:"Name"`
	}

	type SliceViewLevel1 struct {
		Name string `mapper:"Name"`
	}

	type View struct {
		Name string     `mapper:"Name"`
		Sub  ViewLevel1 `mapper:"ModelLevel1"`
	}

	type MapView struct {
		Name      string
		MapInt    map[int]int
		MapString map[string]string
		Map       map[string]ViewLevel1
	}

	type ViewSliceField struct {
		Name string            `mapper:"Name"`
		Dt   time.Time         `mapper:"Dt"`
		Subs []SliceViewLevel1 `mapper:"Subs"`
	}

	type SimpleView struct {
		Name string `mapper:"Name"`
	}

	type ViewDifferentFieldName struct {
		DifferentName string `mapper:"Name"`
	}

	type FlatView struct {
		Name    string `mapper:"Name"`
		SubName string `mapper:"ModelLevel1.Name"`
	}

	type SimpleMapModel struct {
		Name   string
		MapInt map[int]int
	}

	type SimpleTaggedMapView struct {
		TaggedName   string      `mapper:"Name"`
		TaggedMapInt map[int]int `mapper:"MapInt"`
	}

	type TaggedViewLevel2 struct {
		TaggedName string `mapper:"Name"`
	}

	type TaggedViewLevel1 struct {
		TaggedLevel2 TaggedViewLevel2 `mapper:"ModelLevel2"`
		TaggedName   string           `mapper:"Name"`
	}

	tests := []struct {
		name string
		args args
	}{
		// TESTS
		{
			name: "Mapping a map source",
			args: args{
				e: MapModel{
					Name: "teste",
					MapInt: map[int]int{
						1: 2,
						2: 3,
					},
					MapString: map[string]string{
						"key-string":   "value-string",
						"key-string-2": "value-string-2",
					},
					Map: map[string]ModelLevel1{
						"model-1": ModelLevel1{
							Name: "name level 1",
							ModelLevel2: ModelLevel2{
								Name: "name level 2",
							},
						},
					},
				},
				v: &MapView{},
				expected: &MapView{
					Name: "teste",
					MapInt: map[int]int{
						1: 2,
						2: 3,
					},
					MapString: map[string]string{
						"key-string":   "value-string",
						"key-string-2": "value-string-2",
					},
					Map: map[string]ViewLevel1{
						"model-1": ViewLevel1{
							NameLevel1: "name level 1",
							NameLevel2: "name level 2",
						},
					},
				},
			},
		},
		{
			name: "Mapping source to target",
			args: args{
				e: Source{
					Name: "teste",
					List: []string{"1", "2"},
					Sources: []Source{
						Source{
							Name: "teste nivel 1",
						},
					},
				},
				v: &Target{},
				expected: &Target{
					Name: "teste",
					List: []string{"1", "2"},
					Targets: []Target{
						Target{
							Name: "teste nivel 1",
						},
					},
				},
			},
		},
		{
			name: "Mapping equal structure",
			args: args{
				e: Model{
					Name: "teste",
					ModelLevel1: ModelLevel1{
						Name: "Teste",
						ModelLevel2: ModelLevel2{
							Name: "aaa pppp",
						},
					},
				},
				v: &Model{},
				expected: &Model{
					Name: "teste",
					ModelLevel1: ModelLevel1{
						Name: "Teste",
						ModelLevel2: ModelLevel2{
							Name: "aaa pppp",
						},
					},
				},
			},
		},
		{
			name: "View 3 levels",
			args: args{
				e: Model{
					Name: "teste",
					ModelLevel1: ModelLevel1{
						Name: "Teste",
						ModelLevel2: ModelLevel2{
							Name: "aaa pppp",
						},
					},
				},
				v: &View{},
				expected: &View{
					Name: "teste",
					Sub: ViewLevel1{
						NameLevel1: "Teste",
						NameLevel2: "aaa pppp",
					},
				},
			},
		},
		{
			name: "Field Slice View",
			args: args{
				e: ModelSliceField{
					Name: "teste",
					Dt:   timeNow,
					Subs: []ModelLevel1{
						ModelLevel1{
							Name: "SubNameTeste rrrrrr",
						},
					},
				},
				v: &ViewSliceField{},
				expected: &ViewSliceField{
					Name: "teste",
					Dt:   timeNow,
					Subs: []SliceViewLevel1{
						SliceViewLevel1{
							Name: "SubNameTeste rrrrrr",
						},
					},
				},
			},
		},
		{
			name: "Slice View",
			args: args{
				e: []Model{
					Model{
						Name: "teste",
						ModelLevel1: ModelLevel1{
							Name: "SubNameTeste rrrrrr",
						},
					},
				},
				v: &[]View{},
				expected: &[]View{
					View{
						Name: "teste",
						Sub: ViewLevel1{
							NameLevel1: "SubNameTeste rrrrrr",
						},
					},
				},
			},
		},
		{
			name: "Flat View",
			args: args{
				e: Model{
					Name: "teste",
					ModelLevel1: ModelLevel1{
						Name: "SubNameTeste rrrrrr",
					},
				},
				v: &FlatView{},
				expected: &FlatView{
					Name:    "teste",
					SubName: "SubNameTeste rrrrrr",
				},
			},
		},
		{
			name: "Sub View",
			args: args{
				e: Model{
					Name: "teste",
					ModelLevel1: ModelLevel1{
						Name: "SubNameTeste rrrrrr",
					},
				},
				v: &View{},
				expected: &View{
					Name: "teste",
					Sub: ViewLevel1{
						NameLevel1: "SubNameTeste rrrrrr",
					},
				},
			},
		},
		{
			name: "Simple View",
			args: args{
				e: Model{
					Name: "teste",
				},
				v: &SimpleView{},
				expected: &SimpleView{
					Name: "teste",
				},
			},
		},
		{
			name: "From struct has tag",
			args: args{
				e: ViewDifferentFieldName{
					DifferentName: "teste",
				},
				v: &Model{},
				expected: &Model{
					Name: "teste",
				},
			},
		},
		{
			name: "Mapping slices with tags on from",
			args: args{
				e: []ViewDifferentFieldName{
					{DifferentName: "teste"},
				},
				v: &[]Model{},
				expected: &[]Model{
					{Name: "teste"},
				},
			},
		},
		{
			name: "Mapping a map source with tags",
			args: args{
				e: SimpleTaggedMapView{
					TaggedName: "teste",
					TaggedMapInt: map[int]int{
						1: 2,
						2: 3,
					},
				},
				v: &SimpleMapModel{},
				expected: &SimpleMapModel{
					Name: "teste",
					MapInt: map[int]int{
						1: 2,
						2: 3,
					},
				},
			},
		},
		{
			name: "Mapping a nested view with tags",
			args: args{
				e: TaggedViewLevel1{
					TaggedName: "teste_lvl_1",
					TaggedLevel2: TaggedViewLevel2{
						TaggedName: "teste_lvl_2",
					},
				},
				v: &ModelLevel1{},
				expected: &ModelLevel1{
					Name: "teste_lvl_1",
					ModelLevel2: ModelLevel2{
						Name: "teste_lvl_2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.From(tt.args.e).To(tt.args.v)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(tt.args.v, tt.args.expected) {
				t.Error("Error mapping: result is not expected")
			}
		})
	}
}
