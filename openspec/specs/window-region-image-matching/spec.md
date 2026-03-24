## ADDED Requirements

### Requirement: Saved region templates can be used as future image-matching inputs
The system SHALL keep enough saved data to evaluate whether a target window's configured region matches the previously selected reference image.

#### Scenario: Saved template can be reused for matching
- **WHEN** a recognition run uses a previously saved window template
- **THEN** the system can load the saved region position and reference image for matching

### Requirement: Matching is evaluated per configured region
The system SHALL evaluate image matching independently for each configured region.

#### Scenario: Region match is determined independently
- **WHEN** a recognition run checks multiple configured regions
- **THEN** the system determines the match result for each region independently of the others

### Requirement: Matching result distinguishes match and mismatch outcomes
The system SHALL return a clear result indicating whether the current region image matches the saved reference image.

#### Scenario: Matching result reports miss
- **WHEN** the current region image does not match the saved reference image
- **THEN** the system returns a mismatch result for that region

#### Scenario: Matching result reports hit
- **WHEN** the current region image matches the saved reference image
- **THEN** the system returns a match result for that region
