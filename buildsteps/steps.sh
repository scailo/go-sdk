#!/bin/sh

# Run this by going to the root of the repo and running ./buildsteps/steps.sh

rm -rf buf
rm -rf google
python3 buildsteps/truncate.py

goimports -w .