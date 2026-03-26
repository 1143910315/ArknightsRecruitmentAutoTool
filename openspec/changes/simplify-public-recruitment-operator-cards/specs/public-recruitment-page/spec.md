## ADDED Requirements

### Requirement: Public recruitment page shows simplified operator cards in combination results
The system SHALL render each matched operator in the Combination results using a compact card that shows only the operator image and operator name.

#### Scenario: Combination result renders compact operator card
- **WHEN** the public recruitment page displays operators for a matching tag combination
- **THEN** each operator card shows the operator image when a renderable base64 image source is available
- **THEN** each operator card shows the operator name
- **THEN** the card does not show rarity text, profession, origin, or recruitment tag chips

#### Scenario: Combination result keeps card readable without image
- **WHEN** a matching operator does not have a renderable image source
- **THEN** the operator card still shows the operator name
- **THEN** the card remains present in the combination result list
