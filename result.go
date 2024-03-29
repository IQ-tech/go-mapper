package mapper

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/IQ-tech/go-errors"
)

type result struct {
	fieldMap map[string]interface{}
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

	if target.Elem().Kind() == reflect.Map {
		err = r.setMap("", target.Elem())
		if err != nil {
			return errors.Wrap(err)
		}
	} else if target.Elem().Kind() == reflect.Slice {
		err = r.setSlice("", target.Elem())
		if err != nil {
			return errors.Wrap(err)
		}
	} else if target.Elem().Kind() == reflect.Struct {
		err = r.setStruct("", target)
		if err != nil {
			return errors.Wrap(err)
		}
	} else if target.Elem().Kind() == reflect.Interface {
		err = r.To(target.Elem().Interface())
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

func (r result) setMap(rootFieldName string, target reflect.Value) (err error) {

	src, ok := r.fieldMap[rootFieldName]
	if !ok {
		return errors.NewApplicationError("Error mapper: Field is not found")
	}

	source := reflect.ValueOf(src)

	target.Set(reflect.MakeMapWithSize(target.Type(), 0))

	iter := source.MapRange()

	for iter.Next() {

		elementType := reflect.TypeOf(target.Interface()).Elem()

		if elementType.Kind() == reflect.Struct && elementType.String() != timeType {

			targetItem := reflect.New(elementType)
			fieldName := rootFieldName

			if fieldName == "" {
				fieldName = fmt.Sprintf("[%v]", iter.Key())
			} else {
				fieldName = fmt.Sprintf("%s.[%v]", fieldName, iter.Key())
			}

			err = r.setStruct(fieldName, targetItem)
			if err != nil {
				return errors.Wrap(err)
			}

			key := reflect.ValueOf(iter.Key().Interface())
			target.SetMapIndex(key, targetItem.Elem())

		} else {

			key := reflect.ValueOf(iter.Key().Interface())
			rv := reflect.ValueOf(iter.Value().Interface()).Interface()
			item, ok := rv.(reflect.Value)
			if !ok {
				return errors.NewApplicationError("Error mapper: Value from map is not an instance of reflect.Value")
			}

			value, ok := item.Interface().(reflect.Value)
			if !ok {
				return errors.NewApplicationError("Error mapper: Value from map is not an instance of reflect.Value")
			}

			target.SetMapIndex(key, value)
		}

	}

	return nil
}

func (r result) setSlice(rootFieldName string, target reflect.Value) (err error) {

	src, ok := r.fieldMap[rootFieldName]
	if !ok {
		return nil
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

		if tagValue == "-" && foundTag {
			continue
		}

		if rootFieldName != "" && foundTag {
			tagValue = rootFieldName + "." + tagValue
		} else if rootFieldName != "" {
			tagValue = rootFieldName + "." + fieldMetada.Name
		}

		if tagValue == "" && foundTag {
			return errors.NewApplicationError("Error mapper: tag view must have a value when it's present")
		}

		if tagValue == "" {
			tagValue = fieldMetada.Name
		}

		if fieldType == reflect.Map {

			err = r.setMap(tagValue, field)
			if err != nil {
				return errors.Wrap(err)
			}

		} else if fieldType == reflect.Slice {

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
		return nil
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
