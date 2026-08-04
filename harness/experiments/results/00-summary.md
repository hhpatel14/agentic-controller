# Questionnaire Skill Evaluation — Summary

10 open-source applications across 6 ecosystems, each given a vague migration prompt.

## Scorecard

| # | App | Stack | Expected Qs | Covered | Extra Qs | Key Gap |
|---|-----|-------|-------------|---------|----------|---------|
| 1 | Coolstore | Java EE 7 / WebLogic | 4 | 4/4 (100%) | +6 | None — caught audit-lib, JMS, WebLogic stubs, monolith question |
| 2 | PlantsByWebSphere | Java EE / JSF / WebSphere | 3 | 3/3 (100%) | +9 | Missing: no-test detection, admin vs storefront split |
| 3 | DayTrader 7 | Java EE 7 / EJB / JMS | 3 | 3/3 (100%) | +6 | Missing: microservices decomposition, "do nothing" option |
| 4 | IBM Sample App Mod | Java 8 / WebSphere tWAS | 3 | 3/3 (100%) | +4 | Missing: question tiering (blocking vs optional) |
| 5 | Struts Examples | Struts 2 / JSP / Tiles | 3 | 3/3 (100%) | +7 | Missing: OGNL expression migration, scope-gating for 46 modules |
| 6 | Legacy Cycle Store | .NET Framework 4.x / EF | 3 | 3/3 (100%) | +6 | Missing: incomplete repo detection, EF version precision (EF5 not EF6) |
| 7 | BookCatalog | ASP.NET MVC 5 / .NET 4.8 | 2 | 2/2 (100%) | +7 | Missing: question priority tiers, mandatory-step vs decision distinction |
| 8 | Flask-to-FastAPI | Python / Flask / Jinja2 | 4 | 3/4 (75%) | +5 | SQLAlchemy question correctly skipped (no DB in this app) |
| 9 | NodeJS Ecommerce | Express.js / MongoDB / EJS | 3 | 3/3 (100%) | +6 | Missing: dependency chain ordering, security-first vs optional split |
| 10 | Struts 1.3 Legacy | Struts 1.3 → Spring Boot | 3 | 1/3 (33%) | +3 | **Critical**: repo is already migrated — needs pre-flight mismatch detector |

**Overall: 28/32 expected questions covered (87.5%)**

## Cross-Cutting Gaps (patterns that appeared in 3+ apps)

### 1. No question priority tiers
Every app surfaced this. The skill treats all questions equally — but some are **blocking decisions** (target framework, scope) while others are **implementation details** with sensible defaults (config format, test framework). The skill should tier questions:
- **Tier 1 (blocking)**: Must answer before planning can start
- **Tier 2 (design)**: Influences the plan but has a reasonable default
- **Tier 3 (detail)**: Can be deferred to execution

### 2. No pre-flight mismatch detection
App #10 exposed this critically — the repo name says "Struts" but the code is already Spring Boot. The skill should verify that the source tech in the prompt matches what's actually in the code before asking questions. If they don't match, that's the first question.

### 3. No scope-gating for multi-module projects
Apps #3 and #5 (DayTrader with 3 modules, Struts with 46 modules) both needed a "migrate all or subset?" question before anything else. The skill should detect multi-module projects and ask about scope first.

### 4. Missing test coverage detection
Apps #2, #6, #7, #8 all had zero or minimal tests. The skill should flag this and ask: "No tests found — should we add characterization tests before migrating?"

### 5. No "do nothing" option
Apps #3 and #10 could reasonably stay on their current stack. The skill never offers "the current stack is fine, here's why" as an option. For enterprise migrations, this is a valid outcome.

### 6. Mandatory steps presented as questions
App #7 showed this — SDK-style csproj conversion is mandatory for .NET migration, not a decision. The skill should separate **facts** (things that must happen) from **decisions** (things the user chooses).

### 7. No dependency chain ordering
App #9 surfaced this — frontend strategy determines auth approach, which determines session handling. Questions should be ordered by dependency: answer the upstream decision first, then ask downstream ones.

### 8. No incomplete repo detection
App #6 had source files referencing classes that don't exist on disk. The skill should flag incomplete or broken repos before producing a questionnaire.

## Ecosystem Coverage

| Ecosystem | Apps Tested | Detection Quality | Question Quality |
|-----------|------------|-------------------|-----------------|
| Java EE / WebSphere | #1, #2, #3, #4 | Excellent — caught vendor-specific APIs, JNDI, EJB patterns | Good — missed scope-gating |
| Struts | #5, #10 | Good on #5, caught mismatch on #10 | Needs pre-flight check |
| .NET Framework | #6, #7 | Good — caught EF, Windows APIs, membership provider | Needs priority tiers |
| Python / Flask | #8 | Good — correctly skipped inapplicable questions | Solid |
| Node.js / Express | #9 | Good — caught Elasticsearch, Stripe, deprecated deps | Needs dependency ordering |

## Recommended Skill Changes (ranked by impact)

1. **Add pre-flight mismatch check** — verify prompt matches code before asking anything
2. **Add question priority tiers** — blocking → design → detail
3. **Add scope-gating for multi-module projects** — "all or subset?" first
4. **Separate facts from decisions** — mandatory steps are not questions
5. **Detect missing tests** — flag and ask about pre-migration test strategy
6. **Order questions by dependency chain** — upstream decisions before downstream
7. **Offer "do nothing" as a valid option** — not every app needs migration
8. **Detect incomplete repos** — flag missing source files referenced by manifests
