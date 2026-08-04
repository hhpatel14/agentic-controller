# 08 - Flask-to-FastAPI Migration Questionnaire Evaluation

**Vague prompt:** "Migrate this to FastAPI"

---

## Detection Output (the questionnaire.json)

```json
{
  "detection": {
    "source_framework": "Flask 2.2.x",
    "target_framework": "FastAPI (user-specified)",
    "language": "Python",
    "build_system": "pip / requirements.txt (no pyproject.toml, no setup.py)",
    "app_type": "Server-rendered web application with OAuth2/OIDC authentication",
    "dependencies": {
      "flask": ">=2.2.2",
      "authlib": ">=1.2.0 (Flask OAuth integration via authlib.integrations.flask_client)",
      "requests": ">=2.28.2",
      "auth0-python": "==4.1.0"
    },
    "config_mechanism": "configparser (.config INI file) with fallback to hardcoded dict",
    "template_engine": "Jinja2 (5 templates: base.html, nav_bar.html, heels.html, home.html, profile.html)",
    "routing_pattern": "Flask Blueprints (auth_bp, webapp_bp) registered on main app",
    "auth_mechanism": "Auth0 via Authlib OAuth + Flask session-based auth with custom @requires_auth decorator",
    "session_management": "Flask server-side sessions (flask.session)",
    "static_assets": "styles.css, scripts.js served from /static/",
    "database": "None",
    "tests": "None detected",
    "vendored_dependencies": "None",
    "api_style": "HTML-rendering (render_template), not REST/JSON API",
    "async_usage": "None - fully synchronous Flask app",
    "custom_jinja_filters": ["to_pretty_json (utils.py)"],
    "deployment": "flask --app app run --port 4040 --reload (dev server)"
  },
  "decisions_needed": [
    {
      "id": "Q1",
      "question": "Sync or async endpoints?",
      "category": "runtime_model",
      "why_it_matters": "FastAPI supports both sync (def) and async (async def) route handlers. The Flask app is fully synchronous. Auth0 OAuth token exchange is I/O-bound and benefits from async. Choosing async affects how all middleware, dependencies, and external HTTP calls must be written.",
      "options": [
        "Keep endpoints synchronous (def) -- minimal change, FastAPI runs them in a threadpool",
        "Convert to async (async def) -- idiomatic FastAPI, better for I/O-bound OAuth flows, requires httpx instead of requests",
        "Mixed -- async for auth/callback routes, sync for template-rendering routes"
      ],
      "detected_signal": "All Flask routes are sync. The reference fastapi-webapp uses async for login/callback (I/O-bound OAuth), sync for logout. requests library would need to become httpx for async.",
      "default_recommendation": "Mixed: async for OAuth routes, sync for simple template routes"
    },
    {
      "id": "Q2",
      "question": "Use Pydantic models for request/response validation?",
      "category": "data_validation",
      "why_it_matters": "FastAPI's core value proposition is automatic request validation and OpenAPI schema generation via Pydantic. However, this app renders HTML templates and has no JSON API endpoints, so there is little request data to validate beyond session state.",
      "options": [
        "Yes -- add Pydantic models for config, session data, and any future API endpoints",
        "No -- this is a template-rendering app with no JSON APIs, Pydantic adds overhead with little benefit",
        "Partial -- use Pydantic for config/settings (BaseSettings) but not for route parameters"
      ],
      "detected_signal": "No JSON request/response bodies. Config is loaded via configparser dict. Session stores opaque Auth0 token dict. No data validation exists in the Flask app.",
      "default_recommendation": "Partial: use Pydantic BaseSettings for config (.config to .env migration), skip models for routes"
    },
    {
      "id": "Q3",
      "question": "Keep Jinja2 templates or go API-only?",
      "category": "rendering_strategy",
      "why_it_matters": "The Flask app is entirely server-rendered HTML using Jinja2 templates (5 templates with inheritance). FastAPI supports Jinja2 via Jinja2Templates but it is not the idiomatic pattern -- FastAPI is primarily designed for JSON APIs. Keeping templates means adding starlette.templating dependency and maintaining template rendering. Going API-only means building a separate frontend (React, Vue, etc.).",
      "options": [
        "Keep Jinja2 templates -- minimal rewrite, use fastapi.templating.Jinja2Templates with Starlette's TemplateResponse",
        "Go API-only (JSON responses) -- more idiomatic FastAPI, but requires a separate frontend",
        "Hybrid -- serve templates for existing pages, add JSON API endpoints for new features"
      ],
      "detected_signal": "5 Jinja2 templates with inheritance (base.html extends), custom filter (to_pretty_json), static assets. Templates use session data and Flask's request context. Profile template accesses session['userinfo'] directly.",
      "default_recommendation": "Keep Jinja2 templates: this is a small app, rewriting the frontend is disproportionate effort"
    },
    {
      "id": "Q4",
      "question": "How should session management be handled in FastAPI?",
      "category": "auth_session",
      "why_it_matters": "Flask has built-in server-side sessions (flask.session). FastAPI/Starlette has SessionMiddleware (cookie-based, signed but not encrypted, limited size). The Flask app stores the full Auth0 token in the session, which may exceed cookie size limits. This is a critical architectural decision.",
      "options": [
        "Use Starlette SessionMiddleware (cookie-based) -- simple but has 4KB size limit, may not fit full Auth0 token",
        "Use a server-side session backend (e.g., redis, database-backed) -- preserves Flask's behavior, more complex",
        "Store only essential data (id_token, userinfo) in cookie session, not the full token -- reference implementation takes this approach"
      ],
      "detected_signal": "Flask app stores full token dict in session['user']. Reference FastAPI app splits into session['id_token'] and session['userinfo']. Starlette SessionMiddleware is cookie-based with size constraints.",
      "default_recommendation": "Store only id_token and userinfo in cookie session (matches reference implementation pattern)"
    },
    {
      "id": "Q5",
      "question": "How should the Authlib integration be migrated?",
      "category": "auth_library",
      "why_it_matters": "The Flask app uses authlib.integrations.flask_client.OAuth which is Flask-specific. FastAPI/Starlette requires authlib.integrations.starlette_client.OAuth. The API is similar but not identical -- Starlette's OAuth requires passing the Request object explicitly rather than relying on Flask's global request context.",
      "options": [
        "Switch to authlib.integrations.starlette_client.OAuth -- direct equivalent, Authlib supports both",
        "Use a different OAuth library (e.g., fastapi-oauth, python-jose) -- more FastAPI-native but more work",
        "Implement OAuth manually with httpx -- full control but significantly more code"
      ],
      "detected_signal": "authlib>=1.2.0 already in requirements. Authlib provides both flask_client and starlette_client integrations. Reference implementation already uses starlette_client.",
      "default_recommendation": "Switch to authlib.integrations.starlette_client.OAuth -- same library, different integration module"
    },
    {
      "id": "Q6",
      "question": "How should the @requires_auth decorator pattern be migrated?",
      "category": "auth_pattern",
      "why_it_matters": "Flask uses a custom @requires_auth decorator that checks flask.session and redirects. FastAPI's idiomatic equivalent is Depends() with dependency injection. This changes the auth enforcement pattern fundamentally.",
      "options": [
        "Convert to FastAPI Depends() -- idiomatic, testable, composable",
        "Keep as a decorator -- works in FastAPI but loses dependency injection benefits",
        "Use FastAPI middleware for auth checks -- global rather than per-route"
      ],
      "detected_signal": "Custom @requires_auth decorator in auth/decorators.py wraps routes and checks session. Reference FastAPI app uses Depends(protected_endpoint) in route dependencies parameter.",
      "default_recommendation": "Convert to Depends() dependency injection -- the reference implementation already demonstrates this pattern"
    },
    {
      "id": "Q7",
      "question": "Should the configparser (.config INI) pattern be modernized?",
      "category": "configuration",
      "why_it_matters": "The Flask app uses Python's configparser reading a .config INI file. FastAPI apps idiomatically use environment variables, .env files, or Pydantic BaseSettings. The config pattern affects deployment, containerization, and secret management.",
      "options": [
        "Keep configparser as-is -- minimal change, works fine",
        "Migrate to Pydantic BaseSettings with .env file -- idiomatic FastAPI, better validation, environment variable support",
        "Use python-dotenv with os.environ -- simpler than Pydantic but less structured"
      ],
      "detected_signal": "configparser reads .config INI file with sections [WEBAPP] and [AUTH0]. Fallback dict provides defaults. No .env file usage. No environment variable reading.",
      "default_recommendation": "Keep configparser for now; migrate to Pydantic BaseSettings as a follow-up modernization"
    },
    {
      "id": "Q8",
      "question": "Blueprint-to-APIRouter mapping: preserve URL structure?",
      "category": "routing",
      "why_it_matters": "Flask uses Blueprints for modular routing; FastAPI uses APIRouter. The mapping is mostly 1:1 but url_for() calls need to become request.url_for() and route registration syntax changes.",
      "options": [
        "Direct 1:1 mapping (Blueprint -> APIRouter) preserving all URL paths",
        "Restructure URLs to be more RESTful while migrating",
        "Add API versioning prefix (e.g., /api/v1/) during migration"
      ],
      "detected_signal": "Two Blueprints: auth_bp (login/signup/callback/logout) and webapp_bp (/, /home, /profile). Both mounted at url_prefix='/'. Routes use url_for() for cross-blueprint references.",
      "default_recommendation": "Direct 1:1 mapping preserving URLs -- this is a template-rendered app, not a REST API, no need for API versioning"
    }
  ],
  "complexity_assessment": {
    "overall": "Low-Medium",
    "rationale": "Small Flask app (18 Python files across both source and reference, only 7 source files with code). No database, no complex middleware, no background tasks. Main complexity is in OAuth/session handling differences between Flask and Starlette. A reference FastAPI implementation already exists in the repo.",
    "estimated_effort": "1-2 days for an experienced developer",
    "risk_areas": [
      "Session storage size limits when moving from Flask server-side sessions to Starlette cookie sessions",
      "Authlib integration switch (flask_client -> starlette_client) -- subtle API differences in how Request is passed",
      "Jinja2 template context differences -- Flask's global request/session vs explicit passing in Starlette",
      "url_for() -> request.url_for() migration throughout templates and Python code",
      "Custom Jinja2 filter registration mechanism differs between Flask and FastAPI"
    ]
  }
}
```

---

## Questions the Skill SHOULD Ask (from expected list)

1. **Sync or async endpoints?** -- The Flask app is fully synchronous, but FastAPI supports async natively. OAuth token exchange is I/O-bound. This is a fundamental design decision.

2. **Pydantic models for validation?** -- FastAPI's core differentiator is Pydantic-based validation. Even though this app has no JSON APIs, Pydantic BaseSettings could modernize the config pattern.

3. **Keep Jinja2 templates or go API-only?** -- The app renders 5 HTML templates with inheritance and custom filters. This is a critical fork in the migration path.

4. **SQLAlchemy stays?** -- This question assumes a database layer exists. The Flask app has NO database and NO SQLAlchemy dependency. This question is not applicable to this codebase.

---

## Questions the Skill DID Surface (from analysis)

| ID | Question | Maps to Expected? |
|----|----------|-------------------|
| Q1 | Sync or async endpoints? | YES -- direct match |
| Q2 | Use Pydantic models for validation? | YES -- direct match |
| Q3 | Keep Jinja2 templates or go API-only? | YES -- direct match |
| Q4 | How should session management be handled? | NO -- novel, not in expected list |
| Q5 | How should Authlib integration be migrated? | NO -- novel, not in expected list |
| Q6 | How should @requires_auth decorator be migrated? | NO -- novel, not in expected list |
| Q7 | Should configparser be modernized? | NO -- novel, not in expected list |
| Q8 | Blueprint-to-APIRouter mapping: preserve URLs? | NO -- novel, not in expected list |

**Coverage of expected questions: 3 out of 4 (75%)**

---

## Gap Analysis

### Expected question NOT surfaced:

**"SQLAlchemy stays?"** -- This question was correctly NOT surfaced because it is inapplicable. The Flask app has:
- No SQLAlchemy dependency in requirements.txt
- No database models anywhere in the codebase
- No database connection configuration
- No ORM usage whatsoever

The expected question appears to be a generic Flask-to-FastAPI migration question that assumes a typical Flask app with a database. This particular app is an Auth0 demo with no persistence layer. A well-designed questionnaire skill should NOT ask this question for this codebase -- asking it would demonstrate lack of codebase analysis and would waste the user's time.

**Verdict:** The expected list has a false positive. The skill's behavior of omitting this question is CORRECT.

### Novel questions the skill surfaced beyond the expected list:

The skill surfaced 5 additional questions (Q4-Q8) that are highly relevant to this specific migration:

1. **Session management (Q4)** -- CRITICAL gap in the expected list. Flask's session and Starlette's SessionMiddleware are fundamentally different (server-side vs cookie-based). The Auth0 token storage issue is a real migration pitfall.

2. **Authlib integration (Q5)** -- IMPORTANT. The Flask-specific OAuth library must be swapped to the Starlette equivalent. This is a concrete, non-obvious change.

3. **Decorator to Depends() (Q6)** -- IMPORTANT. The @requires_auth pattern is idiomatic Flask; Depends() is idiomatic FastAPI. This is the kind of pattern-level change users need guidance on.

4. **Config modernization (Q7)** -- MODERATE. configparser is not wrong in FastAPI, but Pydantic BaseSettings is the idiomatic approach.

5. **Blueprint-to-APIRouter mapping (Q8)** -- MODERATE. Straightforward but the url_for() differences are a common pitfall.

---

## Skill Improvement Suggestions

### 1. Conditional question filtering based on actual dependencies
The expected question list includes "SQLAlchemy stays?" generically. The skill should maintain a library of potential questions per migration path (Flask->FastAPI) but FILTER them based on what is actually detected in the codebase. The skill correctly did this -- the expected list should be updated, not the skill.

### 2. Prioritize session/auth migration questions
For any web app with authentication, session management migration should be a top-priority question. The expected list misses this entirely. The skill should always surface session-related questions when it detects session usage (`flask.session`, cookies, auth decorators).

### 3. Detect and leverage reference implementations
This repo contains BOTH the Flask source and a FastAPI reference implementation. The skill should detect this and:
- Note that a reference exists
- Use it to inform default recommendations
- Flag differences between the reference and idiomatic best practices
- Adjust question phrasing (e.g., "The repo contains a reference FastAPI implementation -- should we follow that pattern or diverge?")

### 4. Template context migration deserves its own question
When Jinja2 templates are retained, the Flask-to-FastAPI template context mechanism changes significantly:
- Flask: `render_template("x.html", session=session)` with global `request` available
- FastAPI: `templates.TemplateResponse("x.html", {"request": request, ...})`
- The `request` must be explicitly passed in FastAPI
This is a common source of bugs and deserves a dedicated question.

### 5. Static file serving strategy
The Flask app serves static files from `/static/`. FastAPI requires explicit `StaticFiles` mounting (`app.mount("/static", StaticFiles(directory="static"))`). For small apps this is trivial, but for production apps the question of whether to serve statics via the app or via a reverse proxy/CDN matters.

### 6. Improve the expected question list for this test case
The expected questions should be:
- Sync or async endpoints? (KEEP)
- Pydantic models for validation? (KEEP)
- Keep Jinja2 templates or go API-only? (KEEP)
- ~~SQLAlchemy stays?~~ (REMOVE -- not applicable, no database in this app)
- How should Flask sessions be migrated to Starlette? (ADD)
- How should the @requires_auth decorator be migrated? (ADD)
- How should Authlib OAuth integration be switched? (ADD)

### 7. Risk-weighted question ordering
Questions should be ordered by migration risk, not alphabetically or by category. For this app:
1. Session management (highest risk -- size limits, behavioral differences)
2. Authlib integration (subtle API differences, OAuth flow breakage)
3. Jinja2 templates vs API-only (architectural fork)
4. Sync vs async (affects auth flow design)
5. Auth decorator to Depends() (pattern change)
6. Pydantic models (low risk, optional modernization)
7. Config pattern (lowest risk, works either way)
