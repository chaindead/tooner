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
	"github.com/yinxulai/go-jsonrepair/jsonrepair"
)

func convert(logger *log.Logger, line string) string {
	newLine := line

	contents := gjson.Get(line, "result.content").Array()
	if gjson.Get(line, "result.structuredContent").Exists() && len(contents) > 0 {
		var err error
		newLine, err = sjson.Delete(newLine, "result.structuredContent")
		if err != nil {
			logger.Println("!!! remove structured:", newLine)
		}
	}

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
		repaired, err2 := jsonRepair(cData)
		if err2 != nil {
			return "", errors.Wrapf(err, "try json repair(%s)", err2)
		}

		err2 = json.Unmarshal([]byte(repaired), &sData)
		if err2 != nil {
			return "", errors.Wrapf(err, "unmarshal repaired(%s) json:", err2)
		}
	}

	sData = normalizeAny(sData)

	tData, err := toon.Marshal(sData, toon.WithLengthMarkers(true))
	if err != nil {
		return "", errors.Wrap(err, "marshal toon")
	}

	return string(tData), nil
}

func jsonRepair(data string) (string, error) {
	repaired, err := jsonrepair.Repair(data)
	if err != nil {
		return "", errors.Wrap(err, "repair")
	}

	beforeLen := len(data)
	afterLen := len(repaired)
	decrease := float64(beforeLen-afterLen) / float64(beforeLen)

	if decrease > 0.2 {
		return "", fmt.Errorf("repair decrease: %.2f", decrease)
	}

	return repaired, nil
}
