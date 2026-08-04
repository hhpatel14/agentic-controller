## Detection Output (the questionnaire.json)

```json
{
  "detection": {
    "app_name": "DayTrader7",
    "app_description": "Java EE7 stock trading benchmark application for IBM WebSphere/Open Liberty. Online brokerage simulator with buy/sell orders, portfolio management, market summaries, and real-time WebSocket price updates.",
    "source_stack": {
      "language": "Java 8",
      "framework": "Java EE 7 (Full Profile)",
      "app_server": "IBM WebSphere Liberty / Open Liberty",
      "build_tool": "Maven (multi-module: EJB JAR, WAR, EAR)",
      "java_source_version": "1.8",
      "java_target_version": "1.8"
    },
    "modules": [
      {
        "name": "daytrader-ee7-ejb",
        "packaging": "jar (EJB)",
        "description": "Business logic layer: Stateless Session Beans, Singleton EJB, Message-Driven Beans, JPA entities, direct JDBC access"
      },
      {
        "name": "daytrader-ee7-web",
        "packaging": "war",
        "description": "Web tier: Servlets, JSF 2.2, JSP, WebSocket endpoints, CDI beans, REST-like servlet actions"
      },
      {
        "name": "daytrader-ee7",
        "packaging": "ear",
        "description": "Enterprise archive assembling EJB and WAR modules for deployment"
      }
    ],
    "java_ee_features_detected": [
      "EJB 3.2 (Stateless Session Beans: TradeSLSBBean, DirectSLSBBean)",
      "EJB 3.2 (Singleton: MarketSummarySingleton with @Schedule timer)",
      "EJB 3.2 (Message-Driven Beans: DTBroker3MDB for Queue, DTStreamer3MDB for Topic)",
      "JPA 2.1 (persistence.xml with JTA transaction type, 5 entity classes)",
      "JMS 2.0 (QueueConnectionFactory, TopicConnectionFactory, JMSContext)",
      "JSF 2.2 (FacesServlet, .xhtml views, CDI-backed managed beans)",
      "Servlet 3.1 (@WebServlet annotation, async servlet support)",
      "WebSocket 1.1 (@ServerEndpoint for real-time market data)",
      "CDI 1.2 (@Inject, @Observes, CDI Events for JMS-to-WebSocket bridge)",
      "Bean Validation 1.1 (@NotNull on entities)",
      "JSON-P 1.0 (javax.json for WebSocket message formatting)",
      "Concurrency Utilities 1.0 (ManagedThreadFactory, ManagedScheduledExecutorService)",
      "JTA transactions (Container-Managed and UserTransaction for 2-phase commit)",
      "JNDI lookups (InitialContext for DataSource, JMS resources)",
      "IIOP endpoint configured (CORBA interop)"
    ],
    "databases_detected": [
      {
        "type": "Apache Derby (embedded)",
        "usage": "Default development database, bundled via Maven dependency, pre-populated data in resources/data/tradedb",
        "connection": "JNDI jdbc/TradeDataSource via Liberty server.xml"
      },
      {
        "type": "IBM DB2",
        "usage": "Production database option, separate server_db2.xml config",
        "connection": "DB2 Type 4 JDBC driver (db2jcc4.jar vendored in db2jars/), connects to remote DB2 via host/port/dbname env vars",
        "vendored_driver": true
      }
    ],
    "messaging_detected": {
      "provider": "WebSphere JMS (wasJmsServer/wasJmsClient)",
      "queues": ["TradeBrokerQueue (order processing)"],
      "topics": ["TradeStreamerTopic (quote price changes)"],
      "mdbs": ["DTBroker3MDB (queue consumer for order completion)", "DTStreamer3MDB (topic consumer for price change streaming to WebSocket)"],
      "patterns": ["Async order processing via JMS queue", "Pub/Sub quote price streaming via JMS topic", "2-phase commit (JMS + DB in global transaction)"]
    },
    "web_tier_detected": {
      "servlet_based": ["TradeAppServlet", "TradeScenarioServlet", "TradeConfigServlet", "~40 Ping* benchmark servlets"],
      "jsf_views": ["portfolio.xhtml", "tradehome.xhtml", "marketSummary.xhtml", "account.xhtml", "configure.xhtml"],
      "jsp_views": ["tradehome.jsp", "portfolio.jsp", "quote.jsp", "welcome.jsp", "error.jsp", "~10 others"],
      "websocket_endpoints": ["MarketSummaryWebSocket (/marketsummary)"],
      "security": "BASIC authentication with role-based authorization (grp1-grp5, AllAuthenticated)"
    },
    "deployment_detected": {
      "container": "open-liberty:full (Dockerfile), websphere-liberty:kernel-java17 (Containerfile_db2)",
      "kubernetes": "OpenLibertyApplication CRD (deploy/daytrader7-deploy.yaml)",
      "ci": "Travis CI (.travis.yml)"
    },
    "code_metrics": {
      "java_files": 108,
      "ejb_beans": 5,
      "jpa_entities": 5,
      "mdbs": 2,
      "servlets": "~45 (including benchmark primitives)",
      "jsf_beans": "~12",
      "websocket_endpoints": 1,
      "direct_jdbc_dao": 1
    },
    "vendored_dependencies": {
      "db2jars/db2jcc4.jar": "6.3 MB IBM DB2 JDBC driver"
    },
    "complexity_indicators": [
      "Dual data access paths: EJB/JPA and direct JDBC (TradeDirect) switchable at runtime",
      "Multiple order processing modes: synchronous, async JMS, async ManagedThread",
      "2-phase commit transactions (JMS + database)",
      "WebSphere-specific JMS provider (wasJmsServer/wasJmsClient)",
      "IBM-specific deployment descriptors (ibm-ejb-jar-bnd.xml, ibm-web-bnd.xml, ibm-web-ext.xml)",
      "EAR packaging with separate EJB and WAR modules",
      "Hardcoded JNDI lookups throughout",
      "Session-state management in HTTP sessions",
      "CDI Events bridging JMS MDB to WebSocket",
      "EJB @Schedule timer for periodic market summary refresh"
    ]
  },
  "decisions_needed": [
    {
      "id": "Q1",
      "category": "target_framework",
      "question": "When you say 'cloud-native framework', which target do you mean? The main options for Java EE 7 migration are: (a) Quarkus (with MicroProfile + CDI, optimized for containers/GraalVM), (b) Spring Boot (most popular, rich ecosystem), (c) Jakarta EE 10+ on a modern runtime like Open Liberty 23.x or WildFly (lowest friction, keeps EJB/JMS), or (d) MicroProfile-only on Open Liberty (strip EJB, keep CDI). Each has very different implications for the EJB and JMS rewrite effort.",
      "why_it_matters": "DayTrader7 uses 14+ Java EE 7 features. Quarkus and Spring Boot require replacing ALL EJBs, while Jakarta EE preserves them. The migration effort ranges from weeks (Jakarta EE bump) to months (full Quarkus rewrite).",
      "options": ["Quarkus", "Spring Boot", "Jakarta EE 10+ (Open Liberty / WildFly)", "MicroProfile on Open Liberty"],
      "default_if_unasked": "Jakarta EE 10+ on Open Liberty (lowest risk, preserves EJB/JMS patterns)",
      "risk_if_wrong": "high"
    },
    {
      "id": "Q2",
      "category": "ejb_replacement",
      "question": "The application has heavy EJB usage -- 2 Stateless Session Beans (TradeSLSBBean, DirectSLSBBean), 1 Singleton (MarketSummarySingleton with @Schedule), and 2 Message-Driven Beans (DTBroker3MDB, DTStreamer3MDB). If the target is NOT Jakarta EE, what should replace these? Options: (a) CDI beans + @Transactional for session beans, (b) Spring @Service + @Transactional, (c) Keep EJBs (Jakarta EE path). The Singleton with @Schedule is especially tricky -- it needs a replacement scheduler.",
      "why_it_matters": "EJBs provide container-managed transactions, concurrency control (@Lock), timers (@Schedule), and the MDB message listener pattern. Removing them requires finding equivalents for ALL of these capabilities.",
      "options": ["CDI @ApplicationScoped + @Transactional (Quarkus/MicroProfile)", "Spring @Service + @Transactional + @Scheduled", "Keep EJBs (Jakarta EE)", "Mix: CDI for session beans, framework-specific for MDBs"],
      "default_if_unasked": "CDI beans with @Transactional",
      "risk_if_wrong": "high"
    },
    {
      "id": "Q3",
      "category": "jms_replacement",
      "question": "The app uses WebSphere-specific JMS (wasJmsServer/wasJmsClient) with a JMS Queue for async order processing and a JMS Topic for real-time quote price streaming (which bridges to WebSocket via CDI Events). If moving away from WebSphere JMS, what messaging system should replace it? Options: (a) Apache Kafka, (b) RabbitMQ/AMQP, (c) ActiveMQ Artemis (closest JMS drop-in), (d) Cloud-managed messaging (AWS SQS/SNS, GCP Pub/Sub), (e) In-process event bus (if single-instance deployment is acceptable).",
      "why_it_matters": "The 2-phase commit pattern (JMS + DB in a single XA transaction) is tightly coupled to JMS. Kafka and cloud brokers do NOT support XA transactions. This would require redesigning the order processing flow to use the Outbox pattern or eventual consistency.",
      "options": ["ActiveMQ Artemis (JMS-compatible drop-in)", "Apache Kafka (requires outbox pattern redesign)", "RabbitMQ", "Cloud-managed (SQS/SNS, GCP Pub/Sub)", "In-process CDI Events (single-node only)"],
      "default_if_unasked": "ActiveMQ Artemis",
      "risk_if_wrong": "high"
    },
    {
      "id": "Q4",
      "category": "database",
      "question": "The application supports two databases: embedded Derby (development) and DB2 (production, with vendored db2jcc4.jar driver). For the cloud-native target: (a) Keep DB2 as the production database? (b) Migrate to PostgreSQL or another cloud-friendly database? (c) Use a cloud-managed database (AWS RDS, GCP Cloud SQL, Azure Database)? The JPA entities and named queries should be mostly portable, but the direct JDBC code in TradeDirect.java uses raw SQL that may have DB-specific syntax.",
      "why_it_matters": "DB2 licensing costs and operational overhead may conflict with cloud-native goals. However, the TradeDirect class has ~25 hardcoded SQL statements that would need validation against a new database. The JPA entities should be more portable.",
      "options": ["Keep DB2", "Migrate to PostgreSQL", "Cloud-managed database service", "Keep DB2 for now, migrate database later"],
      "default_if_unasked": "Keep DB2 initially, plan PostgreSQL migration separately",
      "risk_if_wrong": "medium"
    },
    {
      "id": "Q5",
      "category": "web_tier",
      "question": "The web tier uses BOTH JSP (Servlet-based flow via TradeAppServlet) AND JSF 2.2 (XHTML views with CDI backing beans), plus a WebSocket endpoint. The user can switch between them at runtime. For the cloud-native target: (a) Keep both rendering technologies? (b) Consolidate to one (which one)? (c) Replace with a REST API + modern SPA frontend? (d) Use server-side templating (Thymeleaf for Spring Boot, Qute for Quarkus)?",
      "why_it_matters": "JSF 2.2 requires significant server-side state and is poorly supported on Quarkus. JSP is deprecated in many modern runtimes. A REST API approach would decouple frontend/backend but requires building a new frontend.",
      "options": ["Keep JSP + JSF (Jakarta EE path only)", "REST API + SPA frontend (React/Angular)", "Server-side templates (Thymeleaf/Qute)", "REST API only (headless, benchmark-focused)"],
      "default_if_unasked": "REST API with the existing JSP/JSF views kept temporarily",
      "risk_if_wrong": "medium"
    },
    {
      "id": "Q6",
      "category": "dual_data_access",
      "question": "DayTrader has TWO parallel data access implementations switchable at runtime: (1) EJB/JPA path (TradeSLSBBean using EntityManager) and (2) direct JDBC path (TradeDirect with raw PreparedStatements, ~25 SQL strings, and manual connection/transaction management). Both implement the TradeServices interface. Should both paths be migrated, or should one be dropped?",
      "why_it_matters": "Maintaining two parallel data access implementations doubles the migration and testing effort. The direct JDBC path exists for benchmark comparison and may not be needed in a cloud-native context. However, dropping it changes the application's benchmarking capabilities.",
      "options": ["Migrate both paths", "Keep JPA only (drop TradeDirect)", "Keep JDBC only (drop EJB/JPA)", "Merge into a single Spring Data / Panache repository"],
      "default_if_unasked": "Keep JPA path only, drop TradeDirect",
      "risk_if_wrong": "low"
    },
    {
      "id": "Q7",
      "category": "transaction_model",
      "question": "The application uses 2-phase commit (XA transactions) coordinating JMS messaging with database operations. The buy/sell flow sends a JMS message and updates the database in a single global transaction using UserTransaction. Should the cloud-native version preserve XA transaction semantics, or is eventual consistency acceptable?",
      "why_it_matters": "XA transactions require an XA-capable transaction manager, JMS broker, and database driver. Many cloud-native messaging systems (Kafka, SQS) do not support XA. Moving to eventual consistency requires implementing compensating transactions or the Outbox pattern.",
      "options": ["Preserve XA (requires XA-capable broker like Artemis)", "Move to eventual consistency (Outbox pattern)", "Synchronous processing only (simplest, may impact throughput)"],
      "default_if_unasked": "Eventual consistency with Outbox pattern",
      "risk_if_wrong": "high"
    },
    {
      "id": "Q8",
      "category": "packaging",
      "question": "The application is currently packaged as an EAR (Enterprise Archive) containing an EJB JAR and a WAR. Cloud-native frameworks typically use a single executable JAR or a thin WAR. Should the multi-module EAR structure be flattened?",
      "why_it_matters": "Quarkus and Spring Boot do not support EAR packaging. Even Jakarta EE runtimes are moving toward WAR-only or executable JAR deployments. Flattening requires merging the EJB and Web modules.",
      "options": ["Flatten to single executable JAR (required for Quarkus/Spring Boot)", "Flatten to single WAR", "Keep EAR (Jakarta EE on Liberty only)"],
      "default_if_unasked": "Flatten to single executable JAR",
      "risk_if_wrong": "low"
    },
    {
      "id": "Q9",
      "category": "deployment_target",
      "question": "The app already has Kubernetes deployment manifests (OpenLibertyApplication CRD) and Dockerfiles. What is the target cloud platform? (a) OpenShift/Kubernetes on-prem, (b) AWS EKS/ECS, (c) GCP GKE/Cloud Run, (d) Azure AKS, (e) Other?",
      "why_it_matters": "The target platform affects container base image selection, managed service integration (databases, messaging), secrets management, and CI/CD pipeline design.",
      "options": ["OpenShift (on-prem or managed)", "AWS (EKS/ECS/Fargate)", "GCP (GKE/Cloud Run)", "Azure (AKS)", "Platform-agnostic Kubernetes"],
      "default_if_unasked": "OpenShift",
      "risk_if_wrong": "medium"
    },
    {
      "id": "Q10",
      "category": "ibm_vendor_lock",
      "question": "The application has several IBM/WebSphere-specific artifacts: ibm-ejb-jar-bnd.xml, ibm-web-bnd.xml, ibm-web-ext.xml, WebSphere JMS provider (wasJmsServer/wasJmsClient), Liberty-specific server.xml configuration, and IIOP endpoint. Should all IBM-specific dependencies be eliminated, or is staying on Liberty (with Jakarta EE update) acceptable?",
      "why_it_matters": "Eliminating all IBM dependencies is the most thorough 'cloud-native' approach but requires the most work. Staying on Open Liberty with a Jakarta EE 10 update preserves most existing code but maintains the Liberty dependency.",
      "options": ["Eliminate all IBM dependencies (full re-platform)", "Stay on Open Liberty with Jakarta EE update", "Hybrid: keep Liberty runtime but replace IBM JMS with standard broker"],
      "default_if_unasked": "Eliminate all IBM dependencies",
      "risk_if_wrong": "medium"
    }
  ],
  "reasoning": {
    "vague_prompt_analysis": "The prompt 'Migrate this to a cloud-native framework' is maximally ambiguous for this application. 'Cloud-native' could mean: (1) just containerizing (already done -- Dockerfile exists), (2) using a cloud-native Java framework (Quarkus, Spring Boot, MicroProfile), (3) adopting cloud-native patterns (12-factor, microservices, externalized config), or (4) breaking into microservices. The term 'framework' suggests option (2), but the specific framework choice cascades into every other decision.",
    "highest_risk_ambiguities": [
      "Target framework choice (Q1) -- affects every other decision",
      "EJB replacement strategy (Q2) -- 5 EJBs including 2 MDBs and a Singleton timer",
      "JMS/messaging replacement (Q3) -- 2-phase commit dependency",
      "Transaction model (Q7) -- XA vs eventual consistency fundamentally changes the architecture"
    ],
    "what_could_go_wrong_if_unasked": [
      "Choosing Quarkus but discovering too late that JSF 2.2 views and EJB @Schedule timers have no equivalent",
      "Replacing WebSphere JMS with Kafka without redesigning the 2-phase commit order flow -- data consistency breaks",
      "Migrating to PostgreSQL but finding that TradeDirect's 25 hardcoded SQL statements use DB2/Derby-specific syntax",
      "Keeping the dual data-access pattern (JPA + raw JDBC) and doubling the migration effort unnecessarily",
      "Assuming the existing Dockerfile/K8s manifests mean the app is already 'cloud-native' when the code still requires a full Java EE app server"
    ]
  }
}
```

## Questions the Skill SHOULD Ask (from expected list)

1. **Cloud-native means what -- Quarkus? Spring Boot? MicroProfile?**
   This is the single most impactful question. "Cloud-native framework" is meaningless without specifying the target. Each option has radically different migration costs for a Java EE 7 application.

2. **The app has heavy EJB + JMS -- replace with what?**
   DayTrader7 has 5 EJBs (2 Stateless, 1 Singleton, 2 MDBs) and deep JMS integration (Queue + Topic + 2-phase commit). These are the hardest components to migrate and the answer depends entirely on the target framework choice.

3. **Database stays the same?**
   The app supports Derby (dev) and DB2 (production) with a vendored 6.3 MB db2jcc4.jar driver. There are two parallel data access paths (JPA and raw JDBC with 25+ hardcoded SQL statements). Database choice affects both paths.

## Questions the Skill DID Surface (from your analysis)

| Expected Question | Corresponding Detected Question(s) | Coverage |
|---|---|---|
| Cloud-native means what -- Quarkus? Spring Boot? MicroProfile? | **Q1 (target_framework)** -- Explicitly asks which framework with 4 options and explains the cascading impact | FULL MATCH |
| The app has heavy EJB + JMS -- replace with what? | **Q2 (ejb_replacement)** -- Covers all 5 EJBs including Singleton timer and MDBs. **Q3 (jms_replacement)** -- Separately and deeply addresses JMS with 5 messaging alternatives. **Q7 (transaction_model)** -- Addresses the 2-phase commit XA dependency that links EJB and JMS | FULL MATCH (exceeded -- split into 3 focused questions) |
| Database stays the same? | **Q4 (database)** -- Covers DB2 vs PostgreSQL vs cloud-managed, notes the raw JDBC SQL portability risk and vendored driver | FULL MATCH |

### Additional Questions Surfaced Beyond Expected List

| Question | Why It Was Surfaced |
|---|---|
| **Q5 (web_tier)** -- JSP + JSF + WebSocket consolidation | Detected dual rendering (JSP and JSF 2.2) plus WebSocket; JSF is poorly supported on Quarkus |
| **Q6 (dual_data_access)** -- Keep both JPA and JDBC paths? | Detected the TradeServices interface with two implementations (TradeSLSBBean and TradeDirect); migrating both doubles effort |
| **Q8 (packaging)** -- EAR to JAR/WAR flattening | EAR packaging is incompatible with Quarkus and Spring Boot |
| **Q9 (deployment_target)** -- Which cloud platform? | Existing K8s manifests use OpenLibertyApplication CRD; target platform affects everything |
| **Q10 (ibm_vendor_lock)** -- Eliminate IBM dependencies? | Found 4 IBM-specific deployment descriptors, Liberty-specific server.xml, and WebSphere JMS provider |

## Gap Analysis (what was missed and why)

### Gaps: NONE for the expected questions

All three expected questions were surfaced with full coverage. In fact, the analysis exceeded expectations by decomposing "The app has heavy EJB + JMS -- replace with what?" into three separate, more actionable questions (Q2: EJB replacement, Q3: JMS replacement, Q7: transaction model), which is arguably better for a questionnaire because each decision point has independent options and risks.

### Minor areas that could have been explored further

1. **Microservices decomposition**: The questionnaire does not ask whether the monolithic DayTrader should be split into microservices (e.g., separate Order Service, Quote Service, Account Service). This is a common "cloud-native" interpretation that was not explicitly surfaced as a question, though the existing questions would naturally lead there.

2. **Observability/monitoring**: No question about replacing Liberty's built-in monitoring with cloud-native observability (OpenTelemetry, Prometheus, Grafana). The app has extensive custom logging (MDBStats, TimerStat) that should be modernized.

3. **Security model migration**: The app uses BASIC auth with XML-configured security roles (grp1-grp5). No question about migrating to OAuth2/OIDC/Keycloak, which is typical for cloud-native applications.

4. **Test strategy**: No existing tests were found in the repo. A question about test coverage requirements before migration would be prudent.

## Skill Improvement Suggestions

### 1. Decompose compound questions (STRENGTH -- already done well)
The skill correctly decomposed "EJB + JMS replacement" into separate questions (Q2, Q3, Q7). This is good practice because the answers are independent -- you could keep EJBs (Jakarta EE) but still replace WebSphere JMS with Artemis.

### 2. Add a "microservices vs monolith" question
For a Java EE 7 application with clear domain boundaries (Order, Quote, Account, Holding), "cloud-native" often implies microservices decomposition. The skill should detect the TradeServices interface and the entity relationships and ask: "Should this remain a monolith, or should it be decomposed into separate services?"

### 3. Add observability/monitoring question
The application has extensive custom instrumentation (MDBStats, TimerStat classes tracking message processing times). The skill should detect these and ask about the cloud-native observability approach.

### 4. Add security modernization question
The skill detected BASIC auth and XML-based security roles but did not surface a question about it. Cloud-native apps typically use token-based auth (JWT/OAuth2). This should be a standard question when security configuration is detected.

### 5. Detect and flag "already partially containerized" status
The skill should note that the app already has Dockerfiles and K8s manifests (OpenLibertyApplication CRD). This contextualizes "cloud-native" -- the user may want MORE than just containerization, since that is already done. The questionnaire should explicitly acknowledge this: "Your app is already containerized on Liberty. By 'cloud-native framework', do you mean changing the runtime framework itself, or adopting cloud-native patterns (config externalization, health checks, etc.) on the existing Liberty runtime?"

### 6. Prioritize questions by cascading impact
Questions should be explicitly ordered by dependency chain. Q1 (target framework) must be answered first because it determines viable answers for Q2-Q10. The skill should enforce this ordering and potentially skip questions that become irrelevant based on earlier answers (e.g., if the user chooses Jakarta EE, Q2 about EJB replacement becomes "keep EJBs" by default).

### 7. Add risk scoring to the overall migration
The questionnaire should include an overall complexity assessment: "Based on detection of 14+ Java EE features, 5 EJBs, 2 MDBs, XA transactions, dual data access paths, JSP + JSF + WebSocket UI, and IBM-specific artifacts, this migration has a COMPLEXITY RATING of HIGH regardless of target framework. Estimated effort: 3-6 months for Jakarta EE bump, 6-12 months for Quarkus/Spring Boot rewrite."

### 8. Surface the "do nothing" option
The questionnaire should explicitly ask: "Open Liberty already supports cloud-native deployment (your Dockerfile and K8s manifests prove this). Would upgrading to Jakarta EE 10 on Open Liberty 23.x (a relatively low-effort update) satisfy your 'cloud-native' requirements, or do you specifically need to leave the Liberty runtime?" This prevents unnecessary large-scale rewrites when a simpler path exists.
