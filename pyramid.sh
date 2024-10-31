#!/bin/bash


if [ $# -ne 1 ]
then
	echo "Usage: $( basename $0 ) ROWS" 1>&2
	exit 1
fi

ROWS=${1}
STARS="*************************************************************"
SPACES="                                                            "
for ROW in $( seq 1 ${ROWS} )
do
	INDENT=$(( ( ${ROWS} - ${ROW} ) / 2 ))
	printf "%.*s" ${INDENT} "${SPACES}"
	printf "%.*s\n" ${ROW} "${STARS}"
done
