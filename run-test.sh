#!/bin/bash
for testfile in tests/*.test;do echo $(basename $testfile);tests/run.rb $testfile;done
