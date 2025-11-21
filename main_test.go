package struct2pflag_test

import (
	"os"
	"testing"

	"github.com/goark/struct2pflag"
	"github.com/spf13/pflag"
)

type ts struct {
	S1 string `pflag:"string option(1)"`
	S2 string `pflag:"S2,string option(2)"`
	S3 string `pflag:"S3,s,string option(3)"`

	I1 int `pflag:"integer option(1)"`
	I2 int `pflag:"I2,integer option(2)"`
	I3 int `pflag:"I3,i,integer option(3)"`

	B1 bool `pflag:"boolean option(1)"`
	B2 bool `pflag:"B2,boolean option(2)"`
	B3 bool `pflag:"B3,b,boolean option(3)"`
}

func TestBindLongname(t *testing.T) {
	ts1 := &ts{}
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
	struct2pflag.Bind(flagSet, ts1)
	err := flagSet.Parse(
		[]string{
			"--s1", "foo",
			"--S2", "bar",
			"--i1", "9",
			"--I2", "8",
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
	if expect := 9; ts1.I1 != expect {
		t.Errorf("expect %#v,but %#v", ts1.I1, expect)
	}
	if expect := 8; ts1.I2 != expect {
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

func TestBindShortname(t *testing.T) {
	ts2 := &ts{}
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
	if expect := 0; ts2.I1 != expect {
		t.Errorf("expect %#v,but %#v", ts2.I1, expect)
	}
	if expect := 0; ts2.I2 != expect {
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

func TestBindErr(t *testing.T) {
	ts3 := &ts{}
	flagSet3 := pflag.NewFlagSet("", pflag.ContinueOnError)
	struct2pflag.Bind(flagSet3, ts3)
	cases := [][]string{
		{"--S1", "should be lower case"},
		{"--s2", "should be upper case"},
		{"-S", "should be lower case"},
		{"--I1", "should be lower case"},
		{"--i2", "should be upper case"},
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
