# GitLab Status

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-gitlab-status?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-gitlab-status/releases)

Update commit status for GitLab repositories.

<details>
<summary>Description</summary>

This Step updates the commit status for a GitLab repository (repo) of your choice with this build's status. Useful anytime you can not or do not want to provide Bitrise write access to your git repo.

### Configuring the Step

1. The **GitLab API base URL** should be `https://gitlab.com/api/v4/` for cloud-hosted GitLabs.
2. In the **GitLab private token** Step input, you need to provide an access token you generated in your Project Settings or User Settings on GitLab.
3. The **Repository URL** input is populated automatically with a variable the value of which is taken from the repository field of the Settings of your app.
4. You can also select a specific branch or tag to post the status to, but it's going to be sent to the default branch unless you change it.
5. The **Commit hash** input is filled in by default with the variable inherited from the **Git Clone** Step.
6. The **Target URL** Step input is the URL of the build, which is forwarded to GitHub as the source of the status.
7. The **Context** Step input, allows you to label the status with a name.
8. The input **Set Specific Status** input has a default value of `auto` which reflects the status of the build, but this input allows you to update the commit with any given status, regardless of the outcome of the build.
9. The **Description** input allows you to provide a short description for the status.
10. Code **Coverage** percent can also be sent from the build to the commit.

### Troubleshooting

If you get a 404 response when running the Step, check your token's scope and validity.
If you use GitLab Enterprise, make sure your API base URL is set to `https://gitlab.local.domain/api/v4`.
If you do not see your status being reflected, double-check **Repository URL** and **Commit hash** input values.

### Useful links

- [GitLab project access tokens](https://docs.gitlab.com/user/project/settings/project_access_tokens/#create-a-project-access-token)
- [GitLab personal access tokens](https://docs.gitlab.com/user/profile/personal_access_tokens/#create-a-personal-access-token)

### Related Steps

- [Git-Clone](https://www.bitrise.io/integrations/steps/git-clone)
- [Build Status Change](https://www.bitrise.io/integrations/steps/build-status-change)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `api_base_url` | API URL for GitLab or GitLab Enterprise  Example: "https://gitlab.example.com/api/v4" | required |  |
| `private_token` | Authorization token for GitLab applications  Generating a [project access token](https://docs.gitlab.com/user/project/settings/project_access_tokens/#create-a-project-access-token): 1. Log in to your GitLab instance. 1. Go to **Project -> Settings > Access Tokens**. 1. Select **Add new token**. 1. In **Token name**, enter a name. 1. Select a role for the token. 1. Select **scope** *api* and *read_api*. 1. Select **Create project access token**.  Alternatively generating a [personal access tokens](https://docs.gitlab.com/user/profile/personal_access_tokens/#create-a-personal-access-token): 1. Log in to your GitLab instance. 1. In the upper-right corner, select your avatar. Select **Edit profile**. 1. In the left sidebar, select **Access > Personal access tokens**. 1. From the **Generate token** dropdown list, select **Legacy token**. 1. In **Token name**, enter a name for the token. 1. In **Expiration date**, enter an expiry date for the token 1. Select **scope** *api* and *read_api*. 1. Select **Generate token**. | required, sensitive |  |
| `repository_url` | The URL for the repository we are working with | required | `$GIT_REPOSITORY_URL` |
| `git_ref` | The name of a repository branch or tag for which the status needs to be reported  In case of a same commit hash on multiple branches, _ref_  will be used so the pipeline status is updated on the correct branch. | required | `$BITRISE_GIT_BRANCH` |
| `commit_hash` | The commit hash for the commit we are working with | required | `$BITRISE_GIT_COMMIT` |
| `target_url` | The target URL to associate with this status. This URL will be linked from the GitLab UI to allow users to easily see the source of the status. |  | `$BITRISE_BUILD_URL` |
| `context` | A string label to differentiate this status from the status of other systems.  If left empty, it will be `default`. |  | `Bitrise` |
| `preset_status` | If set, this Step will set a specific status instead of reporting the current build status.  Can be one of `auto`, `pending`, `running`, `success`, `failed` or `canceled`. If you select `auto`, the Step will send `success` status if the current build status is `success` (no Step failed previously) and `failed` status if the build previously failed. |  | `auto` |
| `description` | The short description of the status.  If left empty, it will be the status of the build. |  |  |
| `coverage` | The test coverage.  Must be a floating point number between 0.0 and 100.0. |  |  |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-gitlab-status/pulls) and [issues](https://github.com/bitrise-steplib/steps-gitlab-status/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
