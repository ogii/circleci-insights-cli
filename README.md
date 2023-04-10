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

Lastly you will need to create a .env file in the following format in the same directory as the executable:

```
CIRCLECI_TOKEN=
API_URL=https://circleci.com/api/v2
```

## Usage/Examples

To get a list of summary metrics for a project on the master branch in a list format:

```
circleci-insights-cli getProjectSummaryMetrics --slug gh/ogii/sampleproject --branch master --format list
```

Output:

![image](https://user-images.githubusercontent.com/640433/230855922-b0f6cad6-c161-46ca-8716-c9fad9a5688b.png)
