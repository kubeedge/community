# SIG Testing Charter

This charter adheres to the conventions described in [KubeEdge Open Governance](https://github.com/kubeedge/community/blob/master/GOVERNANCE.md) and uses the Roles and Organization Management outlined in the governance doc.

## Scope

SIG Testing is interested in effective testing of KubeEdge and automating away project toil. We focus on creating and running tools and infrastructure that make it easier for the community to write and run tests, and to contribute, analyze and act upon test results.

### In scope

- Project CI and workflow automation via tools

  - Configuration management of CI jobs

  - Metrics, reports, and dashboards for tests that execute in CI jobs
  - Test-specific frameworks, tools

- Infrastructure to support running project CI at scale

- Tools, frameworks and libraries that make it possible to write tests against KubeEdge such as e2e

- Enable other SIGs to efficiently protect their own code from defects
  - Pull Request reviews from other SIGs
    - SIG-Testing may be added to any pull request to request testing advice or review
  - Track code test coverage and drive other SIGs continuous improvement
  - Audit existing code for efficiency and correctness
    - Identify inefficient tests and help them execute faster
  - Set consistent quality standards between SIGs
    - Raise this standard over time


### Out of scope
- We are not responsible for writing, fixing nor actively troubleshooting tests for features or subprojects owned by other SIGs
- Maintain all individual tests on behalf of other SIGs
  - Write new automated tests for each feature
  - Fix individual broken tests or the product bugs they point to
  - Approve every pull request which targets the main development branch
  - Track which tests are not automated/automatable in a feature area
  - Monitoring individual test health or their metrics
    - Defining software performance requirements

## Roles and Organization Management

This SIG follows and adheres to the Roles and Organization Management outlined in KubeEdge Open Governance and opts-in to updates and modifications to KubeEdge Open Governance.

### Additional responsibilities of Chairs

- Manage and curate the project boards associated with all sub-projects ahead of every SIG meeting so they may be discussed
- Ensure the agenda is populated 24 hours in advance of the meeting, or the meeting is then cancelled
- Report the SIG status at events and community meetings wherever possible
- Actively promote diversity and inclusion in the SIG
- Uphold the KubeEdge Code of Conduct especially in terms of personal behavior and responsibility
