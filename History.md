
1.4.1 / 2015-11-02
==================

  * Merge branch 'release/v1_4_1' into develop
  * Fixed output filters.

1.4.0 / 2015-10-30
==================

  * Merge branch 'release/v1_4' into develop
  * Introduce output filters to deal with filters not implemented at FA
    level (like hexid and other fields).

1.3.0 / 2015-09-14
==================

  * Merge branch 'release/v1_3' into develop
  * Support all filters as of Flightaware documentation

1.2.0 / 2015-09-11
==================

  * fa-tail has been moved into its own repo
  * Travis-CI support

1.1.0 / 2015-09-01
==================

  * Add support for "flightplan" in addition to "position" events
  * Update documentation.

1.0.0 / 2015-08-25
==================

  * Working version
  * TOML as configuration file format instead of YAML
  * Added fa-tail to see where we are in big files
  * Seek() to near end of file to speed fa-tail up
