#!/usr/bin/env sh

gnuplot -e "
set terminal png size 1200,600;
set output 'mem.png';

set datafile separator whitespace;

set xdata time;
set timefmt '%Y-%m-%d %H:%M:%S';
set format x '%H:%M:%S';

set grid;
set title 'Memory usage (RSS)';
set ylabel 'MB';
set yrange [0:*];

plot 'mem.log' using (strptime('%Y-%m-%d %H:%M:%S', strcol(1).' '.strcol(2))):(\$3/1024) with lines title 'RSS'
"
