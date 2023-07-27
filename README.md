# CircleCI Insights CLI in GO

A wrapper for getting insights data from the CircleCI insights api

## Warning ##

This repository is still in the early stages, and will change significantly.


## Installation

Please clone the the project by running the following command:

```
git clone https://github.com/ogii/circleci-insights-cli.git
```

Next you will need to either build the application with `go build` or install it globally with `go install`.

Lastly you will need to create a .env file in the following format in the directory `~/.goinsightscli/.env`:

```
CIRCLECI_TOKEN=
API_URL=https://circleci.com/api/v2
```
Alternatively, these values can also be set as environment variables.

## Usage/Examples

### Commands

#### get project-workflow-summary

To get a list of summary metrics for a project on the master branch in a list format:

```
project-workflow-summary --slug gh/ogii/circleci-insights-cli --format table
```

Output:

![image](https://user-images.githubusercontent.com/640433/230855922-b0f6cad6-c161-46ca-8716-c9fad9a5688b.png)

Options:

| option             | default             | required                 | sample values                                                    |
|--------------------|---------------------|--------------------------|------------------------------------------------------------------|
| --slug             |                     | yes                      | gh/orgname/repositoryname                                        |
| --branch           | main                | no                       | main                                                             |
| --format           | list                | no                       | list table csv json                                              |
| --reportingWindow  | last-90-days        | no                       | last-7-days last-90-days last-24-hours last-30-days last-60-days |
| --output           | directory of binary | no (only for csv format) | ~/insightsoutput                                                 |


#### get workflow-job-summary

To get the summary metrics for jobs in a specific workflow:

```
circleci-insights-cli get workflow-job-summary --slug gh/ogii/sampleproject --workflow test-testworkflow --branch master --format list
```

| option             | default             | required                 | sample values                                                    |
|--------------------|---------------------|--------------------------|------------------------------------------------------------------|
| --slug             |                     | yes                      | gh/orgname/repositoryname                                        |
| --workflow         |                     | yes                      | workflow1                                                        |
| --branch           | main                | no                       | main                                                             |
| --format           | list                | no                       | list table csv json                                              |
| --reportingWindow  | last-90-days        | no                       | last-7-days last-90-days last-24-hours last-30-days last-60-days |
| --output           | directory of binary | no (only for csv format) | ~/insightsoutput                                                 |