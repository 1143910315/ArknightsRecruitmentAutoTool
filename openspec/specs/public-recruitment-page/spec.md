## ADDED Requirements

### Requirement: User can access a dedicated public recruitment page
The system SHALL provide a dedicated public recruitment page that is reachable from the application's navigation.

#### Scenario: User opens the public recruitment page
- **WHEN** the user selects the public recruitment page from the application navigation
- **THEN** the system displays a page dedicated to public recruitment tag selection and operator filtering

### Requirement: User can choose from the supported recruitment tags
The system SHALL present the supported recruitment tags grouped by category, including profession, deployment position, trait tags, and operator seniority.

#### Scenario: Supported tags are displayed by category
- **WHEN** the public recruitment page is shown
- **THEN** the system displays the supported tags grouped into profession, deployment position, trait tags, and operator seniority sections

### Requirement: User cannot select more than five tags
The system SHALL limit the number of simultaneously selected recruitment tags to five.

#### Scenario: User reaches the selection limit
- **WHEN** the user already has five selected tags and attempts to select another unselected tag
- **THEN** the system keeps the existing five selected tags unchanged
- **THEN** the system shows clear feedback that at most five tags can be selected

#### Scenario: User deselects a chosen tag
- **WHEN** the user deselects one of the currently selected tags
- **THEN** the system removes that tag from the selected set
- **THEN** the user can select a different tag afterward

### Requirement: Page handles missing operator data gracefully
The system SHALL keep the public recruitment page usable even when no local operator data is available.

#### Scenario: No cached operator data is available
- **WHEN** the public recruitment page is opened before operator data has been loaded locally
- **THEN** the system shows an empty or guidance state instead of failing
- **THEN** the system explains that operator data must be loaded before recruitment results can be shown
