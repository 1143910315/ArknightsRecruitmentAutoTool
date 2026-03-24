## ADDED Requirements

### Requirement: Operator cache is stored under the application runtime directory
The system SHALL store operator data JSON and downloaded operator images under a cache directory located relative to the program runtime directory instead of the user cache directory.

#### Scenario: Successful fetch writes cache under runtime directory
- **WHEN** operator data is fetched and cached successfully
- **THEN** the structured operator data file is written under the application runtime directory cache location
- **THEN** downloaded operator images are written under a subdirectory of that runtime cache location

#### Scenario: Cached data is read from runtime directory on page load
- **WHEN** the operator data page is opened and a valid cache exists under the application runtime directory
- **THEN** the application loads operator data from that runtime cache location before attempting any remote request

### Requirement: Existing cache behavior is preserved after moving storage location
The system SHALL preserve local-first loading behavior and wiki source order after moving cache storage to the application runtime directory.

#### Scenario: Runtime directory cache preserves source order
- **WHEN** cached operator data is loaded from the application runtime directory
- **THEN** the operator list is returned in the same order that was parsed from the wiki source
