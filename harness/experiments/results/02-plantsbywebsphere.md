## Detection Output (the questionnaire.json)

```json
{
  "application": {
    "name": "PlantsByWebSphere",
    "description": "E-commerce sample application for plant nursery products, running on WebSphere Liberty",
    "groupId": "net.wasdev.sample",
    "version": "1.0-SNAPSHOT",
    "packaging": "war"
  },
  "source_stack": {
    "language": "Java",
    "java_version": "1.8",
    "platform": "Java EE 7",
    "runtime": "WebSphere Liberty 19.0.0.8",
    "build_tools": ["Maven", "Gradle"],
    "ui_frameworks": ["JSF 2.0 (Facelets)", "JSP", "Raw Servlets"],
    "persistence": {
      "spec": "JPA 1.0",
      "provider": "EclipseLink",
      "transaction_type": "JTA",
      "data_sources": [
        "jdbc/PlantsByWebSphereDataSource (JTA)",
        "jdbc/PlantsByWebSphereDataSourceNONJTA (non-JTA)"
      ]
    },
    "database": "Apache Derby 10.11.1.1 (embedded, in-memory)",
    "dependency_injection": "CDI 1.1 (beans.xml bean-discovery-mode=all)",
    "mail": "JavaMail (javax.mail via @Resource)",
    "security": "BASIC authentication with role SampAdmin",
    "xml_processing": "Xalan 2.7.2"
  },
  "detected_patterns": {
    "cdi_annotations": [
      "@Inject",
      "@Named",
      "@Dependent",
      "@SessionScoped",
      "@Resource"
    ],
    "jpa_annotations": [
      "@Entity",
      "@Table",
      "@Id",
      "@Version",
      "@Transient",
      "@NamedQueries",
      "@PersistenceContext"
    ],
    "web_annotations": [
      "@WebServlet",
      "@Transactional"
    ],
    "jsf_managed_beans": [
      "faces-config.xml managed-bean: pagecode.Help"
    ],
    "websphere_specific_files": [
      "ibm-web-bnd.xml (virtual-host binding)",
      "ibm-web-ext.xml (context-root, reloading, file-serving)",
      "ibm-ejb-jar-bnd.xml (mail resource-ref binding)",
      "server.xml (Liberty features, datasources, keystore, basicRegistry)"
    ],
    "view_files": {
      "xhtml_jsf": 12,
      "jsp": 3,
      "html_static": 4
    },
    "packages": {
      "bean": "Business logic (CatalogMgr, CustomerMgr, BackOrderMgr, ShoppingCartBean, MailerBean, etc.)",
      "jpa": "JPA entities (Inventory, Customer, Order, OrderItem, BackOrder, Supplier, OrderKey)",
      "war": "Web tier (ShoppingBean, AccountBean, AccountServlet, AdminServlet, ImageServlet, etc.)",
      "utils": "Utility classes (Util, ListProperties)"
    },
    "data_seeding": "pbw.properties file parsed at startup (inventory, customers, suppliers)",
    "session_management": "HttpSession-based with session-scoped CDI beans",
    "vendored_dependencies": "None (only gradle-wrapper.jar)"
  },
  "target_stack": "Spring Boot (version unspecified)",
  "migration_prompt": "Modernize this to Spring Boot",
  "ambiguities_requiring_user_input": [
    {
      "id": "Q1",
      "category": "target_version",
      "question": "Which Spring Boot version do you want to target?",
      "context": "The application currently uses Java 8 and Java EE 7. Spring Boot 2.7.x is the last version supporting Java 8. Spring Boot 3.x requires Java 17+ and uses Jakarta EE 9+ (jakarta.* namespace). This is a foundational choice that affects every other decision.",
      "options": [
        {"value": "2.7.x", "tradeoff": "Minimal Java version change (stay on Java 8), javax.* namespace preserved, but Spring Boot 2.7 is EOL since November 2023"},
        {"value": "3.2.x", "tradeoff": "Requires Java 17+ upgrade, all javax.* imports become jakarta.*, long-term support available"},
        {"value": "3.3.x+", "tradeoff": "Latest features, requires Java 17+, longest support runway"}
      ],
      "recommendation": "Spring Boot 3.2.x or 3.3.x with Java 17 -- the javax-to-jakarta migration is mechanical and gives a much longer support lifecycle"
    },
    {
      "id": "Q2",
      "category": "java_version",
      "question": "Which Java version should the migrated application target?",
      "context": "Currently Java 8. Spring Boot 3.x requires Java 17 minimum. Java 21 is the latest LTS.",
      "options": [
        {"value": "8", "tradeoff": "Only compatible with Spring Boot 2.x (EOL)"},
        {"value": "17", "tradeoff": "Minimum for Spring Boot 3.x, LTS, widely adopted"},
        {"value": "21", "tradeoff": "Latest LTS, virtual threads, pattern matching, but may require dependency updates"}
      ],
      "recommendation": "Java 17 for broadest compatibility with Spring Boot 3.x"
    },
    {
      "id": "Q3",
      "category": "ui_framework",
      "question": "What should replace the JSF/JSP UI layer?",
      "context": "The app has a MIXED view layer: 12 JSF Facelets (.xhtml) pages for the storefront (shopping, cart, checkout, account, login, register) using JSF tag libraries (h:form, h:commandLink, h:outputText, c:forEach), 3 JSP pages for admin functions (backorderadmin.jsp, supplierconfig.jsp, error.jsp) with embedded Java scriptlets, and 4 static HTML files. The JSP admin pages contain significant inline Java code and JavaScript. Spring Boot does not support JSF natively.",
      "options": [
        {"value": "Thymeleaf", "tradeoff": "Server-rendered like JSF, natural HTML templates, excellent Spring Boot integration, closest migration path for existing page structure"},
        {"value": "React/Vue SPA + REST API", "tradeoff": "Modern architecture, decoupled frontend, but requires rewriting ALL views from scratch plus building a REST API layer that does not exist today"},
        {"value": "API-only (headless)", "tradeoff": "Backend only, no UI, useful if a separate frontend team will build the UI independently"},
        {"value": "Spring MVC + JSP (keep JSP)", "tradeoff": "Spring Boot still supports JSP with embedded servlet containers (Tomcat), but this is legacy and limits packaging to WAR"}
      ],
      "recommendation": "Thymeleaf -- it is the natural Spring Boot template engine and the JSF page structure maps well to Thymeleaf templates"
    },
    {
      "id": "Q4",
      "category": "dependency_injection",
      "question": "Should CDI beans be converted to Spring DI, or should CDI be retained?",
      "context": "The application uses CDI 1.1 extensively: @Inject for dependency injection (CatalogMgr, CustomerMgr, ShoppingCartBean injected into web beans and servlets), @Named for bean naming, @Dependent and @SessionScoped for scoping, @Resource for JNDI resource injection (mail session). Spring Boot uses its own DI container with @Autowired/@Component/@Service/@Scope. While CDI can technically run alongside Spring, it creates two competing DI containers.",
      "options": [
        {"value": "Convert to Spring DI", "tradeoff": "Clean integration with Spring ecosystem, @Inject->@Autowired, @Named->@Component/@Service, @Dependent->@Scope('prototype'), @SessionScoped->@Scope('session'). Requires touching every bean class but the mapping is mostly mechanical."},
        {"value": "Keep CDI alongside Spring", "tradeoff": "Possible with weld-spring or deltaspike, but adds complexity, two DI containers, potential conflicts, poor ecosystem support"}
      ],
      "recommendation": "Convert to Spring DI -- the CDI annotations have direct Spring equivalents and a single DI container is simpler"
    },
    {
      "id": "Q5",
      "category": "database",
      "question": "Should the embedded Derby database be replaced with a production database?",
      "context": "The app uses Apache Derby 10.11.1.1 as an embedded in-memory database (memory:PLANTSDB). Data is seeded at startup from pbw.properties. This is fine for a demo but unsuitable for production. Spring Boot supports H2 as an embedded alternative and has excellent support for PostgreSQL, MySQL, etc.",
      "options": [
        {"value": "Keep Derby embedded", "tradeoff": "Minimal migration effort, but Derby is uncommon in Spring Boot projects and has limited tooling support"},
        {"value": "H2 embedded (for dev/test)", "tradeoff": "Better Spring Boot integration, widely used for testing, similar embedded model"},
        {"value": "PostgreSQL/MySQL (production)", "tradeoff": "Production-ready, but changes the deployment model; good to combine with H2 for dev/test profiles"}
      ],
      "recommendation": "H2 for development/testing with PostgreSQL for production, using Spring profiles"
    },
    {
      "id": "Q6",
      "category": "jpa_provider",
      "question": "Should the JPA provider switch from EclipseLink to Hibernate?",
      "context": "The app uses EclipseLink-specific properties (eclipselink.ddl-generation, eclipselink.cache.shared.default) in persistence.xml. Spring Boot defaults to Hibernate and has deeper integration with it (spring.jpa.* properties, auto-configuration). The JPA entity code itself is standard and portable.",
      "options": [
        {"value": "Switch to Hibernate (Spring Boot default)", "tradeoff": "Best Spring Boot integration, auto-configuration works out of the box, largest community. Must replace EclipseLink-specific properties with spring.jpa.hibernate.ddl-auto."},
        {"value": "Keep EclipseLink", "tradeoff": "Possible with spring-boot-starter-data-jpa exclusion + manual config, but swimming against the current"}
      ],
      "recommendation": "Switch to Hibernate -- the EclipseLink-specific config is minimal (2 properties) and the entity code is standard JPA"
    },
    {
      "id": "Q7",
      "category": "security",
      "question": "How should authentication and authorization be handled?",
      "context": "The app uses container-managed BASIC authentication with a basicRegistry in Liberty server.xml (hardcoded user 'dev'/'dev'), protecting admin URLs via security-constraint in web.xml with role SampAdmin.",
      "options": [
        {"value": "Spring Security with in-memory users", "tradeoff": "Direct equivalent of current basicRegistry, fastest to migrate"},
        {"value": "Spring Security with database-backed users", "tradeoff": "Production-ready, customers are already in DB (Customer entity)"},
        {"value": "Spring Security with OAuth2/OIDC", "tradeoff": "Modern auth, integrates with identity providers, more complex setup"}
      ],
      "recommendation": "Spring Security with form login and database-backed users, since the Customer entity already exists"
    },
    {
      "id": "Q8",
      "category": "mail",
      "question": "How should email functionality be migrated?",
      "context": "MailerBean uses JavaMail with @Resource injection from JNDI (mail/PlantsByWebSphere). It sends order confirmation emails. Spring Boot has spring-boot-starter-mail with auto-configured JavaMailSender.",
      "options": [
        {"value": "Spring Boot Mail (spring-boot-starter-mail)", "tradeoff": "Simple migration: @Resource Session -> @Autowired JavaMailSender, configure via application.properties"},
        {"value": "Remove email functionality", "tradeoff": "If email is not needed for the migration target"}
      ],
      "recommendation": "Spring Boot Mail -- straightforward mapping"
    },
    {
      "id": "Q9",
      "category": "build_tool",
      "question": "Which build tool should the migrated project use?",
      "context": "The project currently has both pom.xml (Maven) and build.gradle (Gradle) with Liberty-specific plugins. Spring Boot supports both.",
      "options": [
        {"value": "Maven", "tradeoff": "Well-established, spring-boot-starter-parent, familiar to most Java developers"},
        {"value": "Gradle", "tradeoff": "Faster builds, more flexible, Kotlin DSL option, already exists in project"}
      ],
      "recommendation": "Maven -- it is more commonly used with Spring Boot and the existing pom.xml provides a starting point"
    },
    {
      "id": "Q10",
      "category": "packaging",
      "question": "Should the app be packaged as an executable JAR or a WAR?",
      "context": "Currently packaged as a WAR deployed to Liberty. Spring Boot can produce either an executable JAR (embedded Tomcat) or a traditional WAR.",
      "options": [
        {"value": "Executable JAR (embedded Tomcat)", "tradeoff": "Spring Boot default, simplest deployment, self-contained"},
        {"value": "WAR (external server)", "tradeoff": "Familiar deployment model, but loses some Spring Boot auto-configuration benefits"}
      ],
      "recommendation": "Executable JAR -- it is the Spring Boot convention and simplifies deployment"
    },
    {
      "id": "Q11",
      "category": "data_seeding",
      "question": "How should initial data population be handled?",
      "context": "The app reads inventory, customer, and supplier data from pbw.properties at startup via PopulateDBBean and ResetDBBean. Spring Boot offers data.sql/schema.sql, Flyway, and Liquibase for DB initialization.",
      "options": [
        {"value": "data.sql + schema.sql", "tradeoff": "Simple, Spring Boot native, good for demo apps"},
        {"value": "Flyway/Liquibase", "tradeoff": "Version-controlled migrations, production-ready, more setup"},
        {"value": "Keep custom properties loader", "tradeoff": "Port the PopulateDBBean as a CommandLineRunner, preserves existing data format"}
      ],
      "recommendation": "CommandLineRunner with the existing properties format for compatibility, plus data.sql as a simpler alternative"
    },
    {
      "id": "Q12",
      "category": "websphere_bindings",
      "question": "Are there WebSphere-specific runtime behaviors that must be preserved?",
      "context": "The app has ibm-web-bnd.xml (virtual host binding), ibm-web-ext.xml (context root, file serving, reloading settings), and ibm-ejb-jar-bnd.xml (mail resource binding). These are WebSphere-proprietary and have no direct Spring Boot equivalent.",
      "options": [
        {"value": "Drop all IBM-specific config", "tradeoff": "The behaviors (context root, file serving, mail binding) all have Spring Boot equivalents via application.properties. No functional loss."},
        {"value": "Carefully audit each binding", "tradeoff": "Ensures no hidden runtime dependency is missed"}
      ],
      "recommendation": "Drop all IBM-specific config -- the functionality maps cleanly to Spring Boot application.properties (server.servlet.context-path, spring.web.resources.static-locations, spring.mail.*)"
    }
  ]
}
```

## Questions the Skill SHOULD Ask (from expected list)

1. **Which Spring Boot version?** -- The app uses Java 8 + Java EE 7. Spring Boot 2.x supports Java 8 but is EOL. Spring Boot 3.x requires Java 17+ and jakarta.* namespace. This single decision cascades into every other migration choice (namespace changes, dependency versions, feature availability).

2. **Replace JSF/JSP with what (Thymeleaf, React, API-only)?** -- The app has a mixed view layer: 12 JSF Facelets .xhtml files using JSF tag libraries for the main storefront, 3 JSP files with embedded Java scriptlets for the admin interface, and 4 static HTML files. Spring Boot does not natively support JSF. The view layer is the single largest migration surface area in this application.

3. **What about the CDI beans -- keep CDI or switch to Spring DI?** -- The app uses CDI 1.1 throughout: `@Inject` for dependency injection in 10+ classes, `@Named` for bean naming, `@Dependent` and `@SessionScoped` for scoping, `@Resource` for JNDI mail session lookup. These annotations exist in nearly every bean and servlet class. Running CDI alongside Spring DI creates two competing containers.

## Questions the Skill DID Surface (from your analysis)

The questionnaire surfaced all three expected questions plus nine additional ones:

| ID | Question | Maps to Expected? |
|----|----------|--------------------|
| Q1 | Which Spring Boot version? | YES -- exact match to expected question 1 |
| Q2 | Which Java version? | Related to Q1 (downstream dependency of Spring Boot version choice) |
| Q3 | Replace JSF/JSP with what? | YES -- exact match to expected question 2 |
| Q4 | CDI -> Spring DI or keep CDI? | YES -- exact match to expected question 3 |
| Q5 | Replace Derby with production DB? | Additional (not in expected list) |
| Q6 | EclipseLink -> Hibernate? | Additional (not in expected list) |
| Q7 | Authentication approach? | Additional (not in expected list) |
| Q8 | Email migration approach? | Additional (not in expected list) |
| Q9 | Maven or Gradle? | Additional (not in expected list) |
| Q10 | JAR or WAR packaging? | Additional (not in expected list) |
| Q11 | Data seeding approach? | Additional (not in expected list) |
| Q12 | WebSphere-specific bindings? | Additional (not in expected list) |

**Coverage: 3/3 expected questions surfaced (100%)**

## Gap Analysis (what was missed and why)

**No gaps in the expected questions.** All three expected questions were detected and surfaced with accurate context and appropriate options.

### How detection worked for each expected question:

1. **Spring Boot version** -- Detected by reading pom.xml (`maven.compiler.source=1.8`, `javaee-api:7.0`) and understanding the Java 8 vs Java 17+ split in Spring Boot 2.x vs 3.x. The Liberty runtime version (19.0.0.8) and Java EE 7 feature set confirmed the Java EE generation.

2. **JSF/JSP replacement** -- Detected by file type analysis (12 .xhtml, 3 .jsp, 4 .html), reading web.xml (`FacesServlet` mapping to `*.jsf`), faces-config.xml (managed beans), and examining the .xhtml files for JSF tag libraries (`xmlns:h="http://java.sun.com/jsf/html"`, `xmlns:ui="http://java.sun.com/jsf/facelets"`). The JSP files were identified by their `<%@ page %>` directives and embedded Java scriptlets.

3. **CDI vs Spring DI** -- Detected by reading beans.xml (`bean-discovery-mode="all"` for CDI 1.1), and finding CDI annotations across the source code: `@Inject` in ShoppingBean, AccountServlet, AdminServlet; `@Named` on beans; `@Dependent` on CatalogMgr, MailerBean; `@SessionScoped` on ShoppingBean, ShoppingCartBean; `@Resource` on MailerBean.

### Potential improvements to question quality:

- **Q3 (UI framework)** could be split into two sub-questions: one for the storefront (JSF) and one for the admin area (JSP), since they may warrant different approaches (e.g., Thymeleaf for storefront, API+React for admin).
- **Q4 (CDI)** could also mention the `javax.transaction.@Transactional` on ShoppingCartBean, which maps to Spring's `@Transactional` but is a separate concern from DI.
- A question about **servlet migration** is implicit but could be explicit: the app has 3 raw HttpServlets (AccountServlet, AdminServlet, ImageServlet) that use `RequestDispatcher.include()` for page forwarding -- these should become Spring MVC `@Controller` classes.

## Skill Improvement Suggestions

### 1. Auto-detect the Java EE generation and flag Spring Boot version as a gating question
The skill should recognize that Java EE version (7 vs 8) and Java version (8 vs 11 vs 17) together determine which Spring Boot major version is feasible. This should be the FIRST question asked because all other answers depend on it. Detection signal: `javaee-api` version in pom.xml + `maven.compiler.source` value.

### 2. Classify view technologies separately and note mixing
The skill should distinguish between primary and secondary view technologies. This app uses JSF as its primary UI framework but also has JSP pages and raw servlets -- a mixed view layer. The questionnaire should note this mixing explicitly because different parts of the UI may need different migration strategies.

### 3. Detect container-specific artifacts and generate removal checklist
Files like `ibm-web-bnd.xml`, `ibm-web-ext.xml`, `ibm-ejb-jar-bnd.xml`, and `server.xml` are WebSphere-specific. The skill should inventory these, explain what each one does, and map each setting to its Spring Boot equivalent (or flag it as droppable). This reduces the risk of "what does this file do?" uncertainty during migration.

### 4. Detect JNDI lookups and resource bindings
The app uses JNDI for data source (`jdbc/PlantsByWebSphereDataSource`) and mail session (`mail/PlantsByWebSphere`) lookups. Spring Boot replaces JNDI with `application.properties` configuration. The skill should scan for all `@Resource`, `@PersistenceContext`, and JNDI references and list them as items that need Spring Boot property equivalents.

### 5. Detect session management patterns
The app mixes HttpSession direct manipulation (in AccountServlet) with CDI `@SessionScoped` beans (ShoppingCartBean, ShoppingBean). The skill should flag this dual pattern because Spring Boot session management differs: `@SessionScope` beans work but `HttpSession` manipulation in servlets needs to be converted to Spring MVC patterns.

### 6. Detect data seeding patterns
The custom `pbw.properties` file format with pipe-delimited records is non-standard. The skill should detect `@PostConstruct` or `@Startup` data loading patterns and ask whether to preserve the custom format or migrate to standard Spring Boot mechanisms (data.sql, CommandLineRunner, Flyway).

### 7. Ask about test strategy
The build.gradle references JUnit Platform (`useJUnitPlatform()`) but no test files exist. The skill should note the absence of tests and ask whether test creation is part of the migration scope -- migrating without tests is risky.

### 8. Detect the admin vs customer separation
The app has a clear architectural split: customer-facing storefront (JSF/XHTML) vs admin interface (JSP/Servlets with BASIC auth). The skill should identify this split and ask if the admin interface should be migrated to the same technology as the storefront or handled differently (e.g., Spring Boot Admin, separate microservice).

### 9. Prioritize questions by blast radius
Questions should be ordered by how many downstream decisions they affect. The ordering should be: Spring Boot version > Java version > UI framework > DI approach > everything else. Currently the questionnaire lists them but does not explicitly communicate this dependency chain.

### 10. Detect absence of REST API layer
The app has no REST endpoints -- all interactions are through form submissions and JSP/JSF page navigation. If the target is a modern architecture, the skill should ask whether REST APIs should be introduced and whether the migration should split into backend API + frontend client.
