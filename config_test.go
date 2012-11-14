package config

import (
        "fmt"
        "testing"
)

const (
        goodTest = "testdata/test.conf"
        badTest = "testdata/bad.conf"
)

func TestGoodConfig(t *testing.T) {
        fmt.Printf("[+] validating known-good config... ")
        cmap, err := ParseFile(goodTest)
        if err != nil {
                fmt.Printf("\n[!] test failure: %s\n", err.Error())
                t.FailNow()
        } else if len(cmap) != 2 {
                fmt.Printf("\n[!] failed to load config file\n")
                t.FailNow()
        }
        fmt.Println("ok")
}

func TestBadConfig(t *testing.T) {
        fmt.Printf("[+] ensure invalid config file fails... ")
        _, err := ParseFile(badTest)
        if err == nil {
                fmt.Printf("\n[!] parse should have failed, but didn't!\n")
                t.FailNow()
        }
        fmt.Println("ok")
}
