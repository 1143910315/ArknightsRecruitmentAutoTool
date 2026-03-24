## ADDED Requirements

### Requirement: Main navigation uses an Element Plus drawer
The system SHALL replace the fixed left-side menu with an Element Plus Drawer that can be opened from the main page shell and used to switch between top-level pages.

#### Scenario: User opens navigation drawer
- **WHEN** the user clicks the navigation trigger on the home shell
- **THEN** the application opens an Element Plus Drawer containing the available top-level pages
- **THEN** the current active page is visibly identified inside the drawer

#### Scenario: Navigation selection updates current page
- **WHEN** the user selects a page item from the drawer
- **THEN** the application activates that page in the main content area
- **THEN** the drawer may close after selection without changing the chosen page state

### Requirement: Application shell uses a light blue and white visual theme
The system SHALL present the main shell and operator data page using a light blue and white color direction instead of the current warm palette.

#### Scenario: Updated theme is visible on app load
- **WHEN** the application is opened after this change
- **THEN** the page background, panels, and navigation styling use a light blue and white visual palette
- **THEN** the updated theme is applied consistently across the main shell and operator data page
