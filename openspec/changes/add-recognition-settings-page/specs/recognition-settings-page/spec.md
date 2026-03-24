## ADDED Requirements

### Requirement: User can open a dedicated recognition settings page
The system SHALL provide a dedicated recognition settings page that is reachable from the application navigation.

#### Scenario: User navigates to recognition settings
- **WHEN** the user selects the recognition settings page from the application navigation
- **THEN** the system displays the recognition settings workflow for selecting a target window and configuring recognition regions

### Requirement: User can select a target window by pointing at it
The system SHALL allow the user to choose the target window by pointing the mouse at a window.

#### Scenario: Pointing selects the target window
- **WHEN** the user enters window selection mode and points at a window
- **THEN** the system captures that window as the current recognition target
- **THEN** the page shows the selected window's identifying information

### Requirement: Selected window information is displayed with a screenshot preview
The system SHALL show the selected window's title, class name, and a screenshot preview after a target window is chosen.

#### Scenario: Selected window preview is shown
- **WHEN** a target window has been selected successfully
- **THEN** the system displays the window title and class name
- **THEN** the system displays a screenshot preview of that target window
