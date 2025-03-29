#!/bin/sh -l

find . -name '*.ods' -exec java -jar /usr/src/odfvalidator.jar {} \;

exit 0
