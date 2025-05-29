# Database Cleanup Tool

This directory contains the `db_cleanup.go` script and its associated test file, `db_cleanup_test.go` as well as a sample docker file and cronjob manifest. The tool is designed to clean up builds and invocations in the database by executing GraphQL mutations.

## Files

### `db_cleanup.go`
- **Purpose**: Deletes builds and invocations from the database based on a specified time threshold.
- **Features**:
  - Executes two GraphQL mutations:
    - `deleteBuildsBefore`
    - `deleteInvocationsBefore`
  - Supports both HTTP and HTTPS protocols.
  - Optionally loads an SSL certificate file for HTTPS connections.
  - Allows you to specify a configurable lookback period to control which invocations and builds to delete

### `db_cleanup_test.go`
- **Purpose**: Provides unit tests for the `db_cleanup.go` script.
- **Features**:
  - Mocks environment variables and HTTP responses.
  - Tests the `executeGraphQLQuery` function for both builds and invocations.
  - Validates SSL certificate handling.


### `cronjob.yml`
- **Purpose**: Sample Kubernetes cronjob manifest
- **Features**:
  - Schedule - adjust this as needed.  If you create a lot of invocations you probably want to run it more frequently
  - Shows how to provide a certificate and custom environment variables to the cronjob

### `Dockerfile`
- **Purpose**: Allows you to easily build a test container for dev purposes
- **Features**:
  - uses alpine base image


## Environment Variables

The script relies on the following environment variables:

- `BASE_URL`: The base URL of the GraphQL API (e.g., `http://localhost:8081`).
- `HOURS_TO_SUBTRACT`: The number of hours to subtract from the current time to determine the threshold for deletion.
- `TIMEZONE`: The timezone to use for time calculations (e.g., `UTC`).
- `SSL_CERT_FILE`: Path to the SSL certificate file for HTTPS connections (optional).

## Usage

### Running the Script
1. Set the required environment variables.
2. Execute the script:


### Example
Export Environment Variables
```sh
export BASE_URL="https://bbportal.exampleco.io"
export HOURS_TO_SUBTRACT="120"
export TIMEZONE="UTC"
export SSL_CERT_FILE="/path/to/cert.pem"
```
and run the script
 ```bash
go run db_cleanup.go
```
Sample Output
```sh
Deleting builds and invocations starting before 2025-05-01T00:00:00.000Z
Delete Builds Response:
{
  "found": 10,
  "deleted": 8,
  "successful": true
}
Delete Invocations Response:
{
  "found": 5,
  "deleted": 4,
  "successful": true
}
```


### Building the OCI image

You can push the OCI container to your private registry if you like.
1. First modify bb-portal/patches/com_buildbarn_bb_storage/base_image.diff to point to your private registry
2. Make sure your logged into your private registry with `docker login your.private.registry.io`
3. Run the following command:

```sh
bazel run --stamp //tools/db:bb_portal_container_push
```

It does some funky image tagging that doesn't work great outside of CI, but you can easily just use the sha or retag as you like.


### Notes
- Ensure the GraphQL API is accessible and the environment variables are correctly set.
- For HTTPS connections, you can optionally provide a valid SSL certificate file for your private certificate authority.
- There is a dockerfile you can use optionally use to build the with a base for dev purposes
- Tested on postgres production database thats usually around 200GBs in size.  The sql delete queries are blocking, so you may have to adjust the frequency of your cronjob or your lookback window to find what values work best for your environment.  This also works best when you are NOT storing target data in your database as those cascading deletes can fail to work on large builds and large instances.  Its recommended that you do not collect target data for every build anyway and only do it when there is something particular you want to see and explore.