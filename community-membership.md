# KubeEdge Community Membership

**Note :** This document keeps changing based on the status and feedback of KubeEdge Community.

This document gives a brief overview of the KubeEdge community roles with the requirements and responsibilities associated with them.

| Role | Requirements | Responsibilities | Privileges |
| -----| ---------------- | ------------ | -------|
| [Member](#member) | Sponsor from 2 reviewers, active in community, multiple contributions to KubeEdge | Active contributor in the community | KubeEdge GitHub organization Member |
| [Reviewer](#member) | Sponsor from 2 approvers, has good experience and history of review in specific package | Review contributions from other members | Add `lgtm` label to specific PRs |
| [Approver](#approver) | Sponsor from 2 maintainers, highly experienced and knowledge of domain, actively contributed to code and review  | Review and approve contributions from community members | Write access to specific packagies in relevant repository |
| [Maintainer](#maintainer) | Sponsor from 2 owners, shown good technical judgement in feature design/development and PR review | Participate in release planning and feature development/maintenance | Top level write access to relevant repository. Name entry in Maintainers file of the repository |
| [Owner](#owner) | Sponsor from 3 owners, helps drive the overall KubeEdge project | Drive the overall technical roadmap of the project and set priorities of activities in release planning | KubeEdge GitHub organization Admin access |


**Note :** It is mandatory for all KubeEdge community members to follow KubeEdge [Code of Conduct].

## Member

Members are active participants in the community who contribute by authoring PRs,
reviewing issues/PRs or participate in community discussions on slack/mailing list.


### Requirements

- Sponsor from 2 reviewers
- Enabled [two-factor authentication] on their GitHub account
- Actively contributed to the community. Contributions may include, but are not limited to:
    - Authoring PRs
    - Reviewing issues/PRs authored by other community members
    - Participating in community discussions on slack/mailing list
    - Participate in KubeEdge community meetings
- Open an issue against the kubeedge/community repo
- Have your sponsoring reviewers reply confirmation of sponsorship: +1
- Once your sponsors have responded, your request will be reviewed by the org owners


### Responsibilities and privileges

- Member of the KubeEdge GitHub organization
- Can be assigned to issues and PRs and community members can also request their review
- Participate in assigned issues and PRs
- Welcome new contributors
- Guide new contributors to relevant docs/files
- Help/Motivate new members in contributing to KubeEdge

## Reviewer

Reviewers are able to review code for quality and correctness on some part of a subproject. 
They are knowledgeable about both the codebase and software engineering principles.


### Requirements

- member for at least 1 months
- Primary reviewer for at least 5 PRs to the codebase
- Reviewed or merged at least 20 substantial PRs to the codebase
- Knowledgeable about the codebase
- May either self-nominate,  or be nominated by an approver


### Responsibilities and privileges

- Code reviewer status may be a precondition to accepting large code contributions
- Responsible for project quality control
- Focus on code quality and correctness, including testing and refactoring
- May also review for more holistic issues, but not a requirement
- Expected to be responsive to review requests, add `lgtm` label to reviewed PRs

## Approver

Approvers are active members who have good experience and knowledge of the domain.
They have actively participated in the issue/PR reviews and have identified relevant issues during review.


### Requirements

- Sponsor from 2 maintainers
- Member for at least 2 months
- Have reviewed good number of PRs
- Have good codebase knowledge


### Responsibilities and Privileges

- Review code to maintain/improve code quality
- Acknowledge and work on review requests from community members
- May approve code contributions for acceptance related to relevant expertise
- Have 'write access' to specific packages inside a repo, enforced via bot
- Continue to contribute and guide other community members to contribute in KubeEdge project

## Maintainer

Maintainers are approvers who have shown good technical judgement in feature design/development in the past.
Has overall knowledge of the project and features in the project.

### Requirements

- Sponsor from 2 owners
- Approver for at least 2 months
- Nominated by a project owner
- Good technical judgement in feature design/development

### Responsibilities and privileges

- Participate in release planning
- Maintain project code quality
- Ensure API compatibility with forward/backward versions based on feature graduation criteria
- Analyze and propose new features/enhancements in KubeEdge project
- Demonstrate sound technical judgement
- Mentor contributors and approvers
- Have top level write access to relevant repository (able click Merge PR button when manual check-in is necessary)
- Name entry in Maintainers file of the repository
- Participate & Drive design/development of multiple features

## Owner

Owners are maintainers who have helped drive the overall project direction.
Has deep understanding of KubeEdge and related domain and facilitates major agreement in release planning

### Requirements

- Sponsor from 3 owners
- Maintainer for at least 2 months
- Nominated by a project owner
- Not opposed by any project owner
- Helped in driving the overall project

### Responsibilities and Privileges

- Make technical decisions for the overall project
- Drive the overall technical roadmap of the project
- Set priorities of activities in release planning
- Guide and mentor all other community members
- Ensure all community members are following Code of Conduct
- Although given admin access to all repositories, make sure all PRs are properly reviewed and merged
- May get admin access to relevant repository based on requirement
- Participate & Drive design/development of multiple features


## Inactive members

_Members are continuously active contributors in the community._

A core principle in maintaining a healthy community is encouraging active
participation. It is inevitable that people's focuses will change over time and
they are not expected to be actively contributing forever.

Therefore members with an extended period away from the project with no activity
will emeritus or be removed from the KubeEdge GitHub Organizations and will be required to
go through the org membership process again after re-familiarizing themselves
with the current state.


### How inactivity is measured

Inactive members are defined as members of one of the KubeEdge Organizations
with **no** contributions across any organization within 12 months. This is
measured by the CNCF [DevStats project].

**Note:** Devstats does not take into account non-code contributions. If a
non-code contributing member is accidentally removed this way, they may open an
issue to quickly be re-instated.


After an extended period away from the project with no activity
those members would need to re-familiarize themselves with the current state
before being able to contribute effectively.


[Code of Conduct]: https://github.com/kubeedge/community/blob/master/CODE_OF_CONDUCT.md
[two-factor authentication]: https://help.github.com/articles/about-two-factor-authentication
[Devstats project]: https://kubeedge.devstats.cncf.io/
