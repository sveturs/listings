---
allowed-tools: Bash(gh pr:*)
description: Create a GitHub pull request from the current branch to the specified target branch.
---
# pr [target-branch]

Create a GitHub pull request from the current branch to the specified target branch.

## Usage

```
/git:pr           # Create PR to default branch (main/master)
/git:pr main      # Create PR to main branch
/git:pr develop   # Create PR to develop branch
```

## Implementation

1. Run `git status` to check current branch and ensure all changes are committed
2. Run `git log --oneline origin/HEAD..HEAD` to see commits that will be included in PR
3. Analyze the commits to generate a meaningful PR title and description
4. Push the current branch to origin if needed
5. Create PR using `gh pr create` with:
   - Title summarizing the changes
   - Body with detailed description of what was done
   - Target branch (from argument or default)
6. Return the PR URL for the user

## Example Output

```
Creating pull request from 'feature/add-parser' to 'main'...
 Pull request created: https://github.com/owner/repo/pull/123
```

## Important

- Requires GitHub CLI (`gh`) to be installed and authenticated
- Will fail if there are uncommitted changes
- Uses the first argument as target branch, defaults to repo's default branch if not specified
