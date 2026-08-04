# Questionnaire Skill Evaluation: struts-examples

**Vague prompt:** "Migrate this to Spring Boot"

---

## Detection Output (the questionnaire.json)

```json
{
  "detection": {
    "source_framework": "Apache Struts 2",
    "source_version": "7.2.1",
    "target_framework": "Spring Boot",
    "target_version": "unspecified",
    "language": "Java",
    "java_version": "17",
    "build_tool": "Maven",
    "packaging": "WAR (multi-module, 46 submodules)",
    "servlet_spec": "Jakarta EE (Servlet 6.0, JSP 4.0)",
    "project_structure": "multi-module parent POM with 46 independent example WAR submodules",
    "dependencies_detected": {
      "struts2-core": "7.2.1",
      "struts2-tiles-plugin": "7.2.1",
      "struts2-spring-plugin": "7.2.1",
      "struts2-convention-plugin": "used (annotations)",
      "struts2-rest-plugin": "used (rest-angular)",
      "struts2-config-browser-plugin": "7.2.1",
      "struts2-junit-plugin": "7.2.1",
      "struts2-json-plugin": "used (json, json-customize)",
      "log4j2": "2.25.4",
      "jackson": "2.22.0",
      "hibernate-validator": "8.0.1.Final",
      "jakarta.validation": "3.1.1",
      "spring-web": "6.2.12 (in dependencyManagement only)",
      "quarkus-undertow": "3.28.1 (quarkus module)",
      "sitemesh3": "present",
      "apache-shiro": "present (shiro-basic module)",
      "jasperreports": "present",
      "jfreechart": "present"
    },
    "view_technologies": {
      "JSP": "144 files across most modules",
      "FreeMarker": "75 .ftl files (tiles, themes, themes-override, quarkus)",
      "HTML": "several static files",
      "Tiles": "dedicated tiles module with tiles.xml definitions",
      "SiteMesh3": "dedicated sitemesh3 module for page decoration",
      "Angular": "rest-angular module uses AngularJS frontend with REST backend"
    },
    "configuration_style": {
      "xml_based": "46 struts.xml files, web.xml files, tiles.xml, applicationContext.xml",
      "annotation_based": "@Action, @Namespace, @ParentPackage used in annotations and rest-angular modules",
      "mixed": true
    },
    "struts_patterns_used": {
      "ActionSupport_extension": "primary pattern - all action classes extend ActionSupport",
      "ModelDriven": "used in rest-angular (OrderController)",
      "Preparable": "used in preparable-interface, exclude-parameters, crud, shiro-basic, mailreader2",
      "SessionAware": "used in unknown-handler",
      "RestActionSupport": "used in rest-angular for REST endpoints",
      "action_chaining": "dedicated action-chaining module using chain result type",
      "wildcard_method_selection": "dedicated module using *Person pattern in struts.xml",
      "interceptors": "custom interceptors (ShiroUserInterceptor), logger/defaultStack interceptor config",
      "validation_xml": "10+ *-validation.xml files for declarative validation",
      "validation_programmatic": "validate() method override in multiple Register actions",
      "bean_validation": "dedicated module using Hibernate Validator / Jakarta Validation",
      "StrutsParameter_annotation": "@StrutsParameter(depth=N) used for parameter injection",
      "SkipValidation": "@SkipValidation annotation used",
      "type_conversion": "dedicated module for custom type converters",
      "expression_cache": "dedicated module",
      "text_provider": "i18n message resources via getText() and .properties files",
      "themes": "custom FreeMarker themes (KUTheme), theme overrides"
    },
    "test_coverage": {
      "unit_testing_module": true,
      "struts2_junit_plugin": true,
      "junit4": true,
      "test_files_found": "tests in blank, rest-angular, unit-testing modules"
    },
    "security": {
      "shiro": "Apache Shiro integration in shiro-basic module",
      "custom_interceptors": "ShiroUserInterceptor for authentication"
    },
    "deployment": {
      "embedded_server": "Jetty Maven plugin (11.0.18) for development",
      "war_packaging": true,
      "quarkus_variant": "separate quarkus module exists as alternative runtime"
    }
  },
  "decisions_needed": [
    {
      "id": "Q1",
      "category": "view_layer",
      "question": "Replace Tiles + JSP with what view technology in Spring Boot?",
      "context": "The codebase uses 144 JSP files, 75 FreeMarker templates, Apache Tiles for layout composition, and SiteMesh3 for page decoration. Spring Boot supports Thymeleaf (preferred), FreeMarker, Mustache, and even JSP (with limitations). Tiles is effectively dead in the Spring ecosystem.",
      "options": [
        "Thymeleaf (Spring Boot default, natural HTML templates, strong Spring integration)",
        "FreeMarker (already used in 75 templates, reducing rewrite scope for those)",
        "Keep JSP (possible but discouraged in Spring Boot, no embedded JAR support)",
        "React/Angular SPA with REST API (rest-angular module already demonstrates this pattern)"
      ],
      "impact": "high",
      "reason": "Tiles has no direct Spring Boot equivalent. This decision affects every view file in the project (219+ template files) and the entire page layout/composition strategy."
    },
    {
      "id": "Q2",
      "category": "controller_mapping",
      "question": "Struts actions to Spring controllers: 1:1 mapping or redesign?",
      "context": "The codebase has ~197 Java files with diverse action patterns: simple ActionSupport extensions, ModelDriven controllers, RestActionSupport, Preparable interfaces, action chaining, wildcard method selection. Some modules are trivial (one action), others complex (mailreader2, crud, rest-angular). Struts actions bundle state (getters/setters for form data) with controller logic, while Spring controllers are stateless.",
      "options": [
        "1:1 mapping: each Struts action becomes a Spring @Controller with equivalent @RequestMapping methods",
        "Redesign: consolidate related actions into fewer, RESTful Spring controllers grouping by resource",
        "Hybrid: 1:1 for simple modules, redesign for complex ones (rest-angular, crud, mailreader2)"
      ],
      "impact": "high",
      "reason": "Struts actions are stateful (instance variables hold form data per request). Spring controllers are stateless singletons. This fundamental difference requires careful handling of form binding, model attributes, and request-scoped data."
    },
    {
      "id": "Q3",
      "category": "configuration",
      "question": "Keep XML configuration or move to annotation-based Spring Boot configuration?",
      "context": "46 struts.xml files define action mappings, interceptor stacks, and result dispatching via XML. Some modules already use annotation-based configuration (@Action, @Namespace). Spring Boot strongly favors annotation-based configuration with @Controller, @RequestMapping, @GetMapping, etc. There is also an existing Spring integration (spring-struts module) using applicationContext.xml.",
      "options": [
        "Full annotation-based (Spring Boot idiomatic: @Controller, @RequestMapping, @Configuration, component scanning)",
        "Hybrid: annotation-based controllers with some XML for complex bean wiring where needed",
        "Gradual: keep XML config initially using Spring XML imports, migrate to annotations incrementally"
      ],
      "impact": "medium",
      "reason": "46 struts.xml files encode routing, interceptor stacks, and result mappings. These must all be translated regardless of approach, but the target format affects maintainability and developer experience."
    },
    {
      "id": "Q4",
      "category": "scope",
      "question": "Migrate all 46 submodules or select representative subset?",
      "context": "This is a tutorial/examples repository with 46 independent modules demonstrating different Struts 2 features. Many modules overlap in patterns (e.g., form-processing, form-tags, form-validation all handle forms). Some modules demonstrate Struts-specific features that may not have Spring Boot equivalents (expression-cache, themes, themes-override).",
      "options": [
        "Migrate all 46 modules to serve as complete Spring Boot examples",
        "Select ~10-15 representative modules covering all unique patterns and skip duplicates",
        "Migrate core modules, mark Struts-specific ones (themes, expression-cache) as not-applicable"
      ],
      "impact": "high",
      "reason": "This is a tutorial repo, not a production app. Migrating all 46 modules is significant effort. Many demonstrate the same concept. The migration itself could serve as a 'Struts-to-Spring-Boot migration guide' if scoped appropriately."
    },
    {
      "id": "Q5",
      "category": "interceptors",
      "question": "How to translate Struts interceptor stacks to Spring equivalents?",
      "context": "Struts interceptor stacks (logger, exception, timer, defaultStack, custom ShiroUserInterceptor) are central to request processing. Spring Boot has HandlerInterceptors, Servlet Filters, Spring AOP, and Spring Security as equivalents. The interceptors module explicitly demonstrates custom interceptor configuration.",
      "options": [
        "Spring HandlerInterceptors for cross-cutting concerns (logging, timing)",
        "Spring Security filters to replace Shiro interceptors",
        "Spring AOP @Around/@Before advice for method-level concerns",
        "Mix of all three based on the specific interceptor's responsibility"
      ],
      "impact": "medium",
      "reason": "Interceptors define the request processing pipeline. Incorrect mapping can break validation, security, parameter binding, and error handling."
    },
    {
      "id": "Q6",
      "category": "validation",
      "question": "Which validation approach to standardize on in Spring Boot?",
      "context": "The codebase uses three different validation mechanisms: (1) Programmatic validate() method overrides in action classes, (2) XML-based validation via *-validation.xml files (~10 files), (3) Bean Validation (JSR-380) with Hibernate Validator in the bean-validation module. Spring Boot supports all three but favors Bean Validation annotations (@Valid, @NotNull, etc.).",
      "options": [
        "Standardize on Bean Validation annotations (@Valid, @NotEmpty, @Min, etc.) for all modules",
        "Use @Valid for simple cases, custom Validator implementations for complex cross-field validation",
        "Preserve diversity to show different Spring Boot validation approaches in tutorial context"
      ],
      "impact": "medium",
      "reason": "Three different validation strategies exist. Consolidation simplifies migration but loses the tutorial value of showing alternatives."
    },
    {
      "id": "Q7",
      "category": "security",
      "question": "Replace Apache Shiro with Spring Security?",
      "context": "The shiro-basic module uses Apache Shiro for authentication and authorization with a custom ShiroUserInterceptor. Spring Boot has first-party Spring Security integration. Shiro and Spring Security have different programming models.",
      "options": [
        "Replace with Spring Security (idiomatic, better Spring Boot integration)",
        "Keep Apache Shiro (Shiro works with Spring Boot, less rewrite needed)",
        "Skip security module migration (it's a demo, not production)"
      ],
      "impact": "low",
      "reason": "Only one module uses security. Decision mainly affects the shiro-basic example."
    },
    {
      "id": "Q8",
      "category": "rest_api",
      "question": "How to handle the REST plugin migration (rest-angular module)?",
      "context": "The rest-angular module uses Struts2 REST plugin with RestActionSupport, ModelDriven, and convention-based URL mapping (GET/POST/PUT/DELETE mapped to index/create/update/destroy methods). It serves an AngularJS frontend. Spring Boot has native REST support via @RestController.",
      "options": [
        "Migrate to @RestController with @GetMapping/@PostMapping etc. (natural fit)",
        "Use Spring HATEOAS for a more mature REST API",
        "Keep REST controllers simple, modernize Angular frontend separately"
      ],
      "impact": "medium",
      "reason": "The REST plugin has the most direct mapping to Spring Boot @RestController. This module should be one of the easier migrations but the AngularJS frontend is also outdated."
    },
    {
      "id": "Q9",
      "category": "deployment",
      "question": "Target Spring Boot embedded server or keep external WAR deployment?",
      "context": "All modules are WAR-packaged and run via Jetty Maven plugin. Spring Boot defaults to embedded Tomcat/Jetty with executable JAR. A Quarkus variant already exists showing embedded server is acceptable.",
      "options": [
        "Embedded Tomcat/Jetty via Spring Boot JAR (idiomatic, simpler deployment)",
        "WAR packaging for external servlet container (traditional, matches current approach)",
        "Both: support JAR-first with WAR fallback via SpringBootServletInitializer"
      ],
      "impact": "low",
      "reason": "Tutorial/demo project, so embedded JAR is simpler and more aligned with modern Spring Boot practices."
    },
    {
      "id": "Q10",
      "category": "spring_boot_version",
      "question": "Which Spring Boot version to target?",
      "context": "The project already uses Jakarta EE namespaces (jakarta.servlet, not javax.servlet), Java 17, and Spring Web 6.2.12 in dependency management. This means it is already Jakarta EE compatible, which aligns with Spring Boot 3.x (which requires Jakarta EE).",
      "options": [
        "Spring Boot 3.x (latest, requires Jakarta EE -- already satisfied, requires Java 17+ -- already satisfied)",
        "Spring Boot 3.4.x (current stable line)"
      ],
      "impact": "low",
      "reason": "The project's Jakarta EE compatibility makes Spring Boot 3.x the only sensible choice. No javax-to-jakarta migration needed."
    }
  ],
  "reasoning": {
    "project_nature": "This is the official Apache Struts 2 examples repository -- a tutorial/demo collection of 46 independent WAR modules, not a monolithic production application. Migration strategy should account for the educational purpose.",
    "key_complexity_drivers": [
      "46 independent submodules with overlapping but distinct Struts patterns",
      "Three view technologies (JSP, FreeMarker, Tiles) requiring different replacement strategies",
      "Mixed config: XML struts.xml (46 files) vs annotation-based (@Action, @Namespace)",
      "Struts-specific concepts with no direct Spring Boot equivalent (ActionSupport state model, ModelDriven, action chaining, wildcard mappings, OGNL value stack, Struts themes)",
      "Existing Spring integration (spring-struts module) showing partial bridge is already present"
    ],
    "migration_accelerators": [
      "Already on Jakarta EE namespaces -- no javax-to-jakarta migration needed",
      "Already on Java 17 -- meets Spring Boot 3.x requirements",
      "Spring Web 6.2.12 already in dependencyManagement",
      "Bean Validation with Hibernate Validator already present",
      "Jackson already integrated for JSON",
      "REST module already demonstrates controller pattern close to Spring @RestController",
      "Existing spring-struts module shows Spring ApplicationContext integration"
    ]
  }
}
```

---

## Questions the Skill SHOULD Ask (from expected list)

1. **Replace Tiles + JSP with what?** -- The expected question targets the view layer replacement decision. Apache Tiles provides page layout composition (header/menu/body/footer) which has no built-in Spring Boot equivalent. JSP is technically supported but discouraged in Spring Boot. The question is critical because 144 JSP files + 75 FTL files + Tiles layout definitions all need a replacement strategy.

2. **Struts actions to Spring controllers -- 1:1 mapping or redesign?** -- The expected question targets the controller architecture decision. Struts actions are stateful (instance variables hold form data), while Spring controllers are stateless singletons. A 1:1 mapping is simpler but produces non-idiomatic Spring code. A redesign produces better Spring Boot code but requires more analysis.

3. **Keep XML config or annotation-based?** -- The expected question targets configuration style. 46 struts.xml files define action mappings, result types, and interceptor stacks via XML. Spring Boot strongly favors annotations. The project already uses some annotations (@Action, @Namespace) alongside XML.

---

## Questions the Skill DID Surface (from your analysis)

| ID | Question | Maps to Expected? |
|----|----------|-------------------|
| Q1 | Replace Tiles + JSP with what view technology in Spring Boot? | YES -- direct match to expected #1 |
| Q2 | Struts actions to Spring controllers: 1:1 mapping or redesign? | YES -- direct match to expected #2 |
| Q3 | Keep XML configuration or move to annotation-based Spring Boot configuration? | YES -- direct match to expected #3 |
| Q4 | Migrate all 46 submodules or select representative subset? | NEW -- not in expected list |
| Q5 | How to translate Struts interceptor stacks to Spring equivalents? | NEW -- not in expected list |
| Q6 | Which validation approach to standardize on in Spring Boot? | NEW -- not in expected list |
| Q7 | Replace Apache Shiro with Spring Security? | NEW -- not in expected list |
| Q8 | How to handle the REST plugin migration (rest-angular module)? | NEW -- not in expected list |
| Q9 | Target Spring Boot embedded server or keep external WAR deployment? | NEW -- not in expected list |
| Q10 | Which Spring Boot version to target? | NEW -- not in expected list |

---

## Gap Analysis (what was missed and why)

### Coverage of Expected Questions: 3/3 (100%)

All three expected questions were surfaced as the top three highest-impact questions (Q1, Q2, Q3). No expected questions were missed.

### Additional Questions Surfaced: 7

The analysis produced 7 additional questions beyond the expected set. These reflect genuine migration decisions that would be needed for a real migration:

- **Q4 (scope)** is arguably the most important for this specific repo since it is a tutorial collection, not a production app. A real migration team would need to decide scope first.
- **Q5 (interceptors)** is a legitimate Struts-to-Spring mapping concern. Struts interceptor stacks are central to the framework and have no single Spring Boot equivalent.
- **Q6 (validation)** reflects a real complexity: three different validation approaches exist in the codebase and must be reconciled.
- **Q7-Q10** are lower-impact but still real decisions that would arise during planning.

### What Could Be Improved

1. **Question ordering could be better.** Q4 (scope -- migrate all 46 or subset?) should arguably be Q1, since it gates all subsequent decisions. The expected questions assume the "what to migrate" is already decided.

2. **Missing question about OGNL/Value Stack.** Struts 2 uses OGNL expressions in JSP/FTL templates to access the value stack. Spring Boot uses SpEL or Thymeleaf expressions. This is a significant template rewrite concern that was noted in the reasoning but not surfaced as an explicit question.

3. **Missing question about i18n/message resource strategy.** The codebase uses Struts getText() and .properties resource bundles (with locale variants like global_es.properties). Spring Boot has its own MessageSource. This is a cross-cutting migration concern.

4. **Missing question about the Quarkus module.** A Quarkus variant already exists. Should it be preserved, converted, or dropped? This signals the project maintainers are already exploring alternative runtimes.

---

## Skill Improvement Suggestions

### 1. Add a "scope gating" question for multi-module projects
When a project has many modules (especially >10), the skill should always surface a scope question first. This question gates the entire migration plan. Without it, the team might waste effort analyzing modules that will not be migrated.

### 2. Detect template expression languages and surface them
OGNL (Struts 2), SpEL (Spring), Thymeleaf expressions, and FreeMarker expressions are all different. The skill should scan for `<s:property`, `%{...}`, `#request`, `<s:iterator>` and similar Struts tag library usage in JSP/FTL files and flag the expression language migration as a distinct concern. In this codebase, every JSP file uses Struts tag libraries that have no direct Spring equivalent.

### 3. Detect existing target-framework integrations
The spring-struts module already uses `struts2-spring-plugin` with `applicationContext.xml` and `ContextLoaderListener`. The skill should detect this as a bridge/migration accelerator and ask whether to build on the existing Spring integration or start fresh. Similarly, the quarkus module suggests the team is exploring alternatives.

### 4. Surface i18n/localization as a distinct question
When .properties files with locale suffixes exist (e.g., `global_es.properties`, `frontend_de.properties`), the skill should ask about the i18n migration strategy. Struts getText() and Spring MessageSource work differently.

### 5. Detect the "tutorial vs production" nature of the project
The skill should analyze the project structure and README to determine if this is a tutorial/examples repo or a production app. A tutorial repo migration has fundamentally different goals (demonstrate equivalent patterns in the target framework) versus a production app migration (preserve business logic, minimize risk). The presence of 46 independent, non-interconnected modules with names like "helloworld", "basic-struts", "blank" is a strong signal.

### 6. Identify Struts-specific patterns with no direct Spring Boot equivalent
Some Struts 2 patterns have no 1:1 Spring Boot equivalent and should be flagged individually:
- **Action chaining** (chain result type) -- Spring Boot has no equivalent; needs redirect or service-layer orchestration
- **Wildcard method selection** (`*Person` action name patterns) -- Spring Boot uses explicit @RequestMapping
- **ModelDriven interface** -- Spring Boot uses @ModelAttribute
- **Struts themes** (KUTheme, xhtml theme) -- Spring Boot has no form rendering theme system
- **Value stack / OGNL** -- No equivalent; replaced by model attributes and Thymeleaf/FreeMarker expressions

### 7. Ask about the test migration strategy
The codebase uses `struts2-junit-plugin` for testing with JUnit 4. Spring Boot uses `spring-boot-starter-test` with JUnit 5. The skill should ask whether to migrate tests to JUnit 5 + Spring Boot Test and whether MockMvc testing should replace Struts JUnit plugin testing.

### 8. Improve detection output with quantitative metrics
The questionnaire should include quantitative data to help the user understand migration scale:
- Total files to touch: ~197 Java + 144 JSP + 75 FTL + 46 struts.xml + N web.xml = ~470+ files
- Lines of code estimate per language
- Number of unique action classes vs. number of request mappings needed
- Dependency count requiring replacement vs. those that carry over
