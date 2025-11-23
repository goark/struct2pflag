package struct2pflag_test

import (
	"testing"

	"github.com/goark/struct2pflag"
	"github.com/spf13/pflag"
)

type subStruct1a struct {
	Boolean1a bool `pflag:"bool1a,b,Enable boolean 1a mode"`
}
type subStruct1b struct {
	Boolean1b bool `pflag:"bool1b,B,Enable boolean 1b mode"`
}
type subStruct2a struct {
	Boolean2a bool `flag:"bool2a,Enable boolean 2a mode"`
}
type subStruct2b struct {
	Boolean2b bool `flag:"bool2b,Enable boolean 2b mode"`
}

type tsStruct struct {
	Sub0a subStruct1a  // no tag, should be ignored
	Sub0b *subStruct1a `flag:""`  // nil pointer, should be ignored
	Sub1a subStruct1a  `flag:""`  // empty tag, should be processed
	Sub1b subStruct1b  `pflag:""` // empty tag, should be processed
	Sub2a *subStruct2a `flag:""`  // pointer to struct, should be processed
	Sub2b subStruct2b  `pflag:""` // empty tag, should be processed
}

func TestPflagBindStruct(t *testing.T) {
	ts := &tsStruct{Sub2a: &subStruct2a{}} // initialize pointer field
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
	struct2pflag.Bind(flagSet, ts)
	err := flagSet.Parse(
		[]string{
			"--bool1a",
			"--bool1b",
			"--bool2a",
			"--bool2b",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if expect := false; ts.Sub0a.Boolean1a != expect {
		t.Errorf("expect %#v,but %#v", ts.Sub0a.Boolean1a, expect)
	}
	if expect := true; ts.Sub1a.Boolean1a != expect {
		t.Errorf("expect %#v,but %#v", ts.Sub1a.Boolean1a, expect)
	}
	if expect := true; ts.Sub1b.Boolean1b != expect {
		t.Errorf("expect %#v,but %#v", ts.Sub1b.Boolean1b, expect)
	}
	if expect := true; ts.Sub2a.Boolean2a != expect {
		t.Errorf("expect %#v,but %#v", ts.Sub2a.Boolean2a, expect)
	}
	if expect := true; ts.Sub2b.Boolean2b != expect {
		t.Errorf("expect %#v,but %#v", ts.Sub2b.Boolean2b, expect)
	}
}
