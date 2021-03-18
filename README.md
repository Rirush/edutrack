# edutrack
Web-service for tracking lectures with support for Markdown with LaTeX equations

## Development status

- [ ] Fully working storage module
  - [ ] Git mode
- [ ] HTTP API endpoint
- [ ] Fully working front-end

## Directory structure

```
/entries/
-> /[subject-id]/
   -> /[entry-id]/
      -> CONTENT.md
      -> /attachments/
         -> /[attachment-name]

/metadata/
-> /subjects.json
-> /entries.json
```