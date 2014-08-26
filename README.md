tomd
====

tomd: a package that converts stuff **to ma**rkdown

Written to convert .csv files to markdown tables, but csvtomd has too many consonants and not enough vowels for most English speakers, or American speaker in this case. Anyways that kind of consonant/verb ratio frightens me.

So tomd it is!

## What tomd converts.
Not much at the moment, in fact nothing. But it will support:

* csv to markdown tables
  * header support
  * no header support
  * hopefully some kind of support for justification, thinking csv files files can also have template files associated with them, these template files would include:
    * header flag
    * justification of column headers by column name
    * justification of columns (non-header) by:
      * column name (when there are headers only)
      * column numbers

