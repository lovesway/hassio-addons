package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/lovesway/hassio-addons/mq-lightshow/models"
)

type (
	// StringsToStruct represents the entity that converts string interfaces into structs.
	StringsToStruct struct{}
)

const (
	thirtyTwo = 32
	sixtyFour = 64
)

// NewStringsToStruct provides an instance of APIController.
func NewStringsToStruct() StringsToStruct {
	return StringsToStruct{}
}

// Show will convert for a Show model.
func (ss StringsToStruct) Show(in interface{}, out models.Show) (models.Show, error) {
	fieldsIn := reflect.TypeOf(in)
	valuesIn := reflect.ValueOf(in)

	fields := reflect.TypeOf(out)
	values := reflect.ValueOf(out)
	num := fields.NumField()

	var err error

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		_, found := fieldsIn.FieldByName(field.Name)
		if !found {
			continue
		}

		fieldVal := valuesIn.FieldByName(field.Name).String()
		if fieldVal == "" {
			continue
		}

		var (
			fieldValString  string
			fieldValBool    bool
			fieldValInt     int
			fieldValInt32   int32
			fieldValInt64   int64
			fieldValFloat32 float32
			fieldValFloat64 float64
		)

		switch value.Kind() {
		case reflect.String:
			if !value.IsValid() {
				return out, fmt.Errorf("no such field: %s in obj", field.Name)
			}

			fieldValString = fieldVal

		case reflect.Int:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}
		case reflect.Int32:
			fieldValInt, err := strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt32 = int32(fieldValInt)
			log.Debugf("Int32: %v", fieldValInt32)
		case reflect.Int64:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt64 = int64(fieldValInt)
			log.Debugf("Int64: %v", fieldValInt64)
		case reflect.Float32:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, thirtyTwo)
			if err != nil {
				return out, err
			}

			fieldValFloat32 = float32(fieldValFloat)
		case reflect.Float64:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, sixtyFour)
			if err != nil {
				return out, err
			}

			fieldValFloat64 = float64(fieldValFloat)
			log.Debugf("Float64: %v", fieldValFloat64)
		case reflect.Bool:
			fieldValBool, err = strconv.ParseBool(fieldVal)
			if err != nil {
				return out, err
			}
		default:
			log.Errorf("Unsupported type of '%s' in %v", field.Type, field.Name)
		}

		if field.Name == "Name" {
			out.Name = fieldValString
		} else if field.Name == "Topic" {
			out.Topic = fieldValString
		} else if field.Name == "Repeat" {
			out.Repeat = fieldValBool
		} else if field.Name == globalDelay {
			out.GlobalDelay = fieldValFloat32
		} else if field.Name == globalSpeed {
			out.GlobalSpeed = fieldValInt
		} else if field.Name == globalParameter1 {
			out.GlobalParameter1 = fieldValString
		} else if field.Name == out.GlobalParameter2 {
			out.GlobalParameter2 = fieldValString
		}
	}

	return out, err
}

// Cycle will convert for a Cycle model.
func (ss StringsToStruct) Cycle(in interface{}, out models.Cycle) (models.Cycle, error) {
	fieldsIn := reflect.TypeOf(in)
	valuesIn := reflect.ValueOf(in)

	fields := reflect.TypeOf(out)
	values := reflect.ValueOf(out)
	num := fields.NumField()

	var err error

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		_, found := fieldsIn.FieldByName(field.Name)
		if !found {
			continue
		}

		fieldVal := valuesIn.FieldByName(field.Name).String()
		if fieldVal == "" {
			continue
		}

		var (
			fieldValString  string
			fieldValBool    bool
			fieldValInt     int
			fieldValInt32   int32
			fieldValInt64   int64
			fieldValFloat32 float32
			fieldValFloat64 float64
		)

		switch value.Kind() {
		case reflect.String:
			if !value.IsValid() {
				return out, fmt.Errorf("no such field: %s in obj", field.Name)
			}

			fieldValString = fieldVal

		case reflect.Int:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}
		case reflect.Int32:
			fieldValInt, err := strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt32 = int32(fieldValInt)
			log.Debugf("Int32: %v", fieldValInt32)
		case reflect.Int64:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt64 = int64(fieldValInt)
			log.Debugf("Int64: %v", fieldValInt64)
		case reflect.Float32:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, thirtyTwo)
			if err != nil {
				return out, err
			}

			fieldValFloat32 = float32(fieldValFloat)
		case reflect.Float64:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, sixtyFour)
			if err != nil {
				return out, err
			}

			fieldValFloat64 = float64(fieldValFloat)
			log.Debugf("Float64: %v", fieldValFloat64)
		case reflect.Bool:
			fieldValBool, err = strconv.ParseBool(fieldVal)
			if err != nil {
				return out, err
			}
		default:
			log.Errorf("Unsupported type of '%s' in %v", field.Type, field.Name)
		}

		if field.Name == "ShowID" {
			out.ShowID = fieldValInt
		} else if field.Name == "SceneID" {
			out.SceneID = fieldValInt
		} else if field.Name == "SceneCycles" {
			out.SceneCycles = fieldValInt
		} else if field.Name == "EndDelay" {
			out.EndDelay = fieldValFloat32
		} else if field.Name == "LoopInclude" {
			out.LoopInclude = fieldValBool
		} else if field.Name == globalDelay {
			out.GlobalDelay = fieldValFloat32
		} else if field.Name == globalSpeed {
			out.GlobalSpeed = fieldValInt
		} else if field.Name == globalParameter1 {
			out.GlobalParameter1 = fieldValString
		} else if field.Name == globalParameter2 {
			out.GlobalParameter2 = fieldValString
		}
	}

	return out, err
}

// Scene will convert for a Scene model.
func (ss StringsToStruct) Scene(in interface{}, out models.Scene) (models.Scene, error) {
	fieldsIn := reflect.TypeOf(in)
	valuesIn := reflect.ValueOf(in)

	fields := reflect.TypeOf(out)
	values := reflect.ValueOf(out)
	num := fields.NumField()

	var err error

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		_, found := fieldsIn.FieldByName(field.Name)
		if !found {
			continue
		}

		fieldVal := valuesIn.FieldByName(field.Name).String()
		if fieldVal == "" {
			continue
		}

		var (
			fieldValString  string
			fieldValBool    bool
			fieldValInt     int
			fieldValInt32   int32
			fieldValInt64   int64
			fieldValFloat32 float32
			fieldValFloat64 float64
		)

		switch value.Kind() {
		case reflect.String:
			if !value.IsValid() {
				return out, fmt.Errorf("no such field: %s in obj", field.Name)
			}

			fieldValString = fieldVal

		case reflect.Int:
			_, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}
		case reflect.Int32:
			fieldValInt, err := strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt32 = int32(fieldValInt)
			log.Debugf("Int32: %v", fieldValInt32)
		case reflect.Int64:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt64 = int64(fieldValInt)
			log.Debugf("Int64: %v", fieldValInt64)
		case reflect.Float32:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, thirtyTwo)
			if err != nil {
				return out, err
			}

			fieldValFloat32 = float32(fieldValFloat)
			log.Debugf("Float32: %v", fieldValFloat32)
		case reflect.Float64:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, sixtyFour)
			if err != nil {
				return out, err
			}

			fieldValFloat64 = float64(fieldValFloat)
			log.Debugf("Float64: %v", fieldValFloat64)
		case reflect.Bool:
			fieldValBool, err = strconv.ParseBool(fieldVal)
			if err != nil {
				return out, err
			}

			log.Debugf("Bool: %v", fieldValBool)
		default:
			log.Errorf("Unsupported type of '%s' in %v", field.Type, field.Name)
		}

		if field.Name == "Name" {
			out.Name = fieldValString
		}
	}

	return out, err
}

// Group will convert for a Group model.
func (ss StringsToStruct) Group(in interface{}, out models.Group) (models.Group, error) {
	fieldsIn := reflect.TypeOf(in)
	valuesIn := reflect.ValueOf(in)

	fields := reflect.TypeOf(out)
	values := reflect.ValueOf(out)
	num := fields.NumField()

	var err error

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		_, found := fieldsIn.FieldByName(field.Name)
		if !found {
			continue
		}

		fieldVal := valuesIn.FieldByName(field.Name).String()
		if fieldVal == "" {
			continue
		}

		var (
			fieldValString  string
			fieldValBool    bool
			fieldValInt     int
			fieldValInt32   int32
			fieldValInt64   int64
			fieldValFloat32 float32
			fieldValFloat64 float64
		)

		switch value.Kind() {
		case reflect.String:
			if !value.IsValid() {
				return out, fmt.Errorf("no such field: %s in obj", field.Name)
			}

			fieldValString = fieldVal
			log.Debugf("String: %v", fieldValString)
		case reflect.Int:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}
		case reflect.Int32:
			fieldValInt, err := strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt32 = int32(fieldValInt)
			log.Debugf("Int32: %v", fieldValInt32)
		case reflect.Int64:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt64 = int64(fieldValInt)
			log.Debugf("Int64: %v", fieldValInt64)
		case reflect.Float32:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, thirtyTwo)
			if err != nil {
				return out, err
			}

			fieldValFloat32 = float32(fieldValFloat)
		case reflect.Float64:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, sixtyFour)
			if err != nil {
				return out, err
			}

			fieldValFloat64 = float64(fieldValFloat)
			log.Debugf("Float64: %v", fieldValFloat64)
		case reflect.Bool:
			fieldValBool, err = strconv.ParseBool(fieldVal)
			if err != nil {
				return out, err
			}
		default:
			log.Errorf("Unsupported type of '%s' in %v", field.Type, field.Name)
		}

		if field.Name == "Delay" {
			out.Delay = fieldValFloat32
		} else if field.Name == globalDelay {
			out.GlobalDelay = fieldValBool
		} else if field.Name == "Order" {
			out.Order = fieldValInt
		}
	}

	return out, err
}

// Action will convert for an Action model.
func (ss StringsToStruct) Action(in interface{}, out models.Action) (models.Action, error) {
	fieldsIn := reflect.TypeOf(in)
	valuesIn := reflect.ValueOf(in)

	fields := reflect.TypeOf(out)
	values := reflect.ValueOf(out)
	num := fields.NumField()

	var err error

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		_, found := fieldsIn.FieldByName(field.Name)
		if !found {
			continue
		}

		fieldVal := valuesIn.FieldByName(field.Name).String()
		if fieldVal == "" {
			continue
		}

		var (
			fieldValString  string
			fieldValBool    bool
			fieldValInt     int
			fieldValInt32   int32
			fieldValInt64   int64
			fieldValFloat32 float32
			fieldValFloat64 float64
		)

		switch value.Kind() {
		case reflect.String:
			if !value.IsValid() {
				return out, fmt.Errorf("no such field: %s in obj", field.Name)
			}

			fieldValString = fieldVal
		case reflect.Int:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}
		case reflect.Int32:
			fieldValInt, err := strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt32 = int32(fieldValInt)
			log.Debugf("Int32: %v", fieldValInt32)
		case reflect.Int64:
			fieldValInt, err = strconv.Atoi(fieldVal)
			if err != nil {
				return out, err
			}

			fieldValInt64 = int64(fieldValInt)
			log.Debugf("Int64: %v", fieldValInt64)
		case reflect.Float32:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, thirtyTwo)
			if err != nil {
				return out, err
			}

			fieldValFloat32 = float32(fieldValFloat)
			log.Debugf("Float32: %v", fieldValFloat32)
		case reflect.Float64:
			fieldValFloat, err := strconv.ParseFloat(fieldVal, sixtyFour)
			if err != nil {
				return out, err
			}

			fieldValFloat64 = float64(fieldValFloat)
			log.Debugf("Float64: %v", fieldValFloat64)
		case reflect.Bool:
			fieldValBool, err = strconv.ParseBool(fieldVal)
			if err != nil {
				return out, err
			}

			log.Debugf("Bool: %v", fieldValBool)
		default:
			log.Errorf("Unsupported type of '%s' in %v", field.Type, field.Name)
		}

		if field.Name == "Command" {
			out.Command = fieldValString
		} else if field.Name == "Parameter" {
			out.Parameter = fieldValString
		} else if field.Name == "GlobalParameter" {
			out.GlobalParameter = fieldValString
		} else if field.Name == "Order" {
			out.Order = fieldValInt
		}
	}

	return out, err
}
