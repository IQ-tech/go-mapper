package mapper

import (
	"fmt"
	"reflect"
	"strconv"
)

type mapper struct {
}

// From create map form source
func (m *mapper) From(src interface{}) Result {

	retVal := result{}

	// Create Map
	source := reflect.ValueOf(src)

	if source.Kind() == reflect.Map {

		retVal.fieldMap = m.buildMapFromMap("", src)

	} else if source.Kind() == reflect.Slice {

		retVal.fieldMap = m.buildMapFromSlice("", src)

	} else {
		retVal.fieldMap = m.buildMapFromStruct("", src)
	}

	return retVal
}

func (m *mapper) buildMapFromMap(rootFieldName string, src interface{}) (mapField map[string]interface{}) {

	source := reflect.ValueOf(src)

	mapField = make(map[string]interface{})
	srcMap := make(map[interface{}]interface{})

	iter := source.MapRange()

	for iter.Next() {

		fieldName := rootFieldName

		if fieldName == "" {
			fieldName = fmt.Sprintf("[%v]", iter.Key())
		} else {
			fieldName = fmt.Sprintf("%s.[%v]", fieldName, iter.Key())
		}

		if iter.Value().Kind() == reflect.Struct && iter.Value().Type().String() != timeType {

			mapItem := m.buildMapFromStruct(fieldName, iter.Value().Interface())
			for k, v := range mapItem {
				mapField[k] = v
			}

			mapField[fieldName] = mapItem
			key := iter.Key().Interface()
			srcMap[key] = mapItem

		} else {

			itemMap := reflect.ValueOf(iter.Value())
			mapField[fieldName] = itemMap
			key := iter.Key().Interface()
			srcMap[key] = itemMap
		}
	}

	mapField[rootFieldName] = srcMap

	return mapField
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

		tagValue, foundTag := fieldMetada.Tag.Lookup(MapperTag)

		if tagValue == "-" && foundTag {
			continue
		}

		if rootFieldName != "" && foundTag {
			tagValue = rootFieldName + "." + tagValue
		} else if rootFieldName != "" && !foundTag {
			tagValue = rootFieldName + "." + fieldMetada.Name
		} else if rootFieldName == "" && !foundTag {
			tagValue = fieldMetada.Name
		}

		if field.Kind() == reflect.Map {

			mapItem := m.buildMapFromMap(tagValue, field.Interface())

			for k, v := range mapItem {
				mapField[k] = v
			}

		} else if field.Kind() == reflect.Slice {

			mapItem := m.buildMapFromSlice(tagValue, field.Interface())

			for k, v := range mapItem {
				mapField[k] = v
			}

		} else if field.Kind() == reflect.Struct && field.Type().String() != timeType {

			mapStruct := m.buildMapFromStruct(tagValue, field.Interface())

			for k, v := range mapStruct {
				mapField[k] = v
			}

		} else {
			mapField[tagValue] = field
		}
	}

	return mapField
}
