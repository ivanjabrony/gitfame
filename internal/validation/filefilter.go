package validation

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/ivan-jabrony/gitfame/internal/config"
)

type languageData struct {
	Name       string `json:"name"`
	LangType   string `json:"type"`
	Extensions []string
}

func ExtensionCheck(files, extensions []string) []string {
	if extensions == nil {
		return files
	}

	checkedFiles := []string{}
	for _, file := range files {
		for _, extension := range extensions {
			if strings.HasSuffix(file, extension) {
				checkedFiles = append(checkedFiles, file)
				break
			}
		}
	}
	return checkedFiles
}

func LanguagesCheck(files, languages []string) ([]string, error) {
	if languages == nil {
		return files, nil
	}

	var parsedLanguageData []languageData
	neededExtensions := []string{}

	err := json.Unmarshal(config.ExtensionJSON, &parsedLanguageData)
	if err != nil {
		return files, err
	}
	for _, data := range parsedLanguageData {
		for _, language := range languages {
			if strings.EqualFold(data.Name, language) {

				extensionList := data.Extensions
				neededExtensions = append(neededExtensions, extensionList...)
			}
		}
	}

	return ExtensionCheck(files, neededExtensions), err
}

func ExcludeCheck(files, patterns []string) ([]string, error) {
	if patterns == nil {
		return files, nil
	}

	checkedFiles := []string{}

	for _, file := range files {
		isMatchedAny := false
		for _, pattern := range patterns {
			isMatched, err := filepath.Match(pattern, file)
			if err != nil {
				return nil, err
			}

			isMatchedAny = isMatchedAny || isMatched
		}
		if !isMatchedAny {
			checkedFiles = append(checkedFiles, file)
		}
	}

	return checkedFiles, nil
}

func RestrictToCheck(files, patterns []string) ([]string, error) {
	if patterns == nil {
		return files, nil
	}

	checkedFiles := []string{}

	for _, file := range files {
		for _, pattern := range patterns {
			isMatched, err := filepath.Match(pattern, file)

			if err != nil {
				return nil, err
			}

			if isMatched {
				checkedFiles = append(checkedFiles, file)
				break
			}
		}
	}

	return checkedFiles, nil
}
