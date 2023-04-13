# Contributing Guidelines for gogpt

We appreciate your interest in contributing to gogpt! This document provides a set of guidelines to help you contribute to the project effectively.

## Code of Conduct

Please read and follow our [Code of Conduct](./CODE_OF_CONDUCT.md) to ensure a welcoming and inclusive environment for all contributors.

## Getting Started

1. Fork the repository: Click the "Fork" button at the top-right corner of the repository to create your own copy.

2. Clone your fork: Run the following command to clone your forked repository locally:
```bash
git clone <YOUR_FORK_URL>
```
3. Add the original repository as an upstream remote:
```bash
git remote add upstream git@github.com:Makepad-fr/gogpt.git
```
4. Stay up to date with the original repository by regularly syncing your fork:
```bash
git fetch upstream
git checkout master
git merge upstream/main
```

## Reporting Issues

Before submitting a new issue, please search the existing issues to ensure your issue has not already been reported. When submitting an issue, provide as much information as possible, including a clear title and description, steps to reproduce the issue, and any relevant logs or screenshots.

## Submitting a Pull Request

1. Choose an issue to work on: Pick an issue from the issue tracker that you'd like to work on. Leave a comment on the issue to let others know you're working on it.

2. Create a new branch: Create a new branch for your changes using the following naming convention: `slb/<YOUR_INITIALS_OR_USERNAME>/GH-<ISSUE_NUMBER>`

3. Make your changes: Implement your changes, following the project's coding standards and ensuring that your code is clean, well-documented, and well-tested.

4. Commit your changes: Write clear, concise, and descriptive commit messages that explain the purpose of each commit.

5. Rebase and squash: Before submitting your pull request, rebase your branch onto the latest `main` branch and squash any small, related commits.

6. Push your changes to your fork:
```bash
git push origin $(git branch --show-current)
```

7. Submit a pull request: Create a new pull request from your feature branch against the `master` branch of the original repository. In the pull request description, provide an overview of your changes and reference any related issues.

## Review Process

After submitting a pull request, one of the maintainers will review your changes. They may request additional changes or clarification. Please be patient and responsive during the review process.

Once your changes have been approved, they will be merged into the project. Congratulations, and thank you for your contribution!

## Proposing New Contributions

If there are no open issues or you have an idea that is not yet discussed, you can still contribute by following these steps:

1. **Create a new issue**: Describe the improvement, bugfix, or feature you'd like to work on. Provide as much detail as possible, including the motivation behind it, potential implementation approach, and any relevant research or resources.

2. **Discuss with maintainers**: Engage with the maintainers and the community in the issue you created. Collect feedback and refine your idea based on the discussions.

3. **Get approval**: Wait for the maintainers to approve your proposal before starting to work on it. They might have some additional suggestions or requirements.

4. **Follow the steps for submitting a pull request**: Once your proposal is approved, follow the steps outlined in the "Submitting a Pull Request" section to contribute your changes.

Remember, it's essential to engage with the project's maintainers and the community while proposing new contributions to ensure alignment with the project's goals and priorities.


## Additional Resources

- [GitHub Forking Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/forking-workflow)
- [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)

Again, thank you for considering contributing to gogpt!
We look forward to collaborating with you.

Happy coding 😎


