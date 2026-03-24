## ADDED Requirements

### Requirement: Frontend accesses cached operator images through Wails-served asset paths
The system SHALL expose cached local operator images to the frontend through a Wails-managed asset handling path instead of direct local file access.

#### Scenario: Cached operator record includes Wails-served image path
- **WHEN** operator cache is loaded and an operator image exists locally
- **THEN** the backend returns an image path or URL that is resolvable through Wails asset handling
- **THEN** the frontend can render the cached operator image without direct filesystem access

#### Scenario: Missing local image does not break operator data rendering
- **WHEN** an operator record has no locally cached image or the local image cannot be served
- **THEN** the operator data page still renders the operator record
- **THEN** the frontend falls back to its missing-image presentation without requiring direct file reads
