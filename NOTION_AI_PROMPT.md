# 🤖 Промпт для Notion AI

Скопируйте этот текст и вставьте в Notion AI:

---

Create a comprehensive project management database for SveTu Marketplace Platform with the following structure:

## Database Schema:
Create a table with these properties:
- Name (Title)
- Phase (Select): Foundation & Launch, Payments & Delivery, Growth & Optimization, Scale & Innovate, Consolidation & Expansion
- Quarter (Select): Q1 2025, Q2 2025, Q3 2025, Q4 2025
- Month (Select): January through December 2025
- Week (Text): Week 1-2, Week 3-4, Full Month
- Status (Select): 📋 Backlog, 📌 To Do, 🚀 In Progress, 🔍 Review, ✅ Done
- Priority (Select): 🔴 Critical, 🟡 High, 🟢 Medium, 🔵 Low
- Epic (Select): 🏪 Storefronts, 🗺️ Maps, 💳 Payments, 📦 Logistics, 🚗 Automotive, 🌍 i18n, 👥 Admin, 🔧 DevOps, 📱 UX/UI, 🏭 WMS
- Owner (Person): Backend Dev, Frontend Dev, DevOps, Mobile Dev, Full Stack, Security, Business Dev, Team
- Estimate (Number): hours
- Budget (Number): EUR currency
- KPI (Text)
- Start Date (Date)
- End Date (Date)
- Dependencies (Text)
- Risks (Text)
- Tasks (Text): detailed checklist
- Description (Text): full description

## Add these 30 milestones:

### Q1 2025 - Foundation & Launch

1. **Infrastructure Setup**
   - Phase: Foundation & Launch
   - Quarter: Q1 2025
   - Month: January
   - Week: Week 1-2
   - Status: 📌 To Do
   - Priority: 🔴 Critical
   - Owner: DevOps
   - Budget: €5,000
   - Start: 2025-01-01
   - End: 2025-01-14
   - Tasks: Setup production servers; Configure Docker/K8s; SSL certificates; CI/CD pipeline; Backup strategy
   - KPI: 99.9% uptime

2. **Security Audit**
   - Quarter: Q1 2025
   - Month: January
   - Week: Week 1-2
   - Priority: 🔴 Critical
   - Owner: Security
   - Budget: €3,000
   - Start: 2025-01-01
   - End: 2025-01-14
   - Tasks: OWASP Top 10; Penetration testing; DDoS protection; WAF setup

3. **Stripe Integration**
   - Month: January
   - Week: Week 3-4
   - Priority: 🔴 Critical
   - Epic: 💳 Payments
   - Owner: Backend Dev
   - Budget: €2,000
   - Start: 2025-01-15
   - End: 2025-01-31
   - Tasks: Stripe account for Serbia; Checkout Session API; Webhooks; PCI DSS compliance

4. **Monitoring Setup**
   - Month: January
   - Week: Week 3-4
   - Priority: 🟡 High
   - Epic: 🔧 DevOps
   - Owner: DevOps
   - Budget: €1,000
   - Tasks: Prometheus + Grafana; Telegram alerts; Sentry; Uptime monitoring

5. **Post Express Production**
   - Phase: Payments & Delivery
   - Month: February
   - Week: Week 1-2
   - Priority: 🔴 Critical
   - Epic: 📦 Logistics
   - Budget: €2,000
   - Tasks: Production credentials; Webhook integration; Label printing; Real-time tracking

6. **Stripe Connect**
   - Month: February
   - Week: Week 1-2
   - Priority: 🟡 High
   - Epic: 💳 Payments
   - Budget: €3,000
   - Tasks: Seller onboarding; Auto payouts; Commission model; Reporting dashboard

7. **PaySpot Integration**
   - Month: February
   - Week: Week 3-4
   - Priority: 🟡 High
   - Epic: 💳 Payments
   - Budget: €4,000
   - Tasks: API integration; Escrow mechanism; Dispute resolution; Rating system

8. **BEX Express Integration**
   - Month: February
   - Week: Week 3-4
   - Priority: 🟢 Medium
   - Epic: 📦 Logistics
   - Budget: €2,000
   - Tasks: Contract & credentials; API integration; Delivery calculator; Zone mapping

9. **Variant Attributes System**
   - Phase: Storefronts Enhancement
   - Month: March
   - Week: Week 1-2
   - Priority: 🟡 High
   - Epic: 🏪 Storefronts
   - Budget: €3,000
   - Tasks: DB schema update; Admin UI; Attribute inheritance; Inventory tracking; Bulk edit

10. **Import/Export System**
    - Month: March
    - Week: Week 1-2
    - Priority: 🟢 Medium
    - Epic: 🏪 Storefronts
    - Budget: €2,000
    - Tasks: CSV/Excel parser; Import templates; Validation; Error reporting; Scheduled imports

### Q2 2025 - Growth & Optimization

11. **Mobile App Development**
    - Phase: Growth & Optimization
    - Quarter: Q2 2025
    - Month: April
    - Week: Full Month
    - Priority: 🟡 High
    - Epic: 📱 UX/UI
    - Owner: Mobile Dev
    - Budget: €10,000
    - Tasks: React Native setup; Core navigation; Auth flow; Product browsing; Chat integration

12. **PWA Enhancement**
    - Month: April
    - Week: Week 1-2
    - Priority: 🟢 Medium
    - Budget: €2,000
    - Tasks: Offline mode; Push notifications; App-like UX; Install prompts

13. **UX Improvements**
    - Month: April
    - Week: Week 3-4
    - Priority: 🟡 High
    - Epic: 📱 UX/UI
    - Budget: €3,000
    - Tasks: Search autocomplete; Advanced filters; Saved searches; Personalization; Recommendations

14. **Car Section Cleanup**
    - Month: May
    - Week: Week 1-2
    - Priority: 🟢 Medium
    - Epic: 🚗 Automotive
    - Budget: €4,000
    - Tasks: DB optimization; VIN decoder; Make/Model structure; Special filters; Price history

15. **Auto Features**
    - Month: May
    - Week: Week 3-4
    - Priority: 🟢 Medium
    - Epic: 🚗 Automotive
    - Budget: €5,000
    - Tasks: Vehicle history; Insurance calc; Loan calc; Inspection list; 360° viewer

16. **Multi-language Support**
    - Month: June
    - Week: Week 1-2
    - Priority: 🟡 High
    - Epic: 🌍 i18n
    - Budget: €5,000
    - Tasks: Translation system; AI translations; Content moderation; SEO per language; Currency converter

17. **Regional Expansion**
    - Month: June
    - Week: Week 1-2
    - Priority: 🟡 High
    - Owner: Business Dev
    - Budget: €15,000
    - Tasks: Montenegro research; Bosnia entry; Croatia study; Legal compliance; Partnerships

18. **Performance Optimization**
    - Month: June
    - Week: Week 3-4
    - Priority: 🟡 High
    - Epic: 🔧 DevOps
    - Budget: €3,000
    - Tasks: CDN setup; Image optimization; DB indexing; Caching; Load testing

### Q3 2025 - Scale & Innovate

19. **WMS System**
    - Phase: Scale & Innovate
    - Quarter: Q3 2025
    - Month: July
    - Week: Full Month
    - Priority: 🟢 Medium
    - Epic: 🏭 WMS
    - Owner: Team
    - Budget: €20,000
    - Tasks: DDD architecture; Receiving module; Storage locations; Pick & pack; Inventory; Barcode scanning

20. **AI Features**
    - Month: August
    - Week: Full Month
    - Priority: 🟢 Medium
    - Owner: AI Dev
    - Budget: €15,000
    - Tasks: Image recognition; Price recommendations; Fraud detection; Chatbot; Content generation; Forecasting

21. **Business Intelligence**
    - Month: September
    - Week: Full Month
    - Priority: 🟢 Medium
    - Owner: Data Team
    - Budget: €12,000
    - Tasks: Data warehouse; ETL pipelines; Custom dashboards; Predictive analytics; Market insights

### Q4 2025 - Consolidation & Expansion

22. **B2B Marketplace**
    - Phase: Consolidation & Expansion
    - Quarter: Q4 2025
    - Month: October
    - Week: Full Month
    - Priority: 🟡 High
    - Owner: B2B Team
    - Budget: €18,000
    - Tasks: Bulk ordering; Quote system; Net terms; Volume discounts; Corporate accounts; EDI

23. **Social Commerce**
    - Month: November
    - Week: Full Month
    - Priority: 🟢 Medium
    - Owner: Social Team
    - Budget: €10,000
    - Tasks: Reviews & ratings; Social sharing; Influencer program; Referrals; Forums; Live streaming

24. **Platform 2.0 Innovation**
    - Month: December
    - Week: Full Month
    - Priority: 🔵 Low
    - Owner: Innovation Team
    - Budget: €25,000
    - Tasks: Blockchain; Crypto payments; NFT marketplace; Metaverse; AR try-on; Voice commerce; IoT

## Create these Views:

1. **Timeline View** - Gantt chart grouped by Quarter, colored by Priority
2. **Kanban Board** - Grouped by Status, sorted by Priority
3. **Sprint Planning** - Filtered by current month, grouped by Week
4. **Critical Path** - Filtered by Priority = 🔴 Critical
5. **Budget Overview** - Table view sorted by Budget descending
6. **Team Workload** - Grouped by Owner, sorted by Start Date
7. **Risk Management** - Filtered where Risks is not empty
8. **Calendar** - Calendar view by Start Date

## Add these Formulas:

1. **Progress %**:
```
if(prop("Status") == "✅ Done", 100, 
  if(prop("Status") == "🚀 In Progress", 50, 
    if(prop("Status") == "📌 To Do", 25, 0)))
```

2. **Days Remaining**:
```
dateBetween(prop("End Date"), now(), "days")
```

3. **Status Indicator**:
```
if(and(prop("End Date") < now(), prop("Status") != "✅ Done"), 
  "⚠️ Overdue", "✅ On track")
```

4. **Quarter Progress**:
```
round(dateBetween(now(), date(2025, 1, 1), "days") / 365 * 100) + "%"
```

## Add Summary Metrics:

- Total Budget: €150,000
- Target Users: 100,000 by Dec 2025
- Target GMV: €800,000/month
- Target Listings: 55,000
- Team Size: 5-15 people

Make the database interactive, visually appealing, and ready for project management. Add colors, emojis, and make it easy to navigate.

---

## 💡 Дополнительные команды для Notion AI:

После создания базы, можете попросить:

1. "Add a dashboard page with charts showing project progress"
2. "Create a risk matrix based on the Risks column"
3. "Generate a weekly status report template"
4. "Add automation rules for status updates"
5. "Create a budget burndown chart"