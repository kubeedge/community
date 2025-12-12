# POST A JOB

The KubeEdge job center aims to provide a platform for community partners and developers to exchange job opportunities. Partners can post job openings here and provide detailed job information and requirements, while community users can select job positions based on their experience and interests. We sincerely hope that you can find suitable jobs or excellent talents here.

This guide is intended to help you understand how to post your job positions in the community job center and provide guidelines for writing job information documents.

# Requirement for Posting Jobs

- You must be a contributing organization or adopter of the KubeEdge community, including but not limited to having used KubeEdge in your company's public products or public solutions, contributing to the community, and submitting your use cases.
- Your company/organization logo has been added to the [KubeEdge Supporters](https://kubeedge.io/#supporters).

# Writing Guidelines

We hope that you can provide as much detailed information as possible about the job.

1. Please add the following information at the top of the document. This information will be displayed in a card format on the KubeEdge official website's job center homepage, which will help community users quickly understand the key information about the job.
```
title: Job Name
company：Company Name
address： Workplace location
date： Creation Date.（reference format：2023-01-01）
expirydate： Expiry Date.（reference format：2023-01-01）
logo： Company Logo (path to logo image file)
```

2. Document body. Please provide detailed job information, team introduction, job requirements, and contact information in the document body.

# Submit the Job

After completing the job information, please submit your document to the [website](https://github.com/kubeedge/website) repository of KubeEdge. Once it passes community review, it can be displayed in the job center on the official website.

If your document only consists of a company logo image and a document, you can put the image in the `static/img/job-center/` directory, and put the English and Chinese Markdown documents in the `src/pages/job-center/` and `i18n/zh/docusaurus-plugin-content-page/job-center/` directories, respectively. Then submit the materials to the website repository in the form of a pull request.

If you have more materials, you can follow these steps:
1. Create a new directory in `src/pages/job-center/` and `i18n/zh/docusaurus-plugin-content-page/job-center/` under the name of the job you want to submit.
2. Put your English materials in the directory you created under `src/pages/job-center/`, put the Chinese materials in the corresponding newly created directory under `18n/zh/docusaurus-plugin-content-page/job-center/` and put the company logo image in the `tatic/img/job-center directory`
3. Please name your Markdown document as index.md or index.mdx.
4. Please submit your job information to the website repository in the form of a pull request.