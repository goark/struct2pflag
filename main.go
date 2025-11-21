package struct2pflag

import (
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

func Bind(fs *pflag.FlagSet, cfg interface{}) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}
		desc, ok := field.Tag.Lookup("pflag")
		if !ok {
			continue
		}
		switch f.Kind() {
		case reflect.Struct:
			Bind(fs, f.Addr().Interface())
			continue
		case reflect.Pointer:
			if f.Type().Elem().Kind() == reflect.Struct && !f.IsNil() {
				Bind(fs, f.Interface())
				continue
			}
		}
		var shortname, usage string
		longname, rest, ok := strings.Cut(desc, ",")
		if !ok {
			longname = strings.ToLower(field.Name)
			usage = desc
		} else {
			shortname, usage, ok = strings.Cut(rest, ",")
			if !ok {
				shortname = ""
				usage = rest
			}
		}

		switch f.Kind() {
		case reflect.Bool:
			fs.BoolVarP(f.Addr().Interface().(*bool), longname, shortname, f.Bool(), usage)
		case reflect.Int:
			fs.IntVarP(f.Addr().Interface().(*int), longname, shortname, int(f.Int()), usage)
		case reflect.Uint:
			fs.UintVarP(f.Addr().Interface().(*uint), longname, shortname, uint(f.Uint()), usage)
		case reflect.String:
			fs.StringVarP(f.Addr().Interface().(*string), longname, shortname, f.String(), usage)
		}
	}
}

func BindDefault(cfg interface{}) {
	Bind(pflag.CommandLine, cfg)
}
