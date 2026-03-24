## ADDED Requirements

### Requirement: System derives recruitment tags only from complete unique region matches
The system SHALL produce a public recruitment recognition result only when every configured region in the selected recognition template matches exactly one saved state screenshot.

#### Scenario: Every region uniquely matches one state
- **WHEN** a public recruitment recognition run evaluates a template and each configured region matches exactly one saved tagged state screenshot
- **THEN** the system marks the recognition run as successful
- **THEN** the system returns the set of recruitment tags bound to those uniquely matched states

### Requirement: System performs no state-changing output on incomplete or ambiguous region matches
The system SHALL treat the entire recognition run as unusable when any configured region has zero matches or more than one matched state screenshot.

#### Scenario: One region has no matched state
- **WHEN** a recognition run evaluates a template and any configured region matches no saved state screenshot
- **THEN** the system marks the recognition run as failed
- **THEN** the system returns no recognized recruitment tag output for that run

#### Scenario: One region matches multiple states
- **WHEN** a recognition run evaluates a template and any configured region matches more than one saved state screenshot
- **THEN** the system marks the recognition run as failed
- **THEN** the system returns no recognized recruitment tag output for that run

### Requirement: Recognition loop waits 500ms after each completed run before the next run starts
The system SHALL schedule the next public recruitment recognition run no earlier than 500ms after the current run finishes, regardless of success or failure.

#### Scenario: Successful run waits before next retry
- **WHEN** a public recruitment recognition run completes successfully while automatic recognition remains enabled
- **THEN** the system waits at least 500ms after that completion time before starting the next run

#### Scenario: Failed run still waits before next retry
- **WHEN** a public recruitment recognition run completes with a failure result while automatic recognition remains enabled
- **THEN** the system waits at least 500ms after that completion time before starting the next run
