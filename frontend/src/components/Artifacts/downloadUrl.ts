// Mirror of internal/graphql/artifact_download_url.go on the client.
// The BEP File proto doesn't carry digest_function on every emit, so we
// default to sha256 — matching the recorder's behavior.
//
// Patterns:
//
//   bytestream://<host>/<instance>/blobs/<hash>/<size>  → /api/v1/servefile/<instance>/blobs/sha256/file/<hash>-<size>/<name>
//   http(s)://...                                       → returned as-is
//   anything else                                       → ""

export function downloadUrlFor(file: {
	name: string;
	uri?: string;
	digest?: string;
	sizeBytes?: number;
}): string {
	const { name, uri, digest, sizeBytes } = file;
	if (!uri) return "";
	if (uri.startsWith("http://") || uri.startsWith("https://")) {
		return uri;
	}
	if (!uri.startsWith("bytestream://")) return "";
	if (!digest || !sizeBytes) return "";
	let path: string;
	try {
		path = new URL(uri).pathname.replace(/^\//, "");
	} catch {
		return "";
	}
	const idx = path.indexOf("/blobs/");
	const instance = idx > 0 ? path.slice(0, idx) : "";
	const instancePart = instance ? `${instance}/` : "";
	return `/api/v1/servefile/${instancePart}blobs/sha256/file/${digest}-${sizeBytes}/${name}`;
}
