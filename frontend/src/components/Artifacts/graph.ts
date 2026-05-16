import { decompress as zstdDecompress } from "fzstd";

// Parser for the compressed BEP-graph blob returned by
// BazelInvocation.artifactGraph. The payload is a zstd-compressed stream
// of length-prefixed serialized bes.BuildEvent messages, filtered to the
// NamedSetOfFiles and TargetCompleted variants. See
// internal/database/buildeventrecorder/artifact_graph_buffer.go for the
// server-side writer.
//
// We only decode the handful of fields the artifact UI needs, using a
// minimal proto wire-format reader. Field numbers come from
// third_party/bazel/protobuf/build_event_stream.proto:
//
//   BuildEvent.id                = 1 (BuildEventId)
//   BuildEvent.named_set_of_files = 15 (NamedSetOfFiles)
//   BuildEvent.completed         = 8 (TargetComplete)
//
//   BuildEventId.named_set       = 13 (NamedSetOfFilesId)
//   BuildEventId.target_completed = 5 (TargetCompletedId)
//
//   NamedSetOfFilesId.id         = 1 (string)
//   TargetCompletedId.label      = 1 (string)
//   TargetCompletedId.aspect     = 2 (string)
//
//   NamedSetOfFiles.files        = 1 (repeated File)
//   NamedSetOfFiles.file_sets    = 2 (repeated NamedSetOfFilesId)
//
//   File.name                    = 1 (string)
//   File.uri                     = 2 (string, oneof file)
//   File.digest                  = 5 (string)
//   File.length                  = 6 (int64)
//
//   OutputGroup.name             = 1 (string)
//   OutputGroup.file_sets        = 3 (repeated NamedSetOfFilesId)
//   OutputGroup.incomplete       = 4 (bool)
//
//   TargetComplete.output_group  = 2 (repeated OutputGroup)

export interface ParsedFile {
	name: string;
	uri?: string;
	digest?: string;
	sizeBytes?: number;
}

export interface ParsedNamedSet {
	files: ParsedFile[];
	childSetIds: string[];
}

export interface ParsedOutputGroup {
	name: string;
	rootSetIds: string[];
	incomplete: boolean;
}

export interface ParsedTarget {
	label: string;
	aspect: string;
	outputGroups: ParsedOutputGroup[];
}

export interface ParsedGraph {
	sets: Map<string, ParsedNamedSet>;
	targets: ParsedTarget[];
}

export function parseArtifactGraphFromBase64(
	base64Payload: string,
): ParsedGraph {
	const compressed = base64ToBytes(base64Payload);
	const decompressed = zstdDecompress(compressed);
	return parseArtifactGraph(decompressed);
}

export function parseArtifactGraph(stream: Uint8Array): ParsedGraph {
	const reader = new WireReader(stream);
	const sets = new Map<string, ParsedNamedSet>();
	const targets: ParsedTarget[] = [];

	while (reader.hasMore()) {
		const eventLen = reader.uvarint();
		const eventBytes = reader.take(Number(eventLen));
		decodeBuildEvent(eventBytes, sets, targets);
	}

	return { sets, targets };
}

// Walk the file graph rooted at the given set IDs and yield every reachable
// file, terminating on cycles. Useful for rendering output-group contents.
export function* walkSetFiles(
	graph: ParsedGraph,
	rootSetIds: readonly string[],
): Generator<ParsedFile> {
	const visited = new Set<string>();
	const queue: string[] = [...rootSetIds];
	while (queue.length > 0) {
		const id = queue.shift() as string;
		if (visited.has(id)) continue;
		visited.add(id);
		const set = graph.sets.get(id);
		if (!set) continue;
		for (const file of set.files) {
			yield file;
		}
		for (const childId of set.childSetIds) {
			queue.push(childId);
		}
	}
}

// --- internal: minimal wire-format reader -----------------------------------

class WireReader {
	private offset = 0;
	constructor(private readonly buf: Uint8Array) {}

	hasMore(): boolean {
		return this.offset < this.buf.length;
	}

	uvarint(): bigint {
		let result = 0n;
		let shift = 0n;
		while (true) {
			if (this.offset >= this.buf.length) {
				throw new Error("varint past end of buffer");
			}
			const b = this.buf[this.offset++];
			result |= BigInt(b & 0x7f) << shift;
			if ((b & 0x80) === 0) return result;
			shift += 7n;
			if (shift > 70n) throw new Error("varint too long");
		}
	}

	uvarintInt(): number {
		return Number(this.uvarint());
	}

	take(n: number): Uint8Array {
		if (this.offset + n > this.buf.length) {
			throw new Error("length-delimited field past end of buffer");
		}
		const out = this.buf.subarray(this.offset, this.offset + n);
		this.offset += n;
		return out;
	}

	skipField(wireType: number): void {
		switch (wireType) {
			case 0: // varint
				this.uvarint();
				return;
			case 1: // fixed64
				this.offset += 8;
				return;
			case 2: {
				// length-delimited
				const len = this.uvarintInt();
				this.offset += len;
				return;
			}
			case 5: // fixed32
				this.offset += 4;
				return;
			default:
				throw new Error(`unsupported wire type ${wireType}`);
		}
	}
}

function decodeBuildEvent(
	bytes: Uint8Array,
	sets: Map<string, ParsedNamedSet>,
	targets: ParsedTarget[],
): void {
	const reader = new WireReader(bytes);
	let setId: string | undefined;
	let targetLabel: string | undefined;
	let targetAspect = "";
	let isNamedSetPayload = false;
	let isCompletedPayload = false;
	let namedSetBytes: Uint8Array | undefined;
	let completedBytes: Uint8Array | undefined;

	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;

		if (fieldNum === 1 && wireType === 2) {
			// id (BuildEventId)
			const idLen = reader.uvarintInt();
			const idBytes = reader.take(idLen);
			const idResult = decodeBuildEventId(idBytes);
			setId = idResult.setId;
			if (idResult.targetCompleted) {
				targetLabel = idResult.targetCompleted.label;
				targetAspect = idResult.targetCompleted.aspect;
			}
		} else if (fieldNum === 15 && wireType === 2) {
			// named_set_of_files
			isNamedSetPayload = true;
			const len = reader.uvarintInt();
			namedSetBytes = reader.take(len);
		} else if (fieldNum === 8 && wireType === 2) {
			// completed (TargetComplete)
			isCompletedPayload = true;
			const len = reader.uvarintInt();
			completedBytes = reader.take(len);
		} else {
			reader.skipField(wireType);
		}
	}

	if (isNamedSetPayload && setId && namedSetBytes) {
		sets.set(setId, decodeNamedSetOfFiles(namedSetBytes));
	} else if (isCompletedPayload && targetLabel !== undefined && completedBytes) {
		targets.push({
			label: targetLabel,
			aspect: targetAspect,
			outputGroups: decodeTargetComplete(completedBytes),
		});
	}
}

interface DecodedBuildEventId {
	setId?: string;
	targetCompleted?: { label: string; aspect: string };
}

function decodeBuildEventId(bytes: Uint8Array): DecodedBuildEventId {
	const reader = new WireReader(bytes);
	let setId: string | undefined;
	let targetCompleted: { label: string; aspect: string } | undefined;

	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;

		if (fieldNum === 13 && wireType === 2) {
			// named_set (NamedSetOfFilesId)
			const len = reader.uvarintInt();
			const inner = reader.take(len);
			setId = decodeNamedSetOfFilesId(inner);
		} else if (fieldNum === 5 && wireType === 2) {
			// target_completed (TargetCompletedId)
			const len = reader.uvarintInt();
			const inner = reader.take(len);
			targetCompleted = decodeTargetCompletedId(inner);
		} else {
			reader.skipField(wireType);
		}
	}
	return { setId, targetCompleted };
}

function decodeNamedSetOfFilesId(bytes: Uint8Array): string {
	const reader = new WireReader(bytes);
	let id = "";
	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;
		if (fieldNum === 1 && wireType === 2) {
			const len = reader.uvarintInt();
			id = bytesToUtf8(reader.take(len));
		} else {
			reader.skipField(wireType);
		}
	}
	return id;
}

function decodeTargetCompletedId(bytes: Uint8Array): {
	label: string;
	aspect: string;
} {
	const reader = new WireReader(bytes);
	let label = "";
	let aspect = "";
	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;
		if (fieldNum === 1 && wireType === 2) {
			const len = reader.uvarintInt();
			label = bytesToUtf8(reader.take(len));
		} else if (fieldNum === 2 && wireType === 2) {
			const len = reader.uvarintInt();
			aspect = bytesToUtf8(reader.take(len));
		} else {
			reader.skipField(wireType);
		}
	}
	return { label, aspect };
}

function decodeNamedSetOfFiles(bytes: Uint8Array): ParsedNamedSet {
	const reader = new WireReader(bytes);
	const files: ParsedFile[] = [];
	const childSetIds: string[] = [];
	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;
		if (fieldNum === 1 && wireType === 2) {
			const len = reader.uvarintInt();
			files.push(decodeFile(reader.take(len)));
		} else if (fieldNum === 2 && wireType === 2) {
			const len = reader.uvarintInt();
			childSetIds.push(decodeNamedSetOfFilesId(reader.take(len)));
		} else {
			reader.skipField(wireType);
		}
	}
	return { files, childSetIds };
}

function decodeFile(bytes: Uint8Array): ParsedFile {
	const reader = new WireReader(bytes);
	const file: ParsedFile = { name: "" };
	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;
		if (fieldNum === 1 && wireType === 2) {
			const len = reader.uvarintInt();
			file.name = bytesToUtf8(reader.take(len));
		} else if (fieldNum === 2 && wireType === 2) {
			const len = reader.uvarintInt();
			file.uri = bytesToUtf8(reader.take(len));
		} else if (fieldNum === 5 && wireType === 2) {
			const len = reader.uvarintInt();
			file.digest = bytesToUtf8(reader.take(len));
		} else if (fieldNum === 6 && wireType === 0) {
			file.sizeBytes = Number(reader.uvarint());
		} else {
			reader.skipField(wireType);
		}
	}
	return file;
}

function decodeTargetComplete(bytes: Uint8Array): ParsedOutputGroup[] {
	const reader = new WireReader(bytes);
	const groups: ParsedOutputGroup[] = [];
	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;
		if (fieldNum === 2 && wireType === 2) {
			const len = reader.uvarintInt();
			groups.push(decodeOutputGroup(reader.take(len)));
		} else {
			reader.skipField(wireType);
		}
	}
	return groups;
}

function decodeOutputGroup(bytes: Uint8Array): ParsedOutputGroup {
	const reader = new WireReader(bytes);
	let name = "";
	const rootSetIds: string[] = [];
	let incomplete = false;
	while (reader.hasMore()) {
		const tag = reader.uvarintInt();
		const fieldNum = tag >>> 3;
		const wireType = tag & 0x7;
		if (fieldNum === 1 && wireType === 2) {
			const len = reader.uvarintInt();
			name = bytesToUtf8(reader.take(len));
		} else if (fieldNum === 3 && wireType === 2) {
			const len = reader.uvarintInt();
			rootSetIds.push(decodeNamedSetOfFilesId(reader.take(len)));
		} else if (fieldNum === 4 && wireType === 0) {
			incomplete = reader.uvarintInt() !== 0;
		} else {
			reader.skipField(wireType);
		}
	}
	return { name, rootSetIds, incomplete };
}

// --- utilities --------------------------------------------------------------

function base64ToBytes(b64: string): Uint8Array {
	const binary = atob(b64);
	const bytes = new Uint8Array(binary.length);
	for (let i = 0; i < binary.length; i++) {
		bytes[i] = binary.charCodeAt(i);
	}
	return bytes;
}

const utf8Decoder = new TextDecoder("utf-8");
function bytesToUtf8(bytes: Uint8Array): string {
	return utf8Decoder.decode(bytes);
}
