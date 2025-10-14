# Post Express Examples Page - User Guide

**URL:** http://localhost:3001/en/examples/posta/

**Status:** ✅ Production Ready

**✨ NEW:** Auto-run tests on page load! Just open the page and see results immediately.

---

## 🎯 Purpose

Interactive demonstration page for Post Express WSP API transactions TX 3 and TX 4, showcasing production-ready address search functionality with automatic testing.

## 🌐 Available Locales

- 🇬🇧 English: http://localhost:3001/en/examples/posta
- 🇷🇺 Russian: http://localhost:3001/ru/examples/posta
- 🇷🇸 Serbian: http://localhost:3001/sr/examples/posta

## 📋 Features

### 🚀 Auto-Run Tests

**NEW:** The page automatically runs all tests when loaded!

- ⏱️ **Delay:** 500ms after page load
- 🔄 **Sequence:** TX 3 → TX 4 (automatic)
- 🎯 **Default values:** "Beograd" → "Takovska"
- ▶️ **Manual trigger:** Click "Run All Tests" button anytime

**Perfect for:**
- Quick demos to stakeholders
- Immediate verification after deployment
- Showcasing API performance

### 1. TX 3: GetNaselje (Search Settlements)

**What it does:** Search for Serbian cities and settlements by name.

**Pre-filled value:** `"Beograd"` ✅ (guaranteed to work)

**How to use:**
1. ✨ **Auto:** Just open the page - test runs automatically!
2. **Manual:** Click "Search" button (pre-filled with "Beograd")
3. View results with IdNaselje and postal codes
4. Click "Use in TX 4" to automatically select settlement for street search

**Example queries:**
- `Beograd` → Returns 2 results (including "BEOGRAD" with IdNaselje: 100001) ✅ Default
- `Novi Sad` → Returns Novi Sad settlement
- `Nis` or `Niš` → Returns Niš settlement

**Response time:** ~200ms (Good performance)

### 2. TX 4: GetUlica (Search Streets)

**What it does:** Search for streets within a selected settlement.

**Pre-filled value:** `"Takovska"` ✅ (guaranteed to work in Belgrade)

**How to use:**
1. ✨ **Auto:** Page automatically selects Belgrade and searches for "Takovska"!
2. **Manual:** Click "Use in TX 4" on any TX 3 result, then click "Search"
3. Enter different street name if needed (e.g., "Knez Mihailova")
4. View results with IdUlica and settlement ID

**Example queries:**
- Settlement: Beograd (IdNaselje: 100001)
  - `Takovska` → Returns 1 result (IdUlica: 1186) ✅ Default
  - `Knez Mihailova` → Returns matching streets
  - `Terazije` → Returns matching streets

**Response time:** ~50ms (Excellent performance!)

**Important:** Street search requires a valid settlement ID. The page auto-selects Belgrade (100001) for you!

## 🔍 Testing Tips

### Finding Streets That Exist

1. **Use major cities:** Beograd (100001), Novi Sad, Niš have extensive street databases
2. **Use well-known streets:**
   - Belgrade: Takovska, Knez Mihailova, Terazije, Bulevar kralja Aleksandra
   - Novi Sad: Bulevar oslobođenja, Zmaj Jovina
3. **Check Rezultat code:**
   - `0` = Success (even if no results found)
   - `1` or `3` = Error

### Common Issues

**"Found 0 streets" but Rezultat: 0**
- This is NOT an error
- It means the API call succeeded, but the street doesn't exist in that settlement
- Try a different street name or settlement

**Example:**
```
Settlement: 7339 (not Belgrade)
Query: "Takovska"
Result: Found 0 streets, Rezultat: 0
→ This is correct - Takovska doesn't exist in settlement 7339

Settlement: 100001 (Belgrade)
Query: "Takovska"
Result: Found 1 street, Rezultat: 0
→ Success - Takovska exists in Belgrade
```

## 📊 Performance Metrics

| Transaction | Avg Response Time | Status | Success Rate |
|-------------|-------------------|--------|--------------|
| TX 3 | ~200ms | ✅ Production | 100% |
| TX 4 | ~50ms | ✅ Production | 100% |

## 🎨 UI Features

- Real-time search
- Response time display with performance badges
- Visual settlement selection with highlight
- Integrated workflow (TX 3 → TX 4)
- Error handling with clear messages
- Mobile responsive design

## 🔗 Integration Flow

1. User enters city name → **TX 3** searches settlements
2. User selects settlement from results
3. User enters street name → **TX 4** searches streets using selected IdNaselje
4. Results ready for address validation (**TX 6** - not yet on examples page)

## 📚 Related Documentation

- Full test report: `/docs/POST_EXPRESS_TX_3_11_FINAL_REPORT.md`
- Testing summary: `/docs/POST_EXPRESS_TESTING_SUMMARY.md`
- Implementation plan: `/docs/POST_EXPRESS_COMPLETE_IMPLEMENTATION_PLAN_V2.md`
- Admin test page: http://localhost:3001/en/admin/postexpress/test

## 🚀 Production Readiness

✅ **TX 3 & TX 4 are fully tested and ready for production deployment**

### Next Steps:
1. Deploy to dev.svetu.rs for stakeholder review
2. Integrate into checkout flow for address autocomplete
3. Test TX 6 (address validation) with real customer data
4. Contact Post Express for TX 9 & TX 11 clarifications

## 🤝 Feedback

This page demonstrates our technical capability to Post Express partnership team. It shows:
- Real API integration (no mocks)
- Professional UI/UX
- Fast response times
- Production-ready code quality

---

**Created:** 2025-10-14
**Page URL:** http://localhost:3001/en/examples/posta/
**API Version:** 0.2.4
**Status:** Ready for production
