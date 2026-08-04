# Pattern Analysis: What the LLM Figured Out vs. What It Asked

Derived from `interview-questions.md` across 10 applications and 6 ecosystems.

---

## Part 1: What the LLM Could NOT Figure Out (Unknowns by Category)

Across all 10 apps, the unknowns cluster into **7 recurring categories**:

### Category A: Target Ambiguity
*The migration prompt names a direction but not enough to start.*

| App | Unknown |
|-----|---------|
| #3 DayTrader | "Cloud-native framework" — Quarkus? Spring Boot? Micronaut? |
| #7 BookCatalog | "Modernize" — upgrade framework? containerize? rearchitect? |
| #10 Struts Legacy | "Modern Java" — but the app is already Spring Boot 3.5 on Java 21 |
| #4 IBM App Mod | "Liberty" — Open Liberty or WebSphere Liberty? |

**Pattern**: When the prompt says a category instead of a specific technology ("cloud-native", "modernize", "modern Java"), the LLM cannot pick a target without asking.

---

### Category B: Proprietary / Unrecognized Dependencies
*Something in the code the LLM has no training data on.*

| App | Unknown |
|-----|---------|
| #1 Coolstore | `audit-logging-library-1.0.0.jar` — purpose unknown, no source |
| #2 PlantsByWebSphere | IBM deployment descriptors (`ibm-web-bnd.xml`, `ibm-ejb-jar-bnd.xml`) |
| #4 IBM App Mod | `com.ibm.websphere.runtime.ServerName`, `WsnInitialContextFactory`, CORBA/IIOP |
| #6 Legacy Cycle Store | ComponentOne Wijmo controls (`C1.C1Report.4`, `C1.Web.Wijmo.Controls.4`) |

**Pattern**: Vendored JARs, commercial UI components, and vendor-specific runtime APIs are consistently unresolvable. The LLM can name them but cannot determine their purpose or replacement without asking.

---

### Category C: Messaging / Async Replacement Strategy
*The LLM detects the pattern but can't pick the replacement.*

| App | Unknown |
|-----|---------|
| #1 Coolstore | JMS topic `topic/orders` — replace with what? CDI Events, Kafka, Reactive Messaging? |
| #3 DayTrader | 2 MDBs (order broker + market streamer) on JMS queues/topics — what broker? |
| #9 NodeJS Ecommerce | `mongoosastic` Elasticsearch sync — keep, replace with Atlas Search, or drop? |

**Pattern**: The LLM can detect messaging/async usage from annotations and code, but the replacement technology is always a user decision because it depends on infrastructure (what broker do you run?) and architecture (in-process vs external).

---

### Category D: UI / View Layer Decision
*Server-rendered views found — keep, replace, or go API-only?*

| App | Unknown |
|-----|---------|
| #2 PlantsByWebSphere | JSF + XHTML — keep JSF on Spring Boot? Thymeleaf? REST API + SPA? |
| #3 DayTrader | JSP servlets — server-rendered or SPA? |
| #5 Struts Examples | Struts tags + JSP + Tiles — Thymeleaf? Keep JSP? |
| #6 Legacy Cycle Store | Razor views — stay Razor? Move to Blazor? API-only? |
| #7 BookCatalog | Razor CRUD views — server-rendered or rebuild? |
| #8 Flask-to-FastAPI | Jinja2 templates — keep server-rendered or go API-only? |
| #9 NodeJS Ecommerce | EJS templates — NestJS + templates or API + React? |

**Pattern**: This is the most frequent unknown (7/10 apps). The LLM always detects the view technology but never knows whether to preserve server-rendering or switch to API-only. This is always a user/architecture decision.

---

### Category E: Authentication / Security Migration
*Auth exists but the target mechanism is unclear.*

| App | Unknown |
|-----|---------|
| #2 PlantsByWebSphere | BASIC auth with `basicRegistry` — Spring Security? OAuth? |
| #4 IBM App Mod | WebSphere security APIs — Liberty equivalents? |
| #6 Legacy Cycle Store | `SqlMembershipProvider` — ASP.NET Core Identity? External IdP? |
| #8 Flask-to-FastAPI | Auth0 via Authlib Flask integration — JWT? Sessions? |
| #9 NodeJS Ecommerce | Passport.js local + Facebook OAuth — keep Facebook? JWT? |

**Pattern**: The LLM detects auth mechanisms but the replacement is always unclear because it depends on organizational standards (do you use Keycloak? Azure AD? Auth0?) and architecture (stateful sessions vs stateless JWT).

---

### Category F: Database / Persistence Decisions
*DB is detected but production target and migration strategy are unclear.*

| App | Unknown |
|-----|---------|
| #1 Coolstore | What database in production? (Flyway present but no datasource config) |
| #2 PlantsByWebSphere | Derby embedded — clearly not production. What's the real DB? |
| #3 DayTrader | Derby embedded — what production DB? |
| #6 Legacy Cycle Store | EDMX (Database-First EF5) — EF Core has no EDMX. Code-First conversion? |
| #7 BookCatalog | LocalDB — keep SQL Server or switch? |
| #9 NodeJS Ecommerce | MongoDB — keep or switch to PostgreSQL? |

**Pattern**: Dev/embedded databases (Derby, LocalDB, H2) are always detected but the production target is never in the code. ORM migration strategy (EF5 EDMX → EF Core Code-First, EclipseLink → Hibernate) depends on schema complexity the LLM can see but the decision is the user's.

---

### Category G: Scope / Architecture Decisions
*How much of the app to migrate and whether to restructure.*

| App | Unknown |
|-----|---------|
| #1 Coolstore | Keep monolith or split into microservices? |
| #3 DayTrader | EAR with 3 modules — single deployable or split? |
| #5 Struts Examples | 40+ modules — migrate all or subset? |
| #7 BookCatalog | Simple CRUD — is a formal migration even warranted? |
| #10 Struts Legacy | Code is already migrated — what work remains? |

**Pattern**: The LLM can assess app size and complexity but cannot decide scope. Multi-module projects always need a "which parts?" question. Very small or already-migrated apps need a "is this worth it?" question.

---

## Part 2: What the LLM DID Figure Out (no question needed)

Across all 10 apps, the LLM consistently figured out these things **without asking**:

| What | How |
|------|-----|
| **Primary language and file count** | File extension counting |
| **Build tool** | Presence of pom.xml, *.csproj, package.json, requirements.txt |
| **Framework version** | Dependencies in build manifest |
| **App server / runtime** | Config files (server.xml, web.xml, Web.config) and vendor-specific imports |
| **Architecture layers** | Directory structure (model/, service/, controller/, rest/) |
| **Entity / data model** | JPA annotations, EF models, Mongoose schemas |
| **Dependency list** | Build manifest parsing |
| **Vendored / local dependencies** | lib/, vendor/ directory scanning |
| **Java version / .NET version** | Compiler settings in build manifest |
| **Packaging type** | WAR, EAR, JAR, DLL from build config |

**Key insight**: The LLM is reliable at **detection** (what IS the app) but needs to ask about **decisions** (what SHOULD it become).

---

## Part 3: Question Patterns (what the LLM consistently asked)

Across all 10 apps, these question types appeared repeatedly:

### Always Asked (9-10/10 apps)
1. **View layer strategy** — "Keep server-rendered or go API-only?" (7/10)
2. **Target framework version** — "Which specific version?" (6/10)
3. **Database target** — "What database in production?" (6/10)

### Frequently Asked (5-8/10 apps)
4. **Auth migration** — "What auth mechanism in the target?" (5/10)
5. **Proprietary dependency handling** — "What does this do and can it be replaced?" (4/10)
6. **Messaging replacement** — "What broker/pattern in the target?" (3/10)
7. **Scope** — "All modules or subset?" (3/10)

### Sometimes Asked (2-4/10 apps)
8. **Test coverage** — "No tests found, add before migrating?" (4/10)
9. **Session/state management** — "Stateful to stateless?" (3/10)
10. **Deployment target** — "Containers? Cloud? On-prem?" (2/10)

### Rarely Asked (1/10 apps)
11. **"Is this migration needed?"** — Only #10 (already migrated)
12. **Existing migration reference** — Only #8 (fastapi-webapp/ already exists)

---

## Part 4: Summary — 7 Question Categories for the Skill

Based on the patterns above, the questionnaire skill should always probe these 7 areas:

| # | Category | When to ask | When to skip |
|---|----------|-------------|-------------|
| 1 | **Target clarification** | Prompt says a category ("modernize", "cloud-native") not a specific tech | Prompt names exact framework + version |
| 2 | **Proprietary dependency handling** | Vendored jars, system-scoped deps, commercial components detected | No unknown dependencies found |
| 3 | **View layer strategy** | Server-rendered views found (JSP, JSF, Razor, EJS, Jinja2) | App is API-only or has no views |
| 4 | **Auth/security migration** | Any auth mechanism detected | No auth in the app |
| 5 | **Database/persistence target** | Dev/embedded DB found, or ORM migration has breaking changes | Production DB config is explicit and target ORM is clear |
| 6 | **Messaging/async replacement** | JMS, MDB, message queues, event buses detected | No async patterns |
| 7 | **Scope and architecture** | Multi-module project, or app could reasonably not be migrated | Single-module with clear migration path |
