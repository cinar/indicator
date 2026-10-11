@AGENTS.md

## Claude Code Notes

- **Run `task` before every PR.** CI runs `task` (fmt, lint, test, docs) but does not commit the regenerated docs back to the branch. Run it locally and include any resulting `README.md` changes in your commit — otherwise the generated package docs drift out of sync with the source, which they had before being caught and fixed.
- **No promissory or advisory language.** This is an educational tool, not investment advice. Code comments, doc comments, commit messages, and PR descriptions must stay technical and descriptive — what a value measures, how it's computed — never prescriptive (what to do, what returns to expect, buy/sell framing). This is consistent with the disclaimer already in every package's doc comment.
