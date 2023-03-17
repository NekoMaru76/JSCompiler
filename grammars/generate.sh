#!/usr/bin/bash

antlr4 -Dlanguage=Go -visitor -package parser -o ./lib *.g4
cp ./base/*.go ./lib