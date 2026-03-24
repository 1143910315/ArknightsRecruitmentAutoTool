## ADDED Requirements

### Requirement: Region can store multiple tagged state screenshots
The system SHALL allow a single configured recognition region to keep multiple reference screenshots, with each screenshot recorded as a separate state entry bound to one shared recruitment tag.

#### Scenario: User adds multiple states to one region
- **WHEN** the user selects the same configured region and captures more than one reference screenshot for it
- **THEN** the system stores each screenshot as a separate state entry under that region
- **THEN** each state entry stores its selected recruitment tag

### Requirement: Region state management supports review and removal
The system SHALL let the user review existing state entries for a region and remove unwanted state entries before saving the template.

#### Scenario: User reviews region states before saving
- **WHEN** the user opens a configured region in the recognition settings page
- **THEN** the system shows the list of saved state screenshots and their bound recruitment tags for that region
- **THEN** the user can remove an unwanted state entry without deleting the region itself

### Requirement: Template persistence stores region coordinates separately from state screenshots
The system SHALL persist region geometry once per region and persist state screenshot data as nested entries under that region.

#### Scenario: Template is saved with multi-state region data
- **WHEN** the user saves a recognition template containing a region with multiple states
- **THEN** the system stores the region coordinates once for that region
- **THEN** the system stores each tagged screenshot state as a nested entry associated with that region

### Requirement: Existing single-state templates remain readable
The system SHALL keep previously saved single-state region templates readable after the multi-state model is introduced.

#### Scenario: Legacy template is loaded after upgrade
- **WHEN** the system loads a template saved before multi-state region support
- **THEN** the system maps the legacy single reference image into one region state entry
- **THEN** the template remains editable and matchable in the updated recognition settings flow
