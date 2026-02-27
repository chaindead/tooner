package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/toon-format/toon-go"
)

func convert(logger *log.Logger, line string) string {
	newLine := line

	if gjson.Get(line, "result.structuredContent").Exists() {
		// TODO: case if structuredContent but not simple content
		var err error
		newLine, err = sjson.Delete(newLine, "result.structuredContent")
		if err != nil {
			logger.Println("!!! remove structured:", newLine)
		}
	}

	contents := gjson.Get(line, "result.content").Array()
	for i, content := range contents {
		if content.Get("type").String() != "text" {
			continue
		}

		data := content.Get("text").String()
		newData, err := json2toon(data)
		if err != nil {
			logger.Println("!!! convert:", err)
			continue
		}

		pos := fmt.Sprintf("result.content.%d.text", i)
		newLine, err = sjson.Set(newLine, pos, newData)
		if err != nil {
			logger.Println("!!! set:", err)
		}
	}

	return newLine
}

func json2toon(data string) (string, error) {
	cData := strings.TrimSpace(data)

	if len(cData) == 0 {
		return "", errors.New("data is empty")
	}

	if cData[0] != '{' && cData[0] != '[' {
		// simple string
		return data, nil
	}

	var sData any
	if err := json.Unmarshal([]byte(cData), &sData); err != nil {
		return "", errors.Wrap(err, "unmarshal json")
	}

	tData, err := toon.Marshal(sData, toon.WithLengthMarkers(true))
	if err != nil {
		return "", errors.Wrap(err, "marshal toon")
	}

	return string(tData), nil
}
