## ADDED Requirements

### Requirement: Window screenshot is produced from a screen snapshot and window-rectangle crop
The system SHALL capture a full screen snapshot first and then crop the target window image using the target window's screen position and size.

#### Scenario: Visible window is cropped from screen snapshot
- **WHEN** the system needs a screenshot for a target window
- **THEN** it first captures the current screen image
- **THEN** it crops the target window region from that screen image using the window rectangle

### Requirement: Screen-snapshot window capture keeps the existing screenshot call path usable
The system SHALL keep the existing screenshot call path usable for upper-layer recognition flows while changing only the capture source strategy.

#### Scenario: Recognition settings still use the same screenshot entry
- **WHEN** the recognition settings flow requests a target window screenshot
- **THEN** the system returns a screenshot result through the existing screenshot-oriented call path
- **THEN** the caller does not need a separate screenshot workflow for the new capture mode

### Requirement: Screenshot crop handles out-of-bounds window rectangles safely
The system SHALL constrain the crop rectangle to the available screen snapshot bounds before producing the final window image.

#### Scenario: Window rectangle partially exceeds screen bounds
- **WHEN** the target window rectangle extends outside the captured screen image bounds
- **THEN** the system crops only the intersecting visible region
- **THEN** the system avoids invalid out-of-bounds image reads

#### Scenario: Window rectangle has no visible intersection
- **WHEN** the target window rectangle does not intersect the available screen snapshot area
- **THEN** the system returns a clear capture failure instead of an invalid image
