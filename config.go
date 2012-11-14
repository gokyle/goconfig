package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

// ConfigMap is shorthand for the type used as a config struct.
type ConfigMap map[string]map[string]string

var (
	configSection = regexp.MustCompile("^\\s*\\[\\s*(\\w+)\\s*\\]\\s*$")
	configLine    = regexp.MustCompile("^\\s*(\\w+)\\s*=\\s*(.*)\\s*$")
	commentLine   = regexp.MustCompile("^#.*$")
	blankLine     = regexp.MustCompile("^\\s*$")
)

// ParseFile takes the filename as a string and returns a ConfigMap.
func ParseFile(fileName string) (cfg ConfigMap, err error) {
	cfg = make(ConfigMap, 0)
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	buf := bufio.NewReader(file)

	var (
		line           string
		longLine       bool
		currentSection string
	)

	fmt.Println("parsing")
	for {
		err = nil
		lineBytes, isPrefix, err := buf.ReadLine()
		if io.EOF == err {
			break
		} else if err != nil {
			break
		} else if isPrefix {
			line += string(lineBytes)

			longLine = true
			continue
		} else if longLine {
			line += string(lineBytes)
			longLine = false
		} else {
			line = string(lineBytes)
		}
		if commentLine.MatchString(line) {
			continue
		} else if blankLine.MatchString(line) {
			continue
		} else if configSection.MatchString(line) {
			section := configSection.ReplaceAllString(line,
				"$1")
			if section == "" {
				err = fmt.Errorf("invalid structure in file")
				break
			} else if !cfg.SectionInConfig(section) {
				cfg[section] = make(map[string]string, 0)
			}
			currentSection = section
		} else if configLine.MatchString(line) {
			key := configLine.ReplaceAllString(line, "$1")
			val := configLine.ReplaceAllString(line, "$2")
			if key == "" {
				continue
			}
			cfg[currentSection][key] = val
		}
	}
	return
}

// SectionInConfig determines whether a section is in the configuration.
func (c *ConfigMap) SectionInConfig(section string) bool {
	for s, _ := range *c {
		if section == s {
			return true
		}
	}
	return false
}

// ListSections returns the list of sections in the config map.
func (c *ConfigMap) ListSections() (sections []string) {
	for section, _ := range *c {
		sections = append(sections, section)
	}
	return
}
