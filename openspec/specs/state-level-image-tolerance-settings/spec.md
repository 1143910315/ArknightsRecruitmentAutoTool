## ADDED Requirements

### Requirement: Each recognition state image can store its own tolerance value
The system SHALL allow every saved recognition state image to keep an independent image matching tolerance value.

#### Scenario: User configures different tolerance values for two states in one region
- **WHEN** the user edits two saved state images under the same recognition region
- **THEN** the system lets the user set a separate tolerance value for each state image
- **THEN** saving the template preserves both state-specific tolerance values independently

### Requirement: Recognition settings page exposes tolerance input per state image
The system SHALL show the image matching tolerance input as part of each state image editor in the recognition settings page.

#### Scenario: State editor shows tolerance input
- **WHEN** the user views a recognition region state entry in the recognition settings page
- **THEN** the system displays that state's current tolerance value together with its tag and preview image

### Requirement: Legacy state images default to zero tolerance when no value exists
The system SHALL treat state images saved before tolerance support as having a default tolerance value of 0.

#### Scenario: Legacy template state loads without explicit tolerance
- **WHEN** the system loads a previously saved recognition template whose state images do not include a tolerance value
- **THEN** each such state image is assigned a tolerance value of 0
- **THEN** the template remains editable and saveable in the updated recognition settings flow

### Requirement: Invalid state tolerance values are rejected
The system SHALL reject invalid tolerance values for individual state images.

#### Scenario: Negative state tolerance is rejected
- **WHEN** the user attempts to save a state image with a negative tolerance value
- **THEN** the system rejects the save request
- **THEN** the user is informed that the tolerance value must be a valid non-negative number
