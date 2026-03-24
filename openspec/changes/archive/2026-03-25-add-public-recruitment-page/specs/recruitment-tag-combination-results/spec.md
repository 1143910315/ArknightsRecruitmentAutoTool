## ADDED Requirements

### Requirement: System generates all valid non-empty tag combinations from the selected tags
The system SHALL generate every non-empty tag combination from the user's selected tags, up to the configured selection limit.

#### Scenario: Two selected tags produce three combinations
- **WHEN** the user selects the tags `输出` and `近卫`
- **THEN** the system generates the combinations `输出 + 近卫`, `输出`, and `近卫`

#### Scenario: Five selected tags stay within supported combination range
- **WHEN** the user selects five tags
- **THEN** the system generates all non-empty combinations of those five tags
- **THEN** the total number of combinations does not exceed 31

### Requirement: System shows only combinations that match at least one operator
The system SHALL display only the generated tag combinations that match one or more operators.

#### Scenario: Combination with no matching operator is hidden
- **WHEN** a generated tag combination matches no operators in the loaded operator dataset
- **THEN** the system does not display that combination in the results

### Requirement: System matches operators by combination coverage
The system SHALL include an operator in a combination result only when the operator satisfies every tag in that combination.

#### Scenario: Operator appears in full and partial matches it satisfies
- **WHEN** an operator satisfies both `输出` and `近卫`
- **THEN** that operator appears in the `输出 + 近卫` result
- **THEN** that operator also appears in the `输出` result and the `近卫` result

#### Scenario: Operator excluded from non-covered combination
- **WHEN** an operator satisfies `输出` but does not satisfy `近卫`
- **THEN** that operator appears in the `输出` result only
- **THEN** that operator does not appear in the `输出 + 近卫` result

### Requirement: Result ordering remains stable and readable
The system SHALL present combination groups and matched operators in a stable order that preserves usability.

#### Scenario: Combination groups prioritize more specific matches
- **WHEN** the system renders multiple valid tag combinations
- **THEN** combinations with more selected tags are shown before combinations with fewer selected tags

#### Scenario: Operators keep source order inside a combination group
- **WHEN** the system renders operators for a matching tag combination
- **THEN** the operators appear in the same relative order as the loaded operator dataset
