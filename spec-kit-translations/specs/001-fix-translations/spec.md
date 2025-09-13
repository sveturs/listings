# Feature Specification: Complete Translation System Fix

## Problem Statement

Current translation system has only 72.45% coverage instead of the required 95%+. Multiple critical issues:
- 1,357 missing translation keys
- 4,826 untranslated placeholder values showing in production
- Incorrect locale folder (modular-example) being treated as a language
- Russian and Serbian translations showing English placeholders like "marketplace.blackFridayTitle"

## User Scenarios

- As a Russian user, I want to see all interface text in Russian, not English placeholders
- As a Serbian user, I want to see proper Serbian translations in Latin script
- As a developer, I want a clean translation structure without test folders
- As a product owner, I want 100% translation coverage for better user experience

## Requirements

### Functional Requirements

- FR-001: System MUST have 100% translation coverage for ru and sr locales
- FR-002: System MUST NOT show translation keys as text (e.g., "marketplace.blackFridayTitle")
- FR-003: System MUST remove non-locale folders from messages directory
- FR-004: All placeholder values MUST be properly translated
- FR-005: System MUST handle dynamic values correctly (e.g., {count}, {name})

### Non-Functional Requirements

- NFR-001: Translation coverage MUST be at least 95% for production
- NFR-002: All critical modules (marketplace, common, auth) MUST have 100% coverage
- NFR-003: Translation loading MUST NOT exceed 100ms
- NFR-004: JSON files MUST be valid and properly structured

## Acceptance Criteria

### Critical Module Coverage
- Given the marketplace module
- When checking Russian and Serbian translations
- Then all 690+ keys must be properly translated
- And no placeholder text should be visible

### No Visible Translation Keys
- Given any page in the application
- When viewing in Russian or Serbian
- Then no text like "module.key" should be displayed
- And all text should be in the selected language

### Clean Locale Structure
- Given the messages directory
- When listing subdirectories
- Then only valid locale codes (en, ru, sr) should exist
- And no test or example folders should be present

### Dynamic Values
- Given translations with placeholders like {count}
- When rendered with actual values
- Then the placeholders should be replaced correctly
- And the surrounding text should be properly translated

## Out of Scope

- Adding new languages beyond en, ru, sr
- Changing the translation system architecture
- Implementing translation management UI
- Professional translation review (will be done separately)

## Success Metrics

- Translation coverage: >95% (current: 72.45%)
- Missing keys: 0 (current: 1,357)
- Untranslated values: 0 (current: 4,826)
- Visible placeholder texts in UI: 0 (current: hundreds)

## Dependencies

- Access to all English source translations
- Understanding of proper Russian and Serbian terminology
- Knowledge of the application's context and features

## Risks

- Risk of incorrect translations affecting UX
  - Mitigation: Focus on literal translations first, refine later
- Risk of breaking existing functionality
  - Mitigation: Incremental updates with testing
- Risk of JSON structure errors
  - Mitigation: Automated validation after each change

## Notes

The main issue is not missing files but untranslated values within existing files. Many translations are just placeholders copying the key name instead of actual translated text. This needs immediate correction as it severely impacts user experience.