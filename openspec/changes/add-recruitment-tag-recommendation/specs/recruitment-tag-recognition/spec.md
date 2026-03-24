## ADDED Requirements

### Requirement: System recognizes standardized recruitment tags from the screenshot
The system SHALL extract the visible recruitment tags from a valid screenshot and normalize them to the application's standard tag definitions before recommendation begins.

#### Scenario: Recognize all visible tags
- **WHEN** a valid recruitment screenshot is analyzed
- **THEN** the system returns the recognized tags as standardized tag values used by the recommendation engine

#### Scenario: Ignore non-tag text
- **WHEN** OCR detects text outside the recruitment tag set
- **THEN** the system MUST exclude that text from the final recognized tag list

### Requirement: System exposes recognition confidence and exceptions
The system SHALL record whether each recognized tag is reliable enough for direct use and MUST surface ambiguous recognition results to the user.

#### Scenario: Low-confidence tag detected
- **WHEN** a recognized tag does not meet the confidence threshold
- **THEN** the system marks that tag as requiring user confirmation before recommendation is finalized

#### Scenario: Unmatched tag text detected
- **WHEN** OCR output cannot be mapped to any supported standard tag
- **THEN** the system MUST flag the result as unresolved instead of silently inventing a tag value

### Requirement: User can correct recognition results before recommendation
The system SHALL allow the user to review and edit the recognized tag list before the recommendation engine calculates the final result set.

#### Scenario: User fixes a misrecognized tag
- **WHEN** the user replaces an incorrect recognized tag with a valid standard tag
- **THEN** the corrected tag list becomes the input for recommendation

#### Scenario: User confirms recognized tags without changes
- **WHEN** the user accepts the recognized tag list as shown
- **THEN** the system proceeds with recommendation using the confirmed tags
