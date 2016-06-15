// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

type postgresAdapter struct{}

var pgOperators = map[DomainOperator]string{
	OPERATOR_EQUALS:        "= ?",
	OPERATOR_NOT_EQUALS:    "!= ?",
	OPERATOR_LIKE:          "LIKE %?%",
	OPERATOR_NOT_LIKE:      "NOT LIKE %?%",
	OPERATOR_LIKE_PATTERN:  "LIKE ?",
	OPERATOR_ILIKE:         "ILIKE %?%",
	OPERATOR_NOT_ILIKE:     "NOT ILIKE %?%",
	OPERATOR_ILIKE_PATTERN: "ILIKE ?",
	OPERATOR_IN:            "IN (?)",
	OPERATOR_NOT_IN:        "NOT IN (?)",
	OPERATOR_LOWER:         "< ?",
	OPERATOR_LOWER_EQUAL:   "< ?",
	OPERATOR_GREATER:       "> ?",
	OPERATOR_GREATER_EQUAL: ">= ?",
	//OPERATOR_CHILD_OF: "",
}

func (d *postgresAdapter) operatorSQL(do DomainOperator) string {
	return pgOperators[do]
}

var _ dbAdapter = new(postgresAdapter)