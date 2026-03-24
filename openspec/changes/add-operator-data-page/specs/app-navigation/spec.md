## ADDED Requirements

### Requirement: Left navigation menu organizes top-level pages
The system SHALL render a left-side navigation menu on the main application shell so users can switch between top-level pages without mixing unrelated workflows in one view.

#### Scenario: Navigation menu is visible on app launch
- **WHEN** the user opens the application
- **THEN** the main layout shows a left-side navigation menu
- **THEN** the menu lists available top-level pages including `干员数据`

#### Scenario: User switches to operator data page from navigation
- **WHEN** the user clicks the `干员数据` menu item
- **THEN** the application activates the operator data page in the main content area
- **THEN** the selected menu item is visually distinguished from other entries

### Requirement: Active page state remains consistent during in-app actions
The system SHALL keep the current navigation selection stable while the user performs actions inside the active page.

#### Scenario: Fetching operator data does not change active page
- **WHEN** the user starts fetching operator data from the `干员数据` page
- **THEN** the `干员数据` navigation item remains active until the user chooses another page
