# propergol

Properties driver for Go!

This module provides management for `.properties` file.

## Presentation

Properties are mappings of keys to values that are used by an application to
hold configuration or environment data. They are usually stored in an external
text file, parsed as the application loads.

Properties are identified by their key; the latter will be unique per each
`Properties` instance (and should be so per application, too).

Externalizing the properties from the code has multiple advantages: properties
can be modified without having to recompile the code (which can prove extremely
valuable in time-sensitive situations, for example an urgent production patch),
or supplied by users of the application for customized behavior; moreover, it
adheres to the principle of separation of concerns, which facilitates design and
maintenance.

## Reading properties files

Properties are stored in text files. The name of such a file usually ends with
the `.properties` extension. They are defined on dinsticnt lines, in the form:

    key=value

The equals sign separates the key from the value, while conveying the idea of
an assignment (as is its purpose in many programming languages).

The key cannot be empty--functionally, this would rarely make sense; in the rare
cases where it does, some other non-empty key would be at least as sensible, if
not more.
On the other hand, the value is allowed to be empty. In this case, the separator
must still be present. Thus, this is how empty-valued properties are defined:

    key=

Definitions of the same key are silently accepted, and only the last definition
is kept. However, property redefinition is discouraged: why have a first
definition just to overwrite it afterwards?

### Whitespace

Whitepsace (spaces and tabs) before the key and around the separator are not
significant, and are not considered part of the property members. This allows to
present properties in a less compact way:

    key = value

It also makes it possible to align successive properties in a table-like
display:

    fruit     = orange
    vegetable = broccoli
    meat      = mutton
    dairy     = yogurt

Whitespace *after* the value is also discarded, but it is not recommended to
have any, as it has no utility and only takes up unnecessary space. It is
usually there as a result of manual error.

Likewise, blank lines between properties are allowed. They can be useful to
group definitions of semantically-related properties.

### Comments

Lines whose first non-whitespace character is a hash sign are ignored by the
parser. They can be used to give context, precisions or directives to property
definitions, for example:

    # This is the host to which the app connects. IPv4 or IPv6 both handled
    host=127.0.0.1
    # The port. Leave empty to use the default port for the selected protocol
    port=
    # The username. NOTE: for security reasons, the user password is not
    # defined here but sourced from the environment
    username=jean_dupont

Inline comments, or comments on the same line as the property definition, are
not handled. This means that in this case:

    key = value # not a comment

the actual property value is `value # not a comment`.

## Writing properties file

The module takes the assumption that, whereas properties files read by the
module are written by and for humans, properties files where it stores are
produced by machines for machines (for example, as part of an automated process)
and therefore no decoration (whitespace and comments) is output when writing the
properties.
