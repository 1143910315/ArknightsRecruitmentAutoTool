## ADDED Requirements

### Requirement: Combination results encode operator rarity with image border colors
The system SHALL use the operator image border color to communicate rarity in each combination result card.

#### Scenario: Low-rarity operators use neutral border color
- **WHEN** the system renders a 1-star or 2-star operator in a combination result
- **THEN** the operator image border uses `#dedede`

#### Scenario: Mid-rarity operators use mapped border colors
- **WHEN** the system renders a 3-star operator in a combination result
- **THEN** the operator image border uses `#618bf5`
- **WHEN** the system renders a 4-star operator in a combination result
- **THEN** the operator image border uses `#8960ce`

#### Scenario: High-rarity operators use highlighted border colors
- **WHEN** the system renders a 5-star operator in a combination result
- **THEN** the operator image border uses `#f0a94d`
- **WHEN** the system renders a 6-star operator in a combination result
- **THEN** the operator image border uses `#f0e028`
