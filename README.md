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

## Module interface

The module defines a single structure, named `Properties`. Individual properties
are accessed using the method `Get(string) (string, bool)` and defined using
`Set(string, string)`. The methods `Load(io.Reader) error` and
`Store(io.Writer) error` allow to interact with I/O objects (wrapping files,
most of the time) to read (respectively write) property definitions.

More technical details are available in the documentation embedded in [the
source](properties.go).

## Reading properties files

Properties are stored in text files, encoded in UTF-8. The name of such a file
usually ends with the `.properties` extension, although this is not required.
They are defined on distinct lines, in the form:

    key=value

The equals sign separates the key from the value, while conveying the idea of
an assignment (as is its purpose in many programming languages). This is the
only separator that this module recognizes (namely, the colon has no special
meaning).

The key cannot be empty--functionally, this would rarely make sense; in the rare
cases where it does, some other non-empty key would be at least as sensible, if
not more. On the other hand, the value is allowed to be empty. In this case, the
separator must still be present. Thus, this is how empty-valued properties are
defined:

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

Whitespace *after* the value is also silently discarded, but it is not
recommended to have any, as it has no utility and only takes up unnecessary
space. It is usually there as a result of manual error.

Likewise, blank lines between properties are allowed. They can be useful to
group definitions of semantically-related properties.

### Escaping characters

What if one wants to have a property key including an equals sign? Well it is
possible, but the special meaning of the character has to be disabled somehow.
This is done by preceeding the sign with a backslash character, like in the
following example:

    # The key is "key=key" and the value "value"
    key\=key=value

On the other hand, equals signs in the value have no special meaning, because
they occur *after* the actual separator. However, they can still be escaped, to
convey more clearly that they are not the separator:

    # The two following lines are interpreted identically
    url = https://example.com/query?param\=value
    url = https://example.com/query?param=value

Combinations of a backslash and another character are called escape
sequences, and are used to either disable the special meaning of a character
(like the equals sign as separator above) or add a special meaning to a
character that is not otherwise special.
The full list of accepted escape sequences is given below.

|Escape sequence | Result
|----------------|-------
|      `\=`      | A literal equals sign
|      `\\`      | A literal backslash
|      `\n`      | An ASCII newline (LF)
|      `\r`      | An ASCII carriage return (CR)
|      `\t`      | A horizontal tabulation

Note that the escape sequences are oly necessary in properties when read;
properties set using the programmatic interface need not be escaped:

    // Let prop is a Properties object. This statement:
    prop.Set(`key with\=escape sequence`, `value`)
    // sets the value "value" to the property "key with\=escape sequence",
    // not to "key with=escape sequence"

### Line wrapping

If a line length limit is to be enforced, and some properties are longer, it is
possible to break the property definition over several lines, by ending the
initial line with a single back slash character. When parsed, the two lines are
merged back and the backslash, along with any leading whitespace on the next
line, is discarded. This means that any whitespace that is intended to be part
of the reconstituted property has to be placed before the break.
Both key and values can be split over multiple lines this way.

    # In this example, the line length is
    # limited to 40 character
    username_translation.en_US=Username
    # Actual value after reconstruction:
    # “Nome utente”
    username_translation.it_IT=Nome \
                               utente
    # Actual value: “Nom d’utilisateur”
    username_translation.fr_FR=Nom d’utili\
                               sateur

If a property value ends in an escaped backslash, it is not interpreted as line
wrapping:

    # This results in the value “C:\Program Files\”, and not
    # “C:\Program Files\unix_install_dir=/usr/bin/”
    win_install_dir=C:\\Program Files\\
    unix_install_dir=/usr/bin/

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

The module takes the assumption that, while properties files read by the
module are written by and for humans, properties files it writes are produced by
machines for machines (for example, as part of an automated process). Therefore,
no decoration (whitespace and comments) is output when writing the properties.
Moreover, the order in which properties are written is unspecified; in
particular, it may not be the same order in which properties were read.
