package utils

import (
	"io/ioutil"
	"regexp"
	"strings"
)

type Config struct {
	Zones map[string]map[string]string
}

func (config *Config) GetItem(zone string, key string) string {
	Zone, ok := config.Zones[zone]
	if ok == false {
		return ""
	}

	value, ok := Zone[key]
	if ok == false {
		return ""
	}
	return value

}

func (config *Config) Initialize(filename string) {

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	fileString := string(file)
	re := regexp.MustCompile("\\[.*\\]")
	zonesIndex := re.FindAllStringIndex(fileString, -1)
	zonesIndex = append(zonesIndex, []int{len(fileString), len(fileString)})
	var zones []string
	for i := 0; i < len(zonesIndex)-1; i++ {
		zones = append(zones, fileString[zonesIndex[i][0]:zonesIndex[i+1][0]])
	}
	config.Zones = make(map[string]map[string]string)
	for _, zone := range zones {
		_lines := strings.Split(zone, "\n")
		var zoneName string
		var lines []string
		for _, line := range _lines {
			zoneRegex := regexp.MustCompile("(?s)\\[.*]")
			if zoneRegex.MatchString(line) == true {
				zoneName = strings.Trim(line, "\r \n\t[]")
			}
			if strings.Trim(line, "\r \n\t") != "" && strings.Contains(line, "=") {
				lines = append(lines, strings.Trim(line, "\r \n\t"))
			}
		}
		config.Zones[zoneName] = make(map[string]string)
		for _, _item := range lines {
			item := strings.Split(_item, "=")
			key := strings.Trim(item[0], " \r\t")
			value := strings.Trim(item[1], " \r\t")
			config.Zones[zoneName][key] = value
		}
	}
	return
}
