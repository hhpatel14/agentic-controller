# 10 - Struts Legacy Modernization Questionnaire Evaluation

**Repo:** `struts-spring-boot-legacy-code`
**Vague Prompt:** "Migrate this from Struts to modern Java"

---

## Detection Output (the questionnaire.json)

```json
{
  "detection": {
    "source_framework": {
      "declared": "Struts 1.3 (per README and project naming)",
      "actual_code": "Spring Boot 3.5.4 with Thymeleaf",
      "discrepancy": true,
      "explanation": "The repository name and README describe a Struts 1.3 legacy application, but the actual codebase has ALREADY been fully migrated to Spring Boot 3.5.4 + Java 21 + Thymeleaf. There is zero Struts code, no ActionForms, no struts-config.xml, no JSPs. The pom.xml has spring-boot-starter-parent 3.5.4 as its parent and contains no Struts dependencies whatsoever."
    },
    "languages": {
      "java": {
        "version": 21,
        "source_files": 4,
        "packages": ["com.struts.strutsdemoproject", "com.struts.strutsdemoproject.controller", "com.struts.strutsdemoproject.model"]
      },
      "html_thymeleaf": {
        "template_files": 3,
        "templates": ["index.html", "welcome.html", "error.html"]
      }
    },
    "build_system": {
      "tool": "Maven",
      "wrapper": true,
      "pom_parent": "org.springframework.boot:spring-boot-starter-parent:3.5.4"
    },
    "dependencies": {
      "runtime": [
        "spring-boot-starter-web",
        "spring-boot-starter-thymeleaf",
        "lombok"
      ],
      "test": [
        "spring-boot-starter-test"
      ],
      "struts_artifacts": [],
      "vendored_jars": []
    },
    "architecture": {
      "pattern": "Spring MVC (Controller + Model + Thymeleaf views)",
      "controllers": ["UserController"],
      "models": ["User"],
      "views": ["index.html (login form)", "welcome.html (success page)", "error.html (error page)"],
      "services": [],
      "repositories": [],
      "database": "none"
    },
    "configuration": {
      "xml_config_files": 0,
      "struts_config": false,
      "web_xml": false,
      "spring_config": "application.properties (minimal - only spring.application.name)",
      "annotation_driven": true
    },
    "app_complexity": {
      "total_java_files": 4,
      "total_html_files": 3,
      "total_config_files": 1,
      "estimated_loc": 60,
      "complexity": "trivial"
    },
    "critical_finding": "MIGRATION ALREADY COMPLETED. This repo is the OUTPUT of a Struts-to-Spring-Boot migration, not the INPUT. The README explicitly states: 'This project is a modernization of the Struts 1.3 demo application... The codebase has been re-engineered to leverage Java 21, Spring Boot, Thymeleaf, Bootstrap.' The original Struts code is in a separate repository: https://github.com/ShradhaPandey/Struts-1.3-demo-project"
  },
  "decisions_needed": [
    {
      "id": "D1",
      "question": "This codebase is already on Spring Boot 3.5.4 / Java 21 / Thymeleaf. There is no Struts code to migrate. What does 'migrate to modern Java' mean in this context?",
      "options": [
        "Further modernize the existing Spring Boot app (e.g., add REST API, reactive stack, microservices)",
        "Re-migrate from the original Struts 1.3 source repo (github.com/ShradhaPandey/Struts-1.3-demo-project)",
        "Migrate to a different framework entirely (Quarkus, Micronaut, Jakarta EE)",
        "The prompt was applied to the wrong repository"
      ],
      "why_it_matters": "The entire premise of 'Struts to modern Java' is inapplicable -- the Struts migration was already done. Every subsequent decision depends on understanding the actual starting point.",
      "default": null,
      "confidence": "critical"
    },
    {
      "id": "D2",
      "question": "If the intent is to further modernize the current Spring Boot app, what is the target architecture?",
      "options": [
        "Add REST API layer alongside MVC (for SPA or mobile clients)",
        "Convert to reactive stack (Spring WebFlux)",
        "Decompose into microservices",
        "Migrate to Quarkus for cloud-native / GraalVM native compilation",
        "Keep as-is -- no further migration needed"
      ],
      "why_it_matters": "The app is tiny (1 controller, 1 model, 3 templates). The 'right' modernization path depends entirely on the intended production use case.",
      "default": "Keep as-is -- no further migration needed",
      "confidence": "high"
    },
    {
      "id": "D3",
      "question": "What frontend strategy should be used if modernizing further?",
      "options": [
        "Keep Thymeleaf server-side rendering (current state)",
        "Replace with React/Angular/Vue SPA consuming REST APIs",
        "Use HTMX for progressive enhancement over Thymeleaf",
        "No frontend changes needed"
      ],
      "why_it_matters": "The templates already use Thymeleaf + Bootstrap 5 with modern dark-mode, accessibility, and responsive features. A frontend migration is only warranted if the architecture is changing.",
      "default": "Keep Thymeleaf server-side rendering (current state)",
      "confidence": "medium"
    },
    {
      "id": "D4",
      "question": "Does the application need a persistence layer?",
      "options": [
        "No -- the app has no database and operates statelessly",
        "Yes -- add Spring Data JPA with H2/PostgreSQL",
        "Yes -- add a NoSQL store (MongoDB, Redis)"
      ],
      "why_it_matters": "Currently the app has hardcoded validation (username must equal 'Shradha'). A real modernization would likely add user management with a database, but this is a design decision, not a migration one.",
      "default": "No -- the app has no database and operates statelessly",
      "confidence": "medium"
    }
  ],
  "reasoning": {
    "struts_artifacts_found": false,
    "spring_boot_artifacts_found": true,
    "migration_status": "already_completed",
    "original_source_repo": "https://github.com/ShradhaPandey/Struts-1.3-demo-project",
    "prompt_applicability": "The vague prompt 'Migrate this from Struts to modern Java' is fundamentally mismatched with this codebase. A questionnaire skill must detect this mismatch before generating a migration plan, otherwise it will produce nonsensical output."
  }
}
```

---

## Questions the Skill SHOULD Ask (from expected list)

The expected questions assume the repo contains actual Struts 1.3 code:

1. **"Modern Java means what -- Spring Boot? Quarkus? Jakarta EE?"** -- This question probes the vague target specification in the prompt.

2. **"Struts 1.3 ActionForms -- what pattern replaces them?"** -- This question addresses migrating Struts ActionForm classes to modern equivalents (POJOs, DTOs, @ModelAttribute beans, etc.).

3. **"JSP -- what frontend?"** -- This question addresses the view-layer migration from JSP to Thymeleaf, React, etc.

---

## Questions the Skill DID Surface (from your analysis)

1. **D1: "This codebase is already on Spring Boot 3.5.4. There is no Struts code to migrate. What does 'migrate to modern Java' mean in this context?"** -- A meta-question that supersedes all expected questions. The skill correctly detected that the migration has already been completed.

2. **D2: "If the intent is to further modernize, what is the target architecture?"** -- Partially overlaps with expected question #1 (Spring Boot vs Quarkus vs Jakarta EE) but reframed for an already-migrated codebase.

3. **D3: "What frontend strategy should be used if modernizing further?"** -- Overlaps with expected question #3 (JSP -> what frontend?) but correctly notes JSPs are already gone; the question becomes whether to keep Thymeleaf or go further.

4. **D4: "Does the application need a persistence layer?"** -- A new question not in the expected list, surfaced by noticing the app has no database and uses hardcoded validation.

---

## Gap Analysis (what was missed and why)

### Expected questions vs. actual skill output:

| Expected Question | Surfaced? | Analysis |
|---|---|---|
| Modern Java means what -- Spring Boot? Quarkus? Jakarta EE? | Partially (D2) | The skill surfaced this as a secondary question about further modernization, not as the primary target ambiguity question. However, this is arguably correct behavior -- the codebase IS already Spring Boot, so asking "did you mean Spring Boot?" would be redundant. The skill instead asks the right deeper question: "Given it's already Spring Boot, do you want Quarkus or something else?" |
| Struts 1.3 ActionForms -- what pattern replaces them? | Not surfaced | **Correctly omitted.** There are no ActionForms in this codebase. The User model is a Lombok POJO with @Data, and the controller uses @ModelAttribute. Asking about ActionForm replacement would be misleading since the migration is done. |
| JSP -- what frontend? | Partially (D3) | Reframed appropriately. There are no JSPs -- only Thymeleaf templates. The skill asks whether to keep Thymeleaf or move to a SPA framework, which is the correct form of this question for this codebase. |

### Key Gaps:

1. **No gap on ActionForms** -- The expected question about ActionForms is inapplicable. A good skill should NOT ask about migrating ActionForms when none exist. The skill correctly avoided this false positive.

2. **The real gap is at the meta level** -- The expected questions assume the codebase is pre-migration Struts 1.3. The actual codebase is post-migration Spring Boot. A questionnaire skill that asked the three expected questions verbatim would be **wrong** -- it would be generating questions for a codebase it didn't actually analyze.

3. **Missing: should the skill suggest analyzing the original repo?** -- Since the README links to `https://github.com/ShradhaPandey/Struts-1.3-demo-project` as the original Struts 1.3 source, a thorough skill could ask: "Should I analyze the original Struts 1.3 repo instead to plan a migration?"

4. **Missing: packaging/deployment questions** -- The skill did not ask about containerization (Docker/Kubernetes/OpenShift), CI/CD pipeline, or deployment target. For a "modernization" prompt, these are often relevant.

5. **Missing: testing strategy** -- The test suite is a single no-op context-load test. A modernization effort should address test coverage.

---

## Skill Improvement Suggestions

### 1. Add a "Pre-flight Mismatch Detector"
The most important finding from this test case is that the vague prompt ("Migrate from Struts to modern Java") completely contradicts the actual codebase (already on Spring Boot 3.5.4). A questionnaire skill MUST:
- Detect when the stated source technology in the prompt does not match the actual codebase
- Surface this as a blocking question BEFORE generating any migration plan
- Avoid asking migration questions about technology that does not exist in the code

This repo is an excellent adversarial test case: the repo NAME says "struts" and the README describes a Struts 1.3 origin, but the actual code is pure Spring Boot. A naive skill that only reads the README or repo name would generate entirely wrong questions.

### 2. Distinguish "Migration Input" from "Migration Output"
The README explicitly states this is a modernization OF a Struts app. The skill should parse project documentation to determine whether the repo is:
- The source (pre-migration) codebase
- The target (post-migration) codebase
- A hybrid (partially migrated)

### 3. Add Deployment/Infrastructure Questions
When a user says "modern Java," they often implicitly mean modern deployment too. The skill should ask about:
- Containerization targets (Docker, Podman, Buildpacks)
- Orchestration (Kubernetes, OpenShift)
- Cloud-native requirements (health probes, config externalization, 12-factor compliance)

### 4. Add Test Coverage Questions
Detecting that the test suite is essentially empty (1 no-op test for 4 source files) should trigger a question about test strategy -- especially if the goal is modernization, which implies maintainability.

### 5. Handle Trivial Codebases Gracefully
With only 4 Java files and ~60 lines of application logic, a skill should consider whether a "migration plan" is even warranted, or if the answer is simply "this app is too small to need a structured migration -- just rewrite it." The skill should surface complexity metrics and ask whether a formal migration process is appropriate.

### 6. Cross-Reference External Sources
When a README links to an original upstream repo, the skill could offer to analyze that repo instead, since it may be the actual migration source the user intended.
