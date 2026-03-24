## ADDED Requirements

### Requirement: Parsed operator data and images are persisted locally
The system SHALL save parsed operator records and their related images to local storage after a successful fetch from the wiki source.

#### Scenario: Successful fetch writes cache data
- **WHEN** operator data is fetched and parsed successfully from the wiki page
- **THEN** the application stores the structured operator data in a local cache file
- **THEN** the application stores each operator image locally so the cached data can reference local image assets

### Requirement: Operator data page prefers local cached data on load
The system SHALL attempt to load locally cached operator data before requesting remote data when the operator data page is opened.

#### Scenario: Cached data exists when page opens
- **WHEN** the user opens the operator data page and a valid local cache is present
- **THEN** the page displays the cached operator data immediately
- **THEN** the page does not require a remote fetch to show the initial data view

#### Scenario: Cached data is missing when page opens
- **WHEN** the user opens the operator data page and no valid local cache exists
- **THEN** the page falls back to the existing empty state until the user requests data

### Requirement: Display order matches the wiki source order
The system SHALL preserve the original order of `.contentDetail` entries from `https://wiki.biligame.com/arknights/公开招募工具` through parsing, local persistence, and page rendering.

#### Scenario: Parsed data keeps source sequence
- **WHEN** the parser processes multiple `.contentDetail` elements from the wiki HTML
- **THEN** the generated operator records keep the same relative order as the source page

#### Scenario: Cached data renders in preserved order
- **WHEN** the operator data page loads cached records
- **THEN** the page renders the operator list in the same sequence that was parsed from the wiki source
- **THEN** the application does not reorder the list by rarity, name, or other derived fields
