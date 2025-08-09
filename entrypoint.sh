#!/bin/sh -l

find . -name '*.ods' -exec java -jar /usr/src/odfvalidator.jar {} \; >> /odf-errors.log 2>&1

/usr/local/bin/odfvalidatorparser

exit 0
