// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package synthetics

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

const (
	excludedFileTypeCustom    = "custom"
	excludedFileTypeAllExcept = "all_except"
)

func browserCheckV2ExcludedFilesSchema(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: !computed,
		Computed: computed,
		Elem:     browserCheckV2ExcludedFileResource(computed),
	}
}

func browserCheckV2ExcludedFileResource(computed bool) *schema.Resource {
	typeSchema := &schema.Schema{
		Type: schema.TypeString,
	}
	regexSchema := &schema.Schema{
		Type: schema.TypeString,
	}

	if computed {
		typeSchema.Computed = true
		regexSchema.Computed = true
	} else {
		typeSchema.Required = true
		typeSchema.ValidateFunc = validation.StringIsNotEmpty
		regexSchema.Optional = true
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type":  typeSchema,
			"regex": regexSchema,
		},
	}
}

func buildExcludedFilesV2Data(excludedFiles *schema.Set) ([]sc2.ExcludedFile, error) {
	if excludedFiles == nil {
		return make([]sc2.ExcludedFile, 0), nil
	}

	files := excludedFiles.List()
	excludedFilesList := make([]sc2.ExcludedFile, 0, len(files))

	for i, excludedFile := range files {
		data := excludedFile.(map[string]interface{})
		file := sc2.ExcludedFile{
			Type: strings.TrimSpace(data["type"].(string)),
		}

		if regex, ok := data["regex"].(string); ok {
			file.Regex = regex
		}

		if err := validateExcludedFileV2(i, file); err != nil {
			return nil, err
		}

		excludedFilesList = append(excludedFilesList, file)
	}

	return excludedFilesList, nil
}

func validateExcludedFileV2(index int, excludedFile sc2.ExcludedFile) error {
	if excludedFile.Type == "" {
		return fmt.Errorf("excluded_files[%d].type must not be empty", index)
	}

	switch excludedFile.Type {
	case excludedFileTypeCustom, excludedFileTypeAllExcept:
		if strings.TrimSpace(excludedFile.Regex) == "" {
			return fmt.Errorf("excluded_files[%d].regex must be set when type is %s", index, excludedFile.Type)
		}
		if _, err := regexp.Compile(excludedFile.Regex); err != nil {
			return fmt.Errorf("excluded_files[%d].regex is not valid RE2 syntax: %w", index, err)
		}
	default:
		if excludedFile.Regex != "" {
			return fmt.Errorf("excluded_files[%d].regex is only supported when type is custom or all_except", index)
		}
	}

	return nil
}

func flattenExcludedFilesV2Data(excludedFiles []sc2.ExcludedFile) []interface{} {
	if len(excludedFiles) == 0 {
		return make([]interface{}, 0)
	}

	flattened := make([]interface{}, len(excludedFiles))
	for i, excludedFile := range excludedFiles {
		file := map[string]interface{}{
			"type": excludedFile.Type,
		}
		if excludedFile.Regex != "" {
			file["regex"] = excludedFile.Regex
		}
		flattened[i] = file
	}

	return flattened
}
