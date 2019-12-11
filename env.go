package config

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ApplyEnvOverrides apple env overrides
func (c *Configor) ApplyEnvOverrides(prefix string, val interface{}) error {
	return c.applyEnvOverrides(prefix, reflect.ValueOf(val), "")
}

func (c *Configor) applyEnvOverrides(prefix string, spec reflect.Value, structKey string) error {
	element := spec
	value := os.Getenv(prefix)
	if c.Options.ShowLog && len(value) != 0 {
		fmt.Printf("[Config]Loading env %v for field %v\n", prefix, structKey)
	}
	// If spec is a named type and is addressable,
	// check the address to see if it implements encoding.TextUnmarshaler.
	if spec.Kind() != reflect.Ptr && spec.Type().Name() != "" && spec.CanAddr() {
		v := spec.Addr()
		if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
			if len(value) == 0 {
				return nil
			}
			return u.UnmarshalText([]byte(value))
		}
	}
	// If we have a pointer, dereference it
	if spec.Kind() == reflect.Ptr {
		element = spec.Elem()
	}

	switch element.Kind() {
	case reflect.String:
		if len(value) == 0 {
			return nil
		}
		element.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(value) == 0 {
			return nil
		}
		intValue, err := strconv.ParseInt(value, 0, element.Type().Bits())
		if err != nil {
			return fmt.Errorf("Failed to apply %v to %v using type %v and value '%v': %s", prefix, structKey, element.Type().String(), value, err)
		}
		element.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(value) == 0 {
			return nil
		}
		intValue, err := strconv.ParseUint(value, 0, element.Type().Bits())
		if err != nil {
			return fmt.Errorf("Failed to apply %v to %v using type %v and value '%v': %s", prefix, structKey, element.Type().String(), value, err)
		}
		element.SetUint(intValue)
	case reflect.Bool:
		if len(value) == 0 {
			return nil
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to apply %v to %v using type %v and value '%v': %s", prefix, structKey, element.Type().String(), value, err)
		}
		element.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		if len(value) == 0 {
			return nil
		}
		floatValue, err := strconv.ParseFloat(value, element.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to apply %v to %v using type %v and value '%v': %s", prefix, structKey, element.Type().String(), value, err)
		}
		element.SetFloat(floatValue)
	case reflect.Slice:
		if len, err := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_LEN", prefix))); err == nil && len > 0 {
			t := reflect.MakeSlice(reflect.SliceOf(element.Type().Elem()), len, len)
			for i := 0; i < len; i++ {
				f := t.Index(i)
				if err := c.applyEnvOverrides(fmt.Sprintf("%s_%d", prefix, i), f, structKey); err != nil {
					return err
				}
			}
			element.Set(t)
			break
		}
		// If the type is s slice, apply to each using the index as a suffix, e.g. GRAPHITE_0, GRAPHITE_0_TEMPLATES_0 or GRAPHITE_0_TEMPLATES="item1,item2"
		for j := 0; j < element.Len(); j++ {
			f := element.Index(j)
			if err := c.applyEnvOverrides(prefix, f, structKey); err != nil {
				return err
			}

			if err := c.applyEnvOverrides(fmt.Sprintf("%s_%d", prefix, j), f, structKey); err != nil {
				return err
			}
		}

		// If the type is s slice but have value not parsed as slice e.g. GRAPHITE_0_TEMPLATES="item1,item2"
		if element.Len() == 0 && len(value) > 0 {
			rules := strings.Split(value, ",")

			for _, rule := range rules {
				element.Set(reflect.Append(element, reflect.ValueOf(rule)))
			}
		}
	case reflect.Struct:
		typeOfSpec := element.Type()
		for i := 0; i < element.NumField(); i++ {
			field := element.Field(i)
			// Skip any fields that we cannot set
			if !field.CanSet() && field.Kind() != reflect.Slice {
				continue
			}

			structField := typeOfSpec.Field(i)
			fieldName := structField.Name

			configName := structField.Tag.Get("toml")
			if configName == "-" {
				continue // Skip fields with tag `toml:"-"`.
			}

			if configName == "" && structField.Anonymous {
				// Embedded field without a toml tag.
				// Don't modify prefix.
				if err := c.applyEnvOverrides(prefix, field, fieldName); err != nil {
					return err
				}
				continue
			}

			// except tag options, like omitempty or omitzero
			configName = strings.Split(configName, ",")[0]

			// Replace hyphens with underscores to avoid issues with shells
			configName = strings.Replace(configName, "-", "_", -1)

			envKey := strings.ToUpper(configName)
			if prefix != "" {
				envKey = strings.ToUpper(fmt.Sprintf("%s_%s", prefix, configName))
			}

			// If it's a sub-config, recursively apply
			if field.Kind() == reflect.Struct || field.Kind() == reflect.Ptr ||
				field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
				if err := c.applyEnvOverrides(envKey, field, fieldName); err != nil {
					return err
				}
				continue
			}

			value := os.Getenv(envKey)
			// Skip any fields we don't have a value to set
			if len(value) == 0 {
				continue
			}

			if err := c.applyEnvOverrides(envKey, field, fieldName); err != nil {
				return err
			}
		}
	}
	return nil
}
