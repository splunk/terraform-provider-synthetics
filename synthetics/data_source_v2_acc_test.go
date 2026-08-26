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
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"
)

// A note on indexing TypeSet blocks in these tests, since the two halves differ:
//
//   - In HCL, a TypeSet block cannot be indexed, so a fixture is referenced through the
//     top-level resource id (`tonumber(resource.id)`) rather than `test[0].id`.
//   - In assertions, `test.0.name` is correct. helper/schema stores set elements under
//     hash keys, but acceptance tests do not read that state: the SDK reads Terraform's
//     JSON state and rebuilds the flatmap through shimStateFromJson, whose AddSlice
//     numbers every collection sequentially from 0 (state_shim.go). So set elements
//     arrive as `test.0.*` regardless of their hash.
//
// Ordering within a *multi-element* set is still not meaningful, which is why list reads
// go through testAccCheckDataSourceListContains instead of a fixed index. Singleton reads
// assert `<block>.#` is 1 first, so index 0 is unambiguous.

// testAccCheckFixturesDestroyed asserts that every fixture the test created is gone from
// the API after Terraform destroys it. The SDK already fails a test whose destroy errors,
// but that only proves the provider reported success; this re-reads each id and requires a
// 404, so a Delete that silently no-ops cannot pass.
//
// lookups is keyed by resource type. Each lookup receives the raw state id and returns the
// request details of a get for that id.
func testAccCheckFixturesDestroyed(lookups map[string]func(id string) (*sc2.RequestDetails, error)) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			lookup, ok := lookups[rs.Type]
			if !ok {
				continue
			}

			details, err := lookup(rs.Primary.ID)
			if details != nil && details.StatusCode == http.StatusNotFound {
				continue
			}
			if err != nil {
				return fmt.Errorf("verifying %s %s was destroyed: %w", rs.Type, rs.Primary.ID, err)
			}

			return fmt.Errorf("%s %s still exists after destroy", rs.Type, rs.Primary.ID)
		}

		return nil
	}
}

// testAccIntLookup adapts a client getter keyed by an integer id to the string-keyed
// signature testAccCheckFixturesDestroyed expects.
func testAccIntLookup(get func(id int) (*sc2.RequestDetails, error)) func(string) (*sc2.RequestDetails, error) {
	return func(id string) (*sc2.RequestDetails, error) {
		n, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		return get(n)
	}
}

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

// The client reports a missing object as a 404 in RequestDetails *and* a non-nil error,
// so the status check has to come before the error check. This pins that ordering: a
// deleted fixture must pass, and a fixture the API still returns must fail.
func TestCheckFixturesDestroyed(t *testing.T) {
	state := &terraform.State{
		Modules: []*terraform.ModuleState{{
			Path: []string{"root"},
			Resources: map[string]*terraform.ResourceState{
				"synthetics_create_variable_v2.fixture": {
					Type:     "synthetics_create_variable_v2",
					Primary:  &terraform.InstanceState{ID: "42"},
					Provider: "provider.synthetics",
				},
			},
		}},
	}

	lookups := func(details *sc2.RequestDetails, err error) map[string]func(string) (*sc2.RequestDetails, error) {
		return map[string]func(string) (*sc2.RequestDetails, error){
			"synthetics_create_variable_v2": func(string) (*sc2.RequestDetails, error) {
				return details, err
			},
		}
	}

	t.Run("deleted fixture passes even though the client also returns an error", func(t *testing.T) {
		check := testAccCheckFixturesDestroyed(lookups(&sc2.RequestDetails{StatusCode: http.StatusNotFound}, fmt.Errorf("Status Code: 404 Not Found")))
		if err := check(state); err != nil {
			t.Fatalf("expected a destroyed fixture to pass, got: %s", err)
		}
	})

	t.Run("surviving fixture fails", func(t *testing.T) {
		check := testAccCheckFixturesDestroyed(lookups(&sc2.RequestDetails{StatusCode: http.StatusOK}, nil))
		if err := check(state); err == nil {
			t.Fatal("expected a fixture that still exists to fail the destroy check")
		}
	})

	t.Run("unexpected error fails", func(t *testing.T) {
		check := testAccCheckFixturesDestroyed(lookups(&sc2.RequestDetails{StatusCode: http.StatusInternalServerError}, fmt.Errorf("boom")))
		if err := check(state); err == nil {
			t.Fatal("expected an unexpected lookup error to fail the destroy check")
		}
	})

	t.Run("unmapped resource types are ignored", func(t *testing.T) {
		check := testAccCheckFixturesDestroyed(map[string]func(string) (*sc2.RequestDetails, error){})
		if err := check(state); err != nil {
			t.Fatalf("expected unmapped resource types to be skipped, got: %s", err)
		}
	})
}

func TestUniquePrivateLocationIDIsValid(t *testing.T) {
	id := testAccUniquePrivateLocationID("acc")
	if !privateLocationIDPattern.MatchString(id) {
		t.Fatalf("generated private location id %q does not satisfy %s", id, privateLocationIDPattern)
	}
}
