package config

import (
	"testing"
	"time"
)

func TestGetValStr(t *testing.T) {
	t.Setenv("PRIFIX_NAME_STR", "val")
	v := GetString("prifix", "name.str")
	if v != "val" {
		t.Errorf("string dont match [%v]", v)
	}

	t.Setenv("NAME_STR", "val")
	v1 := GetString("", "name.str")
	if v != "val" {
		t.Errorf("string dont match [%v]", v1)
	}

}

func TestGetValInt(t *testing.T) {
	t.Setenv("PRIFIX_NAME_I", "2")
	v := GetInt("prifix", "name.i", 0)
	if v != 2 {
		t.Errorf("string dont match [%v]", v)
	}

	t.Setenv("NAME_I", "3")
	v1 := GetInt("", "name.i", 0)
	if v1 != 3 {
		t.Errorf("string dont match [%v]", v1)
	}

	v2 := GetInt("asdasdasdasdasd", "name.i", 2)
	if v2 != 2 {
		t.Errorf("string dont match [%v]", v2)
	}

}

func TestGetDurationn(t *testing.T) {
	t1 := time.Second * 2
	t.Setenv("PRIFIX_NAME_D", "2s")
	v := GetDuration("prefix", "NAME_D", 2*time.Second)
	if v != t1 {
		t.Errorf("string dont match [%v]", v)
	}

	t2 := time.Second * 3
	v1 := GetDuration("prefix", "NAME_D", 3*time.Second)
	if v1 != t2 {
		t.Errorf("string dont match [%v]", v)
	}

}

func TestBB(t *testing.T) {

	b1 := true
	t.Setenv("NAME_BB", "true")
	b1r := GetBool("", "name.bb")
	if b1 != b1r {
		t.Errorf("string dont match [%v]", b1r)

	}

	b2 := true
	t.Setenv("F_NAME_BB", "true")
	b2r := GetBool("f", "name.bb")
	if b2 != b2r {
		t.Errorf("string dont match [%v]", b2r)

	}

	b3r := GetBool("ff", "name.bb")
	if b3r {
		t.Errorf("string dont match [%v]", b3r)

	}

}
