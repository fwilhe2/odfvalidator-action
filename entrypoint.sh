#!/bin/sh -l

find . -name '*.ods' -exec java -jar odfvalidator.jar {} \;

exit 0
