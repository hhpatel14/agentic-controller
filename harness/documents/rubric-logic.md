# Rubric Logic Diagrams

Five diagrams describing the ingest interview loop rubric logic.

## How to view

**Browser (recommended):** Open `diagrams/rubric-logic.html` — renders all five diagrams with Mermaid.js.

**Individual `.mmd` files** (for editors, CI, or embedding):

| # | File | What it shows |
|---|------|--------------|
| 1 | `diagrams/01-main-flow.mmd` | Complete ingest pipeline: scan, interview rounds, output |
| 2 | `diagrams/02-rubric-categories.mmd` | 6 required + 6 optional categories and their graphify pre-fill relationships |
| 3 | `diagrams/03-coverage-scoring.mmd` | How the 0-100% score is composed (user answer + graphify data + clarity) |
| 4 | `diagrams/04-assumption-levels.mmd` | HIGH/MEDIUM/LOW confidence and how each surfaces in the plan |
| 5 | `diagrams/05-coolstore-walkthrough.mmd` | Sequence diagram of the coolstore Java EE to Quarkus migration |

Render any `.mmd` file with:
```
npx @mermaid-js/mermaid-cli mmdc -i diagrams/01-main-flow.mmd -o 01-main-flow.svg
```

Or paste into [mermaid.live](https://mermaid.live).
