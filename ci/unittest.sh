#!/bin/bash
libs=(
	gallery
	network
	.
)
for lib in ${libs[@]}; do
	cd $lib
	go test
	cd -
done
