package config

import (
	"fmt"
	"strconv"
)

func StoI(str string, i int) int {
	o, err := strconv.Atoi(str)
	if err != nil {
		return i
	}
	return o
}

func ItoS(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func ItoI(i interface{}) int {
	iAreaId, ok := i.(int)
	if ok {
		return iAreaId
	}
	return 0
}

func ItoB(i interface{}) bool {
	iAreaId, ok := i.(bool)
	if ok {
		return iAreaId
	}
	return false
}
