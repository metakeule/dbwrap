dbwrap
======

[![Build Status](https://secure.travis-ci.org/metakeule/dbwrap.png)](http://travis-ci.org/metakeule/dbwrap)

This library offers two utilities for dealing with the Go database/sql package:

  1. a wrapper that can wrap any database driver that is compatible to sql/driver
     and intercept calls to it

  2. a fake driver that does nothing but tracking the queries and values that are
     delivered to him

Why?
----

Use them to

  - debug a database driver
  - run code before each query
  - transform sql before submitting it to database driver
  - test code that emits sql
  - do logging of sql statements
  - do statistics
  - track long running sql queries


How?
----

see examples directory