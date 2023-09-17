package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dattaray-basab/cks-clip-lib/globals"
	"github.com/dattaray-basab/cks-clip-lib/logger"
)

var GetFirstMoveItem = func(templateMap map[string]string) (string, error) {
	moveItems := templateMap[globals.KEY_MOVE_ITEMS]
	moveItemParts := strings.Split(moveItems, ":")
	if len(moveItemParts) == 0 {
		err := fmt.Errorf("no move-item is available")
		return "", err
	}
	firstMoveItem := moveItemParts[0]
	return firstMoveItem, nil
}

var BuildAlterInfrastucture = func(templateMap map[string]string, queryTemplate, controlTemplate string) (globals.SubstitionTemplateT, error){
	var getQueryFilePath = func(templateMap map[string]string) (string, error) { //?1
		dirpath := filepath.Join(templateMap[globals.KEY_TARGET_PATH], globals.TOKENS_DIRNAME, globals.QUERY_DIRNAME)
		if !IsDir(dirpath) {
			err := os.MkdirAll(dirpath, os.ModePerm)
			if err != nil {
				return "", err
			}
		}
		fullAlterRelPath := templateMap[globals.KEY_FULL_ALTER_REL_PATH]
		queryPathName := strings.Replace(fullAlterRelPath, "/", "", -1)
		fName := queryPathName + globals.JSON_EXT
		fPath := filepath.Join(dirpath, fName)
		return fPath, nil
	}

	var getQueryId = func(templateMap map[string]string, queryFilePath string) (string, error) {
		queryFileName := filepath.Base(queryFilePath)
		queryName := queryFileName[:len(queryFileName)-len(globals.JSON_EXT)]
		logger.Log.Debug(queryName)
		suffix := 0
		queryId := "ID_" + strconv.Itoa(suffix)
		fullQueryId := queryName + "." + queryId //??? TODO: check if this is correct

		return fullQueryId, nil
	}

	var alterRecord globals.SubstitionTemplateT

	queryFilePath, err := getQueryFilePath(templateMap)
	if err != nil {
		return alterRecord, err
	}
	moveItemMap, err := GetMoveItemMap(templateMap)
	if err != nil {
		return alterRecord, err
	}

	logger.Log.Debug(moveItemMap)
	fullQueryId, err := getQueryId(templateMap, queryFilePath)
	if err != nil {
		return alterRecord, err
	}
	queryIdParts := strings.Split(fullQueryId, ".")
	shortQueryId := queryIdParts[len(queryIdParts)-1]

	firstFilePath, firstWordInFirstFile, err := GetFirstLineOfFirstFile(templateMap)
	if err != nil {
		return alterRecord, err
	}
	logger.Log.Debug(firstFilePath)
	logger.Log.Debug(firstWordInFirstFile)
	templateMap[globals.KEY_FIRST_WORD_IN_FIRST_FILE] = firstWordInFirstFile

	alterRecord =
		globals.SubstitionTemplateT{
			FullQueryId:  fullQueryId,
			ShortQueryId: shortQueryId,

			MoveItemsInfo:        moveItemMap,
			FirstWordInFirstFile: firstWordInFirstFile,
			FirstFilePath:        firstFilePath,
		}

	contentQuery, error := RunTemplate(queryTemplate, alterRecord)
	if error != nil {
		return alterRecord, error
	}
	err = MakeQueryTokenFile(templateMap, contentQuery, queryFilePath)
	if err != nil {
		return alterRecord, err
	}

	contentControl, error := RunTemplate(controlTemplate, alterRecord)
	if error != nil {
		return alterRecord, error
	}
	err = MakeControlFile(templateMap, contentControl)

	return alterRecord, err
}
