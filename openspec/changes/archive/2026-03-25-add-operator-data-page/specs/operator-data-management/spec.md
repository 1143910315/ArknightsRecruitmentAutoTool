## ADDED Requirements

### Requirement: Operator data page shows empty state before data is loaded
The system SHALL provide an `干员数据` page that starts with no operator records and communicates that data must be fetched explicitly by the user.

#### Scenario: Initial empty state is shown
- **WHEN** the user opens the `干员数据` page before any operator data has been loaded in the current session
- **THEN** the page shows that no operator data is available yet
- **THEN** the page shows a `获取干员数据` action

### Requirement: User can fetch operator data from the Arknights wiki page
The system SHALL fetch operator source HTML by HTTP GET from `https://wiki.biligame.com/arknights/公开招募工具` when the user requests operator data.

#### Scenario: Fetch begins from page action
- **WHEN** the user clicks `获取干员数据`
- **THEN** the application sends a request to fetch operator data from the configured wiki URL
- **THEN** the page shows a loading state until the fetch completes or fails

#### Scenario: Fetch failure is surfaced to the user
- **WHEN** the wiki request fails or returns content that cannot be parsed
- **THEN** the application keeps existing operator data unchanged
- **THEN** the page shows a clear error message describing that operator data could not be loaded

### Requirement: Operator data parser extracts all operators from wiki recruit entries
The system SHALL parse every operator entry represented by a `.contentDetail` element on the source page and convert it into normalized operator data records.

#### Scenario: Parser extracts core fields from a recruit entry
- **WHEN** the source HTML contains a `.contentDetail` element with operator metadata and visible name text
- **THEN** the parser extracts the operator name from `.picText`
- **THEN** the parser extracts rarity from `data-param2`
- **THEN** the parser extracts the recruit-related tags displayed inside `.tags .tagText`

#### Scenario: Parser preserves source metadata needed for later recruitment logic
- **WHEN** the source HTML contains `data-param1` with profession, gender, origin, tags, and availability markers
- **THEN** the parser stores the split metadata fields in a structured operator record
- **THEN** the parser preserves whether the operator is available from public recruitment

### Requirement: Operator data page displays fetched records after a successful load
The system SHALL render the parsed operator records on the `干员数据` page after a successful fetch.

#### Scenario: Successful fetch displays operator list
- **WHEN** operator data is fetched and parsed successfully
- **THEN** the page displays all parsed operator records
- **THEN** each displayed record includes at least the operator name, rarity, and parsed tags
- **THEN** the empty state is replaced by the loaded data view
