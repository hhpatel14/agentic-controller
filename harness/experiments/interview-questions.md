# Migration Interview Questions

## App 1: Coolstore
**Prompt**: "Migrate this to Quarkus"
**What I can figure out from the code**: This is a Java EE 7 monolith (WAR packaging) targeting JBoss/WildFly, compiled with Java 8. It is an e-commerce app with REST endpoints (JAX-RS), CDI beans, a `@Stateful` session bean (`ShoppingCartService`), a `@MessageDriven` bean (`OrderServiceMDB`) listening on a JMS topic (`topic/orders`), and a `@Stateless @Remote` EJB (`ShippingService`) accessed via JNDI lookup. It uses Flyway for DB migrations and has a proprietary system-scoped JAR dependency (`audit-logging-library-1.0.0.jar` from `lib/`, with a v2 also present).
**What I cannot figure out (unknowns)**:
- Whether the `@SessionScoped` cart endpoint and `@Stateful` bean behavior (in-memory cart per user) is intentional or must be preserved with the same session semantics
- What the `audit-logging-library` JAR does, whether its source is available, and whether it has Quarkus-compatible equivalents
- What database is used in production (Flyway is present but no datasource config beyond the JMS topic definition)
- Whether the JMS topic for order processing needs to remain JMS or can be replaced (e.g., with Quarkus reactive messaging, Kafka, etc.)
- What the remote EJB lookup pattern (`ShippingServiceRemote` via JNDI) is connecting to -- is it the same deployment or a separate service?
- Whether there are front-end assets (JSP/HTML) or if this is purely an API backend
**Questions I would ask before starting**:
1. The `ShippingService` is accessed via remote EJB JNDI lookup (`ejb:/ROOT/ShippingService!...`). Is this calling into the same deployment or a separate remote service? In Quarkus, remote EJB is not supported -- should this become a local CDI bean call, or does it need to become a REST client call to a separate microservice?
2. The `OrderServiceMDB` listens on JMS topic `topic/orders`. What messaging broker is used in production (e.g., ActiveMQ Artemis, IBM MQ)? Should this migrate to Quarkus's SmallRye Reactive Messaging (e.g., backed by Kafka or AMQP), or do you need to keep JMS compatibility?
3. There are two versions of `audit-logging-library` in `lib/` (1.0.0 and 2.0.0) but only 1.0.0 is referenced in the POM. Do you have the source code for this library? Does it need to work at build time (Quarkus does a lot at build time), or can it be replaced with a standard logging/audit solution?
4. The `ShoppingCartService` is `@Stateful` and `CartEndpoint` is `@SessionScoped` -- the cart lives in server memory tied to the HTTP session. Quarkus strongly favors stateless services. Should the cart state move to a database or cache (e.g., Redis), or is there an existing external session store?
5. What database does this application use in production, and do you have the Flyway migration scripts (I see Flyway as a dependency but no `db/migration` folder was apparent)? This affects datasource configuration in Quarkus.
6. The POM has `maven.test.skip=true` and there are no test sources visible. Are there any tests elsewhere, or will this migration proceed without a regression test suite?
7. Is there a front-end (JSP pages, HTML) served by this WAR, or is it purely a REST API?
8. The target Java version is currently 1.8. Quarkus 3.x requires Java 17+. Is upgrading to Java 17 or 21 acceptable?

---

## App 2: PlantsByWebSphere
**Prompt**: "Modernize this to Spring Boot"
**What I can figure out from the code**: This is a WebSphere Liberty Java EE 7 application (WAR) using JSF (FacesServlet with `.jsf` URL pattern), JPA with EclipseLink (persistence unit "PBW"), CDI beans (`@SessionScoped`, `@Dependent`), JavaMail (`@Resource Session mailSession`), and BASIC authentication with security roles. The data layer uses Apache Derby embedded (configured in Liberty's `server.xml`), with both JTA and non-JTA datasources. There are IBM-specific deployment descriptors (`ibm-web-bnd.xml`, `ibm-web-ext.xml`, `ibm-ejb-jar-bnd.xml`) and a JSF managed bean configured in `faces-config.xml`. The code has WebSphere Liberty-specific configuration for context roots, virtual hosts, and feature management.
**What I cannot figure out (unknowns)**:
- Whether the JSF front-end should be preserved (migrated to a Spring Boot + JSF setup) or replaced entirely with a different view technology (Thymeleaf, REST API + SPA, etc.)
- What database should be used in production (currently Derby embedded -- clearly not production-grade)
- Whether the JavaMail `mail/PlantsByWebSphere` resource needs to keep working and what SMTP server it connects to
- Whether the security model (BASIC auth with `SampAdmin` role, `basicRegistry` in Liberty) needs to be replicated or replaced
**Questions I would ask before starting**:
1. The UI is built with JSF 2.0 (FacesServlet, `.xhtml` pages, `faces-config.xml` with managed beans). Should the front-end be kept as JSF (Spring Boot can host JSF but it is unusual), or should it be rewritten using Thymeleaf, or replaced with a REST API and separate front-end?
2. The app uses Apache Derby as an embedded in-memory database. What database will be used in the Spring Boot version -- should I configure H2 for development and a production database (PostgreSQL, MySQL, etc.)? If so, which one?
3. The persistence layer uses EclipseLink JPA. Spring Boot defaults to Hibernate. Is it acceptable to switch to Hibernate, or are there EclipseLink-specific features (the DDL generation, shared cache settings) that must be preserved?
4. The `MailerBean` sends order confirmation emails via a JavaMail `@Resource` session. Does this feature need to work in the migrated app? If so, what SMTP configuration should be used?
5. Security is configured with BASIC auth and a `basicRegistry` in `server.xml` (hardcoded user "dev/dev"). What authentication mechanism do you want in Spring Boot -- Spring Security with a database-backed user store, OAuth2/OIDC, or just basic auth for now?
6. There are several IBM WebSphere-specific deployment descriptors (`ibm-web-bnd.xml`, `ibm-web-ext.xml`, `ibm-ejb-jar-bnd.xml`). These configure virtual host bindings, context root, and EJB bindings. Can these all be dropped, or are there operational settings in them that need equivalent Spring Boot configuration?
7. The JPA entities use optimistic locking (`LockModeType.OPTIMISTIC_FORCE_INCREMENT`) and back-order management logic. Should this business logic be preserved exactly as-is, or is this an opportunity to simplify the inventory management flow?

---

## App 3: DayTrader 7
**Prompt**: "Migrate this to a cloud-native framework"
**What I can figure out from the code**: This is a multi-module Java EE 7 benchmark application originally built for WebSphere Liberty, packaged as an EAR with EJB and WAR modules. It uses heavy Java EE features: EJB 3.2 (Stateless session beans with local/remote interfaces, Singleton beans, two Message-Driven Beans), JPA 2.1 with Derby embedded, JMS (queues and topics for order brokering and market streaming), WebSocket 1.1 (real-time market summary push), CDI events, Managed Scheduled Executor Service (Java EE Concurrency), JSP-based web front-end with servlets, and IIOP for EJB remoting. The `TradeAction` class dynamically switches between EJB3 and "Direct" (JDBC) modes at runtime.
**What I cannot figure out (unknowns)**:
- What "cloud-native framework" means -- Quarkus? Spring Boot? Micronaut? MicroProfile? Or a full decomposition into microservices?
- Whether the EAR multi-module structure should become separate deployable microservices or a single app
- Whether the application's benchmarking purpose matters -- is this being migrated for production use or to modernize the benchmark itself?
- How the WebSocket + JMS event-driven architecture for real-time stock prices should be handled
**Questions I would ask before starting**:
1. "Cloud-native framework" is broad. Which framework do you want: Quarkus, Spring Boot, Micronaut, or something else? This fundamentally changes the migration approach for every Java EE feature in the app.
2. The app is packaged as an EAR with separate EJB and WAR modules. Should the migrated app remain a single deployable unit, or should it be split into microservices (e.g., a trade service, a market data service, a web front-end)?
3. DayTrader uses two Message-Driven Beans (`DTBroker3MDB` for order processing, `DTStreamer3MDB` for market streaming) with JMS queues and topics. What messaging system should replace the Liberty embedded JMS -- Kafka, RabbitMQ, ActiveMQ, or cloud-managed messaging (SQS, Pub/Sub)?
4. The app has WebSocket endpoints pushing real-time stock price changes to connected browsers, using `ManagedScheduledExecutorService` for periodic sends. Does this real-time functionality need to be preserved? If splitting into microservices, which service owns the WebSocket connections?
5. `TradeAction` supports two runtime modes: EJB3 (with remote/local SLSB lookup via JNDI) and "Direct" (raw JDBC via `TradeDirect`). Should the migrated version only support one approach, and if so, which pattern -- repository/service layer or direct JDBC?
6. The front-end is JSP-based servlets (`TradeAppServlet`, `TradeScenarioServlet`, `TradeConfigServlet`). Should these be preserved as server-rendered pages, converted to a REST API with a separate SPA, or replaced entirely?
7. The database is Apache Derby embedded. What production database should be targeted? This affects JPA dialect and connection pool configuration.
8. Is this being migrated for actual production use, or is the goal to modernize the benchmark application itself? This affects whether we prioritize performance characteristics or ease of deployment.

---

## App 4: IBM Sample App Mod (ModResorts)
**Prompt**: "Upgrade this to Java 21 and Liberty"
**What I can figure out from the code**: This is a Java EE 7 WAR application ("ModResorts") currently compiled with Java 8, running on WebSphere (traditional, not Liberty). It uses raw `HttpServlet` classes (no framework), JPA via `javax.javaee-api`, CDI (`@Inject`), servlet filters, and has a direct dependency on `com.ibm.websphere.appserver:was_public:9.0.0` (WebSphere traditional 9.x APIs). Critically, the `WeatherServlet` uses `com.ibm.websphere.runtime.ServerName` (a WebSphere traditional-specific class) and `com.ibm.websphere.naming.WsnInitialContextFactory` for JNDI with CORBA/IIOP (`corbaloc:iiop:localhost:2809`). It also registers JMX MBeans. The app calls an external weather API (Weather Underground) and has custom MBean infrastructure.
**What I cannot figure out (unknowns)**:
- Whether "Liberty" means Open Liberty or WebSphere Liberty (commercial), and which version/features
- Whether the WebSphere traditional-specific APIs (`com.ibm.websphere.runtime.ServerName`, `WsnInitialContextFactory`, IIOP JNDI) are actively used in production or are dead code
- What the IIOP/CORBA JNDI connection at port 2809 is for
- Whether the JMX MBeans (`AppInfo`) serve a monitoring purpose that needs to be preserved
**Questions I would ask before starting**:
1. The `WeatherServlet` calls `com.ibm.websphere.runtime.ServerName.getDisplayName()` and `getFullName()` -- these are WebSphere traditional APIs that do not exist in Liberty. Is this code actively used, or can it be removed? If needed, what should replace it (e.g., MicroProfile Config, environment variables)?
2. The `setInitialContextProps()` method configures JNDI with `com.ibm.websphere.naming.WsnInitialContextFactory` and `corbaloc:iiop:localhost:2809` (CORBA/IIOP). What is being looked up via this JNDI connection? Liberty does not support IIOP by default. Can this be replaced with direct CDI injection or a different service discovery mechanism?
3. Do you mean Open Liberty (open source) or WebSphere Liberty (commercial) as the target? This determines feature availability and support model.
4. The `WeatherServlet` registers a custom JMX MBean (`com.acme.modres.mbean:name=appInfo`). Is this MBean used by monitoring tools? Liberty has its own monitoring features (MicroProfile Metrics, Health) -- should the MBean be replaced with MicroProfile equivalents?
5. The POM depends on `com.ibm.websphere.appserver:was_public:9.0.0`. On Liberty, you would use Liberty-specific feature APIs instead. Are there other usages of WebSphere traditional APIs beyond what I found in `WeatherServlet` (e.g., in `AvailabilityCheckerServlet`, the DB layer)?
6. Moving from Java 8 to Java 21 introduces module system considerations and removes `javax.xml.bind`, `javax.annotation`, etc. from the JDK. Also, should the `javax.*` namespace be migrated to `jakarta.*` (Liberty supports both depending on feature version)? This is a critical decision that affects every source file.
7. The `ModResortsCustomerInformation` is `@Inject`ed -- is this backed by a database, or is it a mock/test implementation? This affects the persistence strategy on Liberty.

---

## App 5: Struts Examples
**Prompt**: "Migrate this to Spring Boot"
**What I can figure out from the code**: This is a collection of 40+ independent Struts 2 example modules (Apache Struts 2 v7.2.1), already on Java 17 and using Jakarta namespace (`jakarta.servlet`). Each module is a small standalone WAR application demonstrating a specific Struts feature (form processing, validation, interceptors, JSON, file upload, tiles, etc.). The modules use Struts XML-based action configuration (`struts.xml`), Struts tag libraries in JSP views (`<s:form>`, `<s:textfield>`), and Action classes extending `ActionSupport`. Some modules integrate with Spring (`spring-struts`), JasperReports, Shiro security, and SiteMesh.
**What I cannot figure out (unknowns)**:
- Whether all 40+ modules should be migrated or just a subset
- Whether this is about creating equivalent Spring Boot example applications or consolidating into a single app
- Whether the JSP views should be converted to Thymeleaf or another view technology
- Whether this is a learning exercise or a production migration
**Questions I would ask before starting**:
1. There are 40+ separate Struts 2 example modules (helloworld, form-processing, crud, file-upload, json, tiles, spring-struts, shiro-basic, etc.). Should all of them be migrated, or is there a specific subset you care about? Migrating all 40+ is a large effort.
2. Each module is an independent mini-application. Should each become its own Spring Boot project, or should they be consolidated into a single Spring Boot application with different endpoints?
3. The JSP views use Struts tag libraries (`<s:form>`, `<s:textfield>`, `<s:submit>`, `<s:checkbox>`). Should these be converted to Thymeleaf templates (the Spring Boot convention), kept as JSP (Spring Boot supports it but discourages it), or converted to something else?
4. Struts 2 Action classes (like `Register extends ActionSupport`) use a fundamentally different pattern from Spring MVC controllers -- the action class itself holds form data as fields with getters/setters, and `struts.xml` defines the action-to-view mapping. Is a direct 1:1 translation to `@Controller` + `@ModelAttribute` acceptable, or do you want the migration to adopt Spring Boot idioms more deeply (e.g., `@RestController` with JSON APIs)?
5. Some modules use specific integrations: JasperReports, Apache Shiro security, SiteMesh layout, Spring-Struts integration, and Tiles. Should these integrations be preserved with Spring Boot equivalents, or can they be dropped?
6. The project already uses Java 17 and Jakarta namespace. Is this a teaching/reference project where the goal is to have equivalent Spring Boot examples for documentation purposes?

---

## App 6: Legacy Cycle Store
**Prompt**: "Migrate this to .NET 8"
**What I can figure out from the code**: This is an ASP.NET MVC 4 application targeting .NET Framework 4.5, using the old-style `.csproj` format (MSBuild with `<Reference>` elements and `<Compile Include>`). It uses Entity Framework 5 with an EDMX model (`CycleModel.edmx`) and T4 code generation for the data layer, connecting to a SQL Server RDS instance (AWS). The UI uses Razor views (`.cshtml`) with a layout system. It has references to ComponentOne Wijmo controls (`C1.C1Report.4`, `C1.Web.Wijmo.Controls.4`) -- a commercial UI component library. Authentication uses ASP.NET Forms Authentication with `SqlMembershipProvider`. The data model represents AdventureWorks product catalog (products, categories, subcategories).
**What I cannot figure out (unknowns)**:
- Whether the ComponentOne Wijmo controls are actively used in views and whether replacements exist for .NET 8
- Whether the SQL Server RDS database schema should be preserved or can be redesigned
- Whether Forms Authentication with `SqlMembershipProvider` should become ASP.NET Core Identity or an external identity provider
- Whether there are additional projects in the solution (only one `.csproj` was found)
**Questions I would ask before starting**:
1. The project references ComponentOne Wijmo controls (`C1.C1Report.4`, `C1.Web.Wijmo.Controls.4`, `C1.Web.Wijmo.Controls.Design.4`). These are commercial UI components that may or may not have .NET 8 equivalents. Are these controls actively used in the views? If so, do you have a ComponentOne license for .NET 8, or should they be replaced with something else (e.g., open-source alternatives, plain Bootstrap)?
2. The data layer uses an EDMX model (`CycleModel.edmx`) with T4 code generation (Database-First approach with Entity Framework 5). In .NET 8, Entity Framework Core does not support EDMX. Should this be migrated to EF Core Code-First (requiring manual model creation), or do you want to use a different ORM?
3. The connection string points to a SQL Server on AWS RDS (`sqlrdsdb.xxxxx.us-east-1.rds.amazonaws.com`). Will the database remain on AWS RDS, or is this moving to a different hosting environment? The connection string also has placeholder credentials -- how should secrets be managed in the new version?
4. Authentication uses `SqlMembershipProvider` (the legacy ASP.NET membership system). This does not exist in .NET 8. Should it be replaced with ASP.NET Core Identity, or are you moving to an external identity provider (Azure AD, Auth0, etc.)?
5. The project uses the old-style verbose `.csproj` format with explicit `<Compile Include>` items. The migration to .NET 8 requires the SDK-style project format. Are there any custom MSBuild targets or T4 templates beyond `CycleModel.edmx` that need special handling?
6. The `ProductManager` class uses static methods with LINQ-to-Entities queries against a shared `Common.DataEntities` context. This pattern (static data access with a shared context) is problematic in .NET 8 where `DbContext` should be scoped per request via DI. Is there significant business logic in these static managers that needs careful refactoring, or is this the extent of it?
7. Some LINQ queries in `ProductManager` are commented out (e.g., `GetCategory`). Is the application fully functional in its current state, or is it partially implemented?

---

## App 7: BookCatalog (.NET)
**Prompt**: "Modernize this app"
**What I can figure out from the code**: This is an ASP.NET MVC 5 application on .NET Framework 4.8 (old-style `.csproj`). It is a simple CRUD book catalog using Entity Framework 6 with Code-First approach, connecting to SQL Server LocalDB (`(LocalDB)\MSSQLLocalDB`). It has one controller (`BooksController`) with standard CRUD operations, one model (`Book`), and Razor views (Index, Details, Create, Edit, Delete). The `ApplicationDbContext` includes a `BookCatalogInitializer` that seeds sample data. There are no authentication, authorization, or complex business logic layers. This is part of a "dotnet-modernization-for-beginners" repository, suggesting it is a teaching example.
**What I cannot figure out (unknowns)**:
- What "modernize" means -- target framework (.NET 8? .NET 9?), architecture (keep MVC? move to Blazor? add API layer?), or scope (just upgrade, or rearchitect?)
- Whether the LocalDB dependency should become something else
- Whether this should remain a server-rendered MVC app or become an API + SPA
**Questions I would ask before starting**:
1. "Modernize" is very broad. What is the target? Options include: (a) upgrade to ASP.NET Core on .NET 8 while keeping the MVC pattern, (b) convert to a Blazor Server or Blazor WebAssembly app, (c) split into a REST API backend + React/Angular/Vue SPA frontend, (d) containerize and add cloud-native patterns. Which direction do you want?
2. The app currently uses SQL Server LocalDB (`(LocalDB)\MSSQLLocalDB`). Should the modernized version use a different database (PostgreSQL, SQLite for dev, Azure SQL, etc.), or keep SQL Server?
3. Should Entity Framework 6 be upgraded to Entity Framework Core? This is straightforward for this simple model but changes some APIs (e.g., `DropCreateDatabaseIfModelChanges` initializer does not exist in EF Core).
4. The app has zero authentication/authorization. Should the modernized version add authentication (ASP.NET Core Identity, OAuth/OIDC), or remain open?
5. The current architecture is the simplest possible MVC pattern -- controller directly instantiates `DbContext`, no service layer, no repository pattern, no dependency injection. Should the modernized version introduce proper layering (DI, service layer, repository), or keep it simple since it is a beginner teaching example?
6. This appears to be from a "modernization for beginners" tutorial repository. Is the goal to modernize this as a learning exercise (where the process matters), or for actual production use?
7. There are Razor views for full CRUD (Index, Details, Create, Edit, Delete). Should these remain as server-rendered HTML, or should the UI be rebuilt?

---

## App 8: Flask-to-FastAPI
**Prompt**: "Migrate this to FastAPI"
**What I can figure out from the code**: This is a Flask web application with Auth0 integration for authentication. It uses Authlib for OAuth/OIDC, has two blueprints (`auth_bp` for login/logout/callback, `webapp_bp` for pages), Jinja2 templates for server-rendered HTML, Flask sessions for user state, and a custom `@requires_auth` decorator. Configuration is loaded from a `.config` INI file. The repo already contains a `fastapi-webapp/` directory with what appears to be a completed FastAPI migration. Dependencies are Flask >= 2.2.2, Authlib >= 1.2.0, requests >= 2.28.2, and auth0-python == 4.1.0.
**What I cannot figure out (unknowns)**:
- Whether the existing `fastapi-webapp/` directory is a complete and correct migration or a partial/reference implementation
- Whether the templates/server-rendered approach should be preserved in FastAPI or converted to a REST API
- Whether there is a database layer not visible from the files I read
**Questions I would ask before starting**:
1. There is already a `fastapi-webapp/` directory in this repository that appears to contain a FastAPI version. Is this a complete and working migration, a partial attempt, or a reference implementation? Should I build on it, replace it, or ignore it and start fresh?
2. The Flask app serves HTML templates (Jinja2) with server-side rendering (`heels.html`, `home.html`, `profile.html`). FastAPI can serve Jinja2 templates too, but it is more commonly used as a REST API. Should the migrated version keep server-rendered templates, or become a pure REST API with a separate front-end?
3. Authentication uses Auth0 via Authlib's OAuth integration (Flask-specific `oauth.auth0.authorize_redirect`). FastAPI has different auth patterns. Should the Auth0 integration use the `auth0-python` SDK directly, or a FastAPI-specific OAuth library like `authlib` with Starlette, or FastAPI's built-in security utilities?
4. The Flask app uses server-side sessions (`session['user']`) to track the logged-in user. FastAPI is typically stateless, using JWT tokens or dependency injection for auth. Should the migrated version use stateless JWT-based auth, or preserve server-side sessions (possible with `starlette-session` or similar)?
5. Is there any database or data persistence layer beyond what I see? The app I read only renders templates and handles auth -- is there more business logic, or is this essentially an Auth0 integration demo?
6. The configuration uses Python's `configparser` with a `.config` INI file. Should this be migrated to Pydantic `BaseSettings` (the FastAPI convention) with environment variables, or keep the INI file approach?

---

## App 9: NodeJS Ecommerce
**Prompt**: "Modernize this to NestJS with TypeScript"
**What I can figure out from the code**: This is an Express.js 4.x e-commerce application using JavaScript (no TypeScript). It uses MongoDB with Mongoose, EJS templates for server-rendered views, Passport.js (local + Facebook strategies) for authentication, Stripe for payments, Elasticsearch via `mongoosastic` plugin for product search, and Express sessions stored in MongoDB via `connect-mongo`. The app has routes for products, cart, user management, admin, and payments. It includes security middleware (Helmet, rate limiting, HPP, toobusy-js). Several dependencies are significantly outdated (e.g., `bcrypt-nodejs` is deprecated, Mongoose 5.x, Stripe 5.x, Passport 0.6).
**What I cannot figure out (unknowns)**:
- Whether the EJS server-rendered views should remain server-rendered or become a separate front-end consuming a NestJS API
- Whether MongoDB should remain the database or if this is an opportunity to switch
- Whether Facebook OAuth login is still needed
- Whether the Elasticsearch integration should be preserved
**Questions I would ask before starting**:
1. The app renders server-side HTML using EJS templates with `ejs-mate` layouts. NestJS can serve templates but is primarily an API framework. Should the migrated app be: (a) a NestJS REST/GraphQL API with a separate React/Angular/Vue front-end, (b) NestJS serving server-rendered templates (e.g., Handlebars), or (c) something else?
2. The app uses Mongoose with MongoDB. NestJS has built-in Mongoose integration (`@nestjs/mongoose`), but this is also an opportunity to switch to TypeORM with a relational database. Should MongoDB/Mongoose be preserved, or do you want to change the database?
3. Product search uses `mongoosastic` (a Mongoose plugin that syncs documents to Elasticsearch at `localhost:9200`). This plugin is outdated and unmaintained. Should Elasticsearch search be preserved (using the official Elasticsearch client), replaced with MongoDB Atlas Search, or dropped?
4. Authentication uses Passport.js with local strategy (email/password with `bcryptjs`) and Facebook OAuth (`passport-facebook`). NestJS has `@nestjs/passport` for Passport integration. Should both auth strategies be migrated, or is Facebook login no longer needed?
5. Payments use Stripe with what appears to be a test publishable key hardcoded in the route file (`pk_test_6eUlykjUkKIa4viRDPGNKjwv`). Should the Stripe integration be migrated using the modern Stripe SDK for TypeScript, and should the key management be moved to environment variables / config service?
6. Several dependencies are severely outdated or deprecated: `bcrypt-nodejs` (deprecated, already has `bcryptjs` as well), `faker` 4.x (unmaintained, forked as `@faker-js/faker`), Mongoose 5.x (current is 8.x), Stripe 5.x (current is much higher). Should the migration also update all dependencies to current versions, or just convert to TypeScript/NestJS with the same dependency versions?
7. The app stores sessions in MongoDB using `connect-mongo`. In a NestJS app, should session management continue using server-side sessions, or switch to JWT-based stateless auth?
8. There are security middlewares (Helmet, express-rate-limit, HPP, toobusy-js for server overload detection). NestJS has its own patterns for these (guards, interceptors, middleware). Should all security features be preserved, including the server overload detection and raw body size checking?

---

## App 10: Struts 1.3 Legacy
**Prompt**: "Migrate this from Struts to modern Java"
**What I can figure out from the code**: Despite the repo name "struts-spring-boot-legacy-code", this application is already a Spring Boot 3.5.4 application on Java 21 using Thymeleaf and Lombok. There is no Struts code at all -- the `UserController` is a standard Spring MVC `@Controller` with `@GetMapping` / `@PostMapping`, the `User` model uses Lombok annotations (`@Data`, `@AllArgsConstructor`, `@NoArgsConstructor`), and the POM parent is `spring-boot-starter-parent`. The only "Struts" reference is in the package name (`com.struts.strutsdemoproject`) and the repo name. It has no database, no security framework, and minimal logic (just a username check against a hardcoded value "Shradha").
**What I cannot figure out (unknowns)**:
- Whether the repo name is misleading and this is already the migrated result, or if there was supposed to be legacy Struts code
- What "modern Java" means if the app is already on Java 21 with Spring Boot 3.5
- Whether there is a previous version of this app that actually used Struts
**Questions I would ask before starting**:
1. This application is already running Spring Boot 3.5.4 on Java 21 with Thymeleaf -- there is no Struts code anywhere. The package names reference "struts" (`com.struts.strutsdemoproject`), but the actual implementation is pure Spring MVC. Is there a different version of this application that actually contains Struts code, or is the migration already done?
2. If the migration is already complete, what work remains? Possible cleanup tasks: renaming packages from `com.struts.strutsdemoproject` to something appropriate, adding a database layer, adding proper authentication instead of the hardcoded username check (`"Shradha"`), or adding tests. What is the actual goal?
3. If there is actual Struts legacy code somewhere else that needs to be migrated, can you point me to it? The code in this repository does not match the migration prompt.
4. The `UserController.submitForm` method redirects to `/error` if the username is not "Shradha" -- is this placeholder logic that needs to be replaced with real authentication, or is this intentional for a demo?
5. When you say "modern Java," do you mean adopting modern Java language features (records, sealed classes, pattern matching), adding modern patterns (reactive, virtual threads), or something else? The app is already on a current Java version and framework.
