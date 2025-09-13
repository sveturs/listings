# Technical Implementation Plan: Translation System Fixes

## Architecture Overview

### Current State
- Modular translation system using next-intl
- JSON-based translations in `frontend/svetu/src/messages/{locale}/{module}.json`
- 3 supported locales: en, ru, sr
- ~78% translation coverage (242 missing keys)

### Target State
- >95% translation coverage across all locales
- Automated validation in CI/CD
- Zero MISSING_MESSAGE errors in production
- Consistent translations across all modules

## Technology Stack

### Existing (No Changes)
- **Framework**: Next.js 15.3.2 with App Router
- **i18n Library**: next-intl
- **File Format**: JSON
- **Structure**: Modular (separate files per module)

### New Tools
- **Validation**: Custom Python scripts for analysis
- **CI/CD**: GitHub Actions for translation validation
- **Monitoring**: Sentry integration for missing keys tracking

## Implementation Phases

### Phase 1: Critical Module Fixes (3 days)
**Priority**: HIGH
**Modules**: marketplace, common, admin

1. Add missing translations for marketplace module (45+ keys)
2. Add missing translations for common module (35+ keys)
3. Add missing translations for admin module (30+ keys)
4. Validate JSON structure consistency

**Acceptance**: No MISSING_MESSAGE errors in these modules

### Phase 2: Secondary Module Fixes (3 days)
**Priority**: MEDIUM
**Modules**: cars, auth, map, storefront, misc

1. Add remaining 130+ translations
2. Cross-validate between languages
3. Fix any structural inconsistencies
4. Test with real user scenarios

**Acceptance**: >95% translation coverage achieved

### Phase 3: Quality Assurance (2 days)
**Priority**: HIGH

1. Manual testing of all major user flows
2. Automated testing setup
3. Performance validation
4. User acceptance testing with native speakers

**Acceptance**: All tests passing, native speakers approve

### Phase 4: Automation Setup (2 days)
**Priority**: MEDIUM

1. Create GitHub Actions workflow for validation
2. Setup pre-commit hooks
3. Configure Sentry for production monitoring
4. Document translation guidelines

**Acceptance**: CI/CD pipeline validates all PRs

### Phase 5: Documentation & Training (1 day)
**Priority**: LOW

1. Update developer documentation
2. Create translation guide for contributors
3. Setup translation management process
4. Knowledge transfer session

**Acceptance**: Team can maintain translations independently

## Implementation Details

### File Structure
```
frontend/svetu/src/messages/
├── en/
│   ├── marketplace.json
│   ├── common.json
│   ├── admin.json
│   └── ...
├── ru/
│   ├── marketplace.json (needs 45+ keys)
│   ├── common.json (needs 35+ keys)
│   ├── admin.json (needs 30+ keys)
│   └── ...
└── sr/
    ├── marketplace.json (needs 45+ keys)
    ├── common.json (needs 35+ keys)
    ├── admin.json (needs 30+ keys)
    └── ...
```

### Translation Key Format
```json
{
  "module": {
    "section": {
      "key": "Translation text",
      "nested": {
        "key": "Nested translation"
      }
    }
  }
}
```

### Validation Script Structure
```python
# analyze_translations.py
- Load all translation files
- Compare keys between locales
- Identify missing translations
- Generate report
- Exit with error code if coverage < 95%
```

## Risk Mitigation

### Risk 1: Breaking Existing Functionality
- **Mitigation**: Incremental changes with testing after each module
- **Validation**: Run full E2E tests after each phase

### Risk 2: Poor Translation Quality
- **Mitigation**: Use professional translation services or native speakers
- **Validation**: Review by native speakers before deployment

### Risk 3: Performance Impact
- **Mitigation**: Monitor bundle size and loading time
- **Validation**: Performance tests before/after changes

## Success Metrics

1. **Coverage**: >95% translation keys present in all locales
2. **Errors**: 0 MISSING_MESSAGE errors in production
3. **Performance**: <100ms translation loading time
4. **Quality**: Native speaker approval rating >4/5

## Timeline

- **Total Duration**: 11 working days
- **Start Date**: Immediate
- **Critical Modules Complete**: Day 3
- **Full Coverage**: Day 6
- **Production Ready**: Day 11

## Dependencies

- Access to translation resources/native speakers
- GitHub Actions minutes for CI/CD
- Sentry account for monitoring
- Team availability for review and testing

## Notes

Based on the audit report, the system architecture is sound. The main effort is content addition rather than structural changes. Priority should be given to user-facing modules (marketplace, common) over admin modules.