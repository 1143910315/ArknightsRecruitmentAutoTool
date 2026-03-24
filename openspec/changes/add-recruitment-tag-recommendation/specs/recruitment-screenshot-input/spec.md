## ADDED Requirements

### Requirement: User can provide a recruitment screenshot for analysis
The system SHALL allow the user to start an analysis by supplying an image of the Arknights public recruitment screen through a supported screenshot input method.

#### Scenario: Import screenshot from file
- **WHEN** the user selects a local screenshot file in the desktop application
- **THEN** the system validates that the file can be opened as an image and creates an analysis input from it

#### Scenario: Reject unsupported image input
- **WHEN** the user provides a file that cannot be decoded as an image
- **THEN** the system MUST reject the input and show that analysis cannot continue until a valid screenshot is provided

### Requirement: System validates screenshot readiness before recognition
The system SHALL validate that an analysis input contains enough visible recruitment content to proceed to tag recognition.

#### Scenario: Screenshot is usable
- **WHEN** the provided screenshot contains a readable recruitment tag area
- **THEN** the system marks the screenshot as ready for recognition

#### Scenario: Screenshot is incomplete
- **WHEN** the screenshot does not contain enough recruitment UI content to identify the tag area
- **THEN** the system MUST stop the analysis and ask the user to provide another screenshot
