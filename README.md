[![Build Status](https://travis-ci.org/benjaminbartels/todo.svg?branch=master)](https://travis-ci.org/benjaminbartels/todo)
[![Test Coverage](https://codeclimate.com/github/benjaminbartels/todo/badges/coverage.svg)](https://codeclimate.com/github/benjaminbartels/todo/coverage)
[![Issue Count](https://codeclimate.com/github/benjaminbartels/todo/badges/issue_count.svg)](https://codeclimate.com/github/benjaminbartels/todo)
[![Go Report Card](https://goreportcard.com/badge/github.com/benjaminbartels/todo)](https://goreportcard.com/report/github.com/benjaminbartels/todo)


# ToDo 

Simple ToDo App demonstrating the use of the following stack:

## UI

- Vue.JS frontend built using Vuteify and Vuex 
- Hosted in a AWS S3 bucket behind CloudFront

## Backend

- API hosted using AWS API Gateway with a Lambda function written in Go which 
- Lambda function queries a DynamoDB table

## CI/CD

- Uses TravisCI
- DynamoDB table and IAM executer role deployed via Terraform/Terragrunt modules
- API Gateway, Lambdas and customer domain deployed via Serverless
- Code coverage via gocov and CodeClimate