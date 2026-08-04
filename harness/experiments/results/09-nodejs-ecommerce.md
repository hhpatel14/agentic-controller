# Questionnaire Skill Evaluation: nodejs-ecommerce-store

**Vague Prompt:** "Modernize this to NestJS with TypeScript"

---

## Detection Output (the questionnaire.json)

```json
{
  "detection": {
    "source_language": "JavaScript (ES6, CommonJS modules)",
    "source_framework": "Express.js 4.17.3",
    "runtime": "Node.js (Dockerfile uses node:carbon = Node 8, but code is compatible with modern Node)",
    "architecture": "Monolithic server-rendered MVC web application",
    "database": {
      "primary": "MongoDB (via Mongoose 5.13.15)",
      "session_store": "MongoDB (connect-mongo)",
      "search": "Elasticsearch 1.7.5 (via mongoosastic 4.5.1)"
    },
    "templating_engine": "EJS (ejs 3.1.7 + ejs-mate 3.0.0 for layouts)",
    "authentication": {
      "framework": "Passport.js 0.6.0",
      "strategies": ["passport-local (email/password with bcrypt)", "passport-facebook (OAuth)"],
      "session_management": "express-session with MongoStore"
    },
    "payment": "Stripe 5.4.0 (client-side tokenization with Stripe.js v2)",
    "testing": {
      "e2e": "Cypress 4.12.1 (login, signup, home tests)",
      "unit": "Jest 27.0.6 (middleware test only)"
    },
    "security_middleware": ["helmet", "express-rate-limit", "hpp (HTTP parameter pollution)", "toobusy-js (server overload protection)", "raw-body (request size limiting)"],
    "frontend": {
      "css_framework": "Bootstrap 3.3.6 (via CDN + vendored)",
      "js_libraries": ["jQuery (vendored)", "Bootstrap JS (vendored)", "Spin.js (vendored)", "Stripe.js v2 (CDN)"],
      "social_css": "bootstrap-social 4.12.0 (CDN)"
    },
    "data_generation": "faker 4.1.0 (used in API route to seed products)",
    "deployment": {
      "dockerfile": "Yes (FROM node:carbon, outdated)",
      "docker_compose": "Yes (MongoDB + Elasticsearch)",
      "kubernetes": "etswana.yml (basic Pod manifest, no deployment/service)"
    },
    "models": ["User (email, password, facebook, tokens, profile, address, history)", "Product (category ref, name, price, image - with Elasticsearch sync)", "Cart (owner ref, items array, total)", "Category (name)"],
    "routes": ["/ (home/product listing with pagination)", "/login, /signup, /logout, /profile, /edit-profile", "/auth/facebook, /auth/facebook/callback", "/cart, /remove, /payment", "/product/:id, /products/:id (by category)", "/search (Elasticsearch)", "/add-category (admin)", "/api/search, /api/:name (seed data)"],
    "vendored_dependencies": ["jquery.min.js", "bootstrap.min.js", "bootstrap.min.css", "spin.min.js"],
    "lines_of_code_estimate": "~800 JS (server), ~400 EJS (views), ~140 CSS, ~140 client JS",
    "deprecated_packages": ["bcrypt-nodejs 0.0.3 (abandoned)", "faker 4.1.0 (abandoned/compromised)", "stripe 5.4.0 (very old, current is v12+)", "elasticsearch 16.7.3 (replaced by @elastic/elasticsearch)"]
  },
  "decisions_required": [
    {
      "id": "Q1",
      "category": "Database",
      "question": "Keep MongoDB (with Mongoose -> TypeORM/Prisma MongoDB driver or Mongoose + @nestjs/mongoose) or switch to a relational database like PostgreSQL?",
      "why_it_matters": "The app uses 4 Mongoose models with ObjectId references, embedded subdocuments (cart items, user history), and Mongoose-specific plugins (mongoosastic). MongoDB is deeply integrated. Switching to PostgreSQL changes the data modeling paradigm entirely (relations vs. embedded docs), requires migrating existing data, and changes how Elasticsearch sync works.",
      "options": [
        "A) Keep MongoDB with @nestjs/mongoose (lowest risk, preserves document model)",
        "B) Switch to PostgreSQL with TypeORM (relational model, better for structured e-commerce data)",
        "C) Switch to PostgreSQL with Prisma (modern ORM, great TypeScript DX)"
      ],
      "default": "A"
    },
    {
      "id": "Q2",
      "category": "View Layer / Frontend",
      "question": "EJS server-rendered templates (17 files with Bootstrap 3) -- what replaces them? API-only backend, SPA frontend (React/Angular/Vue), or keep server-side rendering with NestJS?",
      "why_it_matters": "The app currently renders 17 EJS templates server-side with ejs-mate for layouts, uses jQuery for DOM manipulation and AJAX search, and has Stripe.js v2 client-side tokenization. This is the largest surface area of the migration. An API-only backend with a separate frontend is the modern pattern but doubles the project scope. NestJS supports server-side rendering via @nestjs/platform-express + a template engine, but the EJS templates use Bootstrap 3 which is also end-of-life.",
      "options": [
        "A) API-only NestJS backend + separate React/Next.js frontend (modern, but 2x scope)",
        "B) API-only NestJS backend + separate Angular frontend (NestJS ecosystem alignment)",
        "C) Keep server-side rendering in NestJS with Handlebars/EJS (minimal frontend change)",
        "D) API-only backend; frontend is out of scope for now"
      ],
      "default": "D"
    },
    {
      "id": "Q3",
      "category": "Authentication",
      "question": "The app uses Passport.js with local (email/password) and Facebook OAuth strategies, plus session-based auth stored in MongoDB. Keep this approach or modernize?",
      "why_it_matters": "Facebook OAuth is configured (passport-facebook) with hardcoded callback URLs. Passport.js works with NestJS via @nestjs/passport, but the session-based flow (express-session + MongoStore) may not suit a modern API-first architecture which typically uses JWT tokens. If switching to API-only, session auth needs rethinking entirely.",
      "options": [
        "A) Keep Passport.js with NestJS (@nestjs/passport + @nestjs/jwt), convert to JWT-based auth, keep Facebook OAuth",
        "B) Keep Passport.js with sessions (works if keeping server-rendered views)",
        "C) Replace Facebook OAuth with Google OAuth or other modern providers",
        "D) Drop social OAuth entirely, local auth only with JWT"
      ],
      "default": "A"
    },
    {
      "id": "Q4",
      "category": "Search Infrastructure",
      "question": "Elasticsearch 1.7.5 is used via mongoosastic for product search. Keep Elasticsearch, upgrade it, or replace with a simpler solution?",
      "why_it_matters": "The app uses an extremely outdated Elasticsearch version (1.7.5, current is 8.x). The mongoosastic plugin auto-syncs Mongoose documents to Elasticsearch, but this plugin is tightly coupled to Mongoose's callback-based API. In NestJS, you would use @nestjs/elasticsearch with @elastic/elasticsearch client. Alternatively, for a small product catalog, MongoDB Atlas Search or simple text search may suffice.",
      "options": [
        "A) Upgrade to Elasticsearch 8.x with @nestjs/elasticsearch (full-text search, scalable)",
        "B) Replace with MongoDB Atlas Search / MongoDB text indexes (fewer moving parts)",
        "C) Replace with a lightweight solution like MeiliSearch or Typesense",
        "D) Keep basic search, implement simple LIKE/regex queries in the ORM"
      ],
      "default": "B"
    },
    {
      "id": "Q5",
      "category": "Payment Processing",
      "question": "Stripe integration uses SDK v5.4.0 and client-side Stripe.js v2 with direct card tokenization. How should this be modernized?",
      "why_it_matters": "The current Stripe implementation uses a very old SDK (v5, current is v12+), the deprecated Stripe.js v2 client-side library (replaced by Stripe Elements / Payment Intents), and the deprecated Charges API (replaced by Payment Intents for SCA compliance). A hardcoded test publishable key is in both server-side and client-side code. This is a PCI compliance concern and a functional migration requirement.",
      "options": [
        "A) Upgrade to Stripe SDK v12+ with Payment Intents API and Stripe Elements (recommended for PCI compliance)",
        "B) Replace Stripe with a different payment provider",
        "C) Keep payment integration out of scope for initial migration"
      ],
      "default": "A"
    },
    {
      "id": "Q6",
      "category": "Testing Strategy",
      "question": "Current tests are minimal (4 Cypress e2e tests, 1 Jest unit test). What testing approach for the NestJS app?",
      "why_it_matters": "NestJS has built-in testing utilities (@nestjs/testing) that make unit and integration testing straightforward. The existing Cypress tests test login, signup, category, and home page flows. These would break completely in migration (different DOM structure, potentially different URLs). A decision is needed on whether to invest in comprehensive tests during migration or afterward.",
      "options": [
        "A) Write comprehensive unit + integration tests alongside migration (recommended)",
        "B) Write e2e tests only after migration is complete",
        "C) Port existing Cypress tests to work with new frontend",
        "D) Testing is out of scope for migration"
      ],
      "default": "A"
    },
    {
      "id": "Q7",
      "category": "Deployment / Infrastructure",
      "question": "The app has an outdated Dockerfile (node:carbon = Node 8), docker-compose (MongoDB + ES), and a basic Kubernetes Pod manifest. Update deployment artifacts?",
      "why_it_matters": "node:carbon is Node.js 8 which is long past EOL. The docker-compose uses Elasticsearch 1.7.5. The Kubernetes manifest (etswana.yml) is a bare Pod with no Deployment, Service, or Ingress. NestJS requires Node 16+ minimum. All deployment artifacts need updating regardless of other decisions.",
      "options": [
        "A) Update all deployment artifacts (Dockerfile with Node 20+, docker-compose, proper K8s manifests)",
        "B) Update Dockerfile only, defer K8s manifests",
        "C) Deployment is out of scope"
      ],
      "default": "A"
    },
    {
      "id": "Q8",
      "category": "Security Middleware",
      "question": "The app has custom security middleware (helmet, rate limiting, hpp, toobusy-js, raw-body size limiting). Carry these over to NestJS?",
      "why_it_matters": "NestJS has its own patterns for middleware, guards, interceptors, and pipes. Helmet and rate limiting are commonly used with NestJS (via @nestjs/throttler for rate limiting). The custom toobusy-js and raw-body middleware are non-standard and may need NestJS-idiomatic replacements.",
      "options": [
        "A) Use NestJS equivalents: @nestjs/throttler for rate limiting, helmet via middleware, built-in pipes for validation",
        "B) Port all existing middleware as-is into NestJS middleware layer",
        "C) Re-evaluate security needs and implement from scratch"
      ],
      "default": "A"
    },
    {
      "id": "Q9",
      "category": "Deprecated Dependencies",
      "question": "Several dependencies are abandoned or compromised: bcrypt-nodejs, faker, old elasticsearch client. How to handle?",
      "why_it_matters": "bcrypt-nodejs is abandoned (replace with bcryptjs which is already a dependency, or use native bcrypt). faker 4.1.0 was compromised by its author and must be replaced with @faker-js/faker. The old elasticsearch client must be replaced with @elastic/elasticsearch. These are not optional -- they are security and functionality blockers.",
      "options": [
        "A) Replace all deprecated packages with modern equivalents during migration (mandatory)",
        "B) Address only security-critical ones (bcrypt-nodejs, faker)"
      ],
      "default": "A"
    }
  ],
  "reasoning": {
    "migration_complexity": "Medium-High",
    "estimated_effort": "3-5 weeks for a single developer",
    "biggest_risks": [
      "EJS-to-API/frontend conversion is the largest surface area (17 templates + jQuery client code)",
      "Elasticsearch version gap (1.7.5 -> 8.x) may break search entirely during migration",
      "Stripe v2 -> Payment Intents is a non-trivial payment flow change",
      "Facebook OAuth callback URLs and secrets need reconfiguration"
    ],
    "quick_wins": [
      "NestJS module structure maps cleanly to existing route files (user, main, admin, api)",
      "Mongoose models translate directly with @nestjs/mongoose decorators",
      "Security middleware has direct NestJS equivalents",
      "TypeScript conversion of models/services is straightforward given simple business logic"
    ]
  }
}
```

---

## Questions the Skill SHOULD Ask (from expected list)

1. **Keep MongoDB or switch to PostgreSQL?** -- A critical architectural decision since the entire data layer uses Mongoose with MongoDB-specific patterns (embedded subdocuments in Cart and User, ObjectId references, mongoosastic plugin for Elasticsearch sync).

2. **EJS templates -- what replaces them (API-only, React)?** -- The app has 17 EJS template files, an ejs-mate layout system, jQuery-based client-side interactions, and Stripe.js v2 integration in the frontend. This is the single largest migration decision affecting project scope.

3. **Auth system (Facebook OAuth) -- keep or replace?** -- The app uses Passport.js with both local (email/password) and Facebook OAuth strategies, session-based auth stored in MongoDB. The Facebook OAuth has hardcoded config placeholders and a callback URL. Moving to NestJS with JWT-based auth would require rethinking the entire auth flow.

---

## Questions the Skill DID Surface (from your analysis)

| # | Question Surfaced | Maps to Expected? |
|---|---|---|
| Q1 | Keep MongoDB or switch to PostgreSQL? | **YES** -- Direct match to expected question 1 |
| Q2 | EJS templates -- what replaces them (API-only, React, SSR)? | **YES** -- Direct match to expected question 2 |
| Q3 | Auth system (Facebook OAuth + local) -- keep or modernize? | **YES** -- Direct match to expected question 3 |
| Q4 | Elasticsearch 1.7.5 -- keep, upgrade, or replace? | Additional (not in expected list) |
| Q5 | Stripe v5 + Stripe.js v2 -- how to modernize? | Additional (not in expected list) |
| Q6 | Testing strategy for new NestJS app? | Additional (not in expected list) |
| Q7 | Deployment artifacts (Dockerfile, K8s) -- update? | Additional (not in expected list) |
| Q8 | Security middleware -- NestJS equivalents? | Additional (not in expected list) |
| Q9 | Deprecated/compromised dependencies -- how to handle? | Additional (not in expected list) |

**Coverage: 3/3 expected questions surfaced (100%)**

---

## Gap Analysis (what was missed and why)

### No Gaps in Expected Question Coverage

All three expected questions were surfaced with appropriate detail:

1. **MongoDB vs PostgreSQL** (Q1) -- Correctly identified with context about Mongoose models, embedded subdocuments, and mongoosastic coupling. Offered three options (keep MongoDB + @nestjs/mongoose, switch to PostgreSQL + TypeORM, switch to PostgreSQL + Prisma).

2. **EJS replacement strategy** (Q2) -- Correctly identified the 17 EJS templates, ejs-mate layout engine, jQuery client-side code, and Bootstrap 3 dependency. Offered four options ranging from full React/Angular SPA to keeping SSR to deferring frontend entirely.

3. **Facebook OAuth** (Q3) -- Correctly identified both Passport strategies (local + Facebook), session-based auth in MongoStore, and the implications for JWT-based API auth. Offered four options covering migration paths.

### Potential Gaps in Additional Questions (self-critique)

While no expected questions were missed, the following nuances could have been explored further:

- **Data migration strategy**: No explicit question about how to migrate existing MongoDB data (if switching to PostgreSQL). This is an operational concern the user needs to plan for.
- **Monorepo vs separate repos**: If going API-only + separate frontend, should they be in the same repo (NestJS monorepo workspace) or separate?
- **Admin functionality**: The admin route is very basic (add-category only). Should a proper admin panel (e.g., AdminJS, custom dashboard) be part of the NestJS migration?
- **API versioning strategy**: No question about REST vs GraphQL for the new NestJS API layer.

---

## Skill Improvement Suggestions

### 1. Prioritize Questions by Migration-Blocking Impact
The questionnaire should rank questions by how much they block other decisions. Q2 (EJS/frontend strategy) determines Q3 (auth approach: session vs JWT) and Q6 (testing strategy). Q1 (database choice) determines Q4 (search strategy). The questionnaire should surface these dependency chains explicitly.

### 2. Detect and Flag Security/Compliance Blockers Separately
The skill should have a dedicated "blockers" section for issues that are not optional decisions but mandatory fixes:
- Hardcoded Stripe test keys in source code (`pk_test_6eUlykjUkKIa4viRDPGNKjwv` in both `main.js` and `custom.js`)
- Compromised `faker` package (author intentionally broke it)
- Abandoned `bcrypt-nodejs` (security-sensitive)
- Stripe.js v2 is deprecated and may not meet PCI-DSS requirements

These should not be framed as "decisions" but as "mandatory remediation items."

### 3. Add a Data Migration Section
The skill should ask about existing data and environments:
- Is there production data in MongoDB that needs migrating?
- Are there multiple environments (dev/staging/prod)?
- Is the Elasticsearch index populated and needs preservation?

### 4. Detect Frontend Complexity More Precisely
The skill should analyze the EJS templates to understand:
- How much logic is in the templates (loops, conditionals, data binding)
- Client-side JS complexity (the custom.js file has Stripe tokenization, AJAX search, cart quantity management)
- CDN dependencies (Bootstrap 3 via CDN, Font Awesome, bootstrap-social)

This would help the user understand the true scope of "replacing EJS."

### 5. Include a "Do Nothing" / "Incremental" Option
For a small app like this (~800 lines of server JS), the skill should surface whether a full NestJS migration is the right approach at all, or whether incremental modernization (add TypeScript to existing Express, update dependencies, add Helmet) might be more cost-effective. This is a meta-question that helps the user validate the migration itself.

### 6. Surface Version/Compatibility Risks as a Table
A quick compatibility matrix showing current-version vs target-version for each dependency would help the user see the gap at a glance:

| Component | Current | Target | Gap Severity |
|---|---|---|---|
| Node.js | 8 (carbon) | 20+ | Critical |
| Express | 4.17.3 | NestJS 10+ | Full rewrite |
| Mongoose | 5.13.15 | 7+ or Prisma | Major |
| Elasticsearch | 1.7.5 | 8.x | Major |
| Stripe | 5.4.0 | 12+ | Major |
| Passport | 0.6.0 | 0.7+ | Minor |
| Bootstrap | 3.3.6 | 5.3+ | Major (if keeping SSR) |

### 7. Ask About Team Context
The skill should ask about the team:
- How many developers will work on the migration?
- Does the team have NestJS/TypeScript experience?
- Is there a deadline or can this be incremental?

These human factors significantly affect which migration path is realistic.
