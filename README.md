# Buildbarn Portal

Buildbarn Portal is a web service written in React and Go that can display Bazel build output in browsable, digestible fashion.
The service consumes [Build Event Protocol](https://bazel.build/remote/bep) (BEP) data, from local files or as streamed via the [Build Event Service](https://bazel.build/remote/bep#build-event-service) protocol.

Buildbarn Portal groups Bazel invocations into builds based on the observed setting of the `BUILD_URL` environment variable.
If this environment variable is set during an invocation, the invocation will be attributed to a build as identified by the URL.
Buildbarn Portal displays analyzed build output either for individual Bazel invocations, or for builds comprising multiple Bazel invocations.
The grouping of Bazel invocations by build is a differentiating feature of the Buildbarn Portal.

Buildbarn Portal uses a local file-backed database to persist its data so users can continue referencing analyzed results beyond initial viewing.
The service offers basic functionality for browsing and searching among these persisted Bazel invocation results.

## Configuration

The backend is configured via a jsonnet configuration file. Look at `config/portal.jsonnet` for an example configuration
file, and `pkg/proto/configuration/bb_portal/bb_portal.proto` for the configuration schema.

The frontend is configured via environment variables. The following environment variables can be used:

| Variable | Default value | Description |
|----------|---------------|-------------|
| `NEXT_PUBLIC_BES_BACKEND_URL` | `http://localhost:8081` | The URL that the frontend uses to connect to the backend. |
| `NEXT_PUBLIC_BES_GRPC_BACKEND_URL` | `grpc://localhost:8082` | The gRPC URL where the backend listens for the BES. |
| `NEXT_PUBLIC_COMPANY_NAME` | `Example Co` | Used to display the company name in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BEP_UPLOAD` | `true` | Enables the upload of BEP files via the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BES` | `true` | Enables the BES features in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_BUILDS` | `true` | Enables the Builds pages in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_INVOCATIONS` | `true` | Enables the Invocations pages in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TARGETS` | `true` | Enables the Targets pages in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TESTS` | `true` | Enables the Tests pages in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TRENDS` | `true` | Enables the Trends pages in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_BROWSER` | `true` | Enables the Browser pages in the UI. |
| `NEXT_PUBLIC_ENABLED_FEATURES_SCHEDULER` | `true` | Enables the Scheduler pages in the UI. |
| `NEXT_PUBLIC_FOOTER_CONTENT_JSON` | `[{"text": "Buildteam", "href": "https://buildteamworld.slack.com/archives/CD6HZC750", "icon": "slack"}]` | JSON array of links to display in the footer. Each link should have `text`, `href`, and `icon` properties. The `icon` can be a URL or one of `slack`, `github` and `discord`. |

## Setting Up Buildbarn Portal

From `./frontend`, run:

```
npm install
```

## Running the Application

### Running the Backend

From repository root, run:

```
bazel run //cmd/bb_portal -- $PWD/config/portal.jsonnet
```

The backend runs a reverse proxy for the frontend.

### Running the Frontend

From `./frontend`, run:

```
npm run dev
```

### Change where the backend listens

You can run the backend on different bind addresses, but you'll need to update
the frontend too. Modify the backend ports in the config file, and run the frontend:

```
NEXT_PUBLIC_BES_BACKEND_URL=http://localhost:9091 NEXT_PUBLIC_BES_GRPC_BACKEND_URL=grpc://localhost:9092 npm run dev
```

## Using the Application

Go to <http://localhost:3000> or <http://localhost:8081> (which goes through a reverse proxy in the go backend).

> **_NOTE:_** Should the frontend be served on a different url, remember to
> update `frontendProxyUrl` and/or `allowedOrigins` in the configuration.

The home page of the application will appear as follows:

![image](docs/screenshots/home.png)

Bazel invocations known to the service will be listed on the Bazel invocations landing page:

![image](docs/screenshots/bazel-invocations.png)

The problems exhibited during a Bazel invocation will be displayed on a page dedicated to the invocation:

![image](docs/screenshots/bazel-invocation.png)

Builds known to the service will be listed on the builds landing page:

![image](docs/screenshots/builds.png)

From this page, users may navigate to summary views of all invocations associated with a given build:

![image](docs/screenshots/build.png)

### Producing BEP File Examples

From `./bazel-demo`, run:

```
bazel test --keep_going --build_event_json_file=build_events_01.ndjson //...
```

This will produce both a build and a test failure in a single example.

You can then upload the `build_events_01.json` output file to see the results in the application.

### Viewing Build Results From BEP Files

Once you have BEP files produced by Bazel, you can upload them via the application homepage.

### BB-scheduler

BB-portal can show the same information as the web interface from BB-scheduler. To do this, you need to configure the `schedulerServiceConfiguration` in the portal configuration file. The interface can be found under the `Scheduler` tab in the menu.

### BB-browser

BB-portal can show the same information as BB-browser. Everything it can show is visible under the tab `Browser`. To make the browser functionality work, you need to configure the `browserServiceConfiguration` in the portal configuration file. Despite having the name "browser", it is not possible to browse through the content. Instead other parts of Buildbarn will generate links to the browser. To open the content in bb-portal, the the prefix for the links should be `http://url-to-bb-portal/browser/`. After the `/browser/` prefix, the rest of the URL is compatible with urls for bb-browser.

## Using GraphiQL To Explore the GraphQL API

The GraphiQL explorer is available via <http://localhost:8081/graphiql>.

## Generated Code

### Build Event Stream Protocol Buffers

Build portal depends on Build Event Stream protobuf generated code.
This dependency is managed through a Bazel project, `third_party/bazel/_generate-via-bazel/`.
The project itself depends on the original protobuf definitions in Bazel's repository, and patches them in order to produce Go code and to do so with package names specific to Build portal.

> **_NOTE:_** The Bazel project itself uses a specific Bazel version.
> If you use [Bazelisk](https://github.com/bazelbuild/bazelisk), it will automatically activate this version based on the [`third_party/bazel/_generate-via-bazel/.bazelversion`](third_party/bazel/_generate-via-bazel/.bazelversion) file.

For more info, see the [Bazel project's README.md](third_party/bazel/_generate-via-bazel/README.md).
