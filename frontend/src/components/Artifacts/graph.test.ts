import { describe, expect, test } from "vitest";
import { parseArtifactGraph, walkSetFiles } from "./graph";

// Build a length-prefixed stream of serialized BuildEvent messages without
// pulling in a full proto runtime. We hand-encode just the fields the
// parser cares about, using proto3 wire format directly.

function uvarint(n: number): Uint8Array {
	const out: number[] = [];
	let v = n;
	while (v >= 0x80) {
		out.push((v & 0x7f) | 0x80);
		v >>>= 7;
	}
	out.push(v);
	return new Uint8Array(out);
}

function concat(...parts: Uint8Array[]): Uint8Array {
	let len = 0;
	for (const p of parts) len += p.length;
	const out = new Uint8Array(len);
	let off = 0;
	for (const p of parts) {
		out.set(p, off);
		off += p.length;
	}
	return out;
}

function tag(fieldNum: number, wireType: number): Uint8Array {
	return uvarint((fieldNum << 3) | wireType);
}

function lenDelim(fieldNum: number, payload: Uint8Array): Uint8Array {
	return concat(tag(fieldNum, 2), uvarint(payload.length), payload);
}

function stringField(fieldNum: number, value: string): Uint8Array {
	return lenDelim(fieldNum, new TextEncoder().encode(value));
}

function varintField(fieldNum: number, value: number): Uint8Array {
	return concat(tag(fieldNum, 0), uvarint(value));
}

// NamedSetOfFilesId { id = 1 string }
function namedSetOfFilesId(id: string): Uint8Array {
	return stringField(1, id);
}

// TargetCompletedId { label = 1 string, aspect = 2 string }
function targetCompletedId(label: string, aspect = ""): Uint8Array {
	const parts = [stringField(1, label)];
	if (aspect) parts.push(stringField(2, aspect));
	return concat(...parts);
}

// BuildEventId.named_set = 13 (NamedSetOfFilesId)
function buildEventIdNamedSet(id: string): Uint8Array {
	return lenDelim(13, namedSetOfFilesId(id));
}

// BuildEventId.target_completed = 5 (TargetCompletedId)
function buildEventIdTargetCompleted(
	label: string,
	aspect = "",
): Uint8Array {
	return lenDelim(5, targetCompletedId(label, aspect));
}

// File { name=1, uri=2, digest=5, length=6 }
function fileMsg(opts: {
	name: string;
	uri?: string;
	digest?: string;
	length?: number;
}): Uint8Array {
	const parts = [stringField(1, opts.name)];
	if (opts.uri !== undefined) parts.push(stringField(2, opts.uri));
	if (opts.digest !== undefined) parts.push(stringField(5, opts.digest));
	if (opts.length !== undefined) parts.push(varintField(6, opts.length));
	return concat(...parts);
}

// NamedSetOfFiles { files=1 repeated File, file_sets=2 repeated NamedSetOfFilesId }
function namedSetOfFiles(
	files: { name: string; uri?: string; digest?: string; length?: number }[],
	childIds: string[],
): Uint8Array {
	const parts: Uint8Array[] = [];
	for (const f of files) parts.push(lenDelim(1, fileMsg(f)));
	for (const id of childIds) parts.push(lenDelim(2, namedSetOfFilesId(id)));
	return concat(...parts);
}

// OutputGroup { name=1, file_sets=3 repeated NamedSetOfFilesId, incomplete=4 bool }
function outputGroup(
	name: string,
	rootSetIds: string[],
	incomplete = false,
): Uint8Array {
	const parts = [stringField(1, name)];
	for (const id of rootSetIds) parts.push(lenDelim(3, namedSetOfFilesId(id)));
	if (incomplete) parts.push(varintField(4, 1));
	return concat(...parts);
}

// TargetComplete { output_group=2 repeated OutputGroup }
function targetComplete(groups: Uint8Array[]): Uint8Array {
	const parts: Uint8Array[] = [];
	for (const g of groups) parts.push(lenDelim(2, g));
	return concat(...parts);
}

// BuildEvent { id=1 BuildEventId, named_set_of_files=15, completed=8 }
function buildEventNamedSet(setId: string, body: Uint8Array): Uint8Array {
	return concat(lenDelim(1, buildEventIdNamedSet(setId)), lenDelim(15, body));
}

function buildEventTargetCompleted(
	label: string,
	completed: Uint8Array,
	aspect = "",
): Uint8Array {
	return concat(
		lenDelim(1, buildEventIdTargetCompleted(label, aspect)),
		lenDelim(8, completed),
	);
}

// Length-prefix each BuildEvent message into a single stream.
function lengthPrefixedStream(events: Uint8Array[]): Uint8Array {
	const parts: Uint8Array[] = [];
	for (const e of events) {
		parts.push(uvarint(e.length));
		parts.push(e);
	}
	return concat(...parts);
}

describe("parseArtifactGraph", () => {
	test("decodes a simple single-target graph", () => {
		const set0 = namedSetOfFiles(
			[{ name: "hello", uri: "bytestream://x/y/z", digest: "abc", length: 42 }],
			[],
		);
		const completed = targetComplete([outputGroup("default", ["0"])]);

		const stream = lengthPrefixedStream([
			buildEventNamedSet("0", set0),
			buildEventTargetCompleted("//:hello", completed),
		]);

		const graph = parseArtifactGraph(stream);
		expect(graph.targets).toHaveLength(1);
		expect(graph.targets[0].label).toBe("//:hello");
		expect(graph.targets[0].outputGroups).toHaveLength(1);
		expect(graph.targets[0].outputGroups[0].name).toBe("default");
		expect(graph.targets[0].outputGroups[0].rootSetIds).toEqual(["0"]);

		expect(graph.sets.size).toBe(1);
		const files = Array.from(walkSetFiles(graph, ["0"]));
		expect(files).toEqual([
			{ name: "hello", uri: "bytestream://x/y/z", digest: "abc", sizeBytes: 42 },
		]);
	});

	test("walks transitive NamedSetOfFiles children", () => {
		const set0 = namedSetOfFiles(
			[{ name: "root-file" }],
			["1", "2"],
		);
		const set1 = namedSetOfFiles([{ name: "child1-file" }], []);
		const set2 = namedSetOfFiles([{ name: "child2-file" }], ["1"]);

		const stream = lengthPrefixedStream([
			buildEventNamedSet("0", set0),
			buildEventNamedSet("1", set1),
			buildEventNamedSet("2", set2),
		]);

		const graph = parseArtifactGraph(stream);
		const files = Array.from(walkSetFiles(graph, ["0"])).map((f) => f.name);
		expect(files.sort()).toEqual(["child1-file", "child2-file", "root-file"]);
	});

	test("terminates on cycles", () => {
		// set0 references set1, set1 references set0 — must not loop forever.
		const set0 = namedSetOfFiles([{ name: "a" }], ["1"]);
		const set1 = namedSetOfFiles([{ name: "b" }], ["0"]);
		const stream = lengthPrefixedStream([
			buildEventNamedSet("0", set0),
			buildEventNamedSet("1", set1),
		]);

		const graph = parseArtifactGraph(stream);
		const files = Array.from(walkSetFiles(graph, ["0"])).map((f) => f.name);
		expect(files.sort()).toEqual(["a", "b"]);
	});

	test("tolerates target referencing unknown set id", () => {
		const completed = targetComplete([outputGroup("default", ["missing"])]);
		const stream = lengthPrefixedStream([
			buildEventTargetCompleted("//:lonely", completed),
		]);
		const graph = parseArtifactGraph(stream);
		expect(graph.targets).toHaveLength(1);
		expect(Array.from(walkSetFiles(graph, ["missing"]))).toEqual([]);
	});

	test("captures multiple targets with multiple output groups", () => {
		const set0 = namedSetOfFiles([{ name: "bin" }], []);
		const set1 = namedSetOfFiles([{ name: "lib" }], []);
		const stream = lengthPrefixedStream([
			buildEventNamedSet("0", set0),
			buildEventNamedSet("1", set1),
			buildEventTargetCompleted(
				"//:a",
				targetComplete([
					outputGroup("default", ["0"]),
					outputGroup("runfiles", ["1"]),
				]),
			),
			buildEventTargetCompleted(
				"//:b",
				targetComplete([outputGroup("default", ["0"])]),
				"my_aspect",
			),
		]);

		const graph = parseArtifactGraph(stream);
		expect(graph.targets.map((t) => `${t.label}|${t.aspect}`)).toEqual([
			"//:a|",
			"//:b|my_aspect",
		]);
		expect(graph.targets[0].outputGroups.map((g) => g.name)).toEqual([
			"default",
			"runfiles",
		]);
	});

	test("handles empty stream", () => {
		const graph = parseArtifactGraph(new Uint8Array(0));
		expect(graph.targets).toHaveLength(0);
		expect(graph.sets.size).toBe(0);
	});
});
