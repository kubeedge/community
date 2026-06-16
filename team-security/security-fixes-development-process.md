# Security Fixes Development Process

This document describes the private operational process for tracking,
investigating, developing, and preparing security fixes after the KubeEdge
Security Team receives a vulnerability report.

This document supplements, but does not replace, the
[Security Policy](SECURITY.md) and the
[Security Release Process](security-release-process.md). The Security Policy
defines how vulnerabilities are reported and which versions are supported. The
Security Release Process defines the overall response, disclosure, and release
flow. This document focuses on the practical steps for people who participate in
private issue investigation and fix development.

The Security Team may adjust this process when needed based on severity,
reporter expectations, upstream timelines, release manager feedback, or user
impact.

## Goals

- Keep vulnerability analysis and fix development private until coordinated
  disclosure is approved.
- Give security collaborators a clear step-by-step workflow for private
  investigation, fix development, CVE preparation, and release coordination.
- Maintain a clear private audit trail for each vulnerability report.
- Use GitHub Security Advisory comments as the private case log for mapping
  mailing list reports, analysis, fix status, CVE preparation, and release
  coordination.
- Limit access to the minimum set of people required to analyze and fix the
  issue.
- Coordinate advisory publication, fix merge, and patch releases to avoid
  unnecessary exposure windows.

## 1. Report Intake

When the security mailing list receives a vulnerability report, the Security Team
should reply to the reporter as soon as possible.

The reply should:

- Acknowledge that the report has been received.
- State that the community has started analysis.
- Avoid committing to a specific fix or release timeline before assessment is
  complete.
- Ask the reporter to follow coordinated disclosure and not disclose vulnerability
  details publicly before the community completes the fix and public release.

## 2. Private Advisory Setup

The Security Team should create a GitHub Security Advisory draft for the issue.
The advisory must remain private while the issue is being analyzed and fixed.

Use the advisory comments as the private tracking log for the case. Advisory
comments are intended for private case tracking and are not included in the
published advisory, but sensitive report details should still be sanitized where
appropriate. Do not copy sensitive details into public advisory fields until
public disclosure is approved.

Record the following information in advisory comments:

- A link to the KubeEdge security policy used as the process basis for handling
  the report.
- The vulnerability information reported by the reporter, with appropriate
  sanitization while preserving key reproduction details.
- Initial severity assessment, such as the expected CVSS range.
- Expected disclosure window, if one can be estimated.
- Important mailing list context needed to map the report to the advisory and
  track decisions.

## 3. Collaborator Access

Invite only the collaborators required to analyze, reproduce, or fix the issue.
These collaborators do not need to be added to the Security Team.

The Security Team should follow the minimum necessary access principle:

- Grant access only to people needed for this specific case.
- Keep vulnerability discussion in private channels and the advisory comments.
- Avoid sharing exploit details beyond what is needed for analysis and fixing.
- Remove case-specific access when it is no longer needed.

## 4. Vulnerability Analysis

Collaborators analyze and iterate privately in the advisory.

If the issue is not a vulnerability:

- Reply to the reporter with the assessment.
- Close the security advisory.
- End the security fix process for this report.

If the issue is confirmed as a vulnerability:

- Reply to the reporter and confirm that the vulnerability has been accepted.
- Coordinate a public disclosure date with the reporter.
- The Security Team has the final decision on the public disclosure schedule.
- Proceed with private fix development in the private fork associated with the
  advisory.

## 5. Private Fix Development

Security fixes for confirmed vulnerabilities should be developed in the private
fork associated with the GitHub Security Advisory unless the Security Team decides
that a semi-public process is appropriate under the
[Security Release Process](security-release-process.md).

During fix development:

- Keep the fix, reproduction details, and security impact discussion private.
- Track key decisions, review status, testing status, and branch coverage in
  advisory comments.
- Avoid exposing the vulnerability through public issues, public pull requests,
  commit messages, CI logs, or release notes before disclosure is approved.
- Prepare fixes for affected supported release branches according to the
  supported versions policy and the severity and feasibility assessment.
- Confirm whether the fix requires mitigations, configuration guidance, or user
  action before upgrade.

## 6. CVE and Disclosure Preparation

CVE request timing may vary based on the case. The Security Team can request a
CVE after the vulnerability is confirmed and enough information is available, and
should complete the request no later than release preparation. KubeEdge currently
requests CVE IDs through the GitHub Security Advisory UI.

As the fix becomes ready, prepare the release and disclosure materials.

Before the release window:

- Request a CVE through the GitHub Security Advisory UI by using `Request CVE`,
  if this has not already been done.
- Prepare the final public advisory description.
- Include reporter credit when appropriate, but confirm in advance whether the
  reporter wants to be credited publicly.
- Confirm the affected versions and fixed versions.
- Confirm the supported release branches that need patch releases.
- Coordinate advance disclosure to the Private Distributors List when applicable,
  according to the [Security Release Process](security-release-process.md).
- Prepare public announcements and release notes.
- Confirm the planned advisory publish time and fix merge time.

## 7. Fix Merge and Release

Advisory publication and fix merge should be synchronized, or the advisory should
be published slightly before the fix merge, to avoid a window where code is
public but the security advisory is not.

On release day:

- Publish the GitHub Security Advisory at the coordinated disclosure time.
- Merge or cherry-pick the security fixes into the master branch and all affected
  supported release branches.
- Ensure release managers build, publish, and validate the patch releases.
- Publish the release notes and public announcements through the channels listed
  in the [Security Release Process](security-release-process.md).

## 8. Multiple Vulnerabilities

When multiple vulnerabilities are being fixed around the same time, the Security
Team should assess whether they should be released together in a single patch
release batch.

Before batching multiple vulnerabilities:

- Confirm that CVEs have been requested or assigned for each vulnerability that
  requires one.
- Confirm that each vulnerability has a ready fix for every affected supported
  release branch.
- Confirm that public advisory content and announcement content are ready for
  each vulnerability.
- Create one patch release for each supported release line that needs the batch
  of security fixes.

## 9. Case Closure

After public disclosure and patch release:

- Confirm that the advisory, release notes, and announcements point to the fixed
  versions.
- Confirm that reporter credit is correct, if provided.
- Record final release links and CVE identifiers in the advisory comments.
- Close any remaining private tracking items for the case.
