# Development Tasks: Translation System Fixes

## Phase 1: Critical Module Fixes

### TASK-001: Analyze and extract missing marketplace translations
```bash
cd /data/hostel-booking-system/spec-kit/translation-audit
python3 analyze_translations.py --module marketplace --output missing_marketplace.json
```
**Acceptance**: JSON file with all missing marketplace keys generated

### TASK-002: Add marketplace translations for Russian
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages/ru
# Edit marketplace.json and add the missing translations from the analysis
```
**Acceptance**: All marketplace keys present in ru/marketplace.json

### TASK-003: Add marketplace translations for Serbian
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages/sr
# Edit marketplace.json and add the missing translations from the analysis
```
**Acceptance**: All marketplace keys present in sr/marketplace.json

### TASK-004: Add common module translations for Russian
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages/ru
# Edit common.json and add the missing translations
```
**Acceptance**: All common keys present in ru/common.json

### TASK-005: Add common module translations for Serbian
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages/sr
# Edit common.json and add the missing translations
```
**Acceptance**: All common keys present in sr/common.json

### TASK-006: Add admin module translations for Russian
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages/ru
# Edit admin.json and add the missing translations
```
**Acceptance**: All admin keys present in ru/admin.json

### TASK-007: Add admin module translations for Serbian
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages/sr
# Edit admin.json and add the missing translations
```
**Acceptance**: All admin keys present in sr/admin.json

### TASK-008: Validate Phase 1 translations
```bash
cd /data/hostel-booking-system/spec-kit/translation-audit
python3 analyze_translations.py --modules marketplace,common,admin --validate
```
**Acceptance**: No missing keys reported for these modules

## Phase 2: Secondary Module Fixes

### TASK-009: Add cars module translations
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages
# Add missing translations to ru/cars.json and sr/cars.json
```
**Acceptance**: All cars module keys translated

### TASK-010: Add auth module translations
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages
# Add missing translations to ru/auth.json and sr/auth.json
```
**Acceptance**: All auth module keys translated

### TASK-011: Add map module translations
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages
# Add missing translations to ru/map.json and sr/map.json
```
**Acceptance**: All map module keys translated

### TASK-012: Add storefront module translations
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages
# Add missing translations to ru/storefront.json and sr/storefront.json
```
**Acceptance**: All storefront module keys translated

### TASK-013: Add misc module translations
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages
# Add missing translations to ru/misc.json and sr/misc.json
```
**Acceptance**: All misc module keys translated

### TASK-014: Full system validation
```bash
cd /data/hostel-booking-system/spec-kit/translation-audit
python3 analyze_translations.py --full-validation
```
**Acceptance**: >95% translation coverage achieved

## Phase 3: Quality Assurance

### TASK-015: Start development server and test
```bash
cd /data/hostel-booking-system/frontend/svetu
yarn dev -p 3001
# Test all major user flows in ru and sr languages
```
**Acceptance**: No MISSING_MESSAGE errors in console

### TASK-016: Run automated tests
```bash
cd /data/hostel-booking-system/frontend/svetu
yarn test --watchAll=false
```
**Acceptance**: All tests passing

### TASK-017: Check bundle size impact
```bash
cd /data/hostel-booking-system/frontend/svetu
yarn build
# Check build output for bundle sizes
```
**Acceptance**: Bundle size increase <50KB

## Phase 4: Automation Setup

### TASK-018: Create GitHub Actions workflow
```bash
cd /data/hostel-booking-system
cat > .github/workflows/validate-translations.yml << 'EOF'
name: Validate Translations
on: [pull_request]
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-python@v2
      - name: Validate translations
        run: |
          cd spec-kit/translation-audit
          python3 analyze_translations.py --ci-mode
EOF
```
**Acceptance**: Workflow file created

### TASK-019: Setup pre-commit hook
```bash
cd /data/hostel-booking-system
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
cd spec-kit/translation-audit
python3 analyze_translations.py --quick-check
EOF
chmod +x .git/hooks/pre-commit
```
**Acceptance**: Pre-commit hook validates translations

### TASK-020: Configure monitoring
```bash
# Add to frontend/svetu/src/app/[locale]/layout.tsx
# Sentry configuration for missing translations monitoring
```
**Acceptance**: Sentry captures MISSING_MESSAGE errors

## Phase 5: Documentation

### TASK-021: Create translation guide
```bash
cd /data/hostel-booking-system/docs
cat > TRANSLATION_GUIDE.md << 'EOF'
# Translation Guide

## Adding New Translations
1. Add key to en/ module first
2. Add translations to ru/ and sr/
3. Run validation script
4. Test in development

## Key Naming Convention
- Use camelCase for keys
- Group related keys under sections
- Keep nesting max 3 levels deep
EOF
```
**Acceptance**: Guide created and accessible

### TASK-022: Update README
```bash
cd /data/hostel-booking-system
# Add translation section to README.md
```
**Acceptance**: README includes translation info

### TASK-023: Final validation and cleanup
```bash
cd /data/hostel-booking-system/frontend/svetu
yarn format && yarn lint && yarn build
```
**Acceptance**: All commands execute successfully

## Summary

Total Tasks: 23
Estimated Time: 11 working days
Priority Order: TASK-001 through TASK-008 are critical and should be completed first

## Notes

- Tasks can be parallelized within phases
- Native speaker review recommended after TASK-014
- Monitor performance after deployment
- Keep backup of original translation files before changes