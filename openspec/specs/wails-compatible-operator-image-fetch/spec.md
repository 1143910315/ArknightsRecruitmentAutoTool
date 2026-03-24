## ADDED Requirements

### Requirement: Cached operator images are retrievable through a Wails-compatible request path
The system SHALL provide a backend image retrieval mechanism for cached operator images that does not depend on GET-based asset serving.

#### Scenario: Frontend requests cached operator image through supported call path
- **WHEN** the frontend needs to display a locally cached operator image
- **THEN** it can request the image through a Wails-compatible request or binding flow supported in the current runtime
- **THEN** the backend returns image content or an equivalent renderable payload without requiring direct filesystem access

#### Scenario: Missing cached image does not fail the request flow fatally
- **WHEN** the frontend requests a cached operator image that does not exist locally
- **THEN** the backend returns a clear missing-image result that allows the frontend to fall back gracefully

### Requirement: Runtime directory cache remains the source of local images
The system SHALL continue to load cached operator images from the program runtime directory cache location when serving image content through the new request flow.

#### Scenario: Image fetch uses runtime cache location
- **WHEN** the backend resolves a cached operator image request
- **THEN** it reads the image from the runtime directory cache structure rather than a user directory cache
