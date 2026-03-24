## ADDED Requirements

### Requirement: System reuses the captured window instance when it is still valid
The system SHALL prefer the original window instance recorded by the recognition template when that instance still exists and can be captured.

#### Scenario: Original window handle is still usable
- **WHEN** a public recruitment recognition run starts and the template's recorded window handle still identifies a live capturable window
- **THEN** the system uses that window instance for recognition without falling back to other same-title same-class candidates

### Requirement: System resolves a unique candidate when multiple windows share title and class name
The system SHALL use additional template instance metadata beyond title and class name to resolve a single target window when duplicate windows exist.

#### Scenario: Duplicate title and class windows are disambiguated by instance metadata
- **WHEN** multiple visible top-level windows share the template's title and class name
- **THEN** the system compares the template's stored instance metadata against those candidates
- **THEN** the system selects the single candidate that uniquely matches the stored instance metadata

### Requirement: System refuses recognition when no unique window instance can be resolved
The system SHALL not continue into image matching when the target window instance cannot be uniquely resolved.

#### Scenario: No matching window instance exists
- **WHEN** a recognition run cannot find any live window that matches the template's recorded identity and instance metadata
- **THEN** the system stops the run before image matching
- **THEN** the system returns a no-window failure result

#### Scenario: Multiple equivalent candidates remain after resolution
- **WHEN** a recognition run still has more than one equivalent target window candidate after applying title, class name, and stored instance metadata
- **THEN** the system stops the run before image matching
- **THEN** the system returns an ambiguous-window failure result
