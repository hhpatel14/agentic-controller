# Evaluation: Questionnaire Skill on sample-app-mod

**Vague prompt**: "Upgrade this to Java 21 and Liberty"

---

## Detection Output (the questionnaire.json)

```json
{
  "detection": {
    "app_name": "ModResorts",
    "app_type": "Java EE Web Application (WAR)",
    "source_platform": "IBM WebSphere Application Server (tWAS) 9.0",
    "target_platform": "OpenLiberty / WebSphere Liberty",
    "current_java_version": "1.8 (Java 8)",
    "target_java_version": "21",
    "java_ee_version": "Java EE 7 (javax:javaee-api:7.0)",
    "build_tool": "Maven (also has Gradle wrapper)",
    "packaging": "WAR",
    "source_file_count": 26,
    "frameworks_detected": {
      "servlets": true,
      "ejb": true,
      "cdi": true,
      "jmx_mbeans": true,
      "jndi": true,
      "jsp": true,
      "jdbc": true,
      "form_auth": true
    },
    "ibm_websphere_specific": {
      "apis_used": [
        "com.ibm.websphere.runtime.ServerName (WeatherServlet.java - getDisplayName/getFullName)",
        "com.ibm.websphere.naming.WsnInitialContextFactory (WeatherServlet.java - JNDI factory)",
        "com.ibm.websphere.security.WSSecurityHelper (LogoutServlet.java - revokeSSOCookies)",
        "com.ibm.websphere.servlet.response.ResponseUtils (UpperServlet.java - encodeDataString)"
      ],
      "deployment_descriptors": [
        "WebContent/WEB-INF/ibm-web-bnd.xml (virtual host binding)",
        "WebContent/WEB-INF/ibm-web-ext.xml (context-root, JSP reload config)",
        "WebContent/META-INF/ibm-application-bnd.xml (application binding)",
        "WebContent/WEB-INF/ibm-metadata.xml (empty)"
      ],
      "maven_dependency": "com.ibm.websphere.appserver:was_public:9.0.0 (scope: provided)"
    },
    "javax_imports": {
      "servlet": ["javax.servlet.*", "javax.servlet.annotation.*", "javax.servlet.http.*"],
      "ejb": ["javax.ejb.Singleton", "javax.ejb.Startup"],
      "cdi": ["javax.inject.Inject"],
      "jmx": ["javax.management.* (13 classes)"],
      "naming": ["javax.naming.InitialContext", "javax.naming.NamingException"],
      "annotation": ["javax.annotation.Resource"],
      "sql": ["javax.sql.DataSource"]
    },
    "ejb_usage": {
      "classes": ["com.acme.modres.db.ModResortsCustomerInformation"],
      "annotations": ["@Singleton", "@Startup"],
      "purpose": "Database access singleton that initializes on app startup"
    },
    "cdi_usage": {
      "injection_points": ["WeatherServlet injects ModResortsCustomerInformation via @Inject"]
    },
    "security_concerns": {
      "security_manager": "Service.java references System.getSecurityManager() - removed in Java 17+",
      "ssl_utils": "SSLUtils.java has commented-out com.sun.net.ssl imports (not available since Java 11)",
      "form_auth": "web.xml defines FORM auth with j_security_check (currently disabled for demo)"
    },
    "vendored_dependencies": {
      "WebContent/lib/pikaday.js": "Vendored JavaScript date-picker library (42KB)"
    },
    "third_party_dependencies": [
      "com.google.code.gson:gson:2.10.1 (JSON processing)",
      "org.junit.jupiter:junit-jupiter-api:5.10.0 (test)",
      "org.mockito:mockito-core:5.11.0 (test)",
      "org.springframework:spring-test:5.3.20 (test)"
    ]
  },
  "decisions_needed": [
    {
      "id": "D1",
      "question": "This is a runtime upgrade, not a framework rewrite. Keep EJBs or convert to CDI?",
      "context": "ModResortsCustomerInformation uses @Singleton/@Startup EJB annotations and is injected via @Inject (CDI). Liberty supports EJBs (via ejb-3.2 or ejbLite-3.2 features), but since there is only one EJB and it is already injected via CDI, converting to @ApplicationScoped CDI bean is straightforward and reduces the Liberty feature footprint.",
      "options": [
        "Keep @Singleton/@Startup EJB - requires ejbLite-3.2 or ejb-3.2 Liberty feature",
        "Convert to @ApplicationScoped @Initialized(ApplicationScoped.class) CDI bean - lighter, only needs cdi-2.0 feature"
      ],
      "recommendation": "Convert to CDI - only one EJB class, already CDI-injected, reduces required Liberty features",
      "impact": "low"
    },
    {
      "id": "D2",
      "question": "Jakarta namespace migration: stay on javax (Java EE 8 features) or move to jakarta (Jakarta EE 9+/10)?",
      "context": "The entire codebase uses javax.* namespace (35+ import statements across all source files). Liberty supports both: javaee-8.0 convenience feature keeps javax namespace; jakartaee-9.1 or jakartaee-10.0 require jakarta namespace. Java 21 works with either, but Jakarta EE 10 is the forward-looking choice.",
      "options": [
        "Stay on javax namespace using Liberty javaee-8.0 features - zero code changes for namespace, fastest migration",
        "Migrate to jakarta namespace using Jakarta EE 9.1 (jakartaee-9.1) - moderate effort, better long-term support",
        "Migrate to jakarta namespace using Jakarta EE 10 (jakartaee-10.0) - most modern, best long-term investment"
      ],
      "recommendation": "Stay on javax (javaee-8.0 features) for initial migration, then plan a follow-up Jakarta namespace migration",
      "impact": "high"
    },
    {
      "id": "D3",
      "question": "How should the IBM WebSphere-specific APIs be replaced?",
      "context": "Four IBM-specific API usages must be replaced: (1) ServerName.getDisplayName/getFullName for environment discovery, (2) WsnInitialContextFactory for JNDI, (3) WSSecurityHelper.revokeSSOCookies for logout, (4) ResponseUtils.encodeDataString for output encoding. Liberty has some equivalents; others need standard Java EE alternatives.",
      "options": [
        "Replace with Liberty-specific equivalents where available, standard Java EE for the rest",
        "Replace all with pure Java EE / standard library equivalents for maximum portability"
      ],
      "recommendation": "Use standard Java EE equivalents for portability: standard JNDI, HttpServletRequest.logout(), and manual HTML encoding",
      "impact": "medium"
    },
    {
      "id": "D4",
      "question": "Should the application be containerized for deployment?",
      "context": "The current app is a WAR deployed to tWAS. Liberty is cloud-native and commonly deployed as a container image. The app has no external service dependencies in its current demo configuration (DB connection is commented out).",
      "options": [
        "Deploy as WAR to a standalone Liberty server (traditional deployment)",
        "Containerize with a Dockerfile/Containerfile using icr.io/appcafe/open-liberty or UBI-based Liberty image",
        "Containerize and deploy to OpenShift/Kubernetes"
      ],
      "recommendation": "Containerize - Liberty is designed for container deployment and this is the standard modernization path",
      "impact": "medium"
    },
    {
      "id": "D5",
      "question": "What Liberty features should be configured in server.xml?",
      "context": "Liberty uses a feature-based architecture. The app needs: servlet-4.0 (or 3.1), jsp-2.3, cdi-2.0 (or ejb-3.2 if keeping EJBs), jndi-1.0, and potentially jdbc-4.2 and appSecurity-3.0. The ibm-web-ext.xml context-root (/resorts) must be moved to server.xml.",
      "options": [
        "Use javaee-8.0 convenience feature (includes everything but is heavier)",
        "Specify individual features for a minimal footprint"
      ],
      "recommendation": "Start with individual features for faster startup and smaller image",
      "impact": "low"
    },
    {
      "id": "D6",
      "question": "How should the JMX DynamicMBean (AppInfo) be handled on Liberty?",
      "context": "The app registers a custom DynamicMBean (AppInfo) with the platform MBeanServer. Liberty supports JMX via the monitor-1.0 or localConnector-1.0 features, but the MBean pattern is unusual for Liberty apps. The MBean reads config from ops.json and provides increaseMaxLimit/resetMaxLimit operations.",
      "options": [
        "Keep the JMX MBean pattern and add localConnector-1.0 Liberty feature",
        "Refactor to use Liberty config variables or MicroProfile Config for the configurable operations"
      ],
      "recommendation": "Keep JMX for initial migration (it works on Liberty), refactor later if needed",
      "impact": "low"
    },
    {
      "id": "D7",
      "question": "The SecurityManager is used in Service.java - how should this be handled for Java 21?",
      "context": "java.lang.SecurityManager was deprecated for removal in Java 17 (JEP 411) and the checkMemberAccess call is already commented out. The code still calls System.getSecurityManager() which returns null on Java 17+ and is removed in Java 24.",
      "options": [
        "Remove SecurityManager references entirely (the check is already a no-op)",
        "Replace with an alternative access control mechanism"
      ],
      "recommendation": "Remove the SecurityManager code - it is already effectively dead code",
      "impact": "low"
    }
  ],
  "migration_complexity": "medium",
  "estimated_file_changes": 8,
  "risk_areas": [
    "IBM WebSphere API replacements (4 classes affected)",
    "javax-to-jakarta namespace if chosen (all 26 Java files)",
    "SecurityManager removal (1 class)",
    "Liberty server.xml creation and feature selection",
    "IBM deployment descriptor removal/replacement"
  ]
}
```

---

## Questions the Skill SHOULD Ask (from expected list)

1. **"This is a runtime upgrade, not a framework rewrite. Keep EJBs or convert to CDI?"**
   - Why this matters: The app uses `@Singleton`/`@Startup` EJB on `ModResortsCustomerInformation`, which is injected via CDI `@Inject` into `WeatherServlet`. Liberty supports EJBs but a CDI-only approach is lighter. The user said "upgrade" not "rewrite," so this is a scope-setting question.

2. **"Jakarta namespace migration needed?"**
   - Why this matters: Every Java source file uses `javax.*` imports (Java EE 7). Liberty can run Java 21 with either `javax` (javaee-8.0 features) or `jakarta` (jakartaee-9.1/10.0 features). This is the single biggest scope decision for the migration.

3. **"Containerize?"**
   - Why this matters: Liberty is cloud-native and typically deployed in containers. The user just said "Liberty" without specifying deployment model. This determines whether a Dockerfile, server.xml, and container build pipeline are needed.

---

## Questions the Skill DID Surface (from analysis)

| # | Question | Maps to Expected? |
|---|----------|-------------------|
| D1 | Keep EJBs or convert to CDI? | YES - matches expected #1 exactly |
| D2 | Jakarta namespace migration: stay on javax or move to jakarta? | YES - matches expected #2 exactly |
| D3 | How should IBM WebSphere-specific APIs be replaced? | NO - additional question (not in expected list but critical) |
| D4 | Should the application be containerized? | YES - matches expected #3 exactly |
| D5 | What Liberty features should be configured in server.xml? | NO - additional question (Liberty-specific detail) |
| D6 | How should the JMX DynamicMBean be handled on Liberty? | NO - additional question (niche but relevant) |
| D7 | SecurityManager removal for Java 21? | NO - additional question (Java version-specific) |

**Coverage: 3/3 expected questions surfaced (100%)**

---

## Gap Analysis (what was missed and why)

### No gaps in expected question coverage
All three expected questions were surfaced by the analysis. The skill correctly identified:

1. **EJB vs CDI (D1)**: Detected `@Singleton`/`@Startup` in `ModResortsCustomerInformation.java` and `@Inject` in `WeatherServlet.java`. The skill correctly framed this as a runtime-upgrade scope decision rather than assuming a rewrite.

2. **Jakarta namespace (D2)**: Detected the `javax:javaee-api:7.0` dependency and all `javax.*` imports across the codebase. The skill correctly identified that Liberty offers both paths (javaee-8.0 vs jakartaee-9.1/10.0) and that this is high-impact.

3. **Containerize (D4)**: Recognized that Liberty is cloud-native and the user's vague prompt didn't specify deployment model.

### Additional questions surfaced beyond expectations
The skill surfaced four additional questions (D3, D5, D6, D7) that are all legitimate concerns:

- **D3 (WebSphere API replacement)** is arguably the most critical technical question -- the app has four distinct IBM-specific API usages that MUST be changed. This was not in the expected list but is essential.
- **D5 (Liberty features/server.xml)** is a necessary implementation detail for any tWAS-to-Liberty migration.
- **D6 (JMX MBean)** is niche but the app has a real DynamicMBean implementation that needs consideration.
- **D7 (SecurityManager)** is a Java 21-specific concern that would cause runtime warnings or failures.

### Potential risk: over-questioning
The skill generated 7 questions total. For a small application (26 Java files, ~8 requiring changes), this is reasonable. However, D5, D6, and D7 could be collapsed into implementation details rather than user-facing decision points, since they have clear "right answers."

---

## Skill Improvement Suggestions

### 1. Prioritize and tier the questions
Not all questions deserve equal weight. The skill should separate:
- **Blocking decisions** (must answer before any migration work): D1 (EJB scope), D2 (namespace), D4 (containerize)
- **Implementation details** (have sensible defaults, can be auto-decided): D3, D5, D6, D7

Present blocking decisions first. For implementation details, state the recommended default and only ask if the user wants to override.

### 2. Detect the "upgrade vs rewrite" framing from the prompt
The prompt says "Upgrade this to Java 21 and Liberty" -- the word "upgrade" strongly implies keeping the existing architecture. The skill should detect this signal and frame all questions accordingly (e.g., "Since you said 'upgrade,' we recommend keeping the existing patterns where possible. Do you agree?"). This would collapse D1 into a confirmation rather than an open question.

### 3. Surface the WebSphere API replacement question more prominently
D3 (IBM WebSphere API replacement) is not in the expected list but is arguably more important than the containerization question. The four IBM-specific API usages (`ServerName`, `WsnInitialContextFactory`, `WSSecurityHelper`, `ResponseUtils`) are hard blockers -- the app literally will not compile without them being replaced. The skill should flag these as mandatory changes, not optional decisions.

### 4. Auto-detect the "right" Jakarta namespace answer based on context
For a small app on Java EE 7 targeting Java 21 + Liberty, the practical answer is almost always "stay on javax for the initial migration" (using javaee-8.0 Liberty features). The skill should present this as the default and ask if the user wants to also tackle the namespace migration in the same pass.

### 5. Add a complexity/effort estimate per question
Each question should include an estimated effort impact so the user can make informed trade-offs. For example:
- "Keep EJBs" = 0 additional files changed; "Convert to CDI" = 1 file changed, 5 minutes
- "Stay on javax" = 0 namespace changes; "Move to jakarta" = 26 files, find-and-replace

### 6. Detect vendored/local dependencies and flag compatibility
The skill found `pikaday.js` in `WebContent/lib/` but did not raise it as a question. While this JS library is not affected by a Java/Liberty migration, the skill should at least confirm it is front-end only and not coupled to any server-side rendering that might change.

### 7. Check for existing Liberty branches or migration artifacts
The README reveals that `liberty-java8` and (implied) `liberty-java21` branches already exist. A smart questionnaire skill should detect these and ask: "Reference implementations exist on other branches. Should we use them as a migration guide?" This dramatically reduces ambiguity.
