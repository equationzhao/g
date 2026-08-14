---
name: Flag request (rewrite)
about: Propose a new CLI flag or enum value for the rewrite on rewrite/v1
title: "[flag] "
labels: rewrite, flag-request
---

Walk the five gates yourself. Maintainers will close this if a gate fails or if the name is in `docs/rejected-flags.md`.

### Gate 1 — Is this listing?

Can `find`, `du`, `stat`, `xargs`, or a pipe already do this?

- [ ] No, listing itself cannot express it.

If yes, use that tool instead.

### Gate 2 — New value or new dimension?

- [ ] New value on an existing flag (`--sort=`, `--time=`, `--format=`, `--color=`, …)
- [ ] New dimension (will consume budget)

Existing dimension / proposed value:

### Gate 3 — Does it have to be a flag?

Could this be a theme JSON key or a config-only key?

- [ ] Users must switch it **per invocation**.

### Gate 4 — Precedent

- GNU ls:
- eza:
- lsd:
- Proposed name (must match GNU if GNU has it):
- Short letter (check `docs/flag-registry.md`; GNU-reserved letters are not available):

### Gate 5 — Budget

`parse.Budget` is 40. A new dimension needs a cut **or** a Key Decision that revises the cap.

- Flag I would remove (if any):
- Or KD I will write:

### Interaction sketch

Which of the 12 dimensions is this? How does it interact with `--format`, `-l`, `-0`, and json?

### Alternative I considered and rejected
