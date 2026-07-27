# Ingest Loop Protocol — Migration Interview Step

## 1. Overview & Trigger

### When to invoke
Any time a migration prompt is received — e.g. "migrate this application from Java EE to Quarkus." This step runs **after** detect (graphify) and **before** plan.

### Pre-requisite
`graphify-out/` must exist in the project directory. If it doesn't, run `mig-graphify` first to build the knowledge graph.

### Purpose
Gather enough context from the user — through a structured, batched interview loop — to produce a high-quality migration plan. The LLM reads the knowledge graph to pre-fill what it already knows, then asks about what it doesn't: custom libraries, design preferences, replacement strategies, and constraints.

### Pipeline position
```
detect (graphify) ──► INGEST (this step) ──► plan ──► execute ──► verify
                      ▲                       │
                      │  reads graphify-out/   │  consumes ingest-context.md
                      │  interviews user       │
                      │  produces context doc  │
```

---

## 2. Rubric Categories

### Required (must reach ≥80% confidence before planning)

| # | Category | What to gather | Pre-fillable from graphify? |
|---|----------|---------------|---------------------------|
| R1 | **Source stack confirmation** | Language, framework, app server, version, build tool | Yes — confirm what graphify detected |
| R2 | **Target stack details** | Target framework, version, specific extensions or features desired | No — user must specify |
| R3 | **Custom/proprietary dependencies** | Local jars, internal libraries, vendor-specific APIs the LLM has no training data for | Partially — graphify detects dependency nodes with unknown ecosystem mapping; user explains purpose and migration strategy |
| R4 | **Messaging & async patterns** | JMS topics/queues, MDBs, event buses, scheduled tasks → what replacement strategy (reactive messaging, CDI events, Kafka, direct calls) | Partially — graphify detects JMS/MDB usage; user picks replacement |
| R5 | **Data persistence strategy** | Database type (current and target), ORM approach, schema migration tool, dev vs prod DB strategy | Partially — graphify detects persistence.xml, DataSource config, entity classes |
| R6 | **Entry point & dependency order** | Application entry point, which files/modules to migrate first based on dependency graph topology | Yes — from graph node degree, leaf-to-hub ordering |

### Optional (ask if interview rounds remain after required categories are filled)

| # | Category | What to gather |
|---|----------|---------------|
| O1 | **Authentication & authorization** | Current auth mechanism (JAAS, Keycloak, custom), target auth approach |
| O2 | **External integrations** | REST clients, SOAP services, LDAP, mail, third-party APIs |
| O3 | **Deployment target** | Container platform, CI/CD pipeline, environment constraints |
| O4 | **Design preferences** | Keep monolith vs split, naming conventions, code style, module boundaries |
| O5 | **Team context** | Team size, familiarity with target stack, training budget |
| O6 | **Timeline & constraints** | Hard deadlines, phased rollout, budget, regulatory requirements |

---

## 3. Loop Protocol

### ROUND 0 — SCAN (no user interaction)

```
1. Read graphify-out/GRAPH_REPORT.md
2. Read graphify-out/graph.json (nodes, edges, communities)
3. Pre-fill categories:
   - R1 (Source stack): extract language, framework, app server from graph metadata
   - R5 (Data persistence): extract DB config, entity classes, ORM from graph
   - R6 (Entry point): identify entry point from graph topology
     (node with @ApplicationPath, main(), or highest fan-out)
   - R6 (Dependency order): compute migration order from leaf nodes → hub nodes
4. Identify unknowns:
   - Nodes referencing packages not in any known ecosystem → flag for R3
   - JMS/MDB/event annotations detected → flag for R4
   - System-scoped or local jar dependencies → flag for R3
5. Compute initial coverage score for all 6 required categories
6. Present findings summary to user:
   "I've scanned your application. Here's what I found: [summary].
    I have some questions before I can create a solid migration plan."
```

### ROUND 1..N — INTERVIEW (batched, rubric-driven)

```
FOR each round (up to max 6):

  1. PICK: Select the 3-4 lowest-confidence required categories
  2. FORMULATE: Generate one specific question per selected category
     - If graphify pre-filled partial data, frame as confirmation:
       "I see you're using JMS with a topic called 'orders'. 
        For Quarkus, would you prefer CDI Events (simplest for in-process),
        SmallRye Reactive Messaging, or Kafka?"
     - If no data exists, ask directly:
       "I found a dependency on audit-logging-library-1.0.0.jar which I'm
        not familiar with. What does it do, and should we find a replacement
        or remove it?"
  3. PRESENT: Show questions as a numbered batch (max 4)
  4. RECEIVE: User responds (may answer all, skip some, say "assume X")
  5. UPDATE:
     - Parse answers into category scores
     - For skipped questions: create assumption with MEDIUM confidence
     - For "I don't know": create assumption with LOW confidence
     - For concrete answers: set category to HIGH confidence
  6. LOG: Append Q&A to the raw log section of ingest-context.md
  7. CHECK TERMINATION: evaluate exit conditions (see below)
  
  IF all required categories ≥ 80%:
    Proceed to FINAL ROUND
  ELSE IF round count = max rounds (6):
    Document remaining gaps as assumptions, proceed to OUTPUT
  ELSE:
    Continue to next round
```

### FINAL ROUND — FREE-FORM (1-2 questions, open-ended)

```
After all required categories reach ≥ 80%, ask:

"Before I start planning, is there anything else I should know?
 For example:
 - Custom build steps or deployment scripts
 - Known gotchas or workarounds in the current codebase
 - Things that surprised previous developers
 - Specific features you want to add or remove during migration"

This round catches edge cases the rubric missed.
If the user says "no" or "that's it" → proceed to OUTPUT.
If the user provides new info → log it under Design Decisions, update
relevant category scores, proceed to OUTPUT.
```

### TERMINATION — exit the loop when ANY condition is met

| Condition | What happens |
|-----------|-------------|
| **(a) Coverage threshold** | All 6 required categories score ≥ 80% confidence AND free-form round completed |
| **(b) Max rounds reached** | 6 interview rounds completed — summarize what's known, document all unknowns as assumptions |
| **(c) User says stop** | User explicitly says "enough", "go ahead", "start planning", "skip the rest" — document remaining unknowns as assumptions |

---

## 4. Coverage Scoring

Each required category (R1–R6) carries a confidence score from 0% to 100%.

### Scoring components

| Component | Points | Criteria |
|-----------|--------|----------|
| **User answer** | 0–40% | 0 = not asked yet, 20 = vague/partial answer, 40 = concrete specific answer |
| **Graphify data** | 0–30% | 0 = no graph data, 15 = partial detection, 30 = confirmed by graph topology |
| **No ambiguity** | 0–30% | 0 = multiple interpretations possible, 15 = minor assumptions needed, 30 = fully unambiguous |

### Score thresholds

| Score | Meaning | Action |
|-------|---------|--------|
| **100%** | Fully resolved with concrete details | No questions needed |
| **80–99%** | Answered, minor details assumed | Good enough to plan — document assumptions |
| **50–79%** | Partially answered, significant gaps | Must ask more questions |
| **0–49%** | Not addressed or highly ambiguous | Priority for next round's questions |

### After each round

```
COVERAGE REPORT (internal, shown to user as a brief status):

Category         Score   Status
─────────────────────────────────
R1 Source stack    90%   ✓ Confirmed
R2 Target stack    40%   ✗ Need details
R3 Custom deps     60%   ~ Partially known
R4 Messaging       30%   ✗ Not discussed
R5 Data persist    80%   ✓ Pre-filled + confirmed
R6 Entry point    100%   ✓ From graph

Overall: 4/6 categories below threshold
Action: Continue interview — ask about R2, R4, R3
```

---

## 5. Assumption Handling

When a category cannot be fully resolved — user skips, says "I don't know", or max rounds are reached:

### Step 1: State the assumption explicitly
> "I'll assume CDI Events as the JMS replacement since the messaging is in-process only."

### Step 2: Tag with confidence level

| Level | When to use | Risk in plan |
|-------|------------|-------------|
| **HIGH** | LLM is confident based on graphify data + common patterns; user didn't contradict | Low risk — mention in plan as a default |
| **MEDIUM** | Reasonable guess but multiple valid options exist; user said "whatever makes sense" | Medium risk — include as a decision point in plan |
| **LOW** | No data, no user input, pure guess | High risk — flag prominently in plan, may need revisiting |

### Step 3: Log in the Assumptions table of ingest-context.md

### Step 4: Surface in plan
- HIGH assumptions: footnote in the relevant plan section
- MEDIUM assumptions: explicit callout box in the plan
- LOW assumptions: listed as risks/blockers that need user decision before execution

---

## 6. Output Format

The ingest loop produces a single file: **`ingest-context.md`** in the project's working directory (alongside `graphify-out/`).

```markdown
# Migration Ingest Context

| Field | Value |
|-------|-------|
| Generated | YYYY-MM-DD HH:MM |
| Source application | <absolute path> |
| Migration direction | <source> → <target> |
| Rounds completed | N of 6 max |
| Termination reason | coverage threshold / max rounds / user requested |

---

## Source Stack
- **Language**: Java 8
- **Framework**: Java EE 7 (JAX-RS, EJB, JPA, JMS, CDI)
- **App server**: WebLogic / WildFly (detected WebLogic stubs)
- **Build tool**: Maven
- **Packaging**: WAR
- **Confirmed by**: graphify scan + user confirmation in Round 1

## Target Stack
- **Framework**: Quarkus 3.x
- **Key extensions**: resteasy-jackson, hibernate-orm, flyway, arc
- **Java version**: 17
- **Packaging**: JAR (uber-jar)
- **Specified by**: user in Round 1

## Custom Dependencies

| Dependency | Location | Purpose | Migration Strategy | Confidence |
|-----------|----------|---------|-------------------|-----------|
| audit-logging-library-1.0.0.jar | lib/ | File-based audit logging for orders | Remove — replace with java.util.logging | MEDIUM |

## Messaging Patterns

| Current Pattern | Files | Replacement | Rationale |
|----------------|-------|-------------|-----------|
| JMS Topic (topic/orders) | ShoppingCartOrderProcessor, OrderServiceMDB, InventoryNotificationMDB | CDI Events (@Observes) | In-process only, no external broker needed |
| @MessageDriven MDB | OrderServiceMDB | @ApplicationScoped + @Observes | Quarkus has no EJB container |
| WebLogic JNDI JMS | InventoryNotificationMDB | Direct CDI injection | JNDI not available in Quarkus |

## Data Persistence

| Aspect | Current | Target |
|--------|---------|--------|
| Database | JNDI DataSource (java:jboss/datasources/CoolstoreDS) | H2 (dev), PostgreSQL (prod) |
| ORM | JPA 2.1 (Hibernate) | Quarkus Hibernate ORM |
| Schema migration | Flyway 4.x (manual startup bean) | Quarkus Flyway extension (auto) |
| EntityManager | CDI producer in Resources.java | Quarkus built-in @Inject |

## Entry Point & Migration Order

**Entry point**: `RestApplication.java` (@ApplicationPath("/services"))

**Migration order** (from graph leaf nodes to hubs):
1. Model classes (POJOs, entities) — no dependencies on services
2. Utility classes (Transformers, Producers) — used by services
3. Service classes (leaf services first: ShippingService, PromoService)
4. Service classes (hub services: CatalogService, ProductService, OrderService)
5. Service classes (orchestrators: ShoppingCartService, ShoppingCartOrderProcessor)
6. Message-driven beans (OrderServiceMDB, InventoryNotificationMDB)
7. REST endpoints (CartEndpoint, OrderEndpoint, ProductEndpoint)
8. Configuration (pom.xml, application.properties, persistence.xml → delete)
9. Cleanup (remove WebLogic stubs, StartupListener, DataBaseMigrationStartup)

## Design Decisions
- Keep as monolith (no microservices split)
- Replace all javax.* with jakarta.* (Quarkus 3.x requirement)
- Remove WebLogic-specific code entirely (not behind a flag)
- Use in-process eventing (CDI Events) rather than external broker

## Assumptions

| # | Assumption | Confidence | Category | Risk if wrong |
|---|-----------|-----------|----------|--------------|
| A1 | H2 for dev, PostgreSQL for prod | HIGH | Data Persistence | Would need different JDBC driver + config |
| A2 | CDI Events replace JMS (no external broker needed) | MEDIUM | Messaging | If external consumers exist, need Kafka/AMQP instead |
| A3 | No authentication changes needed | MEDIUM | Auth (optional) | May need quarkus-oidc if auth exists |
| A4 | audit-logging-library can be removed | LOW | Custom Deps | If audit trail is compliance-required, need alternative |

---

## Raw Q&A Log

### Round 0 — Scan
Graphify detected: Java EE 7 WAR, WebLogic stubs, JPA entities, JMS messaging,
Flyway migrations, 3 REST endpoints, 1 proprietary jar dependency.

### Round 1
**Q1** [R2]: What version of Quarkus are you targeting, and are there specific
extensions you need (e.g., reactive, security, health checks)?
**A1**: Quarkus latest, just the basics for now.

**Q2** [R3]: I found `audit-logging-library-1.0.0.jar` in `lib/` — it's used by
OrderService for file-based audit logging. What does this library do, and should
we find a Quarkus-compatible replacement or remove it?
**A2**: It's an internal library, just remove it and use standard logging.

**Q3** [R4]: Your app uses JMS with a topic called `topic/orders` for order
processing. For Quarkus, do you want: (a) CDI Events (simplest, in-process only),
(b) SmallRye Reactive Messaging, or (c) Kafka/AMQP for external consumers?
**A3**: CDI Events is fine, it's all in-process.

### Round 2
...
```

---

## 7. Guardrails

### Hard limits

| Guardrail | Value | Rationale |
|-----------|-------|-----------|
| Max interview rounds | 6 | Prevents infinite loops; 6 rounds × 4 questions = 24 questions max |
| Max questions per round | 4 | Prevents user fatigue; keeps each round answerable in 1-2 minutes |
| Max total questions | 24 | 6 rounds × 4 questions — absolute ceiling |

### Quality rules

| Rule | How enforced |
|------|-------------|
| **No duplicate questions** | Track asked questions by category+topic; never re-ask what's already answered |
| **No trivial questions** | If graphify answered it with ≥80% confidence, don't ask — briefly confirm instead |
| **Confirm, don't interrogate** | When graphify pre-filled a category, frame as "I see X — is that right?" not "What is your X?" |
| **Respect user pace** | If user gives short answers, don't push — take what's given, assume the rest |
| **Progressive depth** | Start with broad required categories, narrow to specifics in later rounds |

### Escape hatches

| User says | LLM response |
|-----------|-------------|
| "skip" / "next" | Skip that question, create MEDIUM-confidence assumption, move on |
| "I don't know" | Create LOW-confidence assumption, flag as risk in plan, move on |
| "assume whatever" | Create MEDIUM-confidence assumption based on common patterns, move on |
| "enough" / "go ahead" / "start planning" | End interview immediately, document all unresolved categories as assumptions |
| "go back" / "actually..." | Reopen a previously-answered category, update the score |

### Anti-patterns to avoid

| Anti-pattern | Why it's bad | What to do instead |
|-------------|-------------|-------------------|
| Asking all 24 questions upfront | Overwhelming; user abandons | Batch 3-4 per round, stop when coverage is met |
| Re-asking what graphify already knows | Wastes user time, feels dumb | Pre-fill from graph, confirm briefly |
| Asking about things that won't affect the plan | Scope creep | Stick to categories that directly influence migration decisions |
| Continuing after user says stop | Ignoring user autonomy | Immediately terminate, document assumptions |
| Not showing progress | User doesn't know when it'll end | Show coverage status after each round |
