13 Dec 2025

# **Measuring the Impact of Agentic AI on Enterprise Software Development Life Cycle (SDLC)**

## **Methods Discussion**

### **Thoughts on metrics proposed in [preliminary research proposal](https://docs.google.com/presentation/d/1aJJIrQvcwrt3Vt1r7OuIOHDG-QHW7WKm/edit?usp=sharing&ouid=112480188286379590350&rtpof=true&sd=true):**

**Code duplication**: While conceptually appealing, should be explicitly excluded or treated with caution on methodological grounds. In particular, code duplication metrics that are restricted to detection within a single repository are poorly aligned with modern enterprise software development practices. In environments dominated by microservices and multi-repository architectures, meaningful duplication frequently occurs across repositories rather than within them. Repository-local duplication measures therefore risk being both incomplete and potentially misleading, while still imposing non-trivial compute and processing costs. Given these limitations, duplication does not constitute a reliable or interpretable signal in this setting and could be excluded from the core analysis.

**Code Coverage**: Test coverage metrics must be interpreted with caution in enterprise settings. CI systems such as Jenkins do not compute coverage themselves but collect reports produced by language-specific tools during individual build executions. Coverage is therefore calculated only for the source code present in the build workspace and only for tests executed within that specific pipeline, resulting in a repository-local view of testing activity.

In many enterprise environments, particularly those based on microservices, substantial portions of the test suite such as integration, contract, end-to-end, and performance tests are maintained and executed outside the primary code repository and often in separate pipelines or CI systems. These tests are not reflected in repository-level coverage reports unless coverage artifacts are explicitly aggregated, which is uncommon in practice. As a result, coverage measurements may systematically underrepresent actual testing effort and vary across teams for reasons unrelated to software quality or developer productivity.

| Proposed to Be Removed |  |
| :---- | :---- |
| Repository-local code duplication | Misses cross-repository duplication common in microservice architectures; incomplete and costly to compute if done right |
| Test coverage (as primary metric) | Reflects only repository-local tests executed in a single CI pipeline; often excludes integration and end-to-end testing. |
| Deep static or architectural metrics | Require full commit history scanning or semantic code analysis; high compute and deployment overhead. |

## **Proposal for Discussion \- Experimental Design Framework**

**Objective:** Isolate the causal effect of AI usage on software delivery velocity, review costs, and code quality.  
**Data Source:** Git Repositories & AI Tool Telemetry (Git-native only).

### **Table 1: Independent Variables (Inputs & Controls)**

Variables used to define the treatment (AI) and control for task complexity.

| Metric Category | Metric Name | Calculation / Definition | Scientific Rationale |
| :---- | :---- | :---- | :---- |
| **Treatment** | **AI Usage Intensity** | AI Usage API calls | Defines the "dose" of AI. Allows correlation of high-usage vs. low-usage days, individuals, or shas |
| **Control (Size)** | **PR Volume** | Total LoC (Added \+ Removed) per PR. Potential transformation: `log(LoC + 1)`. | **Essential Covariate.** Normalizes for task size. Log-transformation may be required as code changes follow a power law. |
| **Control (Complexity)** | **PR Scatter** | Count of unique files modified in the PR. | **Essential Covariate.** Modifying 10 files is cognitively harder than modifying 1 file, even if line count is identical. |
| **Control (Context)** | **Greenfield Index** | % of PR lines belonging to files created \<X days ago. (Scale 0-1). | AI performs differently on new code (Greenfield) vs. legacy maintenance (Brownfield). |
| **Control (Env)** | **Repo Maturity** | Repository Age (Days since init) \+ Primary Language \+ Total Repo Size | **Structural Control.** Controls for technical debt (Age), ecosystem verbosity (Language), and codebase scale (Size). May want to split into 2-3. |

Problem: **multiple simultaneous AI tools** (Copilot, Cursor, Claude Code adopted concurrently)

### **Table 2: Dependent Variables: Velocity (Speed)**

Measures throughput and cycle time. Hypothesis: AI reduces Coding Lead Time but may impact Review Lead Time.

| Metric Category | Metric Name | Calculation / Definition | Hypothesis |
| :---- | :---- | :---- | :---- |
| **Cycle Time** | **Coding Lead Time** | PR Open Timestamp \- First Commit Timestamp | **Drafting Speed.** AI should reduce this. |
| **Cycle Time** | **Pickup Time** | First Review Timestamp \- PR Open Timestamp | **Queue Latency.** Mostly a control variable. High pickup time usually indicates team busyness, not AI impact. |
| **Cycle Time** | **Review Lead Time** | Merge Timestamp \- First Review Timestamp | **Iteration Speed.** AI may increase this if the code is hard to verify. |
| **Throughput** | **Volume Throughput** | Count of Merged PRs / Active Developer / Week | Captures raw output volume (velocity). |
| **Efficiency** | **Merge Rate** | Merged PRs / (Merged PRs \+ Closed PRs) | Proxy for **Success Rate**. A drop indicates AI is generating "waste" (bad code) that is discarded. |

### **Table 3: Dependent Variables: Review Costs (Friction)**

Measures the human effort required to process AI-generated code. High costs here negate velocity gains.

| Metric Category | Metric Name | Calculation / Definition | Hypothesis |
| :---- | :---- | :---- | :---- |
| **Review Load** | **Review Density** | Total Review Comments / PR Volume (LoC) \[or PR count\] | **Critical Metric.** An increase suggests AI shifts labor from the author (writing) to the reviewer (correcting). |
| **Friction** | **Iteration Count** | Count of "Review Requested" to "New Commit" cycles. | Measures "ping-pong." High iterations imply AI code looks plausible but fails detailed scrutiny. |
| **Rework** | **Review Rework Ratio** | Total LoC Changed During Review / Total LoC in First Draft | **Correction Intensity.** Measures how much the code had to change to pass review. High ratio \= "First Draft" was poor. Could also analyze this per time period rather than per PR |
| **Scope** | **Scope Creep** | (Final LoC \- Initial LoC) / Final LoC | **Omission.** High creep suggests AI missed requirements, forcing humans to add them during review. |
| **Coordination** | **Reviewer Count** | Count of unique users who commented or approved per PR. | **Overhead.** Does AI code require more eyeballs to verify? |

### Table 4: Dependent Variables: Quality (Stability)

Measures the robustness of the output. Did the speed come at the cost of reliability?

| Metric Category | Metric Name | Calculation / Definition | Hypothesis |
| :---- | :---- | :---- | :---- |
| **Stability** | **Revert Rate** | % of Merged PRs reverted within X days. | The standard proxy for "catalytic failure" or breaking changes in production. |
| **Longevity** | **Code Survival Rate** | % of lines added in Month M still present in Month M+X. 21/60 days? | **Validity Test.** If AI code is written fast but deleted/rewritten shortly after, it was "waste," not productivity. |
| **Defects** | **Hotfix Follow-up** | % of PRs followed by a fix-PR to the same files within Xh | Indicates "leaky" code where the primary PR introduced a bug requiring immediate remediation. |
| **Code Health** | **Complexity Delta** | (Post-Merge Complexity \- Pre-Merge Complexity) / PR Volume | **Tech Debt.** AI often defaults to verbose or convoluted logic, increasing maintenance burden even if the code works. |
| **Code Health** | **Static Analysis Delta** | (New Warnings Introduced \- Old Warnings Fixed) / PR Volume | **Code Hygiene.** Captures subtle defects or "code smells" that don't break the build but degrade quality. |
| **Code Health** | **“Code Inflation”** | Average LoC (Added \+ Removed) per Merged PR. Also a covariate | **Technical Debt.** If this increases while "PRs per Developer" stays flat, AI is likely generating verbose, unoptimized code (Bloat). |

### **Potentially Include:**

1. Data from JIRA

