package mapper

import (
	"reflect"
	"regexp"
	"strconv"

	"github.com/IQ-tech/go-errors"
)

// MapperTag standards tag name for view mapping
const MapperTag string = "mapper"
const timeType string = "time.Time"

var regPoint = regexp.MustCompile(`\.`)

// Result holds result of mapper
type Result interface {
	Merge(src interface{}) Result
	To(tgr interface{}) error
}

// Mapper holds mapping operations
type Mapper interface {
	From(src interface{}) (retVal Result)
}

type mapper struct {
}

type result struct {
	fieldMap map[string]interface{}
}

// New returns an instance of mapper
func New() Mapper {
	return &mapper{}
}

// From create map form source
func (m *mapper) From(src interface{}) Result {

	retVal := result{}

	// Create Map
	source := reflect.ValueOf(src)

	if source.Kind() == reflect.Slice {

		retVal.fieldMap = m.buildMapFromSlice("", src)

	} else {
		retVal.fieldMap = m.buildMapFromStruct("", src)
	}

	return retVal
}

// Merge mixes maps from source
func (r result) Merge(source interface{}) (retVal Result) {
	// TODO
	return retVal
}

// To sets target form map value
func (r result) To(trg interface{}) (err error) {

	target := reflect.ValueOf(trg)

	if target.Kind() != reflect.Ptr {
		return errors.NewApplicationError("Error mapper: must receive a pointer, but received " + target.Kind().String())
	}

	if target.Elem().Kind() == reflect.Slice {
		err = r.setSlice("", target.Elem())
		if err != nil {
			return errors.Wrap(err)
		}
	} else if target.Elem().Kind() == reflect.Struct {
		err = r.setStruct("", target)
		if err != nil {
			return errors.Wrap(err)
		}
	} else {
		err = r.setField("", target)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}

func (r result) setSlice(rootFieldName string, target reflect.Value) (err error) {

	src, ok := r.fieldMap[rootFieldName]
	if !ok {
		return errors.NewApplicationError("Error mapper: Field is not found")
	}

	source := reflect.ValueOf(src)

	if source.Len() == 0 {
		return nil
	}

	target.Set(reflect.MakeSlice(target.Type(), source.Len(), source.Cap()))

	for i := 0; i < source.Len(); i++ {
		sourceItem := source.Index(i)
		targetItem := target.Index(i)
		var value reflect.Value

		// Corrigir para os casos de lista de structs
		if targetItem.Kind() == reflect.Struct && targetItem.Type().String() != timeType {

			fieldName := rootFieldName

			if fieldName == "" {
				fieldName = strconv.Itoa(i)
			} else {
				fieldName = fieldName + "." + strconv.Itoa(i)
			}

			err = r.setStruct(fieldName, targetItem.Addr())
			if err != nil {
				return errors.Wrap(err)
			}

		} else {
			item, ok := sourceItem.Interface().(reflect.Value)
			if !ok {
				return errors.NewApplicationError("Error mapper: Value from list is not an instance of reflect.Value")
			}

			value, ok = item.Interface().(reflect.Value)
			if !ok {
				return errors.NewApplicationError("Error mapper: Value from list is not an instance of reflect.Value")
			}

			targetItem.Set(value)
		}

	}

	return nil
}

func (r result) setStruct(rootFieldName string, trg reflect.Value) (err error) {

	target := trg.Elem()

	for i := 0; i < target.NumField(); i++ {

		fieldMetada := target.Type().Field(i)

		tagValue, foundTag := fieldMetada.Tag.Lookup(MapperTag)

		field := target.Field(i)
		fieldType := field.Kind()

		if rootFieldName != "" && foundTag {
			tagValue = rootFieldName + "." + tagValue
		} else if rootFieldName != "" {
			tagValue = rootFieldName + "." + fieldMetada.Name
		}

		if tagValue == "-" && foundTag {
			continue
		}

		if tagValue == "" && foundTag {
			return errors.NewApplicationError("Error mapper: tag view must have a value when it's present")
		}

		if tagValue == "" {
			tagValue = fieldMetada.Name
		}

		if fieldType == reflect.Slice {

			err = r.setSlice(tagValue, field)
			if err != nil {
				return errors.Wrap(err)
			}

		} else if fieldType == reflect.Struct && field.Type().String() != timeType {

			err = r.setStruct(tagValue, field.Addr())
			if err != nil {
				return errors.Wrap(err)
			}

		} else {

			err = r.setField(tagValue, field)
			if err != nil {
				return errors.Wrap(err)
			}
		}
	}

	return nil
}

func (r result) setField(fieldName string, field reflect.Value) error {

	item, ok := r.fieldMap[fieldName]
	if !ok {
		return errors.NewApplicationError("Error mapper: Field is not found")
	}

	value, ok := item.(reflect.Value)
	if !ok {
		return errors.NewApplicationError("Error mapper: Value from map is not a reflect.Value")
	}

	if !field.IsValid() && !field.CanSet() {
		return errors.NewApplicationError("Error mapper: Field is not valid")
	}

	if field.Kind() != value.Kind() {
		return errors.NewApplicationError("Mapper error: different type field mapping")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(value.Int())
	case reflect.Float32, reflect.Float64:
		field.SetFloat(value.Float())
	default:
		field.Set(value)
	}

	return nil
}

func (m *mapper) buildMapFromSlice(rootFieldName string, src interface{}) (mapField map[string]interface{}) {

	source := reflect.ValueOf(src)
	mapField = make(map[string]interface{})
	list := []interface{}{}

	for i := 0; i < source.Len(); i++ {

		sourceItem := source.Index(i)

		fieldName := rootFieldName

		if fieldName == "" {
			fieldName = strconv.Itoa(i)
		} else {
			fieldName = fieldName + "." + strconv.Itoa(i)
		}

		if sourceItem.Kind() == reflect.Struct && sourceItem.Type().String() != timeType {

			mapItem := m.buildMapFromStruct(fieldName, sourceItem.Interface())
			for k, v := range mapItem {
				mapField[k] = v
			}

			list = append(list, mapItem)

		} else {
			itemSlice := reflect.ValueOf(sourceItem)
			mapField[fieldName] = itemSlice
			list = append(list, itemSlice)
		}
	}

	mapField[rootFieldName] = list

	return mapField
}

func (m *mapper) buildMapFromStruct(rootFieldName string, item interface{}) (mapField map[string]interface{}) {

	itemStruct := reflect.ValueOf(item)
	mapField = make(map[string]interface{})

	for i := 0; i < itemStruct.NumField(); i++ {

		field := itemStruct.Field(i)

		fieldMetada := itemStruct.Type().Field(i)

		fieldName := rootFieldName

		if fieldName == "" {
			fieldName = fieldMetada.Name
		} else {
			fieldName = fieldName + "." + fieldMetada.Name
		}

		if field.Kind() == reflect.Slice {

			mapItem := m.buildMapFromSlice(fieldName, field.Interface())

			for k, v := range mapItem {
				mapField[k] = v
			}

		} else if field.Kind() == reflect.Struct && field.Type().String() != timeType {

			mapStruct := m.buildMapFromStruct(fieldName, field.Interface())

			for k, v := range mapStruct {
				mapField[k] = v
			}

		} else {
			mapField[fieldName] = field
		}
	}

	return mapField
}
