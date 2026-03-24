## ADDED Requirements

### Requirement: User can define multiple regions on the selected window screenshot
The system SHALL allow the user to define multiple labeled regions on the selected window screenshot.

#### Scenario: User creates multiple regions
- **WHEN** the user draws more than one region on the selected window screenshot
- **THEN** the system keeps each region as a separate configurable entry
- **THEN** each region can be labeled independently

### Requirement: Region positions are saved relative to the window screenshot
The system SHALL save each configured region using coordinates relative to the selected window screenshot.

#### Scenario: Region data stores relative position
- **WHEN** the user saves a configured region
- **THEN** the system stores the region position relative to the selected window screenshot rather than only as absolute screen coordinates

### Requirement: Region reference image and label are persisted together
The system SHALL persist each region's label and reference image together with its position data.

#### Scenario: Region configuration is saved
- **WHEN** the user saves the recognition template
- **THEN** the system stores each region's label
- **THEN** the system stores a reference image for each region
- **THEN** the saved data remains available for future recognition runs
