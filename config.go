package goconfig

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

var DefaultSection = "default"

// ParseFile takes the filename as a string and returns a ConfigMap.
func ParseFile(fileName string) (cfg ConfigMap, err error) {
	var file *os.File

	cfg = make(ConfigMap, 0)
	file, err = os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	buf := bufio.NewReader(file)

	var (
		line           string
		longLine       bool
		currentSection string
		lineBytes      []byte
		isPrefix       bool
	)

	for {
		err = nil
		lineBytes, isPrefix, err = buf.ReadLine()
		if io.EOF == err {
			err = nil
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
			if currentSection == "" {
				currentSection = DefaultSection
				if !cfg.SectionInConfig(currentSection) {
					cfg[currentSection] = make(map[string]string, 0)
				}
			}
			key := configLine.ReplaceAllString(line, "$1")
			val := configLine.ReplaceAllString(line, "$2")
			if key == "" {
				continue
			}
			cfg[currentSection][key] = val
		} else {
			err = fmt.Errorf("invalid config file")
			break
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

// WriteFile writes out the configuration to a file.
func (c *ConfigMap) WriteFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, section := range c.ListSections() {
		sName := fmt.Sprintf("[ %s ]\n", section)
		_, err = file.Write([]byte(sName))
		if err != nil {
			return
		}

		for k, v := range (*c)[section] {
			line := fmt.Sprintf("%s = %s\n", k, v)
			_, err = file.Write([]byte(line))
			if err != nil {
				return
			}
		}
		_, err = file.Write([]byte{0x0a})
		if err != nil {
			return
		}
	}
	return
}

// AddSection creates a new section in the config map.
func (c *ConfigMap) AddSection(section string) {
	if nil != (*c)[section] {
		(*c)[section] = make(map[string]string, 0)
	}
}

// AddKeyVal adds a key value pair to a config map.
func (c *ConfigMap) AddKeyVal(section, key, val string) {
	if "" == section {
		section = DefaultSection
	}

	if nil == (*c)[section] {
		c.AddSection(section)
	}

	(*c)[section][key] = val
}

// Retrieve the value from a key map.
func (c *ConfigMap) GetValue(section, key string) (val string, present bool) {
	if c == nil {
		return
	}

	if section == "" {
		section = DefaultSection
	}

	cm := *c
	_, ok := cm[section]
	if !ok {
		return
	}

	val, present = cm[section][key]
	return
}
