package common

import (
	"fmt"
	"testing"
)

func TestIP(t *testing.T) {
	IP := GetIPFromNetWork("en0")
	fmt.Println(IP)
}

func TestFile(t *testing.T) {
	PrintFilesName("/Users/ywaz/test", "lo0")
}
