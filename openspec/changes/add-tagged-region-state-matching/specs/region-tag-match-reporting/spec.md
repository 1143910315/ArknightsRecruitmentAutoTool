## ADDED Requirements

### Requirement: Region matching reports matched tag states
The system SHALL report which tagged state screenshots matched successfully for each configured region during recognition matching.

#### Scenario: Region returns matched tag states
- **WHEN** a recognition run evaluates a configured region that contains multiple tagged state screenshots
- **THEN** the system compares the current cropped region image against each saved state screenshot
- **THEN** the match result includes the recruitment tag of every state that matched successfully for that region

### Requirement: Region matching preserves region-level context when no state matches
The system SHALL return a region-level result even when none of the saved state screenshots match.

#### Scenario: Region returns no matched tag state
- **WHEN** a recognition run evaluates a configured region and no saved state screenshot matches
- **THEN** the system returns that region in the result set
- **THEN** the system marks that region as having no matched tag state

### Requirement: Recognition settings can display tag-oriented match feedback
The system SHALL expose enough match result detail for the recognition settings page to show which tag states matched in each region.

#### Scenario: UI receives tag-oriented match detail
- **WHEN** the recognition settings page requests a template match result
- **THEN** the system returns region results together with the matched state identifiers and their recruitment tags
- **THEN** the page can render which tag image matched successfully for each region
