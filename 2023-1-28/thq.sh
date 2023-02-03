#!/bin/sh

INPUT=$1
OUTPUT=$2

usage() {
	echo "usage: $0 input output"
	return 0
}


main() {
	if [ -z ${INPUT} ] || [ -z ${OUTPUT} ]; then
		usage
		exit 0
	fi

	if [ ! -f ${INPUT} ]; then
		echo "input file ${INPUT} not exists"
		exit 0
	fi

	cat $INPUT | awk '
	BEGIN{FS=" "}
	{
		if (NF==6) {
			m1[$3]++;
			m2[$3]+=$5
			m3[$3]+=$6
		}
	}
	END{
		for (key in m1) {
			printf("%d, %.6f, %.6f\n",key, m2[key]/m1[key], m3[key]/m1[key])
		}
	}
	' | awk '
	BEGIN{FS=","}
	{
		printf("%d\t%g\t%g\n", $1, $2, $3)
	}
	' | sort -k1n > $OUTPUT
}

main