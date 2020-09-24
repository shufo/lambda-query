# lambda-query

A serverless DB querying client with Lambda

## Overview

An architecture overview

```
+------------------+        +-----------------+          +-----------+
|                  |        |                 |          |           |
|   lambda-query   +------->+ lambda function +---------->    RDS    |
|                  |        |                 |          |           |
+------------------+        +-----------------+          +-----------+
                     invoke                      query
```

lambda-query invokes lambda function with given query and get response and display result with formatted text

This architecture pros and cons

- Pros
  - no more bastion server(!)
    - no requirement for hardening bastion server
    - you only require manage IAM accounts
  - no requirement for keeping db connection while operation
- Cons
  - setup
  - no GUI leads some overhead or operation error sometimes

## Installation

Install client

```bash
$ go get -u github.com/shufo/lambda-query
```

Upload lambda function to your Lambda

```bash
$ cd lambda_function
$ pip install -r requirements.txt -t .
$ zip -r lambda_function.zip *
```

see terraform resource [example](./example/main.tf) for deployment example.

## Usage

```bash
$ lambda-query -f lambda_function -q "select * from users" --format table

+----+------+---------------------+-------------------+----------+----------------+------------+------------+
| id | name | email               | email_verified_at | password | remember_token | created_at | updated_at |
| 1  | foo  | bar@example.com     |                   | pass     |                |            |            |
| 5  | foo  | bartest@example.com |                   |          |                |            |            |
+----+------+---------------------+-------------------+----------+----------------+------------+------------+
```

## Options

|               name |                   description |                 default |
| -----------------: | ----------------------------: | ----------------------: |
| `--function`, `-f` |          Lambda function name |                      - |
|    `--query`, `-q` |                       A query |                      - |
|    `--limit`, `-l` |      Result limit per request | default: `0` (no limit) |
|         `--format` |   specify format [table, csv] |          default: `csv` |
|  `--timeout`, `-t` | max execution time in seconds |           default: `60` |
|   `--output`, `-o` |              output file path |                      - |
|  `--verbose`, `-v` |             Show verbose logs |        default: `false` |

## Example


- Querying table that has many records (e.g. querying million records) 

```bash
$ lambda-query -f function_name \
  -q "select * from users" \
  --limit 10000
```

Add `limit` option will limit records per request (result with many records will timeout or occurs error with Lambda limitation (response size is limited to 6MB))

- Output result to CSV file

```bash
$ lambda-query -f function_name \
  -q "select * from users" \
  --limit 10000 \
  --format csv
  --output result.csv

$ cat result.csv
id,name,email,email_verified_at,password,remember_token,created_at,updated_at
1,foo,bar@example.com,,pass,,,
5,foo,bartest@example.com,,,,,
```

- Querying DB with [aws-vault](https://github.com/99designs/aws-vault) 

```bash
$ aws-lambda foo exec -- lambda-query -f function_name \
  -q "select * from users" \
  --limit 10000 \
  --format csv
  --output result.csv
```

## Motivation

RDS has no serverless querying solution like [Aurora Serverless Query API](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/data-api.html) yet. So I created this client.
