_Use this email template for pre-disclosing security vulnerabilities to distributors-announce._

TO: `cncf-kubeedge-distrib-announce@lists.cncf.io`

SUBJECT: `[EMBARGOED] $CVE: $SUMMARY`

---

### EMBARGOED

The information contained in this email is **[under embargo](../private-distributors-list.md#embargo-Policy)** until the scheduled public disclosure on **$DATE, at 9AM PT**.

_Additional details on the embargo conditions._
- _If a patch is provided, can it be deployed?_
- _Is the patch itself under embargo?_

### Issue Details

A security issue was discovered in KubeEdge where $ACTOR may be able to $DO_SOMETHING. <optional> KubeEdge are only affected if $CONDITION </end optional>

This issue has been rated **$SEVERITY** (link to CVSS calculator https://www.first.org/cvss/calculator/3.1) (optional: $SCORE), and assigned **$CVE_NUMBER**

_Additional background and high level description of the vulnerability._

### Affected Components and Configurations

_How to determine if a cluster is impacted. Include:_
- _Vulnerable configuration details_
- _Commands that indicate whether a component, version or configuration is used_

#### Affected Versions

- $COMPONENT $VERSION_RANGE_1
- $COMPONENT $VERSION_RANGE_2 ...
- ...

### Mitigations

_If a patch is provided, describe it here._

_(If fix has side effects)_ **Fix impact:** details of impact.

_(If additional steps required after upgrade)_
**ACTION REQUIRED:** The following steps must be taken to mitigate this
vulnerability: ...

_(If possible):_ Prior to upgrading, this vulnerability can be mitigated by ...

### Detection

_How can exploitation of this vulnerability be detected?_

If you find evidence that this vulnerability has been exploited, please contact  [cncf-kubeedge-security@lists.cncf.io](mailto:cncf-kubeedge-security@lists.cncf.io)

Thank You,

$PERSON on behalf of the KubeEdge Security Team
