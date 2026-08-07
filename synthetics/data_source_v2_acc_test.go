// Copyright 2026 Splunk, Inc.
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
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccUniqueName builds a collision-resistant fixture name so the acceptance
// suite can be re-run and run in parallel against the same live org.
func testAccUniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// testAccUniquePrivateLocationID builds a unique private location id. Location ids
// are validated against `\Aprivate-[a-z\-]*[a-z]\z`, so the usual numeric UnixNano
// suffix is rejected; the timestamp is encoded as lowercase letters instead.
func testAccUniquePrivateLocationID(prefix string) string {
	var suffix strings.Builder
	for n := time.Now().UnixNano(); n > 0; n /= 26 {
		suffix.WriteByte(byte('a' + n%26))
	}
	return fmt.Sprintf("private-%s-%s", prefix, suffix.String())
}

// testAccCheckDataSourceListContains asserts that a collection returned by a list
// data source includes an element whose `field` equals `expected`.
//
// These collections are org-wide, so the fixture's position depends on whatever else
// exists in the account. Matching on any index keeps the assertion about presence
// rather than ordering.
func testAccCheckDataSourceListContains(dataSourceName, collection, field, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source %s not found", dataSourceName)
		}

		if count := rs.Primary.Attributes[collection+".#"]; count == "" || count == "0" {
			return fmt.Errorf("%s returned no %s entries", dataSourceName, collection)
		}

		for index := range testAccDataSourceCollectionIndexes(rs.Primary.Attributes, collection) {
			if rs.Primary.Attributes[fmt.Sprintf("%s.%s.%s", collection, index, field)] == expected {
				return nil
			}
		}

		return fmt.Errorf("%s did not include a %s entry with %s = %q", dataSourceName, collection, field, expected)
	}
}

// testAccDataSourceCollectionIndexes returns the element indexes present for a
// collection in flatmap state, ignoring the `.#` count key and any deeper nesting.
func testAccDataSourceCollectionIndexes(attributes map[string]string, collection string) map[string]struct{} {
	indexes := map[string]struct{}{}
	prefix := collection + "."
	for key := range attributes {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		rest := strings.TrimPrefix(key, prefix)
		if index, _, found := strings.Cut(rest, "."); found {
			indexes[index] = struct{}{}
		}
	}
	return indexes
}

// testAccCheckDataSourceAttrsAbsent asserts that none of the given attribute names
// appear anywhere in a data source's state, at any nesting depth. Used to prove
// metadata-only data sources never surface certificate or private key material.
func testAccCheckDataSourceAttrsAbsent(dataSourceName string, fields ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source %s not found", dataSourceName)
		}

		for key := range rs.Primary.Attributes {
			for _, field := range fields {
				if key == field || strings.HasSuffix(key, "."+field) {
					// Deliberately does not log the value.
					return fmt.Errorf("%s exposed forbidden attribute %s", dataSourceName, key)
				}
			}
		}

		return nil
	}
}

// testAccCheckDataSourceCollectionNotEmpty asserts a list data source returned at
// least one element. Used for account-global reads (devices, locations) where no
// fixture can be created.
func testAccCheckDataSourceCollectionNotEmpty(dataSourceName, collection string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source %s not found", dataSourceName)
		}

		if count := rs.Primary.Attributes[collection+".#"]; count == "" || count == "0" {
			return fmt.Errorf("%s returned no %s entries", dataSourceName, collection)
		}

		return nil
	}
}

// testAccCheckDataSourceElemFieldsNonEmpty asserts that at least one element of a
// collection has a meaningful value for every named field. Used for account-global
// reads where exact values are not ours to control but the flattened shape must
// still be validated.
//
// A numeric zero counts as unset: an unpopulated int renders as "0" in flatmap state,
// and none of the fields this is used on (ids, viewport dimensions) are legitimately 0.
func testAccCheckDataSourceElemFieldsNonEmpty(dataSourceName, collection string, fields ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source %s not found", dataSourceName)
		}

		if count := rs.Primary.Attributes[collection+".#"]; count == "" || count == "0" {
			return fmt.Errorf("%s returned no %s entries", dataSourceName, collection)
		}

		for index := range testAccDataSourceCollectionIndexes(rs.Primary.Attributes, collection) {
			complete := true
			for _, field := range fields {
				value := rs.Primary.Attributes[fmt.Sprintf("%s.%s.%s", collection, index, field)]
				if value == "" || value == "0" {
					complete = false
					break
				}
			}
			if complete {
				return nil
			}
		}

		return fmt.Errorf("%s had no %s element with all of %v set to a non-empty value", dataSourceName, collection, fields)
	}
}

func TestUniquePrivateLocationIDIsValid(t *testing.T) {
	id := testAccUniquePrivateLocationID("acc")
	if !privateLocationIDPattern.MatchString(id) {
		t.Fatalf("generated private location id %q does not satisfy %s", id, privateLocationIDPattern)
	}
}
