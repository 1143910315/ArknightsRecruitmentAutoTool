## ADDED Requirements

### Requirement: System evaluates valid recruitment tag combinations
The system SHALL compute the valid single-tag, two-tag, and three-tag combinations from the confirmed tag list and match them against the local recruitment data set.

#### Scenario: Generate recommendation candidates
- **WHEN** the user confirms a set of recognized recruitment tags
- **THEN** the system returns every supported tag combination that can be evaluated by the recruitment rules

#### Scenario: Unsupported combination produces no candidate
- **WHEN** a tag combination does not match any recruitable operator set in the local data
- **THEN** the system MUST omit that combination from the recommended result list

### Requirement: System explains recommendation outcomes
The system SHALL show the user why a tag combination is valuable, including candidate operators and the combination's rarity or exclusivity implications.

#### Scenario: High-value combination found
- **WHEN** a combination guarantees a high-rarity or special-value outcome
- **THEN** the system highlights the combination and states the reason for the recommendation

#### Scenario: Multiple combinations remain viable
- **WHEN** more than one valid combination exists
- **THEN** the system shows each combination with its candidate operators and value explanation

### Requirement: System orders results by recommendation priority
The system SHALL sort recommendation results by the value and certainty of the recruitment outcome rather than by arbitrary display order.

#### Scenario: Guaranteed high-rarity result outranks broader result
- **WHEN** one combination guarantees a better rarity outcome than another combination
- **THEN** the higher-value combination appears earlier in the recommendation list

#### Scenario: Similar rarity results are differentiated by certainty
- **WHEN** two combinations lead to the same top rarity band
- **THEN** the system ranks the combination with the narrower or more distinctive operator pool higher
