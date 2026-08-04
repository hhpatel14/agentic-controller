# Legacy Cycle Store MVC App - Migration Questionnaire Analysis

**Source Application**: AdventureWorksMVC (Legacy Cycles) - ASP.NET MVC 4 / .NET Framework 4.5
**Vague Migration Prompt**: "Migrate this to .NET 8"

---

## Detection Output (the questionnaire.json)

```json
{
  "questionnaire": {
    "metadata": {
      "app_name": "AdventureWorksMVC (Legacy Cycles)",
      "source_framework": "ASP.NET MVC 4.0 (.NET Framework 4.5)",
      "target_framework": ".NET 8",
      "solution_file": "AdventureWorksMVC_2013.sln",
      "project_type": "ASP.NET MVC Web Application (Library output, IIS-hosted)",
      "visual_studio_version": "Visual Studio 2013",
      "complexity": "medium",
      "estimated_risk": "high"
    },
    "detection": {
      "framework_stack": {
        "runtime": ".NET Framework 4.5",
        "web_framework": "ASP.NET MVC 4.0.20710.0",
        "web_api": "ASP.NET Web API 4.0.20710.0",
        "view_engine": "Razor 2.0 (System.Web.Razor)",
        "orm": "Entity Framework 5.0.0 (Database-First via EDMX)",
        "json": "Newtonsoft.Json 13.0.1",
        "project_format": "Legacy .csproj (MSBuild ToolsVersion 4.0, packages.config)"
      },
      "database": {
        "provider": "System.Data.SqlClient (SQL Server)",
        "connection_string_format": "EDMX metadata-style (EntityClient provider)",
        "database_name": "CYCLE_STORE",
        "schema": "Production",
        "host": "AWS RDS (sqlrdsdb.xxxxx.us-east-1.rds.amazonaws.com)",
        "tables": ["Production.Product", "Production.ProductCategory", "Production.ProductSubcategory"],
        "entity_model": "CycleModel.edmx (EDMX Database-First with T4 code generation)"
      },
      "authentication": {
        "mode": "Forms Authentication",
        "membership_provider": "SqlMembershipProvider (ASP.NET Membership)",
        "cookie_name": "ADVENTUREWORKS.AUTH",
        "login_url": "~/Home/Login"
      },
      "vendored_dependencies": {
        "c1_report": "C1.C1Report.4 (ComponentOne/GrapeCity reporting - no HintPath, likely GAC or local)",
        "c1_wijmo": "C1.Web.Wijmo.Controls.4 (ComponentOne Wijmo UI controls - no HintPath)",
        "c1_wijmo_design": "C1.Web.Wijmo.Controls.Design.4 (Wijmo design-time support - no HintPath)"
      },
      "windows_specific_apis": [
        "System.Web.HttpApplication (Global.asax)",
        "System.Web.UI.Page (WebResource.cs helper - used for GetWebResourceUrl)",
        "System.Web.UI.HtmlControls (imported in SiteLayoutController)",
        "System.Web.Security.SqlMembershipProvider (ASP.NET Membership)",
        "System.Drawing (referenced in csproj)",
        "System.EnterpriseServices (referenced in csproj)",
        "System.Data.Services / System.Data.Services.Client / System.Data.Services.Design (WCF Data Services)",
        "System.ServiceModel.Web (WCF WebHttp)",
        "System.Data.Entity (EF legacy assembly)",
        "System.Web.DynamicData",
        "System.Web.Entity",
        "System.Web.ApplicationServices",
        "System.Runtime.Interop.ComVisible (AssemblyInfo)"
      ],
      "hosting": {
        "server": "IIS / IIS Express",
        "pipeline_mode": "Integrated (with classic mode fallback handlers in Web.config)",
        "local_url": "http://localhost:55185/",
        "deployment": "Web Deploy (pubxml publish profile present)"
      },
      "code_generation": {
        "t4_templates": ["CycleModel.Context.tt", "CycleModel.tt"],
        "edmx_designer": "CycleModel.Designer.cs (auto-generated from EDMX)",
        "entities_generated": ["Product.cs", "ProductCategory.cs", "ProductSubcategory.cs"]
      },
      "razor_views": [
        "Views/Shared/_SiteLayout.cshtml (master layout)",
        "Views/Home/Default.cshtml",
        "Views/SiteLayout/HeaderLayout.cshtml",
        "Views/SiteLayout/ContentLayout.cshtml",
        "Views/Error/Default.cshtml"
      ],
      "architectural_patterns": {
        "data_access": "Static manager classes (ProductManager, CategoryManager) with static Common.DataEntities property creating new DbContext per call",
        "controllers": "MVC Controllers + Web API (separate route configs)",
        "no_dependency_injection": true,
        "global_asax_startup": true,
        "child_actions": "Html.RenderAction used for partial view composition"
      }
    },
    "decisions": [
      {
        "id": "D1",
        "category": "ORM Migration",
        "question": "Entity Framework 5.0 (Database-First with EDMX) is used. EF Core does not support EDMX/Database-First designer. Should we migrate to EF Core Code-First with reverse-engineered models (scaffold from existing DB), or use EF Core with manual model definitions matching the current schema?",
        "options": [
          "Scaffold EF Core models from existing CYCLE_STORE database using dotnet ef dbcontext scaffold",
          "Manually rewrite entity classes as Code-First with Fluent API or Data Annotations matching current EDMX schema",
          "Keep EF6 temporarily (.NET 8 does not support EF6 -- this is NOT viable)"
        ],
        "default": "Scaffold EF Core models from existing database",
        "impact": "high",
        "reasoning": "The EDMX file defines 3 entities (Product, ProductCategory, ProductSubcategory) with relationships. The T4-generated entity classes are simple POCOs already, so migration to EF Core Code-First is straightforward. However, the EDMX contains explicit column mappings (nvarchar lengths, money types, precision/scale), navigation properties, and schema metadata (Production schema) that must be preserved. The connection string format also changes from EntityClient metadata style to standard SqlClient."
      },
      {
        "id": "D2",
        "category": "Schema Breaking Changes",
        "question": "The current EDMX maps to SQL Server 'money' type for StandardCost/ListPrice (Precision 19, Scale 4). EF Core maps decimal differently. The 'Production' schema prefix is used for all tables. Are there any stored procedures, triggers, or views in the database that depend on the current schema that could break during migration?",
        "options": [
          "Schema is fully captured in the SQL file and EDMX -- no additional DB objects to worry about",
          "There are additional database objects (stored procs, views, triggers) not visible in the repo that need audit",
          "The database is shared with other applications and schema changes are not permitted"
        ],
        "default": "Audit the live database before finalizing EF Core model configuration",
        "impact": "high",
        "reasoning": "The SQL file (CYCLE_STORE_Schema_data.sql) shows table definitions only. The EDMX uses UseLegacyProvider=true. EF Core handles money/decimal mapping differently and the Production schema must be explicitly configured. Without auditing the live database, we cannot know if there are views, stored procedures, or other objects that depend on the current table structure."
      },
      {
        "id": "D3",
        "category": "View Layer Strategy",
        "question": "The app uses ASP.NET MVC 4 Razor views (5 .cshtml files) with Html.RenderAction for partial composition. Should the Razor views be kept as server-side rendered views in ASP.NET Core MVC, migrated to Razor Pages, or should the front-end move to Blazor or a separate SPA with a Web API backend?",
        "options": [
          "Keep as ASP.NET Core MVC with Razor views (minimal view changes, straightforward migration)",
          "Convert to Razor Pages (better for page-focused scenarios, requires restructuring controllers)",
          "Move to Blazor Server (modern .NET UI, but significant rewrite of view logic)",
          "Move to Blazor WebAssembly with Web API backend (full SPA approach, major rewrite)",
          "Convert to API-only backend with separate SPA frontend (React/Angular/Vue)"
        ],
        "default": "Keep as ASP.NET Core MVC with Razor views",
        "impact": "high",
        "reasoning": "The app has only 5 Razor views with simple model binding and Html.RenderAction calls. The views use @model, @ViewBag, and @foreach -- all compatible with ASP.NET Core Razor with minor syntax updates. Html.RenderAction must become ViewComponents or partial views with <partial> tag helper. The ContentLayout.cshtml already uses 'asp-for' tag helper syntax mixed with MVC 4 patterns, suggesting partial modernization was attempted. The simplest path is keeping MVC+Razor and updating the incompatible helpers."
      },
      {
        "id": "D4",
        "category": "Third-Party UI Controls",
        "question": "The project references ComponentOne/GrapeCity Wijmo controls (C1.C1Report.4, C1.Web.Wijmo.Controls.4) which are Windows/.NET Framework-only WebForms controls. These have NO .NET 8 equivalent. What should replace them?",
        "options": [
          "Replace with GrapeCity's modern Wijmo 5 (JavaScript-based, compatible with any backend)",
          "Replace with alternative .NET 8-compatible UI libraries (Telerik, DevExpress, Syncfusion)",
          "Remove reporting/rich-UI features and use plain HTML/CSS/Bootstrap",
          "Evaluate if these controls are actually used in the current codebase (references exist but may be unused)"
        ],
        "default": "Evaluate usage first, then replace with Wijmo 5 JavaScript or alternative",
        "impact": "high",
        "reasoning": "The C1/Wijmo references have no HintPath pointing to NuGet packages, suggesting they were GAC-installed or locally referenced. The current source files do not show explicit usage of Wijmo controls in the committed code, but the second csproj variant references additional views (Products/Index, ShoppingCart/OpenCart, etc.) that likely use these controls. These WebForms-based server controls cannot run on .NET 8. GrapeCity's modern Wijmo 5 is a JavaScript library that works with any backend."
      },
      {
        "id": "D5",
        "category": "Authentication Migration",
        "question": "The app uses ASP.NET Membership (SqlMembershipProvider) with Forms Authentication. This system is completely removed in .NET 8. What should replace it?",
        "options": [
          "Migrate to ASP.NET Core Identity (recommended replacement, requires new Identity tables)",
          "Migrate to external identity provider (Azure AD, Auth0, Keycloak)",
          "Implement custom cookie authentication with ASP.NET Core authentication middleware",
          "If no active users exist, remove authentication entirely and re-add later"
        ],
        "default": "Migrate to ASP.NET Core Identity",
        "impact": "high",
        "reasoning": "SqlMembershipProvider stores user credentials in a specific schema (aspnet_Membership tables). ASP.NET Core Identity uses a different schema. If the app has existing users, a data migration strategy is needed to move hashed passwords (noting that the password format is 'Hashed' which uses SHA1 -- not compatible with ASP.NET Core Identity's PBKDF2). The connection string references 'ASPNetDB' which is separate from the main CYCLE_STORE database."
      },
      {
        "id": "D6",
        "category": "Windows-Specific API Dependencies",
        "question": "Multiple Windows/.NET Framework-only APIs are referenced: System.Drawing, System.EnterpriseServices, System.Data.Services (WCF Data Services), System.ServiceModel.Web (WCF), System.Web.UI (WebForms Page class used in WebResource.cs helper), System.Web.DynamicData. Are these actively used or vestigial references?",
        "options": [
          "Audit each reference for active usage and remove unused ones; replace active ones with .NET 8 equivalents",
          "Assume all are vestigial (typical in VS2013 project templates) and remove them all",
          "Some are actively used and require specific replacement strategies"
        ],
        "default": "Audit and remove unused; WebResource.cs helper definitely needs rewriting",
        "impact": "medium",
        "reasoning": "System.Web.UI.Page is actively used in WebResource.cs to generate WebResource.axd URLs -- this pattern does not exist in ASP.NET Core. System.Drawing is referenced but no image manipulation code is visible. System.EnterpriseServices, WCF Data Services, and WCF WebHttp appear to be default VS2013 template references. The WebResource.cs helper must be rewritten to use ASP.NET Core's static file middleware or embedded resource serving."
      },
      {
        "id": "D7",
        "category": "Project Structure",
        "question": "The solution has a legacy .csproj format with packages.config and contains two copies of the .csproj file (one at project root, one in Business/ directory -- likely a version control artifact). The project uses Global.asax for startup. Should we convert to SDK-style .csproj with a single project, or split into multiple projects?",
        "options": [
          "Single SDK-style .csproj with Program.cs/Startup.cs replacing Global.asax",
          "Split into Web project + Business/Data Access class library",
          "Split into Web API + separate frontend project"
        ],
        "default": "Single SDK-style .csproj (app is small enough)",
        "impact": "medium",
        "reasoning": "The app has only 19 .cs files, 5 views, and 3 entities. Splitting would add unnecessary complexity. The Business/ folder can remain as a namespace boundary within a single project. Global.asax Application_Start must be replaced with Program.cs top-level statements and builder.Services / app.Use middleware pipeline."
      },
      {
        "id": "D8",
        "category": "Data Access Pattern",
        "question": "The app uses a static Common.DataEntities property that creates a new DbContext on every property access (no disposal, no dependency injection). This is an anti-pattern. Should DI be introduced during migration?",
        "options": [
          "Introduce proper DI with scoped DbContext lifetime (ASP.NET Core standard pattern)",
          "Keep static pattern but add proper disposal (minimal change, still anti-pattern)",
          "Full repository pattern with DI"
        ],
        "default": "Introduce DI with scoped DbContext (required by ASP.NET Core conventions)",
        "impact": "medium",
        "reasoning": "ASP.NET Core strongly favors constructor injection. The static manager classes (ProductManager, CategoryManager) would need refactoring to accept DbContext via constructor injection or method parameters. The current pattern leaks DbContext instances (no using/Dispose). This is a good opportunity to fix a correctness bug alongside the migration."
      },
      {
        "id": "D9",
        "category": "Hosting Model",
        "question": "The app targets IIS (with IIS Express for development) and uses ISAPI handlers for extensionless URLs. What hosting model should be used in .NET 8?",
        "options": [
          "Kestrel standalone (cross-platform, recommended for .NET 8)",
          "Kestrel behind IIS reverse proxy (if staying on Windows/IIS infrastructure)",
          "Kestrel behind nginx/Apache (for Linux deployment)",
          "Container-based deployment (Docker/Kubernetes)"
        ],
        "default": "Depends on deployment target -- Kestrel standalone or behind reverse proxy",
        "impact": "medium",
        "reasoning": "The current Web.config has IIS-specific handler registrations (ISAPI, integrated mode). These are replaced by ASP.NET Core's built-in middleware. The database is on AWS RDS, suggesting cloud deployment is possible. The publish profile (pubxml) suggests Web Deploy to IIS was the deployment method."
      },
      {
        "id": "D10",
        "category": "Configuration Migration",
        "question": "All configuration is in Web.config (XML). .NET 8 uses appsettings.json. How should configuration be migrated?",
        "options": [
          "Convert all settings to appsettings.json with environment-specific overrides",
          "Use environment variables for sensitive values (connection strings, credentials)",
          "Use Azure Key Vault or AWS Secrets Manager for secrets"
        ],
        "default": "appsettings.json with environment variables for secrets",
        "impact": "low",
        "reasoning": "The Web.config contains connection strings with embedded credentials (currently placeholder xxxxx values), entityFramework section, authentication settings, and custom errors. All must be migrated to appsettings.json format. The EntityClient connection string format must change to a standard SqlClient connection string for EF Core."
      }
    ]
  }
}
```

---

## Questions the Skill SHOULD Ask (from expected list)

1. **EF6 to EF Core -- any breaking schema changes?** The expected question probes whether the Entity Framework migration will cause schema-level breakage in the database (e.g., column type mapping differences, naming convention changes, missing navigation properties, precision/scale mismatches).

2. **Razor views stay or move to Blazor/API?** The expected question asks whether the existing Razor view layer should remain as server-rendered MVC views, be modernized to Blazor, or be replaced entirely with an API backend + SPA frontend.

3. **Windows-specific APIs used?** The expected question probes whether the codebase depends on APIs that only run on Windows/.NET Framework and would not work on .NET 8's cross-platform runtime.

---

## Questions the Skill DID Surface (from your analysis)

1. **D1 (ORM Migration)**: Directly addresses EF5 Database-First/EDMX to EF Core Code-First migration. Covers the need to scaffold or rewrite entity models. **MATCHES expected question #1 (partially).**

2. **D2 (Schema Breaking Changes)**: Explicitly asks about schema breaking changes -- money type mapping, Production schema prefix, stored procedures/triggers/views. **MATCHES expected question #1 (directly).**

3. **D3 (View Layer Strategy)**: Asks whether Razor views should stay as MVC Razor, convert to Razor Pages, move to Blazor (Server or WebAssembly), or become an API+SPA. **MATCHES expected question #2 (directly).**

4. **D4 (Third-Party UI Controls)**: Identifies ComponentOne/Wijmo WebForms controls as incompatible with .NET 8. This is a Windows-specific dependency. **PARTIALLY matches expected question #3.**

5. **D5 (Authentication Migration)**: Identifies SqlMembershipProvider removal. This is both a Windows-specific API concern and a distinct migration decision. **PARTIALLY matches expected question #3.**

6. **D6 (Windows-Specific API Dependencies)**: Explicitly catalogs System.Drawing, System.EnterpriseServices, WCF Data Services, System.Web.UI, System.Web.DynamicData. **MATCHES expected question #3 (directly).**

7. **D7 (Project Structure)**: Addresses legacy csproj to SDK-style and Global.asax replacement. **ADDITIONAL -- not in expected list but necessary.**

8. **D8 (Data Access Pattern)**: Flags the static DbContext anti-pattern. **ADDITIONAL -- not in expected list but important for correctness.**

9. **D9 (Hosting Model)**: Addresses IIS to Kestrel transition. **ADDITIONAL -- not in expected list but relevant.**

10. **D10 (Configuration Migration)**: Web.config to appsettings.json. **ADDITIONAL -- not in expected list but standard migration concern.**

---

## Gap Analysis (what was missed and why)

### Coverage of Expected Questions

| Expected Question | Covered? | Decision IDs | Notes |
|---|---|---|---|
| EF6 to EF Core -- any breaking schema changes? | YES | D1, D2 | Covered thoroughly. D1 handles the ORM migration strategy (EDMX to Code-First). D2 specifically probes schema breaking changes (money type, Production schema, stored procs). The questionnaire correctly identifies that EF5 (not EF6) is in use, which is actually an even older version. |
| Razor views stay or move to Blazor/API? | YES | D3 | Covered with 5 options ranging from keep-as-is to full SPA rewrite. Correctly notes Html.RenderAction must become ViewComponents. |
| Windows-specific APIs used? | YES | D4, D5, D6 | Covered across three decisions. D6 is the direct hit; D4 (Wijmo) and D5 (SqlMembershipProvider) are specific instances elevated to their own decisions due to high impact. |

### Gaps and Weaknesses

1. **No explicit question about EF version accuracy**: The detection says EF 5.0.0 (correct per packages.config) but decisions sometimes say "EF6" in shorthand. The skill should be precise: EF 5 to EF Core is an even bigger jump than EF6 to EF Core because EF5's DbContext API has fewer features.

2. **Missing question about the dual-csproj discrepancy**: The Business/ subdirectory contains a second copy of the .csproj with additional views (Products/Index, ShoppingCart/OpenCart, ShoppingCart/CheckOut, ShoppingCart/Shipping, ShoppingCart/ReviewOrder, ShoppingCart/OrderComplete, Home/Login, Home/Logout) and additional models (ReviewOrderModel, ProductsModel) that are NOT present in the committed source. This suggests the repo is an incomplete snapshot. The skill should ask: "The project references views and models (ShoppingCart, Products, Login) that are not present in the repository. Are these files missing from version control, or has the app been intentionally stripped down?"

3. **Missing question about SQL Server version compatibility**: The EDMX uses ProviderManifestToken="2008" (SQL Server 2008). .NET 8 + EF Core may have different minimum SQL Server version requirements. The skill should ask whether the production SQL Server version is compatible.

4. **No question about test coverage**: There are zero test projects in the solution. The skill should ask whether tests should be written as part of the migration to validate functional equivalence.

5. **No question about the AWS RDS hosting**: The database is on AWS RDS. The skill should ask whether the target deployment environment changes (e.g., moving from IIS on-prem + RDS to a fully cloud-native deployment like ECS/EKS + RDS, or Azure App Service + Azure SQL).

6. **The `asp-for` tag helper in ContentLayout.cshtml is a red flag the skill did not deeply analyze**: The view uses `<div asp-for="@category.Name">` which is ASP.NET Core Tag Helper syntax in what is supposed to be an MVC 4 view. This indicates either the code was partially migrated already, or it is invalid markup that happens to render as a plain attribute. The skill should flag this as a potential inconsistency.

---

## Skill Improvement Suggestions

### 1. Detect Incomplete Repositories
The skill should diff the .csproj `<Compile>` and `<Content>` items against actual files on disk. When referenced files are missing (like the ShoppingCart views, Products views, Login/Logout views, and ReviewOrderModel/ProductsModel), the questionnaire should raise a blocking question: "N files referenced in the project are missing from the repository. This migration cannot be fully planned without them."

### 2. EF Version Precision
When the packages.config says `EntityFramework 5.0.0`, the questionnaire must say "EF5" not "EF6." EF5 has notable differences from EF6 (e.g., enum support limitations, spatial type handling, migration support). The migration path from EF5 to EF Core is strictly harder than EF6 to EF Core.

### 3. Probe the Database Layer More Deeply
The skill should:
- Parse the SQL schema file to identify constraints, indexes, foreign keys
- Check for stored procedures, views, triggers mentioned in any SQL files
- Identify the ProviderManifestToken version and flag SQL Server compatibility concerns
- Note the connection string uses `MultipleActiveResultSets=True` which has different behavior implications in EF Core

### 4. Identify Anti-Patterns as Migration Risks
The static `Common.DataEntities` pattern that creates a new `DbContext` on every property access without disposal is a memory leak. The skill should not just note this as a "nice to fix" but as a **migration risk** because EF Core's behavior with undisposed contexts differs from EF5/EF6 (EF Core tracks more state, making leaks more impactful).

### 5. Vendored/GAC Dependencies Need Stronger Red Flags
The C1/Wijmo references have no HintPath, meaning they were likely installed in the GAC or as a local reference. The skill should:
- Flag that these binaries are not in the repository
- Warn that GAC is not available in .NET 8
- Ask whether the team has licenses for modern GrapeCity products
- Check if the controls are actually used in the committed source code (they may be vestigial)

### 6. Cross-Reference Views for Tag Helper / WebForms Control Usage
The skill should scan .cshtml files for:
- `@Html.RenderAction` (must become ViewComponents in ASP.NET Core)
- `asp-*` tag helpers (valid in Core, invalid in MVC 4 -- indicates code inconsistency)
- Any WebForms-style `<asp:*>` controls or Wijmo `<C1*>` controls
- `@Url.Content("~/...")` (still works in Core but Static Files middleware is preferred)

### 7. Authentication Migration Needs Password Hash Compatibility Check
The SqlMembershipProvider uses SHA1 hashing (passwordFormat="Hashed"). ASP.NET Core Identity uses PBKDF2. The skill should explicitly ask: "Do existing user accounts need to be preserved? If so, a custom IPasswordHasher that supports legacy SHA1 hashes will be needed, or a forced password reset is required."

### 8. Configuration Completeness Check
The Web.config references a connection string named "ASPNetDB" for the membership provider, but this connection string is not defined in the committed Web.config. The skill should flag missing configuration entries as migration risks.

### 9. Dual .csproj Detection
The repository has two .csproj files with the same ProjectGuid but different content (the Business/ variant has more views and models). The skill should detect and flag this as a potential version control problem that must be resolved before migration begins.

### 10. Lifecycle/Support Consideration
The skill should note that .NET 8 is an LTS release (supported until November 2026) and ask whether the team should target .NET 8 or plan for .NET 9+ given the migration timeline.
