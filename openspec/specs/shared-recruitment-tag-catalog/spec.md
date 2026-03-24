## ADDED Requirements

### Requirement: System exposes a shared recruitment tag catalog
The system SHALL provide a shared catalog of public recruitment tags that can be reused across recognition settings and other recruitment-related features.

#### Scenario: Recognition settings loads shared tag options
- **WHEN** the recognition settings flow needs tag options for a region state
- **THEN** the system returns the shared recruitment tag catalog instead of a page-local hardcoded list

### Requirement: Shared recruitment tag catalog preserves standard tag groupings
The system SHALL organize the shared recruitment tags into the standard groups of profession, deployment position, trait tag, and operator seniority.

#### Scenario: Tag groups are returned with their allowed values
- **WHEN** a caller requests the recruitment tag catalog
- **THEN** the system returns the standard groups with these values
- **THEN** profession contains 近卫、狙击、重装、医疗、辅助、术师、特种、先锋
- **THEN** deployment position contains 近战位、远程位
- **THEN** trait tag contains 控场、爆发、治疗、支援、费用回复、输出、生存、群攻、防护、减速、削弱、快速复活、位移、召唤、支援机械、元素
- **THEN** operator seniority contains 新手、资深干员、高级资深干员
