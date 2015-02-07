# EXTENDS for go templates.

This code is a minimal 1-file implementation for adding `extends` support to go templates.
The `preprocess` functions works on both `html/template` and `text/template`.

### USAGE:

1. Include `extends.go` into your view rendering package.
2. In your template, use `{{ extends "super.template" }}` for extension.
3. Load your template with `ParseFilesWithExtends`
4. Render your template by executing the top-most template.
5. Enjoy!

### HOW IT WORKS:

Since GO template engine will sometimes embed parent template as part of the child
template, this makes it difficult to use a recursive approach, so the following approach
is used instead:

1. The `preprocess` function walks the template parse tree looking for `{{extends}}`
2. This list is then compared to all the currently loaded template names.
3. If an `extends`-ion have not been loaded yet, it is loaded.
4. Repeat 1-3 for newly loaded templates until there are no new templates to be loaded.

