import fs from "node:fs";
import path from "node:path";
import type { ResolvedConfig } from "vite";

interface ManifestChunk {
  file?: string;
  src?: string;
  isEntry?: boolean;
  imports?: string[];
  css?: string[];
}

interface RouteAssetChunks {
  path: string;
  js: string[];
  css: string[];
}

type ViteManifest = Record<string, ManifestChunk>;

// Normalizes the route file path
const strippedRoute = (file: string) => {
  return file
    .replace(/^.*src\/routes\//, "")
    .replace(/\?tsr-split=component$/, "")
    .replace(/\.tsx$/, "")
    .replace(/\//g, ".");
};

// This generates the MUX compatible path for each route
function generatePath(file: string): string {
  const route = strippedRoute(file)
    .replace(/^.*index\.html$/, "index.html")
    .replace(/\.html$/, "");

  const segments = route.split(".");
  const pathSegments: string[] = [];

  for (const segment of segments) {
    if (segment === "index") {
      continue;
    }

    //Handle dynamic routes
    if (segment === "$") {
      pathSegments.push("{rest...}");
      continue;
    }

    // Handle dynamic route parameters such as $invocationID
    if (segment.startsWith("$")) {
      pathSegments.push(`{${segment.slice(1)}}`);
      continue;
    }

    pathSegments.push(segment);
  }
  return `/${pathSegments.join("/")}`;
}

// Recursively fetches all the chunks for a given route file
function fetchImportRooutes(routeFile: string, manifest: ViteManifest) {
  const visited = new Set<string>();
  const js = new Set<string>();
  const css = new Set<string>();

  function visitRoute(key: string) {
    if (visited.has(key)) return;
    visited.add(key);

    const chunk = manifest[key];
    if (!chunk) return;

    if (chunk.file) {
      js.add(chunk.file);
    }

    for (const c of chunk.css ?? []) {
      css.add(c);
    }

    // Add all the imports of the current chunk to the visit queue
    for (const imp of chunk.imports ?? []) {
      visitRoute(imp);
    }
  }

  visitRoute(routeFile);

  return {
    js: [...js],
    css: [...css],
  };
}

function generateRouteKeys(
  routeFile: string,
  manifest: ViteManifest,
): RouteAssetChunks {
  const routeId = strippedRoute(routeFile);
  const parts = routeId.split(".");

  const routeChunks: RouteAssetChunks = {
    js: [],
    css: [],
    path: "",
  };

  // Always start the layout chain with __root
  const routeIds = ["__root"];

  if (routeId === "index") {
    routeIds.push("index");
  } else if (routeId !== "__root") {
    // Generates the layout chain for the route
    // For example, "foo.bar.baz" would be
    // ["__root", "foo", "foo.bar", "foo.bar.baz"]
    let prefix = "";
    for (const segment of parts) {
      prefix = prefix ? `${prefix}.${segment}` : segment;
      routeIds.push(prefix);
    }
  }

  // Collect the JS/CSS for the root, all parent layouts, and the route itself
  for (const parentId of routeIds) {
    const parentManifestKey = Object.keys(manifest).find(
      (key) => strippedRoute(key) === parentId,
    );

    if (!parentManifestKey) continue;

    const parentChunks = fetchImportRooutes(parentManifestKey, manifest);
    routeChunks.js.push(...parentChunks.js);
    routeChunks.css.push(...parentChunks.css);
  }

  return routeChunks;
}

function generateRouteManifest(data: ViteManifest): RouteAssetChunks[] {
  // Collect the routes
  const routeKeys = Object.keys(data).filter(
    (file) =>
      file.includes("src/routes/") && file.endsWith(".tsx?tsr-split=component"),
  );

  // To handle keys for index file
  const mainEntryKey = Object.keys(data).find((key) => data[key].isEntry);
  const mainEntryChunks = mainEntryKey
    ? fetchImportRooutes(mainEntryKey, data)
    : { js: [], css: [] };

  const routeMap = new Map<string, { js: Set<string>; css: Set<string> }>();

  for (const routeFile of routeKeys) {
    const path = generatePath(routeFile);
    const routeSpecificChunks = generateRouteKeys(routeFile, data);

    let routeStart = routeMap.get(path);

    // Create new entry if it doesn't exist
    if (!routeStart) {
      routeStart = {
        js: new Set<string>(),
        css: new Set<string>(),
      };
      routeMap.set(path, routeStart);
    }

    // Add the main entry chunks to the route's JS/CSS sets
    for (const js of mainEntryChunks.js) routeStart.js.add(js);
    for (const css of mainEntryChunks.css) routeStart.css.add(css);

    // Add the route-specific chunks to the route's JS/CSS sets
    for (const js of routeSpecificChunks.js) routeStart.js.add(js);
    for (const css of routeSpecificChunks.css) routeStart.css.add(css);
  }

  const routeManifest: RouteAssetChunks[] = Array.from(routeMap.entries()).map(
    ([path, sets]) => ({
      path,
      js: Array.from(sets.js),
      css: Array.from(sets.css),
    }),
  );

  return routeManifest;
}

let config: ResolvedConfig;
export function goRouteManifestPlugin() {
  return {
    name: "go-route-manifest",
    configResolved(resolved: ResolvedConfig) {
      config = resolved;
    },

    writeBundle(options: { dir: string }) {
      const outDir = path.resolve(config.root, config.build.outDir);
      const manifestPath = path.join(outDir, ".vite", "manifest.json");

      if (!fs.existsSync(manifestPath)) {
        throw new Error("Manifest.json not found");
      }
      const content = fs.readFileSync(manifestPath, "utf-8");
      const data: ViteManifest = JSON.parse(content);
      const routeTree = generateRouteManifest(data);

      fs.writeFileSync(
        path.resolve(options.dir, "route-manifest.json"),
        JSON.stringify(routeTree, null, 2),
      );
    },
  };
}
