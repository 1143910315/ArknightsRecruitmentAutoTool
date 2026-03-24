## ADDED Requirements

### Requirement: Recognition settings page exposes an image matching tolerance field
The system SHALL allow the user to configure an image matching tolerance value in the recognition settings workflow.

#### Scenario: User edits template image tolerance
- **WHEN** the user configures a recognition template in the recognition settings page
- **THEN** the system shows an image matching tolerance input together with the template settings
- **THEN** the user can set the tolerance value before saving the template

### Requirement: Recognition template persistence stores the configured tolerance value
The system SHALL persist the configured image matching tolerance value as part of the saved recognition template.

#### Scenario: Saved template keeps tolerance value
- **WHEN** the user saves a recognition template with a configured image matching tolerance value
- **THEN** the system stores that tolerance value with the template metadata
- **THEN** loading the template later restores the same tolerance value in the recognition settings page

### Requirement: Recognition settings rejects invalid tolerance values
The system SHALL prevent invalid image matching tolerance values from being saved.

#### Scenario: Negative tolerance is rejected
- **WHEN** the user attempts to save a recognition template with a negative tolerance value
- **THEN** the system rejects the save request
- **THEN** the user is informed that the tolerance value must be a valid non-negative number

#### Scenario: Legacy template without tolerance is loaded
- **WHEN** the system loads a recognition template saved before tolerance support existed
- **THEN** the system assigns that template a default tolerance value of 0
- **THEN** the template remains editable and matchable in the updated recognition settings flow
