## ADDED Requirements

### Requirement: Region image matching uses tolerance-based absolute pixel difference
The system SHALL compare each template image pixel against the corresponding screenshot pixel by absolute difference and treat the pixel as matched only when the difference stays below the configured tolerance value.

#### Scenario: Pixel difference stays within tolerance
- **WHEN** a region matching run compares a template pixel and screenshot pixel whose absolute per-channel difference is less than the configured tolerance value
- **THEN** the system treats that pixel position as matched

#### Scenario: Pixel difference exceeds tolerance
- **WHEN** a region matching run compares a template pixel and screenshot pixel whose absolute per-channel difference is greater than or equal to the configured tolerance value
- **THEN** the system treats that pixel position as mismatched

### Requirement: Zero tolerance preserves exact-match behavior
The system SHALL preserve the existing exact image matching behavior when the configured tolerance value is 0.

#### Scenario: Zero tolerance requires exact equality
- **WHEN** a region matching run uses a template whose configured tolerance value is 0
- **THEN** the system reports a match only when every compared pixel is exactly equal to the template image

### Requirement: Tolerance-based matching remains region-scoped and deterministic
The system SHALL keep applying the tolerance rule independently for each configured region and return deterministic match results for the same inputs.

#### Scenario: Same region inputs produce stable tolerance-based result
- **WHEN** the system runs region image matching multiple times with the same template image, screenshot image, and tolerance value
- **THEN** the system returns the same match result each time for that region
