## ADDED Requirements

### Requirement: Frontend renders cached images from backend-provided image payloads
The system SHALL render cached operator images in the frontend using the new Wails-compatible image retrieval result instead of direct GET asset URLs.

#### Scenario: Cached operator image is rendered in operator list
- **WHEN** cached operator data includes a locally available image reference
- **THEN** the frontend resolves that image through the new backend-supported flow
- **THEN** the operator list displays the image for the corresponding operator

#### Scenario: Frontend falls back when no image payload is available
- **WHEN** the frontend cannot obtain a renderable image payload for a cached operator image
- **THEN** the operator record still renders in the list
- **THEN** the missing-image fallback UI is shown for that record
