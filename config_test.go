package goconfig

import (
        "fmt"
        "testing"
)

// FailWithError is a utility for dumping errors and failing the test.
func FailWithError(t *testing.T, err error) {
        fmt.Println("failed")
        if err != nil {
                fmt.Println("[!] ", err.Error())
        }
        t.FailNow()
}

func TestGoodConfig(t *testing.T) {
        testFile := "testdata/test.conf"
        fmt.Printf("[+] validating known-good config... ")
        cmap, err := ParseFile(testFile)
        if err != nil {
                FailWithError(t, err)
        } else if len(cmap) != 2 {
                FailWithError(t, err)
        }
        fmt.Println("ok")
}

func TestGoodConfig2(t *testing.T) {
        testFile := "testdata/test2.conf"
        fmt.Printf("[+] validating second known-good config... ")
        cmap, err := ParseFile(testFile)
        if err != nil {
                FailWithError(t, err)
        } else if len(cmap) != 1 {
                FailWithError(t, err)
        } else if len(cmap["default"]) != 3 {
                FailWithError(t, err)
        }
        fmt.Println("ok")
}

func TestBadConfig(t *testing.T) {
        testFile := "testdata/bad.conf"
        fmt.Printf("[+] ensure invalid config file fails... ")
        _, err := ParseFile(testFile)
        if err == nil {
                FailWithError(t, err)
        }
        fmt.Println("ok")
}
