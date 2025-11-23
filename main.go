package struct2pflag

import (
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

type tagType int // tag type

const (
	TAG_NONE  tagType = iota // no tag
	TAG_FLAG                 // flag tag
	TAG_PFLAG                // pflag tag
)

// Bind binds the exported fields of a configuration struct to flags on the given pflag.FlagSet.
//
// Bind expects cfg to be a pointer to a struct. It walks the struct's exported fields and, for each
// field that carries either a "pflag" or "flag" struct tag, registers a corresponding flag on fs.
// Fields without those tags or unexported fields are ignored.
//
// Tag formats:
//   - flag:  "longname,usage"
//     If the tag does not contain a comma the tag value is treated as the usage text and the
//     long flag name defaults to the lowercase field name.
//   - pflag: "longname,shortname,usage"
//     If shortname is omitted (no second comma) shortname is set to the empty string and the
//     remainder is treated as usage. If the tag doesn't contain a comma at all the behavior is
//     the same as for "flag".
//
// Recursion:
//   - If a field is itself a struct, Bind is called recursively on that struct value.
//   - If a field is a pointer to a struct, the pointer is followed and Bind is called only if the
//     pointer is non-nil (nil pointers are skipped).
//
// Supported field kinds:
//   - bool, int, uint, string
//
// For each supported field Bind registers the corresponding fs.*VarP binding using the field's
// address and the field's current value as the flag default.
//
// Notes:
//   - Only exported struct fields are considered for binding.
//   - The function uses the field's current value as the default for the registered flag.
//   - shortname for pflag entries may be supplied as the second tag component; if omitted an
//     empty short name is used.
func Bind(fs *pflag.FlagSet, cfg interface{}) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := range v.NumField() {
		f := v.Field(i)
		field := t.Field(i)

		if !field.IsExported() { // unexported field, skip
			continue
		}
		var desc string
		var ok bool
		var t tagType
		if desc, ok = field.Tag.Lookup("pflag"); ok {
			t = TAG_PFLAG
		} else if desc, ok = field.Tag.Lookup("flag"); ok {
			t = TAG_FLAG
		} else {
			continue
		}
		switch f.Kind() {
		case reflect.Struct:
			Bind(fs, f.Addr().Interface()) // Bind the struct field
			continue
		case reflect.Pointer:
			if f.Type().Elem().Kind() == reflect.Struct && !f.IsNil() {
				Bind(fs, f.Interface()) // Bind the struct pointer field
				continue
			}
		}
		var longname, shortname, usage string
		switch t {
		case TAG_FLAG:
			if longname, usage, ok = strings.Cut(desc, ","); !ok {
				longname = strings.ToLower(field.Name)
				usage = desc
			}
		case TAG_PFLAG:
			var rest string
			if longname, rest, ok = strings.Cut(desc, ","); !ok {
				longname = strings.ToLower(field.Name)
				usage = desc
			} else if shortname, usage, ok = strings.Cut(rest, ","); !ok {
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

// BindDefault registers flags for the given configuration value on the
// default pflag.CommandLine flag set. It is a convenience wrapper around Bind
// that saves callers from having to specify the FlagSet explicitly — pass a
// pointer to your config struct and BindDefault will create the corresponding
// command-line flags on the global pflag.CommandLine.
func BindDefault(cfg interface{}) {
	Bind(pflag.CommandLine, cfg)
}
