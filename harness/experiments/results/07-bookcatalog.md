# BookCatalog (.NET Modernization for Beginners) - Questionnaire Skill Evaluation

**Repo:** `/Users/hitpatel/agentic-controller/harness/experiments/repos/dotnet-modernization-for-beginners`
**Vague Prompt:** "Modernize this app"

---

## Detection Output (the questionnaire.json)

```json
{
  "detection": {
    "app_name": "BookCatalog.Web",
    "app_type": "ASP.NET MVC 5 Web Application",
    "source_language": "C# 7.3",
    "source_framework": ".NET Framework 4.8",
    "source_runtime": "IIS / IIS Express",
    "build_system": "MSBuild (legacy .csproj format, non-SDK style)",
    "package_manager": "NuGet (packages.config)",
    "database": "SQL Server LocalDB (via Entity Framework 6.4.4, SqlClient provider)",
    "hosting_model": "IIS-hosted (System.Web pipeline, Global.asax lifecycle)",
    "architecture": "Monolithic MVC (single project, no service layer separation)",
    "test_coverage": "None detected (no test projects in solution)",
    "ci_cd": "GitHub Actions workflows present (.github/workflows/)",
    "repo_purpose": "Educational course repo: teaches .NET Framework 4.8 -> .NET 10 migration using GitHub Copilot modernization agent",
    "secondary_project": "SimpleLegacyApp - console app also targeting net48 with BinaryFormatter (deprecated/removed API)",
    "dependencies": [
      { "name": "EntityFramework", "version": "6.4.4", "category": "ORM", "migration_impact": "high - must migrate to EF Core" },
      { "name": "Microsoft.AspNet.Mvc", "version": "5.2.9", "category": "Web Framework", "migration_impact": "high - must migrate to ASP.NET Core MVC" },
      { "name": "Microsoft.AspNet.Razor", "version": "3.2.9", "category": "View Engine", "migration_impact": "high - Razor syntax changes in ASP.NET Core" },
      { "name": "Microsoft.AspNet.WebPages", "version": "3.2.9", "category": "Web Framework", "migration_impact": "high - removed in ASP.NET Core" },
      { "name": "Microsoft.Web.Infrastructure", "version": "2.0.0", "category": "Web Infrastructure", "migration_impact": "high - not available in ASP.NET Core" },
      { "name": "Newtonsoft.Json", "version": "13.0.3", "category": "Serialization", "migration_impact": "low - works on .NET 8+, but System.Text.Json is the default" }
    ],
    "blocking_apis": [
      "System.Web.HttpContext (not available in ASP.NET Core)",
      "System.Web.Mvc.Controller (replaced by Microsoft.AspNetCore.Mvc.Controller)",
      "System.Web.Routing (replaced by ASP.NET Core routing)",
      "System.Data.Entity.DbContext (replaced by Microsoft.EntityFrameworkCore.DbContext)",
      "Global.asax / HttpApplication (replaced by Program.cs / Startup.cs)",
      "System.Runtime.Serialization.Formatters.Binary.BinaryFormatter (removed in .NET 8+)",
      "System.Configuration.ConfigurationManager (replaced by IConfiguration in ASP.NET Core)",
      "DropCreateDatabaseIfModelChanges (EF6 initializer - no equivalent in EF Core)"
    ],
    "config_concerns": [
      "Web.config XML-based config must migrate to appsettings.json",
      "Connection string uses LocalDB with MDF file attach - needs migration strategy",
      "Assembly binding redirects (handled automatically by .NET Core runtime)",
      "IIS-specific handler/module configuration in system.webServer section"
    ],
    "complexity_assessment": "Low-Medium: Small app (3 source files, 1 controller, 1 model, 1 DbContext), well-structured MVC pattern, no complex integrations, no authentication/authorization, no external service calls, no messaging or caching layers"
  },
  "ambiguities": [
    {
      "id": "AMB-001",
      "category": "target_framework",
      "question": "Modernize to what? What is the target .NET version?",
      "why_it_matters": "The repo's own README targets .NET 10, but the user said only 'modernize.' Possible targets: .NET 8 (LTS, stable), .NET 9, or .NET 10 (latest). Each has different support timelines and feature sets. .NET 8 is the current LTS; .NET 10 is what this course teaches.",
      "options": [".NET 8 (LTS - recommended for production)", ".NET 9", ".NET 10 (latest, matches repo intent)"],
      "default_if_unasked": ".NET 8 (LTS)",
      "confidence_without_answer": 0.4
    },
    {
      "id": "AMB-002",
      "category": "modernization_scope",
      "question": "What does 'modernize' mean to you? Framework upgrade only, or also containerize, re-platform to cloud, or re-architect?",
      "why_it_matters": "Modernization is a spectrum. At minimum it means upgrading from .NET Framework 4.8 to .NET 8+. But it could also mean: containerizing with Docker, deploying to Azure/AWS/OpenShift, splitting into microservices, or adding cloud-native patterns (health checks, observability, managed databases). Each path has vastly different scope.",
      "options": [
        "Framework upgrade only (.NET Framework 4.8 -> .NET 8+)",
        "Framework upgrade + containerization (add Dockerfile)",
        "Framework upgrade + cloud deployment (Azure App Service, AWS, etc.)",
        "Framework upgrade + re-architecture (microservices, API-first)",
        "Full modernization (all of the above)"
      ],
      "default_if_unasked": "Framework upgrade only",
      "confidence_without_answer": 0.3
    },
    {
      "id": "AMB-003",
      "category": "database_strategy",
      "question": "How should the database layer be modernized?",
      "why_it_matters": "The app uses Entity Framework 6 with SQL Server LocalDB (file-based MDF). Migration to EF Core is required for .NET 8+, but the database hosting strategy also matters: keep LocalDB for dev? Move to full SQL Server? Use Azure SQL? Use a different database entirely (PostgreSQL, SQLite)?",
      "options": [
        "Migrate to EF Core with SQLite (simplest for dev/demo)",
        "Migrate to EF Core with SQL Server (production-like)",
        "Migrate to EF Core with Azure SQL (cloud-native)",
        "Migrate to EF Core with PostgreSQL (open source)"
      ],
      "default_if_unasked": "Migrate to EF Core with SQLite (simplest)",
      "confidence_without_answer": 0.5
    },
    {
      "id": "AMB-004",
      "category": "serialization_strategy",
      "question": "The SimpleLegacyApp uses BinaryFormatter which is removed in modern .NET. Should it be migrated to System.Text.Json, or is it out of scope?",
      "why_it_matters": "BinaryFormatter is a security risk and removed in .NET 8+. The SimpleLegacyApp is a separate project in this repo that demonstrates this exact issue. The user may or may not consider it in scope for 'modernize this app.'",
      "options": [
        "Migrate to System.Text.Json (recommended)",
        "Out of scope (focus only on BookCatalog.Web)",
        "Remove SimpleLegacyApp entirely"
      ],
      "default_if_unasked": "Migrate to System.Text.Json",
      "confidence_without_answer": 0.6
    },
    {
      "id": "AMB-005",
      "category": "deployment_target",
      "question": "Where will the modernized app be deployed?",
      "why_it_matters": "Deployment target affects technology choices. The repo's Chapter 04 covers Azure App Service deployment. But 'modernize' might mean deploying to OpenShift, Kubernetes, or just running locally. The deployment target influences whether to add Dockerfiles, health check endpoints, managed identity, or cloud-specific configuration.",
      "options": [
        "Local development only (no deployment target)",
        "Azure App Service (matches repo's course path)",
        "Container platform (Docker/Kubernetes/OpenShift)",
        "On-premises IIS (upgraded runtime only)"
      ],
      "default_if_unasked": "Local development only",
      "confidence_without_answer": 0.3
    },
    {
      "id": "AMB-006",
      "category": "project_scope",
      "question": "Which project(s) should be modernized? The BookCatalog.Web app, the SimpleLegacyApp, or both?",
      "why_it_matters": "The repo contains two separate .NET projects: BookCatalog.Web (the main ASP.NET MVC app) and SimpleLegacyApp (a console demo). They have different modernization paths. The user may mean one or both.",
      "options": [
        "BookCatalog.Web only (primary app)",
        "Both BookCatalog.Web and SimpleLegacyApp",
        "SimpleLegacyApp only"
      ],
      "default_if_unasked": "BookCatalog.Web only",
      "confidence_without_answer": 0.6
    },
    {
      "id": "AMB-007",
      "category": "testing_strategy",
      "question": "Should tests be added as part of modernization? There are currently none.",
      "why_it_matters": "The solution has zero test projects. Adding tests before migration provides a safety net to verify the modernized app behaves identically. Adding them after gives tests that target the new framework. Or the user may not care about tests for what is fundamentally a course/demo app.",
      "options": [
        "Add tests before migration (safety net approach)",
        "Add tests after migration (target new framework)",
        "No tests needed (this is a demo/learning app)"
      ],
      "default_if_unasked": "No tests needed",
      "confidence_without_answer": 0.5
    },
    {
      "id": "AMB-008",
      "category": "csproj_format",
      "question": "The project uses legacy non-SDK-style .csproj. Should the migration include converting to SDK-style .csproj format?",
      "why_it_matters": "SDK-style .csproj is required for .NET 8+ and is dramatically simpler. This conversion is a necessary step but worth confirming the user understands it will change the project file format significantly.",
      "options": [
        "Yes, convert to SDK-style (required for .NET 8+)",
        "The user should be aware this is mandatory, not optional"
      ],
      "default_if_unasked": "Yes (mandatory)",
      "confidence_without_answer": 0.95
    },
    {
      "id": "AMB-009",
      "category": "frontend_modernization",
      "question": "Should the Razor views and CSS be modernized beyond what is strictly required for the framework upgrade?",
      "why_it_matters": "The app uses Razor views with basic CSS. ASP.NET Core Razor is compatible but has differences (_ViewImports.cshtml replaces Web.config in Views folder, tag helpers replace HTML helpers). Beyond the required changes, the user might want to modernize to Razor Pages, Blazor, or add a modern CSS framework.",
      "options": [
        "Minimal changes (just make Razor views work in ASP.NET Core)",
        "Adopt tag helpers and modern Razor conventions",
        "Re-write frontend with Blazor",
        "Keep MVC but add modern CSS (Bootstrap 5, Tailwind)"
      ],
      "default_if_unasked": "Minimal changes",
      "confidence_without_answer": 0.6
    }
  ],
  "recommendations": {
    "likely_migration_path": ".NET Framework 4.8 ASP.NET MVC 5 -> ASP.NET Core MVC on .NET 8 (LTS) or .NET 10",
    "estimated_effort": "Small (1-2 days for experienced developer, given app simplicity)",
    "key_transformations": [
      "Convert legacy .csproj to SDK-style .csproj targeting net8.0 or net10.0",
      "Replace Global.asax with Program.cs (minimal hosting model)",
      "Replace System.Web.Mvc with Microsoft.AspNetCore.Mvc",
      "Migrate Entity Framework 6 to Entity Framework Core",
      "Replace Web.config with appsettings.json",
      "Update Razor views for ASP.NET Core conventions",
      "Replace packages.config with PackageReference in .csproj",
      "Replace BinaryFormatter with System.Text.Json in SimpleLegacyApp",
      "Replace HttpContext.Current with dependency-injected IHttpContextAccessor"
    ],
    "risk_level": "Low - small, well-structured app with no complex dependencies"
  }
}
```

---

## Questions the Skill SHOULD Ask (from expected list)

The expected questions for this maximally vague prompt are:

1. **Modernize to what? .NET 8? Containerize? Re-platform to cloud?**
   - The prompt "Modernize this app" gives zero information about the target. The skill must disambiguate between a framework version upgrade (.NET 8 LTS vs .NET 9 vs .NET 10), containerization (adding Docker support), cloud re-platforming (Azure, AWS, OpenShift), or re-architecture (microservices, Blazor frontend, etc.).

2. **This prompt is maximally vague -- skill must clarify the target.**
   - The skill should recognize that "modernize" is not a specific action. It is an umbrella term that could mean any combination of: upgrade runtime, adopt new patterns, containerize, deploy to cloud, add observability, improve security posture, add tests, modernize the frontend, or restructure the architecture.

---

## Questions the Skill DID Surface (from analysis)

The questionnaire produced **9 distinct ambiguity questions** covering:

| ID | Category | Question Summary | Matches Expected? |
|----|----------|-----------------|-------------------|
| AMB-001 | target_framework | "Modernize to what? What is the target .NET version?" | **YES** - directly matches expected question #1 (.NET 8?) |
| AMB-002 | modernization_scope | "What does 'modernize' mean? Framework upgrade, containerize, cloud, or re-architect?" | **YES** - directly matches expected question #1 (Containerize? Re-platform to cloud?) and #2 (must clarify target) |
| AMB-003 | database_strategy | "How should the database layer be modernized?" | Extends beyond expected - good additional depth |
| AMB-004 | serialization_strategy | "BinaryFormatter is removed - should SimpleLegacyApp be migrated?" | Extends beyond expected - identifies specific technical blocker |
| AMB-005 | deployment_target | "Where will the modernized app be deployed?" | **YES** - matches expected question #1 (Re-platform to cloud?) |
| AMB-006 | project_scope | "Which project(s) should be modernized?" | Extends beyond expected - important since repo has 2 projects |
| AMB-007 | testing_strategy | "Should tests be added? There are currently none." | Extends beyond expected - relevant modernization concern |
| AMB-008 | csproj_format | "Convert to SDK-style .csproj?" | Informational - this is mandatory, not truly a question |
| AMB-009 | frontend_modernization | "Should Razor views be modernized beyond minimum?" | Extends beyond expected - reasonable scope question |

---

## Gap Analysis (what was missed and why)

### What the Skill Got Right

1. **Core ambiguity fully captured.** The two expected questions are both directly addressed:
   - AMB-001 asks about target .NET version (.NET 8, .NET 9, .NET 10)
   - AMB-002 asks what "modernize" means (upgrade only vs. containerize vs. cloud vs. re-architect)
   - AMB-005 asks about deployment target (Azure, containers, on-prem)
   
2. **Went beyond the expected questions** with 7 additional domain-specific questions that demonstrate genuine understanding of the codebase:
   - Database migration strategy (EF6 -> EF Core, database provider choice)
   - BinaryFormatter security concern (removed API in modern .NET)
   - Multi-project scope (two projects in the repo)
   - Testing gap (zero tests currently)
   - Frontend modernization depth (Razor views, CSS)

3. **Correct detection of all blocking APIs and dependencies.** The detection output accurately identifies every System.Web dependency, EF6, BinaryFormatter, and the legacy .csproj format.

### What the Skill Could Improve

1. **AMB-008 is not a real question.** Converting to SDK-style .csproj is mandatory for .NET 8+, not optional. The skill correctly notes "confidence_without_answer: 0.95" and "mandatory," but including it as a question creates noise. It should be stated as a fact in the detection output, not presented as a decision point.

2. **Missing: authentication/authorization modernization question.** The app has no auth today, but modernization often includes adding it. The skill could ask whether security features should be added as part of modernization (e.g., ASP.NET Core Identity, OAuth, Azure AD).

3. **Missing: observability/monitoring question.** Modern apps typically add health checks, structured logging (Serilog/OpenTelemetry), and metrics. The skill does not ask whether these cloud-native patterns should be adopted.

4. **Missing: API-first question.** The current app is server-rendered MVC. A common modernization path is to add a REST API layer (Web API controllers) alongside or instead of the MVC views. This was not asked.

5. **Repo context not fully leveraged.** The repo is an educational course with a pre-defined migration path (.NET Framework 4.8 -> .NET 10, then Azure deployment). The skill detected this in the README but could have used it more prominently to suggest: "This repo already defines a migration path in its course structure. Should we follow that path, or do you have different goals?"

6. **Priority ordering missing.** The 9 questions are not prioritized. AMB-001 and AMB-002 are blocking -- nothing can proceed without answers. AMB-003 through AMB-009 are secondary. The questionnaire should indicate which questions must be answered first.

---

## Skill Improvement Suggestions

### 1. Prioritize Questions into Tiers

Split questions into:
- **Tier 1 (Blocking):** Must be answered before any work begins. For this case: target framework version (AMB-001) and modernization scope (AMB-002).
- **Tier 2 (Important):** Affect the migration plan but have reasonable defaults. For this case: database strategy, deployment target, project scope.
- **Tier 3 (Nice-to-have):** Can be deferred or defaulted without significant risk. For this case: testing, frontend depth, serialization strategy.

### 2. Remove Non-Questions

AMB-008 (SDK-style .csproj) is not a decision -- it is a requirement. Move mandatory transformations to the detection/facts section. Only present genuine decision points as questions.

### 3. Add "Repo Context Awareness" Heuristic

When a repo's README or documentation explicitly states a migration target (as this one does: ".NET Framework 4.8 to .NET 10"), the skill should surface this as a strong default and ask whether the user wants to follow the documented path or diverge. This reduces question count for repos with clear intent.

### 4. Add Cloud-Native Patterns Question

For any web application modernization, add a standard question about cloud-native patterns:
- Health check endpoints (`/health`, `/ready`)
- Structured logging (Serilog, OpenTelemetry)
- Configuration via environment variables (12-factor)
- Containerization (Dockerfile, docker-compose)

### 5. Add Security Posture Question

Ask whether security should be added or upgraded as part of modernization, especially when the current app has no authentication (as is the case here).

### 6. Detect and Surface "This Is a Learning Repo" Signal

The repo name (`dotnet-modernization-for-beginners`), the course structure (chapters 00-04), and the README all indicate this is educational material, not a production application. The skill should detect this and adjust its questions accordingly -- for example, suggesting that the user may want to follow the pre-defined course path rather than define a custom migration.

### 7. Confidence Scoring Calibration

The confidence scores are generally reasonable but AMB-002 (modernization scope) at 0.3 is appropriately low -- "modernize" truly is maximally ambiguous. However, AMB-006 (project scope) at 0.6 may be too high given that the user said "this app" (singular) but the repo contains two projects. A lower score (0.4) would better reflect the genuine ambiguity.

### 8. Add a "What I Already Know" Summary

Before asking questions, the skill should present a brief summary of what it detected and is confident about. This gives the user context and lets them correct any misdetections before answering scope questions. For example: "I detected a .NET Framework 4.8 ASP.NET MVC 5 web application called BookCatalog.Web with Entity Framework 6 and SQL Server LocalDB. Is this the application you want to modernize?"
