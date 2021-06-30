module example

go 1.16

require gee v0.0.0

// using 'replace' statement to make 'gee' point to './gee'
replace gee => ./gee
