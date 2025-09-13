# Development Tasks: Complete Translation System Fix

## Phase 1: Cleanup and Structure Fix

### TASK-001: Create backup of current translations
```bash
cd /data/hostel-booking-system
mkdir -p backups/translations-$(date +%Y%m%d-%H%M%S)
cp -r frontend/svetu/src/messages/* backups/translations-$(date +%Y%m%d-%H%M%S)/
```
**Acceptance**: Backup created with all current translation files

### TASK-002: Remove invalid locale folder
```bash
cd /data/hostel-booking-system/frontend/svetu/src/messages
rm -rf modular-example
```
**Acceptance**: Only en, ru, sr folders remain

### TASK-003: Create translation fix script
```bash
cd /data/hostel-booking-system/spec-kit-translations/scripts
cat > fix-placeholder-translations.py << 'EOF'
# Script to replace placeholder translations with real ones
# [Implementation will be added in next task]
EOF
chmod +x fix-placeholder-translations.py
```
**Acceptance**: Script file created and executable

## Phase 2: Fix marketplace module translations

### TASK-004: Analyze marketplace placeholders
```bash
cd /data/hostel-booking-system/spec-kit-translations
python3 scripts/real-translation-audit.py | grep marketplace > marketplace-issues.txt
```
**Acceptance**: List of all marketplace translation issues identified

### TASK-005: Fix Russian marketplace translations
```python
# Fix all placeholder values in ru/marketplace.json
# Replace values like "marketplace.blackFridayTitle" with actual Russian text
# Example fixes:
# "blackFridayTitle": "marketplace.blackFridayTitle" -> "blackFridayTitle": "Черная пятница"
# "blackFridaySubtitle": "marketplace.blackFridaySubtitle" -> "blackFridaySubtitle": "Скидки до 70% на избранные товары"
```
**Acceptance**: No placeholders remain in ru/marketplace.json

### TASK-006: Fix Serbian marketplace translations  
```python
# Fix all placeholder values in sr/marketplace.json
# Use Latin script for Serbian
# Example fixes:
# "blackFridayTitle": "marketplace.blackFridayTitle" -> "blackFridayTitle": "Crni petak"
# "blackFridaySubtitle": "marketplace.blackFridaySubtitle" -> "blackFridaySubtitle": "Popusti do 70% na odabrane proizvode"
```
**Acceptance**: No placeholders remain in sr/marketplace.json

## Phase 3: Fix common module translations

### TASK-007: Fix Russian common module
```python
# Fix all placeholder values in ru/common.json
# Focus on navigation, buttons, and UI elements
# Examples:
# "viewDetails": "common.viewDetails" -> "viewDetails": "Подробнее"
# "addToCart": "common.addToCart" -> "addToCart": "В корзину"
```
**Acceptance**: All common UI elements properly translated to Russian

### TASK-008: Fix Serbian common module
```python
# Fix all placeholder values in sr/common.json
# Examples:
# "viewDetails": "common.viewDetails" -> "viewDetails": "Prikaži detalje"
# "addToCart": "common.addToCart" -> "addToCart": "Dodaj u korpu"
```
**Acceptance**: All common UI elements properly translated to Serbian

## Phase 4: Fix auth module translations

### TASK-009: Fix Russian auth module
```python
# Fix authentication-related translations in ru/auth.json
# Examples:
# "signIn": "auth.signIn" -> "signIn": "Войти"
# "register": "auth.register" -> "register": "Регистрация"
# "forgotPassword": "auth.forgotPassword" -> "forgotPassword": "Забыли пароль?"
```
**Acceptance**: All auth flows work in Russian

### TASK-010: Fix Serbian auth module
```python
# Fix authentication-related translations in sr/auth.json
# Examples:
# "signIn": "auth.signIn" -> "signIn": "Prijavite se"
# "register": "auth.register" -> "register": "Registracija"
```
**Acceptance**: All auth flows work in Serbian

## Phase 5: Fix remaining critical modules

### TASK-011: Fix cart module translations
```bash
# Fix both ru/cart.json and sr/cart.json
# Handle dynamic placeholders correctly
# Example: "itemsCount": "{count} товаров" (Russian)
```
**Acceptance**: Shopping cart fully translated

### TASK-012: Fix storefront module translations
```bash
# Fix both ru/storefronts.json and sr/storefronts.json
# Focus on store management terminology
```
**Acceptance**: Storefront management fully translated

### TASK-013: Fix orders module translations
```bash
# Fix both ru/orders.json and sr/orders.json
# Include order statuses and actions
```
**Acceptance**: Order management fully translated

## Phase 6: Automated fix for remaining modules

### TASK-014: Create bulk translation script
```python
# Create script to automatically translate remaining placeholders
# Use translation mapping from already fixed modules
# Apply consistent terminology
```
**Acceptance**: Script can process all remaining modules

### TASK-015: Run bulk translation fix
```bash
cd /data/hostel-booking-system/spec-kit-translations
python3 scripts/bulk-fix-translations.py --locales ru,sr --validate
```
**Acceptance**: All modules processed, no placeholders remain

## Phase 7: Validation and Testing

### TASK-016: Run comprehensive validation
```bash
cd /data/hostel-booking-system/spec-kit-translations
python3 scripts/real-translation-audit.py > final-audit.txt
grep "Coverage" final-audit.txt
```
**Acceptance**: Coverage >95%, 0 untranslated values

### TASK-017: Test UI in Russian
```bash
# Open browser at http://localhost:3001/ru
# Check main pages for any visible placeholders
# Test: Homepage, Login, Product listing, Cart, Checkout
```
**Acceptance**: No English or placeholder text visible

### TASK-018: Test UI in Serbian
```bash
# Open browser at http://localhost:3001/sr
# Check main pages for any visible placeholders
# Test: Homepage, Login, Product listing, Cart, Checkout
```
**Acceptance**: No English or placeholder text visible

### TASK-019: Fix any remaining issues
```bash
# Based on testing, fix any found issues
# Re-run validation
# Commit fixes
```
**Acceptance**: All issues resolved

### TASK-020: Create final report
```bash
cd /data/hostel-booking-system/spec-kit-translations
cat > reports/FINAL_TRANSLATION_FIX_REPORT.md << 'EOF'
# Final Translation Fix Report
## Results
- Coverage: [actual]%
- Missing Keys: 0
- Untranslated Values: 0
## Changes Made
[List of all modules fixed]
EOF
```
**Acceptance**: Comprehensive report documenting all changes

## Summary

**Total Tasks**: 20
**Estimated Time**: 3-4 days
**Priority**: Fix in order (TASK-001 through TASK-020)

## Critical Success Factors

1. NO placeholder text visible in UI (like "marketplace.blackFridayTitle")
2. Coverage must exceed 95%
3. All dynamic placeholders ({count}, {name}) handled correctly
4. Russian uses Cyrillic, Serbian uses Latin script
5. Consistent terminology across modules

## Notes

- Focus on replacing placeholder values, not just adding keys
- Test after each module to catch issues early
- Keep backups at each phase
- Use existing good translations as reference for consistency