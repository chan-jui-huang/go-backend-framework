#!/bin/bash

printf "Enter your project name:"
read PROJECT_NAME
PROJECT_NAME=${PROJECT_NAME:=go-backend}
mkdir $PROJECT_NAME

printf "Enter the framework version what you want (release page: https://github.com/chan-jui-huang/go-backend-framework/releases ):"
read VERSION

curl -L https://github.com/chan-jui-huang/go-backend-framework/archive/refs/tags/v$VERSION.tar.gz | tar -zxv -C $PROJECT_NAME --strip-components=1
cd $PROJECT_NAME
rm LICENSE
cp .env.example .env
cp .env.test.example .env.test
cp .env.staging.example .env.staging
cp .air.exmaple.toml .air.toml

OLD_MODULE_NAME=github.com/chan-jui-huang/go-backend-framework
printf "Enter your go module name (\e[33mex: github.com/author/my-project\e[0m):"
read NEW_MODULE_NAME

go mod edit -module $NEW_MODULE_NAME
case $OSTYPE in
darwin*)  find . -type f -name "*.go" -exec sed -i "" "s,$OLD_MODULE_NAME,$NEW_MODULE_NAME,g" {} \; ;; 
linux*)   find . -type f -name "*.go" -exec sed -i "s,$OLD_MODULE_NAME,$NEW_MODULE_NAME,g" {} \; ;;
esac
