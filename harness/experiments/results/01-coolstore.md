## Detection Output (the questionnaire.json)

```json
{
  "migration_prompt": "Migrate this to Quarkus",
  "source_application": {
    "name": "coolstore-monolith",
    "groupId": "com.redhat.coolstore",
    "artifactId": "monolith",
    "version": "1.0.0-SNAPSHOT",
    "packaging": "war",
    "java_source_level": "1.8",
    "java_target_level": "1.8"
  },
  "detected_technologies": {
    "framework": "Java EE 7 (javax.* namespace)",
    "application_server": "JBoss EAP / WildFly (with WebLogic compatibility stubs)",
    "build_tool": "Maven",
    "dependencies": [
      { "groupId": "javax", "artifactId": "javaee-web-api", "version": "7.0", "scope": "provided" },
      { "groupId": "javax", "artifactId": "javaee-api", "version": "7.0", "scope": "provided" },
      { "groupId": "org.jboss.spec.javax.jms", "artifactId": "jboss-jms-api_2.0_spec", "version": "2.0.0.Final", "scope": "compile" },
      { "groupId": "org.flywaydb", "artifactId": "flyway-core", "version": "4.1.2", "scope": "compile" },
      { "groupId": "org.jboss.spec.javax.rmi", "artifactId": "jboss-rmi-api_1.0_spec", "version": "1.0.2.Final", "scope": "compile" },
      { "groupId": "com.enterprise", "artifactId": "audit-logging-library", "version": "1.0.0", "scope": "system", "systemPath": "${project.basedir}/lib/audit-logging-library-1.0.0.jar" }
    ],
    "java_ee_apis_used": [
      "EJB (@Stateless, @Stateful, @Singleton, @MessageDriven, @Remote, @Startup)",
      "JPA/Hibernate (persistence.xml, @Entity, EntityManager, CriteriaBuilder)",
      "CDI (@Inject, @Produces, @Dependent, @SessionScoped, beans.xml)",
      "JAX-RS (@Path, @GET, @POST, @DELETE, @Produces, @ApplicationPath)",
      "JMS (JMSContext, Topic, MessageListener, TopicConnection, TopicSession)",
      "JTA (jta-data-source in persistence.xml)",
      "JNDI (InitialContext lookups for EJB Remote and JMS)",
      "RMI (PortableRemoteObject.narrow)",
      "Servlet API (JSP pages, web.xml)",
      "javax.json (Json.createObjectBuilder, JsonReader, JsonWriter)"
    ],
    "messaging": {
      "type": "JMS",
      "destination": "topic/orders",
      "producers": ["ShoppingCartOrderProcessor - publishes cart JSON to JMS Topic via JMSContext"],
      "consumers": [
        "OrderServiceMDB - @MessageDriven, subscribes to topic/orders, persists Order to DB",
        "InventoryNotificationMDB - WebLogic-style manual JNDI/TopicConnection subscriber, checks inventory thresholds"
      ]
    },
    "persistence": {
      "type": "JPA 2.1",
      "persistence_unit": "primary",
      "datasource": "java:jboss/datasources/CoolstoreDS",
      "entities": ["Order", "OrderItem", "CatalogItemEntity", "InventoryEntity"],
      "db_migration": "Flyway 4.1.2 (triggered via @Singleton @Startup EJB)"
    },
    "authentication": {
      "type": "Keycloak",
      "realm": "eap",
      "client": "eap-app",
      "sso_enabled": true,
      "realm_export_present": true
    },
    "frontend": {
      "type": "AngularJS 1.x",
      "package_manager": "Bower",
      "ui_framework": "PatternFly",
      "delivery": "JSP (index.jsp serves the SPA shell)"
    },
    "vendored_libraries": {
      "lib_directory": "lib/",
      "jars_found": [
        "audit-logging-library-1.0.0.jar (24KB, currently referenced in pom.xml as system scope)",
        "audit-logging-library-2.0.0.jar (19KB, present but not referenced)"
      ],
      "usage_in_code": "OrderService.java imports FileSystemAuditLogger, AuditConfiguration, AuditLoggingException"
    },
    "weblogic_artifacts": {
      "stub_classes": [
        "weblogic.application.ApplicationLifecycleListener (abstract base class stub)",
        "weblogic.application.ApplicationLifecycleEvent (event class stub)",
        "weblogic.i18n.logging.NonCatalogLogger (logger stub)"
      ],
      "usage": "StartupListener extends ApplicationLifecycleListener with postStart/preStop lifecycle hooks",
      "jndi_references": [
        "weblogic.jndi.WLInitialContextFactory (in InventoryNotificationMDB)",
        "t3://localhost:7001 (WebLogic provider URL)"
      ]
    },
    "ejb_remote": {
      "interfaces": ["ShippingServiceRemote"],
      "implementations": ["ShippingService (@Stateless @Remote)"],
      "lookup_pattern": "JNDI via WildFlyInitialContextFactory in ShoppingCartService"
    },
    "architectural_pattern": "monolith",
    "tests": "none (maven.test.skip=true, no test sources found)"
  },
  "ambiguities_requiring_user_input": [
    {
      "id": "Q1",
      "category": "target_version",
      "question": "Which Quarkus version do you want to target?",
      "why_it_matters": "Quarkus 2.x uses javax.* namespace (closer to current code); Quarkus 3.x uses jakarta.* namespace (requires all import rewrites but is the future-proof choice). This fundamentally changes the scope of the migration.",
      "options": ["Quarkus 2.16.x LTS (javax.*)", "Quarkus 3.x latest (jakarta.*)"],
      "default_suggestion": "Quarkus 3.x (jakarta.*) - despite larger diff, it avoids a second migration later"
    },
    {
      "id": "Q2",
      "category": "vendored_dependency",
      "question": "What should be done about the audit-logging-library JAR (com.enterprise:audit-logging-library)?",
      "why_it_matters": "This is a system-scoped vendored JAR (lib/audit-logging-library-1.0.0.jar). System scope dependencies are not supported in Quarkus builds. The code uses FileSystemAuditLogger in OrderService.java. A v2.0.0 jar also exists in lib/ but is not referenced. Konveyor rules (rules/rule.yaml) describe a migration path from v1.x to v2.x with API changes (builder pattern removal, async logging, StreamableAuditLogger).",
      "options": [
        "Upgrade to v2.0.0 (already in lib/), adapt API calls per rule.yaml guidance, publish to a Maven repository",
        "Rewrite the audit logging inline using Quarkus-native logging (e.g., quarkus-logging-json, OpenTelemetry)",
        "Keep the library but move from system scope to a proper Maven repository dependency",
        "Remove audit logging entirely for now"
      ],
      "default_suggestion": "Upgrade to v2.0.0, publish to Maven repo, and adapt OrderService.java per the Konveyor rules"
    },
    {
      "id": "Q3",
      "category": "messaging_strategy",
      "question": "How should JMS messaging be replaced for Quarkus?",
      "why_it_matters": "The app uses JMS extensively: ShoppingCartOrderProcessor publishes to topic/orders; OrderServiceMDB and InventoryNotificationMDB consume from it. Quarkus does not support traditional JMS MDBs. The InventoryNotificationMDB additionally uses WebLogic-specific JNDI for JMS connection setup.",
      "options": [
        "SmallRye Reactive Messaging with an in-memory connector (simplest, keeps monolith, no external broker)",
        "SmallRye Reactive Messaging with Apache Kafka (adds external dependency but is production-grade)",
        "SmallRye Reactive Messaging with AMQP/ActiveMQ Artemis (closest to JMS semantics)",
        "Quarkus Artemis JMS extension (quarkus-artemis-jms) to keep JMS API with minimal code changes",
        "Replace with direct CDI events (@Observes) since this is a monolith and messages stay in-process"
      ],
      "default_suggestion": "CDI events (@Observes/@ObservesAsync) for the monolith case; SmallRye Reactive Messaging with Kafka if splitting into microservices"
    },
    {
      "id": "Q4",
      "category": "architecture",
      "question": "Should the application remain a monolith or be split into microservices?",
      "why_it_matters": "The artifact is literally named 'monolith' and packages all domains (catalog, cart, orders, shipping, inventory, promotions) in one WAR. Quarkus supports both patterns, but splitting affects messaging strategy, data model, deployment topology, and migration scope significantly.",
      "options": [
        "Keep as a single Quarkus application (monolith-first, simplest migration)",
        "Split into microservices along domain boundaries (catalog, cart/order, shipping, inventory)",
        "Migrate as monolith first, then incrementally extract services later"
      ],
      "default_suggestion": "Migrate as monolith first, then extract microservices incrementally"
    },
    {
      "id": "Q5",
      "category": "java_version",
      "question": "What Java version should be targeted?",
      "why_it_matters": "Source is currently Java 8. Quarkus 3.x requires Java 17 minimum and recommends 21. The audit-logging-library v2.x rules reference Java 21+. This affects language features available during migration.",
      "options": ["Java 17 (minimum for Quarkus 3.x)", "Java 21 (LTS, recommended, required by audit-logging-library v2.x)"],
      "default_suggestion": "Java 21 (aligns with audit-logging-library v2.x requirement and is current LTS)"
    },
    {
      "id": "Q6",
      "category": "ejb_remote",
      "question": "How should the EJB Remote interface (ShippingService) be handled?",
      "why_it_matters": "ShippingService uses @Remote and is looked up via JNDI in ShoppingCartService. Quarkus has no EJB container. The JNDI lookup uses WildFlyInitialContextFactory.",
      "options": [
        "Convert to a CDI bean with direct @Inject (simplest if keeping monolith)",
        "Convert to a REST client/service if splitting into microservices",
        "Use gRPC for service-to-service communication"
      ],
      "default_suggestion": "Convert to CDI bean with @ApplicationScoped and direct injection"
    },
    {
      "id": "Q7",
      "category": "session_state",
      "question": "How should session-scoped state be handled?",
      "why_it_matters": "CartEndpoint is @SessionScoped and ShoppingCartService is @Stateful EJB (maintaining cart state in memory). Quarkus favors stateless services. This affects horizontal scaling.",
      "options": [
        "Move cart state to database/Redis with a cart ID lookup",
        "Use quarkus-undertow to keep HTTP session support",
        "Redesign as stateless REST with client-side or database-backed cart storage"
      ],
      "default_suggestion": "Move cart state to database with cart ID, making the service stateless"
    },
    {
      "id": "Q8",
      "category": "frontend",
      "question": "What should be done with the AngularJS 1.x frontend?",
      "why_it_matters": "The frontend uses AngularJS 1.x (EOL since Dec 2021), Bower (deprecated), and PatternFly. It is served via JSP. Quarkus does not natively serve JSPs.",
      "options": [
        "Serve as static assets from src/main/resources/META-INF/resources (minimal change)",
        "Modernize to a current framework (React, Angular 17+, or PatternFly 5 with React)",
        "Separate the frontend entirely as an independent SPA with its own build"
      ],
      "default_suggestion": "Serve as static assets initially; modernize frontend as a separate effort"
    },
    {
      "id": "Q9",
      "category": "authentication",
      "question": "Should Keycloak integration use the Quarkus OIDC extension?",
      "why_it_matters": "The app has a Keycloak realm export and keycloak.json for the 'eap' realm. Quarkus has native OIDC support via quarkus-oidc that replaces the old Keycloak adapter approach.",
      "options": [
        "Use quarkus-oidc extension (recommended, native integration)",
        "Keep Keycloak JavaScript adapter for frontend-only auth",
        "Defer auth migration, remove SSO temporarily"
      ],
      "default_suggestion": "Use quarkus-oidc for backend, keep Keycloak JS adapter for frontend"
    },
    {
      "id": "Q10",
      "category": "weblogic_removal",
      "question": "How should WebLogic-specific code be handled?",
      "why_it_matters": "The codebase includes WebLogic stubs (ApplicationLifecycleListener, ApplicationLifecycleEvent, NonCatalogLogger) and WebLogic JNDI references (WLInitialContextFactory, t3://localhost:7001). These have no equivalent in Quarkus.",
      "options": [
        "Replace with Quarkus lifecycle events (@Observes StartupEvent/ShutdownEvent) and remove stubs",
        "Remove entirely if the lifecycle hooks are not critical"
      ],
      "default_suggestion": "Replace with Quarkus @Observes StartupEvent/@Observes ShutdownEvent"
    }
  ]
}
```

## Questions the Skill SHOULD Ask (from expected list)

1. **Which Quarkus version?** -- The prompt "Migrate this to Quarkus" does not specify a version. Quarkus 2.x (javax.*) vs 3.x (jakarta.*) is the single most impactful version decision because it determines whether all javax.* imports must be rewritten to jakarta.*, effectively doubling the migration diff.

2. **What about the audit-logging-library jar?** -- The pom.xml declares a `system`-scoped dependency on `com.enterprise:audit-logging-library:1.0.0` pointing to `lib/audit-logging-library-1.0.0.jar`. A v2.0.0 jar also exists in `lib/`. The `OrderService.java` imports and uses `FileSystemAuditLogger`, `AuditConfiguration`, and `AuditLoggingException` from this library. The Konveyor `rules/rule.yaml` describes a detailed v1-to-v2 migration path including API changes. System-scope dependencies are not supported in standard Quarkus builds, so this must be addressed.

3. **JMS replacement strategy?** -- The app has three JMS-dependent classes: `ShoppingCartOrderProcessor` (producer using `JMSContext` and `@Resource Topic`), `OrderServiceMDB` (standard `@MessageDriven` consumer), and `InventoryNotificationMDB` (WebLogic-style manual JNDI/TopicConnection consumer). Quarkus has no EJB container and thus no `@MessageDriven` support. Multiple replacement strategies exist with very different trade-offs.

4. **Keep monolith or split?** -- The artifact is literally named `coolstore-monolith`. It packages 6+ domain concerns (catalog, cart, orders, shipping, inventory, promotions) into a single WAR. Whether to keep it as one Quarkus app or decompose into microservices changes the migration scope by an order of magnitude.

## Questions the Skill DID Surface (from your analysis)

| # | Question ID | Topic | Matches Expected? |
|---|-------------|-------|-------------------|
| 1 | Q1 | Which Quarkus version (2.x vs 3.x)? | YES -- matches expected #1 |
| 2 | Q2 | What about the audit-logging-library JAR? | YES -- matches expected #2 |
| 3 | Q3 | JMS replacement strategy? | YES -- matches expected #3 |
| 4 | Q4 | Keep monolith or split into microservices? | YES -- matches expected #4 |
| 5 | Q5 | Target Java version (17 vs 21)? | BONUS -- not in expected list but critical |
| 6 | Q6 | EJB Remote (ShippingService) replacement? | BONUS -- not in expected list but important |
| 7 | Q7 | Session-scoped state handling? | BONUS -- not in expected list but important |
| 8 | Q8 | AngularJS 1.x frontend strategy? | BONUS -- not in expected list but relevant |
| 9 | Q9 | Keycloak/OIDC migration approach? | BONUS -- not in expected list but relevant |
| 10 | Q10 | WebLogic stub removal strategy? | BONUS -- not in expected list but necessary |

**Score: 4/4 expected questions surfaced (100%), plus 6 additional relevant questions.**

## Gap Analysis (what was missed and why)

### Against Expected Questions: No Gaps

All four expected questions were surfaced:

1. **Quarkus version (Q1)** -- Detected by observing that the source uses `javax.*` namespace throughout, and knowing that Quarkus 2.x vs 3.x differ on the `javax` vs `jakarta` namespace boundary. This is the most fundamental version question.

2. **audit-logging-library (Q2)** -- Detected through three converging signals: (a) the `system` scope dependency in pom.xml with explicit `systemPath` pointing to `lib/`, (b) the actual jar files in `lib/` (both v1.0.0 and v2.0.0), (c) import statements in `OrderService.java`, and (d) the Konveyor migration rules in `rules/rule.yaml` that describe the v1-to-v2 upgrade path.

3. **JMS replacement (Q3)** -- Detected by finding three classes that depend on JMS APIs: `ShoppingCartOrderProcessor` (producer), `OrderServiceMDB` (@MessageDriven consumer), and `InventoryNotificationMDB` (manual WebLogic-style JNDI subscriber). The diversity of JMS usage patterns makes this a multi-faceted question.

4. **Monolith vs split (Q4)** -- Detected from (a) the artifact name `monolith`, (b) the `coolstore.json` config with `"MONOLITH": true`, and (c) code analysis showing 6+ tightly-coupled domain services all in one WAR.

### Potential Gaps a Weaker Skill Version Might Have

- **Missing the vendored jar**: A skill that only parses `pom.xml` dependencies without checking `scope=system` and `systemPath`, or without scanning `lib/` directories, would miss the audit-logging-library entirely.
- **Missing the second jar version**: Simply finding the pom dependency shows v1.0.0, but a directory listing of `lib/` reveals v2.0.0 is also present -- this changes the question from "how to replace" to "should we upgrade to the version already available."
- **Missing WebLogic artifacts**: The WebLogic stubs are in the source tree itself (`src/main/java/weblogic/`), not in a dependency. A skill that only examines pom.xml would not catch this.
- **Missing the Konveyor rules**: The `rules/` directory contains migration guidance that directly informs the audit-logging-library question. A skill that only looks at `src/` would miss this context.
- **Conflating JMS patterns**: The two MDBs use different JMS patterns (standard @MessageDriven vs WebLogic JNDI). A skill that treats them as identical would produce a less nuanced question.

## Skill Improvement Suggestions

### 1. Mandatory: Scan for vendored/local dependencies beyond pom.xml
The skill must check for `lib/`, `vendor/`, and `libs/` directories AND cross-reference with `system` scope dependencies in the build manifest. The audit-logging-library is the prime example: it only appears in pom.xml as `scope=system` with a `systemPath`, and the actual jars live in `lib/`. A pure pom.xml parser would either miss it or not understand its significance.

**Implementation**: After parsing pom.xml, find all `<scope>system</scope>` entries, resolve their `systemPath`, verify the file exists, and flag for user questioning. Also independently scan for `lib/`, `vendor/`, `libs/` directories and cross-reference any jars found against declared dependencies.

### 2. Mandatory: Detect ALL application-server-specific artifacts in source code
WebLogic, WildFly, and other server-specific code may be vendored as source stubs (not jar dependencies). The skill should grep for known server-specific package prefixes (`weblogic.*`, `com.ibm.websphere.*`, `oracle.adf.*`, etc.) in the Java source tree.

**Implementation**: Run a package-prefix scan across all `.java` files for known proprietary namespaces. The coolstore has `weblogic.application.ApplicationLifecycleListener` and `weblogic.i18n.logging.NonCatalogLogger` as actual source files in the tree.

### 3. Important: Cross-reference migration rules if present
If a `rules/` directory or Konveyor ruleset exists in the repo, the skill should parse it to enrich its questions. In this case, `rules/rule.yaml` contains 5 detailed migration rules for the audit-logging-library upgrade that directly inform the vendored-dependency question.

**Implementation**: Check for `rules/`, `*.windup.xml`, `ruleset.yaml`, or similar migration rule files. Parse them to extract source/target labels and enrich the corresponding questionnaire items.

### 4. Important: Distinguish between different flavors of the same technology
The two JMS consumers use fundamentally different patterns -- one is a standard Java EE `@MessageDriven` bean, the other is a manual WebLogic-style JNDI subscriber. The skill should distinguish these because they require different migration approaches.

**Implementation**: When detecting a technology (e.g., JMS), classify the usage patterns (annotation-driven vs programmatic, standard vs vendor-specific) and surface these distinctions in the question rationale.

### 5. Nice-to-have: Detect absence of tests
The pom.xml has `maven.test.skip=true` and there are zero test source files. The skill should flag this and ask whether the user wants test generation as part of the migration, since migrating without tests is high-risk.

**Implementation**: Check for `src/test/`, `maven.test.skip`, `skipTests`, and test dependencies (JUnit, TestNG, Arquillian). If none found, add a question about test strategy.

### 6. Nice-to-have: Detect frontend technology and its delivery mechanism
The frontend is AngularJS 1.x served via JSP. Since Quarkus does not support JSP, and AngularJS 1.x is EOL, both the frontend framework and its delivery mechanism are migration concerns. A skill should detect the frontend stack separately from the backend.

**Implementation**: Scan `webapp/` for JS framework signatures (angular.js, react, vue), package managers (bower.json, package.json), and serving mechanisms (JSP, Thymeleaf, static HTML). Cross-reference with backend technology to identify incompatibilities.

### 7. Nice-to-have: Flag session state implications
The `@SessionScoped` endpoint and `@Stateful` EJB represent an implicit architecture decision (in-memory session state). Since Quarkus favors stateless designs and these patterns affect horizontal scaling, the skill should surface this proactively.

**Implementation**: Detect `@SessionScoped`, `@Stateful`, `HttpSession` usage, and `<distributable/>` in web.xml. If found, add a question about state management strategy in the target platform.
