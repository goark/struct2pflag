package struct2pflag_test

import (
	"os"
	"testing"

	"github.com/goark/struct2pflag"
	"github.com/spf13/pflag"
)

type tsPflag struct {
	S1 string `pflag:"string option(1)"`      // name omitted
	S2 string `pflag:"S2,string option(2)"`   // long name only
	S3 string `pflag:"S3,s,string option(3)"` // long and short names
	s4 string `pflag:"s4,s,string option(4)"` // unexported field, should be ignored
	S5 string `other:"string option(5)"`      // no pflag/flag tag, should be ignored

	I1 int  `pflag:"integer option(1)"`              // name omitted
	I2 uint `pflag:"UI2,unsigned integer option(2)"` // long name only
	I3 int  `pflag:"I3,i,integer option(3)"`         // long and short names

	B1 bool `pflag:"boolean option(1)"`      // name omitted
	B2 bool `pflag:"B2,boolean option(2)"`   // long name only
	B3 bool `pflag:"B3,b,boolean option(3)"` // long and short names
}

func TestPflagBindLongname(t *testing.T) {
	ts1 := &tsPflag{}
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
	struct2pflag.Bind(flagSet, ts1)
	err := flagSet.Parse(
		[]string{
			"--s1", "foo",
			"--S2", "bar",
			"--i1", "9",
			"--UI2", "8",
			"--b1",
			"--B2",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if expect := "foo"; ts1.S1 != expect {
		t.Errorf("expect %#v,but %#v", ts1.S1, expect)
	}
	if expect := "bar"; ts1.S2 != expect {
		t.Errorf("expect %#v,but %#v", ts1.S2, expect)
	}
	if expect := ""; ts1.S3 != expect {
		t.Errorf("expect %#v,but %#v", ts1.S3, expect)
	}
	if expect := ""; ts1.s4 != expect {
		t.Errorf("expect %#v,but %#v", ts1.s4, expect)
	}
	if expect := ""; ts1.S5 != expect {
		t.Errorf("expect %#v,but %#v", ts1.S5, expect)
	}
	if expect := 9; ts1.I1 != expect {
		t.Errorf("expect %#v,but %#v", ts1.I1, expect)
	}
	if expect := uint(8); ts1.I2 != expect {
		t.Errorf("expect %#v,but %#v", ts1.I2, expect)
	}
	if expect := 0; ts1.I3 != expect {
		t.Errorf("expect %#v,but %#v", ts1.I3, expect)
	}
	if expect := true; ts1.B1 != expect {
		t.Errorf("expect %#v,but %#v", ts1.B1, expect)
	}
	if expect := true; ts1.B2 != expect {
		t.Errorf("expect %#v,but %#v", ts1.B2, expect)
	}
	if expect := false; ts1.B3 != expect {
		t.Errorf("expect %#v,but %#v", ts1.B3, expect)
	}
}

func TestPflagBindShortname(t *testing.T) {
	ts2 := &tsPflag{}
	flagSet2 := pflag.NewFlagSet("", pflag.ContinueOnError)
	struct2pflag.Bind(flagSet2, ts2)
	err := flagSet2.Parse(
		[]string{
			"-s", "foo",
			"-i", "9",
			"-b",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if expect := ""; ts2.S1 != expect {
		t.Errorf("expect %#v,but %#v", ts2.S1, expect)
	}
	if expect := ""; ts2.S2 != expect {
		t.Errorf("expect %#v,but %#v", ts2.S2, expect)
	}
	if expect := "foo"; ts2.S3 != expect {
		t.Errorf("expect %#v,but %#v", ts2.S3, expect)
	}
	if expect := ""; ts2.s4 != expect {
		t.Errorf("expect %#v,but %#v", ts2.s4, expect)
	}
	if expect := ""; ts2.S5 != expect {
		t.Errorf("expect %#v,but %#v", ts2.S5, expect)
	}
	if expect := 0; ts2.I1 != expect {
		t.Errorf("expect %#v,but %#v", ts2.I1, expect)
	}
	if expect := uint(0); ts2.I2 != expect {
		t.Errorf("expect %#v,but %#v", ts2.I2, expect)
	}
	if expect := 9; ts2.I3 != expect {
		t.Errorf("expect %#v,but %#v", ts2.I3, expect)
	}
	if expect := false; ts2.B1 != expect {
		t.Errorf("expect %#v,but %#v", ts2.B1, expect)
	}
	if expect := false; ts2.B2 != expect {
		t.Errorf("expect %#v,but %#v", ts2.B2, expect)
	}
	if expect := true; ts2.B3 != expect {
		t.Errorf("expect %#v,but %#v", ts2.B3, expect)
	}
}

func TestPflagBindErr(t *testing.T) {
	ts3 := &tsPflag{}
	flagSet3 := pflag.NewFlagSet("", pflag.ContinueOnError)
	struct2pflag.Bind(flagSet3, ts3)
	cases := [][]string{
		{"--S1", "should be lower case"},
		{"--s2", "should be upper case"},
		{"-S", "should be lower case"},
		{"--s4", "should be ignored"},
		{"--S5", "should be ignored"},
		{"--I1", "should be lower case"},
		{"--ui2", "should be upper case"},
		{"-I", "should be lower case"},
		{"--B1"},
		{"--b2"},
		{"-B"},
	}
	stderrSaved := os.Stderr
	devnull, err := os.Create(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = devnull.Close() }()
	for _, case1 := range cases {
		os.Stderr = devnull
		err := flagSet3.Parse(case1)
		os.Stderr = stderrSaved
		if err == nil {
			t.Errorf("expect error, but succeeded for %#v", case1)
		}
	}
}
