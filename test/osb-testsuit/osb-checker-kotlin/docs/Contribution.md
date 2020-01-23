## Table of Contents
- [Description](../README.md#description)
- [Getting Started](../README.md#getting-started)
    - [Build Application](../README.md#build-application)
    - [Basic Run Configuration](../README.md#basic-run-configuration)
- [Changes](../README.md#changes)
- [Usage](docs/Usage.md)
    - [Declaring Test Runs](docs/Usage.md##declaring-test-runs)
    - [Parameters](docs/Usage.md#parameters)
    - [Originating Identity](docs/Usage.md#originating-identity)
    - [Declaring Services](docs/Usage.md#declaring-services)
- Test Classes
    - [Catalog](docs/CatalogTest.md)
    - [Provision](docs/ProvisionTests.md)
        - [Test Procedure](docs/ProvisionTests.md#test-procedure)
        - [Version specific Tests](docs/ProvisionTests.md#version-specific-tests)
        - [Example Output](docs/ProvisionTests.md#example-output)
    - [Binding](docs/BindingTests.md#binding)
        - [Test Procedure](docs/BindingTests.md#test-procedure)
        - [Version specific Tests](docs/BindingTests.md#version-specific-tests)
        - [Example Output](docs/BindingTests.md#example-output)
    - [Authentication](docs/AuthenticationTests.md)   
    - [Contract](docs/ContractTest.md)
- [Contribution](#Contribution)
# Contribution

## Create an Issue

In case you find a bug, don't understand the documentation or have a question about the project you can create an issue.

When creating an issue you should follow these tips:

- Check if a similar issue *already exists*.
- Describe your problem *as clear as possible*. What was your expected outcome and what happened instead?
- Name your *system details*, for example what operation system or library you are using.
- Paste your *error or logs* in the issue. Make sure to wrap it in three backticks to render it automatically.

## Pull Request

If you are able to patch a bug or add a feature, you can create a pull request with the code to contribute. But first of make sure you understand the license. Once you created the pull request, the maintainer(s) can check out your code and decide whether or not to pull in your changes.

When creating a pull request you should follow these tips:

- [Fork](https://guides.github.com/activities/forking/) the *repository* and *clone* it locally. Connect your local repository to the original.
- *Pull changes* as often as possible to stay up to date, so that merge conflicts will be less likely.
- [Branch](https://guides.github.com/introduction/flow/) for your changes.
- Run your changes against *existing tests*, or create new ones. 
- Follow the *style of the project*, to make it easier for the maintainer to merge.

If you are asked to make some changes to your request, add more commits and push them to your branch. These changes will automatically go into your existing pull request.
