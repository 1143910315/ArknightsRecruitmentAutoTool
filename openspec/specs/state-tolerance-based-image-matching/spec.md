## ADDED Requirements

### Requirement: Each state image is matched using its own tolerance value
The system SHALL compare a recognition state image against the current screenshot region using that state image's own configured tolerance value.

#### Scenario: Two states in one region use different tolerance values
- **WHEN** the system evaluates two saved state images in the same recognition region
- **THEN** each state image is compared using its own tolerance value
- **THEN** a match or mismatch for one state does not change the tolerance used for the other state

### Requirement: State-level tolerance uses absolute pixel difference per comparison
The system SHALL treat a pixel position as matched only when the absolute per-channel difference between the state image pixel and screenshot pixel is less than that state's configured tolerance value.

#### Scenario: Pixel difference stays within a state's tolerance
- **WHEN** a state image comparison finds that every compared channel difference at a pixel position is less than that state's tolerance value
- **THEN** the system treats that pixel position as matched for that state image

#### Scenario: Pixel difference exceeds a state's tolerance
- **WHEN** a state image comparison finds that any compared channel difference at a pixel position is greater than or equal to that state's tolerance value
- **THEN** the system treats that pixel position as mismatched for that state image

### Requirement: Zero tolerance keeps exact-match behavior for that state image
The system SHALL preserve exact pixel equality behavior for any state image whose configured tolerance value is 0.

#### Scenario: State image with zero tolerance requires exact equality
- **WHEN** the system compares a saved state image whose tolerance value is 0
- **THEN** that state image matches only when every compared pixel is exactly equal to the screenshot region
