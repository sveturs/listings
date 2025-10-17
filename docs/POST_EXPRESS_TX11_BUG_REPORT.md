# Post Express TX 11 Bug Report - Database Error

**Date:** 2025-10-14
**Transaction:** TX 11 - PostarinaPosiljke (Calculate Postage)
**Status:** ❌ **Post Express API Bug** - Needs to be reported to Post Express

---

## 🐛 Problem Description

TX 11 (Calculate Postage) returns a **database error from Post Express API**:

```
"Column 'PREVOD_SR' does not belong to table Prevodi."
```

**This is a bug in Post Express WSP API** - their database schema has an issue with the `PREVOD_SR` column.

---

## 📊 Test Data

**Request sent to Post Express API:**

```json
{
  "TransactionId": "11",
  "Data": {
    "IdRukovanje": 71,
    "PostanskiBrojOdlaska": "11000",
    "PostanskiBrojDolaska": "21000",
    "Masa": 500,
    "Otkupnina": 0,
    "Vrednost": 0,
    "PosebneUsluge": "PNA"
  }
}
```

**Response from Post Express:**

```json
{
  "Rezultat": 3,
  "StrOut": null,
  "StrRezultat": "{\"Poruka\":\"Column 'PREVOD_SR' does not belong to table Prevodi.\",\"PorukaKorisnik\":\"Column 'PREVOD_SR' does not belong to table Prevodi.\",\"Info\":null}"
}
```

---

## 🔍 Backend Logs

```
DEBUG: WSP API Request - transaction_id: b7e4ddf8-7f4c-4593-b7ae-33d52c28b565, type: 11
DEBUG: WSP API Response - status_code: 200, execution_time_ms: 20
DEBUG: WSP API Raw Response Body: {"Rezultat":3,"StrOut":null,"StrRezultat":"{\"Poruka\":\"Column 'PREVOD_SR' does not belong to table Prevodi.\",\"PorukaKorisnik\":\"Column 'PREVOD_SR' does not belong to table Prevodi.\",\"Info\":null}"}
ERROR: WSP transaction failed - Rezultat: 3, Poruka: unknown error
ERROR: Failed to calculate postage - error: CalculatePostage failed: unknown error
```

---

## ✅ Our Implementation Status

**Our side is correct:**
- ✅ Request format matches WSP API specification
- ✅ All required fields are present (IdRukovanje, PostanskiBrojOdlaska, PostanskiBrojDolaska, Masa)
- ✅ Data types are correct (integers, strings)
- ✅ PascalCase field naming is correct
- ✅ API endpoint is correct (http://212.62.32.201/WspWebApi/transakcija)
- ✅ **Backend parses StrRezultat and extracts detailed Post Express errors**
- ✅ **Backend returns detailed errors via API** (`{"error": "Column 'PREVOD_SR' does not belong to table Prevodi."}`)
- ✅ **Frontend now displays detailed Post Express error messages** (not generic placeholders)

**The bug is on Post Express side** - their database has a schema issue.

---

## 🎯 Tested Services (IdRukovanje)

All services fail with the same database error:

| IdRukovanje | Service Name | Status |
|------------|--------------|---------|
| 29 | PE_Danas_za_sutra_12 | ❌ DB Error |
| 30 | PE_Danas_za_danas | ❌ DB Error |
| 55 | PE_Danas_za_odmah | ❌ DB Error |
| 58 | PE_Danas_za_sutra_19 | ❌ DB Error |
| 59 | PE_Danas_za_odmah_Bg | ❌ DB Error |
| 71 | PE_Danas_za_sutra_isporuka | ❌ DB Error |
| 85 | Isporuka_na_paketomatu | ❌ DB Error |

---

## 📝 Action Items

### ✅ Done (Our Side)

1. ✅ Fixed frontend error handling for better error messages
2. ✅ Added optional chaining (`?.`) to prevent "Cannot read properties of undefined" errors
3. ✅ **Modified backend to extract and return detailed Post Express errors from `StrRezultat`**
   - Backend `parseWSPResponse()` parses `StrRezultat` JSON when `Rezultat != 0`
   - Extracts `Poruka` or `PorukaKorisnik` from nested error structure
   - Returns detailed error immediately in `TransactionResponse.ErrorMessage`
4. ✅ **Modified backend `CalculatePostage()` to return error details without throwing**
   - Changed logic to return response structure with `Rezultat: 3` and detailed `Poruka`
   - Handler now receives and returns detailed error message to frontend
5. ✅ **Fixed frontend error display in TX 11 modal**
   - Changed line 506 from `response?.data?.message` to `response?.data?.error || response?.data?.message`
   - Frontend now correctly displays: "Column 'PREVOD_SR' does not belong to table Prevodi."
   - Also fixed TX 6 and TX 9 modals for consistency
6. ✅ Documented the issue in this report

### ⏳ Waiting (Post Express Side)

1. ❌ **Contact Post Express support** - Report the database bug:
   - Error: `Column 'PREVOD_SR' does not belong to table Prevodi.`
   - Transaction: TX 11 (PostarinaPosiljke)
   - Occurs for all service IDs (IdRukovanje)
   - Their database schema needs to be fixed

2. ❌ **Wait for Post Express fix** - Until they fix their database, TX 11 will not work

---

## 🔄 Workaround

**None available** - This is a fundamental database error on Post Express side.

**Alternative options:**
1. Use TX 73 (B2B Manifest) - Create full shipments instead of just calculating postage
2. Contact Post Express for manual postage calculations
3. Use fixed postage rates from documentation

---

## 📧 Contact Post Express

**Email:** support@postexpress.rs (or their B2B support contact)

**Subject:** TX 11 Database Error - Column 'PREVOD_SR' Missing

**Message Template:**

```
Dear Post Express Support,

We are integrating with your WSP API (version 0.2.4) and encountering a database error
when calling Transaction 11 (PostarinaPosiljke - Calculate Postage).

Error details:
- Transaction: TX 11
- Rezultat: 3
- Error message: "Column 'PREVOD_SR' does not belong to table Prevodi."
- Occurs for all service IDs (IdRukovanje: 29, 30, 55, 58, 59, 71, 85)

Our request format is correct according to the API specification. This appears to be
a database schema issue on your side.

Could you please investigate and fix this issue?

Test request we are sending:
{
  "IdRukovanje": 71,
  "PostanskiBrojOdlaska": "11000",
  "PostanskiBrojDolaska": "21000",
  "Masa": 500,
  "Otkupnina": 0,
  "Vrednost": 0,
  "PosebneUsluge": "PNA"
}

Thank you for your assistance.

Best regards,
SVETU.rs Development Team
```

---

## 📌 Related Documentation

- TX 3 (GetNaselje): ✅ Working
- TX 4 (GetUlica): ✅ Working
- TX 6 (ProveraAdrese): ⚠️ Needs real address data for testing
- TX 9 (ProveraDostupnostiUsluge): ⚠️ Needs testing
- TX 11 (PostarinaPosiljke): ❌ **Post Express Database Bug**
- TX 73 (B2B Manifest): ✅ Working

**Main Report:** `/docs/POST_EXPRESS_TX_3_11_FINAL_REPORT.md`

---

## 🎯 Complete Error Flow (Working!)

1. **Post Express API** returns database error in `StrRezultat`:
   ```json
   {
     "Rezultat": 3,
     "StrRezultat": "{\"Poruka\":\"Column 'PREVOD_SR' does not belong to table Prevodi.\",...}"
   }
   ```

2. **Backend** (`parseWSPResponse()`) extracts detailed error:
   - Parses `StrRezultat` JSON
   - Extracts `Poruka`: "Column 'PREVOD_SR' does not belong to table Prevodi."
   - Returns in `TransactionResponse.ErrorMessage`

3. **Backend** (`CalculatePostage()`) creates response with error:
   - Returns `PostageCalculationResponse{Rezultat: 3, Poruka: "..."}` (not error)
   - Handler checks `Rezultat != 0` and returns 500 error

4. **Backend API** returns to frontend:
   ```json
   {"error": "Column 'PREVOD_SR' does not belong to table Prevodi."}
   ```

5. **Frontend** displays detailed error:
   - Extracts `response?.data?.error`
   - Shows actual Post Express database error to user
   - ✅ No more generic "Failed to calculate postage" message!

---

**Updated:** 2025-10-14 21:30
**Status:** Waiting for Post Express to fix their database
**Our Status:** ✅ All error handling fixed and working correctly
