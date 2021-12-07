#!/bin/sh
NUMBER=0
for i in 1 2 3
do
	JSONPATH="{range .items[$NUMBER]}{@.metadata.name}:{range @.status.conditions[?(@.type=='Ready')]}{@.type}={@.status};{end}{end}"
	STATUS=$(kubectl get nodes -o jsonpath="$JSONPATH")
	CHECK="Ready=True"
	if [[ "$STATUS" == *"$CHECK"* ]];
	then
		echo "${STATUS%%:*} is ready"
	else
		echo "${STATUS%%:*} is not ready"
	fi
	((NUMBER++))
done
