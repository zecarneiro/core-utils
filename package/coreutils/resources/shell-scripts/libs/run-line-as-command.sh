#!/usr/bin/env bash

declare file="$1"
declare section="$2"
declare prefix_sufix_key="######"
declare allKey="$prefix_sufix_key ALL $prefix_sufix_key"
declare section="$prefix_sufix_key $section $prefix_sufix_key"
declare canRun=false
while read line; do
    if [[ "${line}" == "${prefix_sufix_key}"* ]]; then
        if [[ "${line}" == "${allKey}"* ]]||[[ "${line}" == "${section}"* ]]; then
            canRun=true
        else
            canRun=false
        fi
    fi
    if [[ "${canRun}" == "true" ]]&&[[ -n "${line}" ]]; then            
        if [[ "${line}" != "${prefix_sufix_key}"* ]]; then
            evalc "$line" 
        fi
    fi
done <"$file"
