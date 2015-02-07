# EXTENDS for go templates.

This code is a minimal 1-file implementation for adding `extends` support to go templates.
The `preprocess` functions works on both `html/template` and `text/template`.

### USAGE:

1. Paste the contents of `preprocess`, `listExtends`, `listNodeExtends` into your view
   rendering package.
2. Before parsing the first template, adds a placeholder function `extends` that simply
   returns an empty string (`""`)
3. In your template, use `{{ extends "super.template" }}` for extension.
4. Before rendering your template, call `preprocess(template)` to process the extensions.
5. Enjoy!

### HOW IT WORKS:

Since GO template engine will sometimes embed parent template as part of the child
template, this makes it difficult to use a recursive approach, so the following approach
is used instead:

1. The `preprocess` function walks the template parse tree looking for all `extends`
   command by using the `listExtends` function.
2. This list is then compared to all the currently loaded template names.
3. If an `extends`-ed template have not been loaded yet, it is loaded into the template.
4. Repeat 1-3 until there are no new templates to be loaded.

