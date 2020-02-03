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

	type ViewSliceField struct {
		Name string            `mapper:"Name"`
		Dt   time.Time         `mapper:"Dt"`
		Subs []SliceViewLevel1 `mapper:"Subs"`
	}

	type SimpleView struct {
		Name string `mapper:"Name"`
	}

	type FlatView struct {
		Name    string `mapper:"Name"`
		SubName string `mapper:"ModelLevel1.Name"`
	}

	tests := []struct {
		name string
		args args
	}{
		// TESTS
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
