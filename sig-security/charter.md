# SIG Security Charter

This charter describes the scope of the SIG Security. 

## Scope

SIG Security is responsible for the design, implementation, and maintenance of features in
KubeEdge that cover horizontal security initiatives, including regular security audits, the vulnerability management process, cross-cutting security documentation, and security community management. As a process-oriented SIG, it does not directly own KubeEdge component code. Instead, SIG Security focuses on improving the security of the KubeEdge project across all components.

SIG Security continues to manage the third-party security audits, while serving a wider mission of advocating for security-related structural or systemic issues and default configuration settings, managing the non-embargoed (public) vulnerability process, defining the bug bounty, creating official KubeEdge Hardening Guides and security documents, and serving as a public relations contact point for KubeEdge security.

### In scope

#### Vulnerability Management Process

Work with the Security Team to define the processes for fixing and disclosing vulnerabilities, as outlined in [security-release-process](../security-team/security-release-process.md). For example:

- When the private fix & release process is invoked
- How vulnerabilities are rated
- The scope of the bug bounty
- Post-announcement follow-ups, such as additional fixes, mitigations, preventions or documentation after a vulnerability is made public
- Distributor announcement policies, such as timelines, criteria for joining the list, etc.
- How, when and where vulnerabilities are announced

#### Security Community Management and Outreach

Provide an entry point to the KubeEdge community for new security-minded contributors, as well as a meeting point to discuss security themes and issues within KubeEdge, including:

- Discuss and design security features to harden KubeEdge.
- Answer security questions from inexperienced users (that don't know what SIG to go to), and identify common questions or issues as areas for improvement.
- Provide an "entry point" for new contributors interested in security. Route these new contributors to other SIGs when they have more specific goals.

#### Horizontal Security Documentation

Author and maintain cross-cutting security documentation, such as hardening guides and security benchmarks. Seek out and coordinate with experts in other SIGs for input on the documentation. In-scope documentation includes:

- Hardening guides and best practices
- Security benchmarks
- Improving documentation to address common misunderstandings or questions
- Threat models

#### Security Audit

Manage recurring security audits and follow up on issues. Coordinate vendors to perform the audit and publish the findings. Follow up on issues with the affected SIG and help coordinate resolution, which can include:

- Helping to prioritize the fixes, possibly by recruiting from SIG Security (while acknowledging that the ultimate authority in deciding whether and how to fix an issue lies with the responsible SIG).
- Documenting mitigations, workarounds, or caveats, especially when the responsible SIG decides not to fix a reported issue.

### Out of scope

SIG Security’s scope does not include:

- Vulnerability response, , including:

  - Embargoed vulnerability management
  - Bug bounty submission triage and management
  - Non-public vulnerability collection, triage, and disclosure

  Notes: Please contact the Security Team to report a vulnerability using [these instructions](../security-team/report-a-vulnerability.md).

- Any projects outside of the KubeEdge project.

- Cloud provider-specific or distributor-specific hardening guides.

- Recommendations or endorsements of specific commercial product vendors or cloud providers.

## Roles and Organization Management

This SIG follows and adheres to the Roles and Organization Management outlined in KubeEdge Open Governance and opts-in to updates and modifications to KubeEdge Open Governance.

### Additional responsibilities of Chairs

None defined at this time.

### Additional responsibilities of Tech Leads

Security Documents and Documentation Tech Leads will be responsible for maintaining the official KubeEdge project Security Hardening Guide.