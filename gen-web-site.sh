#!/bin/bash
# Build the website: 
docker run --rm -it -v ${PWD}:/docs squidfunk/mkdocs-material build
