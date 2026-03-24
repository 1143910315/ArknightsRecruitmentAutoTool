## ADDED Requirements

### Requirement: Public recruitment page exposes an automatic recognition toggle
The system SHALL provide a control on the public recruitment page that lets the user enable or disable template-driven automatic recognition.

#### Scenario: User enables automatic recognition
- **WHEN** the user turns on automatic recognition on the public recruitment page
- **THEN** the system starts the public recruitment recognition loop for the configured template

#### Scenario: User disables automatic recognition
- **WHEN** the user turns off automatic recognition on the public recruitment page
- **THEN** the system stops scheduling further recognition runs

### Requirement: Successful recognition replaces the page's selected recruitment tags
The system SHALL update the public recruitment page's selected tags from the recognized recruitment tags when an automatic recognition run succeeds.

#### Scenario: Recognition success updates selected tags and results
- **WHEN** automatic recognition is enabled and a recognition run returns a successful tag set
- **THEN** the system replaces the page's selected recruitment tags with the recognized tag set
- **THEN** the public recruitment combination results recompute from that updated selected tag set

### Requirement: Unusable recognition runs leave the page selection unchanged
The system SHALL keep the current public recruitment page selection unchanged when a recognition run does not produce a usable tag set.

#### Scenario: Failed recognition does not overwrite current selection
- **WHEN** automatic recognition is enabled and a recognition run returns no recognized tag output
- **THEN** the system keeps the page's existing selected recruitment tags unchanged
- **THEN** the current combination results remain based on the unchanged selection
