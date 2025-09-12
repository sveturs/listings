# Technical Implementation Plan: Complete Translation System Fix

## Current State Analysis

### Issues Identified
- **Coverage**: 72.45% (target: >95%)
- **Missing Keys**: 1,357 across all modules
- **Untranslated Values**: 4,826 placeholders showing as text
- **Invalid Locale**: "modular-example" folder treated as locale
- **Critical Modules Affected**: marketplace, common, auth, storefront

### Root Causes
1. Previous "fixes" only added keys but not actual translations
2. Many values are just the key name repeated (e.g., "blackFridayTitle": "marketplace.blackFridayTitle")
3. Test/example folders mixed with production locales
4. No validation process to catch untranslated placeholders

## Technology Stack

- **Framework**: Next.js 15.3.2
- **i18n Library**: next-intl
- **File Structure**: JSON modules per locale
- **Languages**: en (source), ru, sr (targets)

## Implementation Phases

### Phase 1: Cleanup and Structure Fix (Day 1 - 2 hours)
1. Remove "modular-example" folder from messages directory
2. Audit all JSON files for structural issues
3. Create backup of current state
4. Set up validation scripts

**Deliverables**:
- Clean locale structure with only en, ru, sr
- Backup in `/backups/translations-{timestamp}/`
- Validation script ready

### Phase 2: Critical Module Translations (Day 1-2 - 8 hours)
Fix modules that directly impact user experience:

1. **marketplace module** (690+ keys)
   - Fix all placeholder values in ru and sr
   - Translate Black Friday campaign texts
   - Fix category and product-related texts

2. **common module** (590+ keys)
   - Fix navigation and footer texts
   - Translate error messages
   - Fix form labels and buttons

3. **auth module** (150+ keys)
   - Fix login/registration texts
   - Translate validation messages
   - Fix security-related texts

**Deliverables**:
- 100% real translations for critical modules
- No placeholders visible in UI

### Phase 3: Secondary Module Translations (Day 2-3 - 6 hours)
Fix remaining high-usage modules:

1. **storefront** - Store management texts
2. **cart** - Shopping cart and checkout
3. **orders** - Order management
4. **profile** - User profile texts
5. **map** - Map interface texts

**Deliverables**:
- All user-facing modules fully translated
- Coverage increased to >85%

### Phase 4: Complete Coverage (Day 3 - 4 hours)
1. Fix all remaining modules
2. Handle dynamic placeholders correctly
3. Ensure consistency across translations

**Deliverables**:
- >95% translation coverage
- All placeholders properly handled

### Phase 5: Validation and Testing (Day 3-4 - 4 hours)
1. Run automated validation scripts
2. Manual UI testing in all languages
3. Fix any remaining issues
4. Performance testing

**Deliverables**:
- Zero untranslated placeholders
- All tests passing
- Performance benchmarks met

## Implementation Strategy

### Translation Approach
1. Use English text as reference
2. For Russian: Professional business terminology
3. For Serbian: Latin script, regional terminology
4. Keep technical terms consistent across modules

### File Processing Method
```python
for locale in ['ru', 'sr']:
    for module in all_modules:
        1. Load en/{module}.json as reference
        2. Load {locale}/{module}.json
        3. For each key in English:
           - If value in locale is placeholder: translate
           - If value is missing: add translation
           - If value is English: translate
        4. Validate JSON structure
        5. Save with proper formatting
```

### Quality Checks
- No key names as values
- No English text in non-English locales
- Proper handling of {variables}
- Consistent terminology
- Valid JSON syntax

## Risk Mitigation

1. **Backup Strategy**
   - Full backup before starting
   - Incremental backups after each phase
   - Git commits after each module

2. **Testing Strategy**
   - Automated validation after each file
   - UI spot checks after each module
   - Full regression test at end

3. **Rollback Plan**
   - Keep original files in backup
   - Can restore any module independently
   - Git history for granular rollback

## Success Metrics

| Metric | Current | Target | 
|--------|---------|--------|
| Coverage | 72.45% | >95% |
| Missing Keys | 1,357 | 0 |
| Untranslated Values | 4,826 | 0 |
| Visible Placeholders | Hundreds | 0 |
| Load Time | Unknown | <100ms |

## Timeline

- **Total Duration**: 3-4 days
- **Day 1**: Cleanup + Critical modules (marketplace, common)
- **Day 2**: Critical (auth) + Secondary modules  
- **Day 3**: Complete coverage + Validation
- **Day 4**: Testing + Fixes + Documentation

## Tools Required

- Python for automation scripts
- JSON validator
- Translation audit script
- UI testing via browser
- Git for version control

## Notes

The key insight is that the previous "fix" added translation keys but filled them with placeholders instead of actual translations. This implementation will replace ALL placeholder values with proper translations, achieving true 95%+ coverage.