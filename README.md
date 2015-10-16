transmogrifier
==============
[![Build Status](https://travis-ci.org/mohae/transmogrifier.png)](https://travis-ci.org/mohae/transmogrifier)

Transmogrifier is named after one of Calvin's most complex inventions. This one a lot simpler, more reliable, and does not have a dinosaur option.

## About
A transmogrifier is something that changes things to a different shape or form.  In this case it is for transmogrifying data between formats.

For an example implementation see [mog](https://github.com/mohae/mog), a cli tool for creating a []GitHub Flavored Markdown](https://github.com/articles/github-flavored-markdown) table out of CSV data.

## Supported Transformations
Currently, only one transformation is supported:  CSV to GitHub Flavored Markdown Table.

### CSV -> MD Table
The input CSV data can start with an optional header row. If this row does not exist, a format file must exist.  The separator for the CSV data can be specified if it is something other than a comma.  For CSV data that comes from a file and is written to a file. the resulting output file with the MD is saved in the same directory as the source, using the same filename.  The orginal extension is replaced with `.md`.

#### Format file
A format file can be used to specify both the column names and the formatting to be applied to the column.  This includes column justification and column text transformations: ____italics_____, ____bold___, and ~~strikethrough~~.  If a format file is not used, the first row of the data must be the column names.  The format file is expected to be `filename.fmt` and is expected to be in the same directory as the data.

## License
This is licensed under the MIT license. Please view the LICENSE file for more information.

