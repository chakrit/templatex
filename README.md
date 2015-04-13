# EXTENDS for go templates.

[![Build Status](https://travis-ci.org/chakrit/templatex.svg)](https://travis-ci.org/chakrit/templatex)

This code is a minimal 1-file implementation of a thin layer on top of `html/template`
that provides the following features:

* An `{{extends ""}}` directive to allow template to extends another template.
* All templates' names are now its relative path to the main executable working directory
  which allows templates in different subfolders to share the same name without conflict.
* Proper template composition.

The `extends.go` file also works for `text/template`, simply change the import on the top
of the file to use it.

### USAGE:

1. Include `extends.go` into your view rendering package.
2. Define a top-level template `templates/base.template` as a base to extend from.
3. Define a child template `templates/subfolder/child.template`.
4. Adds `{{ extends "../base.template" }}` to make it extends form the base template.
5. Render your template by executing the top-most template name: `templates/base.template`

An example is provided in the `main.go` file.

### HOW IT WORKS:

There are quite a few shortcomings inside the normal template packages that made this
harder than it should be:

1. When you use `ParseFiles`, templates are named by their filename, excluding all paths
   information. So if you accidentally names any template file the same, you will get an
   error.
2. This also pose a problem when we try to resolve template extensions from within another
   template as it might not be in the same folder as the original template or the working
   directory.
3. Although you can add custom template functions, it will not have access to the current
   parsing tree so you are not allowed to modify it on-the-fly.
4. Due to 3., rendering another template file is possible, but not in the same rendering
   context. This prevents any attempt at proper template composition.
4. Due to all of the above, it is then required that you build a rigid convention
   around how your template files are named to achieve any level of composition. For
   example maintaining a large, separate `map[string][]*Template` of all your templates
   and their dependencies.

Thus, we work around this with the following:

1. Name templates according to their absolute file paths.
2. Avoid using the built-in `ParseFiles` and its problematic parsing logic directly.
3. Provide a no-op `extends` template function as our marker.
4. Manually walks through each template's parse tree as they're loaded looking for
   `{{extends}}` fragments, resolving both the extender and the extended's to their
   absolute paths as we go.
5. Loads the extended template and adds it to the top-level template parse tree.
6. Repeat 3-5 for the recently loaded template until we have no more templates to load.
7. Always render from the top-most template (usually the only file without any
   `{{extends}}` in it).

With the above steps we achieve the following:

* Template consumers only need to care about the most specialized template they need.
* Each template specifies their own requirements inside the template file itself without
  requiring any other module's involvements.

### TODO:

* Proper tests.
* Benchmarks and optimizations.
* Extract common interface between html/template and text/template
* Better error messages. (provide parsing context in all errors.)

