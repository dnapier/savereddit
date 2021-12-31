#!/bin/bash

cd dashboard || exit
npm run build
cd ..
go run .
