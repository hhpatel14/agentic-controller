# 10 Migration Test Cases

Diverse open-source applications across 6 ecosystems for evaluating the questionnaire skill.
Each test case uses an intentionally vague migration prompt to test whether the skill asks the right follow-up questions.

---

## Test Cases

### #1 Coolstore
- **Repo**: [konveyor-ecosystem/coolstore](https://github.com/konveyor-ecosystem/coolstore)
- **Stack**: Java EE 7 / WebLogic / WAR
- **Vague Migration Prompt**: "Migrate this to Quarkus"

### #2 PlantsByWebSphere
- **Repo**: [WASdev/sample.plantsbywebsphere](https://github.com/WASdev/sample.plantsbywebsphere)
- **Stack**: Java EE / JSF / JSP / CDI / WebSphere
- **Vague Migration Prompt**: "Modernize this to Spring Boot"

### #3 DayTrader 7
- **Repo**: [WASdev/sample.daytrader7](https://github.com/WASdev/sample.daytrader7)
- **Stack**: Java EE 7 / EJB / JMS / JPA / WebSphere
- **Vague Migration Prompt**: "Migrate this to a cloud-native framework"

### #4 IBM Sample App Mod
- **Repo**: [IBM/sample-app-mod](https://github.com/IBM/sample-app-mod)
- **Stack**: Java 8 / WebSphere traditional / EJB / Servlets
- **Vague Migration Prompt**: "Upgrade this to Java 21 and Liberty"

### #5 Struts Legacy Mailreader
- **Repo**: [apache/struts-examples](https://github.com/apache/struts-examples)
- **Stack**: Struts 2 / JSP / XML config / Tiles
- **Vague Migration Prompt**: "Migrate this to Spring Boot"

### #6 Legacy Cycle Store
- **Repo**: [aws-samples/legacy-cycle-store-mvc-app](https://github.com/aws-samples/legacy-cycle-store-mvc-app)
- **Stack**: .NET Framework 4.x / MVC / Entity Framework
- **Vague Migration Prompt**: "Migrate this to .NET 8"

### #7 BookCatalog (.NET)
- **Repo**: [microsoft/dotnet-modernization-for-beginners](https://github.com/microsoft/dotnet-modernization-for-beginners)
- **Stack**: ASP.NET MVC 5 / .NET Framework 4.8
- **Vague Migration Prompt**: "Modernize this app"

### #8 Flask-to-FastAPI Tutorial
- **Repo**: [jtemporal/flask-to-fastapi](https://github.com/jtemporal/flask-to-fastapi)
- **Stack**: Python / Flask / Jinja2 / SQLite
- **Vague Migration Prompt**: "Migrate this to FastAPI"

### #9 NodeJS Ecommerce Store
- **Repo**: [mrmodise/nodejs-ecommerce-store](https://github.com/mrmodise/nodejs-ecommerce-store)
- **Stack**: Node.js / Express.js / MongoDB / EJS
- **Vague Migration Prompt**: "Modernize this to NestJS with TypeScript"

### #10 Struts 1.3 to Spring Boot
- **Repo**: [SarahRSilva/struts-spring-boot-legacy-code](https://github.com/SarahRSilva/struts-spring-boot-legacy-code)
- **Stack**: Struts 1.3 / JSP / XML / Java 8
- **Vague Migration Prompt**: "Migrate this from Struts to modern Java"

---

## Expected Questions (what the questionnaire skill should surface)

### #1 Coolstore — "Migrate this to Quarkus"
- Which Quarkus version?
- What about the `audit-logging-library` jar?
- JMS replacement strategy?
- Keep monolith or split?

### #2 PlantsByWebSphere — "Modernize this to Spring Boot"
- Which Spring Boot version?
- Replace JSF/JSP with what (Thymeleaf, React, API-only)?
- What about the CDI beans — keep CDI or switch to Spring DI?

### #3 DayTrader 7 — "Migrate this to a cloud-native framework"
- Cloud-native means what — Quarkus? Spring Boot? MicroProfile?
- The app has heavy EJB + JMS — replace with what?
- Database stays the same?

### #4 IBM Sample App Mod — "Upgrade this to Java 21 and Liberty"
- This is a runtime upgrade, not a framework rewrite. Keep EJBs or convert to CDI?
- Jakarta namespace migration needed?
- Containerize?

### #5 Struts Mailreader — "Migrate this to Spring Boot"
- Replace Tiles + JSP with what?
- Struts actions to Spring controllers — 1:1 mapping or redesign?
- Keep XML config or annotation-based?

### #6 Legacy Cycle Store — "Migrate this to .NET 8"
- EF6 to EF Core — any breaking schema changes?
- Razor views stay or move to Blazor/API?
- Windows-specific APIs used?

### #7 BookCatalog — "Modernize this app"
- Modernize to what? .NET 8? Containerize? Re-platform to cloud?
- This prompt is maximally vague — skill must clarify the target.

### #8 Flask to FastAPI — "Migrate this to FastAPI"
- Sync or async endpoints?
- Pydantic models for validation?
- Keep Jinja2 templates or go API-only?
- SQLAlchemy stays?

### #9 NodeJS Ecommerce — "Modernize this to NestJS with TypeScript"
- Keep MongoDB or switch to PostgreSQL?
- EJS templates — what replaces them (API-only, React)?
- Auth system (Facebook OAuth) — keep or replace?

### #10 Struts 1.3 — "Migrate this from Struts to modern Java"
- Modern Java means what — Spring Boot? Quarkus? Jakarta EE?
- Struts 1.3 ActionForms — what pattern replaces them?
- JSP — what frontend?

---

## Evaluation Criteria

For each test case, evaluate the questionnaire skill output on:

1. **Detection accuracy** — did it correctly identify the source stack, frameworks, and dependencies?
2. **Unknowns surfaced** — did it flag proprietary/unfamiliar libraries and patterns?
3. **Target clarity** — when the prompt is vague (e.g., "modernize this app"), did the skill ask or make a reasonable choice?
4. **Decision quality** — are the priorities, constraints, and scope appropriate for the detected app?
5. **Reasoning depth** — does the reasoning explain WHY, not just WHAT?
6. **Missing questions** — compare against the expected questions above. What did the skill miss?

## Results Tracking

| # | App | Detection | Unknowns | Target | Decisions | Reasoning | Missing Qs | Notes |
|---|-----|-----------|----------|--------|-----------|-----------|------------|-------|
| 1 | Coolstore | | | | | | | |
| 2 | PlantsByWebSphere | | | | | | | |
| 3 | DayTrader 7 | | | | | | | |
| 4 | IBM Sample App Mod | | | | | | | |
| 5 | Struts Mailreader | | | | | | | |
| 6 | Legacy Cycle Store | | | | | | | |
| 7 | BookCatalog | | | | | | | |
| 8 | Flask to FastAPI | | | | | | | |
| 9 | NodeJS Ecommerce | | | | | | | |
| 10 | Struts 1.3 | | | | | | | |
