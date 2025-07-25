// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v2.6.1
//   protoc               v3.19.1
// source: opentelemetry/proto/common/v1/common.proto

/* eslint-disable */
import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";

export const protobufPackage = "opentelemetry.proto.common.v1";

/**
 * AnyValue is used to represent any type of attribute value. AnyValue may contain a
 * primitive value such as a string or integer or it may contain an arbitrary nested
 * object containing arrays, key-value lists and primitives.
 */
export interface AnyValue {
  stringValue?: string | undefined;
  boolValue?: boolean | undefined;
  intValue?: string | undefined;
  doubleValue?: number | undefined;
  arrayValue?: ArrayValue | undefined;
  kvlistValue?: KeyValueList | undefined;
  bytesValue?: Uint8Array | undefined;
}

/**
 * ArrayValue is a list of AnyValue messages. We need ArrayValue as a message
 * since oneof in AnyValue does not allow repeated fields.
 */
export interface ArrayValue {
  /** Array of values. The array may be empty (contain 0 elements). */
  values: AnyValue[];
}

/**
 * KeyValueList is a list of KeyValue messages. We need KeyValueList as a message
 * since `oneof` in AnyValue does not allow repeated fields. Everywhere else where we need
 * a list of KeyValue messages (e.g. in Span) we use `repeated KeyValue` directly to
 * avoid unnecessary extra wrapping (which slows down the protocol). The 2 approaches
 * are semantically equivalent.
 */
export interface KeyValueList {
  /**
   * A collection of key/value pairs of key-value pairs. The list may be empty (may
   * contain 0 elements).
   * The keys MUST be unique (it is not allowed to have more than one
   * value with the same key).
   */
  values: KeyValue[];
}

/**
 * KeyValue is a key-value pair that is used to store Span attributes, Link
 * attributes, etc.
 */
export interface KeyValue {
  key: string;
  value: AnyValue | undefined;
}

/**
 * InstrumentationScope is a message representing the instrumentation scope information
 * such as the fully qualified name and version.
 */
export interface InstrumentationScope {
  /** An empty instrumentation scope name means the name is unknown. */
  name: string;
  version: string;
  /**
   * Additional attributes that describe the scope. [Optional].
   * Attribute keys MUST be unique (it is not allowed to have more than one
   * attribute with the same key).
   */
  attributes: KeyValue[];
  droppedAttributesCount: number;
}

/**
 * A reference to an Entity.
 * Entity represents an object of interest associated with produced telemetry: e.g spans, metrics, profiles, or logs.
 *
 * Status: [Development]
 */
export interface EntityRef {
  /**
   * The Schema URL, if known. This is the identifier of the Schema that the entity data
   * is recorded in. To learn more about Schema URL see
   * https://opentelemetry.io/docs/specs/otel/schemas/#schema-url
   *
   * This schema_url applies to the data in this message and to the Resource attributes
   * referenced by id_keys and description_keys.
   * TODO: discuss if we are happy with this somewhat complicated definition of what
   * the schema_url applies to.
   *
   * This field obsoletes the schema_url field in ResourceMetrics/ResourceSpans/ResourceLogs.
   */
  schemaUrl: string;
  /**
   * Defines the type of the entity. MUST not change during the lifetime of the entity.
   * For example: "service" or "host". This field is required and MUST not be empty
   * for valid entities.
   */
  type: string;
  /**
   * Attribute Keys that identify the entity.
   * MUST not change during the lifetime of the entity. The Id must contain at least one attribute.
   * These keys MUST exist in the containing {message}.attributes.
   */
  idKeys: string[];
  /**
   * Descriptive (non-identifying) attribute keys of the entity.
   * MAY change over the lifetime of the entity. MAY be empty.
   * These attribute keys are not part of entity's identity.
   * These keys MUST exist in the containing {message}.attributes.
   */
  descriptionKeys: string[];
}

function createBaseAnyValue(): AnyValue {
  return {
    stringValue: undefined,
    boolValue: undefined,
    intValue: undefined,
    doubleValue: undefined,
    arrayValue: undefined,
    kvlistValue: undefined,
    bytesValue: undefined,
  };
}

export const AnyValue: MessageFns<AnyValue> = {
  encode(message: AnyValue, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.stringValue !== undefined) {
      writer.uint32(10).string(message.stringValue);
    }
    if (message.boolValue !== undefined) {
      writer.uint32(16).bool(message.boolValue);
    }
    if (message.intValue !== undefined) {
      writer.uint32(24).int64(message.intValue);
    }
    if (message.doubleValue !== undefined) {
      writer.uint32(33).double(message.doubleValue);
    }
    if (message.arrayValue !== undefined) {
      ArrayValue.encode(message.arrayValue, writer.uint32(42).fork()).join();
    }
    if (message.kvlistValue !== undefined) {
      KeyValueList.encode(message.kvlistValue, writer.uint32(50).fork()).join();
    }
    if (message.bytesValue !== undefined) {
      writer.uint32(58).bytes(message.bytesValue);
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): AnyValue {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAnyValue();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.stringValue = reader.string();
          continue;
        }
        case 2: {
          if (tag !== 16) {
            break;
          }

          message.boolValue = reader.bool();
          continue;
        }
        case 3: {
          if (tag !== 24) {
            break;
          }

          message.intValue = reader.int64().toString();
          continue;
        }
        case 4: {
          if (tag !== 33) {
            break;
          }

          message.doubleValue = reader.double();
          continue;
        }
        case 5: {
          if (tag !== 42) {
            break;
          }

          message.arrayValue = ArrayValue.decode(reader, reader.uint32());
          continue;
        }
        case 6: {
          if (tag !== 50) {
            break;
          }

          message.kvlistValue = KeyValueList.decode(reader, reader.uint32());
          continue;
        }
        case 7: {
          if (tag !== 58) {
            break;
          }

          message.bytesValue = reader.bytes();
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AnyValue {
    return {
      stringValue: isSet(object.stringValue) ? globalThis.String(object.stringValue) : undefined,
      boolValue: isSet(object.boolValue) ? globalThis.Boolean(object.boolValue) : undefined,
      intValue: isSet(object.intValue) ? globalThis.String(object.intValue) : undefined,
      doubleValue: isSet(object.doubleValue) ? globalThis.Number(object.doubleValue) : undefined,
      arrayValue: isSet(object.arrayValue) ? ArrayValue.fromJSON(object.arrayValue) : undefined,
      kvlistValue: isSet(object.kvlistValue) ? KeyValueList.fromJSON(object.kvlistValue) : undefined,
      bytesValue: isSet(object.bytesValue) ? bytesFromBase64(object.bytesValue) : undefined,
    };
  },

  toJSON(message: AnyValue): unknown {
    const obj: any = {};
    if (message.stringValue !== undefined) {
      obj.stringValue = message.stringValue;
    }
    if (message.boolValue !== undefined) {
      obj.boolValue = message.boolValue;
    }
    if (message.intValue !== undefined) {
      obj.intValue = message.intValue;
    }
    if (message.doubleValue !== undefined) {
      obj.doubleValue = message.doubleValue;
    }
    if (message.arrayValue !== undefined) {
      obj.arrayValue = ArrayValue.toJSON(message.arrayValue);
    }
    if (message.kvlistValue !== undefined) {
      obj.kvlistValue = KeyValueList.toJSON(message.kvlistValue);
    }
    if (message.bytesValue !== undefined) {
      obj.bytesValue = base64FromBytes(message.bytesValue);
    }
    return obj;
  },

  create(base?: DeepPartial<AnyValue>): AnyValue {
    return AnyValue.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<AnyValue>): AnyValue {
    const message = createBaseAnyValue();
    message.stringValue = object.stringValue ?? undefined;
    message.boolValue = object.boolValue ?? undefined;
    message.intValue = object.intValue ?? undefined;
    message.doubleValue = object.doubleValue ?? undefined;
    message.arrayValue = (object.arrayValue !== undefined && object.arrayValue !== null)
      ? ArrayValue.fromPartial(object.arrayValue)
      : undefined;
    message.kvlistValue = (object.kvlistValue !== undefined && object.kvlistValue !== null)
      ? KeyValueList.fromPartial(object.kvlistValue)
      : undefined;
    message.bytesValue = object.bytesValue ?? undefined;
    return message;
  },
};

function createBaseArrayValue(): ArrayValue {
  return { values: [] };
}

export const ArrayValue: MessageFns<ArrayValue> = {
  encode(message: ArrayValue, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    for (const v of message.values) {
      AnyValue.encode(v!, writer.uint32(10).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): ArrayValue {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseArrayValue();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.values.push(AnyValue.decode(reader, reader.uint32()));
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ArrayValue {
    return {
      values: globalThis.Array.isArray(object?.values) ? object.values.map((e: any) => AnyValue.fromJSON(e)) : [],
    };
  },

  toJSON(message: ArrayValue): unknown {
    const obj: any = {};
    if (message.values?.length) {
      obj.values = message.values.map((e) => AnyValue.toJSON(e));
    }
    return obj;
  },

  create(base?: DeepPartial<ArrayValue>): ArrayValue {
    return ArrayValue.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<ArrayValue>): ArrayValue {
    const message = createBaseArrayValue();
    message.values = object.values?.map((e) => AnyValue.fromPartial(e)) || [];
    return message;
  },
};

function createBaseKeyValueList(): KeyValueList {
  return { values: [] };
}

export const KeyValueList: MessageFns<KeyValueList> = {
  encode(message: KeyValueList, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    for (const v of message.values) {
      KeyValue.encode(v!, writer.uint32(10).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): KeyValueList {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseKeyValueList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.values.push(KeyValue.decode(reader, reader.uint32()));
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): KeyValueList {
    return {
      values: globalThis.Array.isArray(object?.values) ? object.values.map((e: any) => KeyValue.fromJSON(e)) : [],
    };
  },

  toJSON(message: KeyValueList): unknown {
    const obj: any = {};
    if (message.values?.length) {
      obj.values = message.values.map((e) => KeyValue.toJSON(e));
    }
    return obj;
  },

  create(base?: DeepPartial<KeyValueList>): KeyValueList {
    return KeyValueList.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<KeyValueList>): KeyValueList {
    const message = createBaseKeyValueList();
    message.values = object.values?.map((e) => KeyValue.fromPartial(e)) || [];
    return message;
  },
};

function createBaseKeyValue(): KeyValue {
  return { key: "", value: undefined };
}

export const KeyValue: MessageFns<KeyValue> = {
  encode(message: KeyValue, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      AnyValue.encode(message.value, writer.uint32(18).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): KeyValue {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseKeyValue();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.value = AnyValue.decode(reader, reader.uint32());
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): KeyValue {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? AnyValue.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: KeyValue): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = AnyValue.toJSON(message.value);
    }
    return obj;
  },

  create(base?: DeepPartial<KeyValue>): KeyValue {
    return KeyValue.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<KeyValue>): KeyValue {
    const message = createBaseKeyValue();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? AnyValue.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseInstrumentationScope(): InstrumentationScope {
  return { name: "", version: "", attributes: [], droppedAttributesCount: 0 };
}

export const InstrumentationScope: MessageFns<InstrumentationScope> = {
  encode(message: InstrumentationScope, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.version !== "") {
      writer.uint32(18).string(message.version);
    }
    for (const v of message.attributes) {
      KeyValue.encode(v!, writer.uint32(26).fork()).join();
    }
    if (message.droppedAttributesCount !== 0) {
      writer.uint32(32).uint32(message.droppedAttributesCount);
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): InstrumentationScope {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstrumentationScope();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.version = reader.string();
          continue;
        }
        case 3: {
          if (tag !== 26) {
            break;
          }

          message.attributes.push(KeyValue.decode(reader, reader.uint32()));
          continue;
        }
        case 4: {
          if (tag !== 32) {
            break;
          }

          message.droppedAttributesCount = reader.uint32();
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): InstrumentationScope {
    return {
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      version: isSet(object.version) ? globalThis.String(object.version) : "",
      attributes: globalThis.Array.isArray(object?.attributes)
        ? object.attributes.map((e: any) => KeyValue.fromJSON(e))
        : [],
      droppedAttributesCount: isSet(object.droppedAttributesCount)
        ? globalThis.Number(object.droppedAttributesCount)
        : 0,
    };
  },

  toJSON(message: InstrumentationScope): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.version !== "") {
      obj.version = message.version;
    }
    if (message.attributes?.length) {
      obj.attributes = message.attributes.map((e) => KeyValue.toJSON(e));
    }
    if (message.droppedAttributesCount !== 0) {
      obj.droppedAttributesCount = Math.round(message.droppedAttributesCount);
    }
    return obj;
  },

  create(base?: DeepPartial<InstrumentationScope>): InstrumentationScope {
    return InstrumentationScope.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<InstrumentationScope>): InstrumentationScope {
    const message = createBaseInstrumentationScope();
    message.name = object.name ?? "";
    message.version = object.version ?? "";
    message.attributes = object.attributes?.map((e) => KeyValue.fromPartial(e)) || [];
    message.droppedAttributesCount = object.droppedAttributesCount ?? 0;
    return message;
  },
};

function createBaseEntityRef(): EntityRef {
  return { schemaUrl: "", type: "", idKeys: [], descriptionKeys: [] };
}

export const EntityRef: MessageFns<EntityRef> = {
  encode(message: EntityRef, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.schemaUrl !== "") {
      writer.uint32(10).string(message.schemaUrl);
    }
    if (message.type !== "") {
      writer.uint32(18).string(message.type);
    }
    for (const v of message.idKeys) {
      writer.uint32(26).string(v!);
    }
    for (const v of message.descriptionKeys) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): EntityRef {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntityRef();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.schemaUrl = reader.string();
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.type = reader.string();
          continue;
        }
        case 3: {
          if (tag !== 26) {
            break;
          }

          message.idKeys.push(reader.string());
          continue;
        }
        case 4: {
          if (tag !== 34) {
            break;
          }

          message.descriptionKeys.push(reader.string());
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EntityRef {
    return {
      schemaUrl: isSet(object.schemaUrl) ? globalThis.String(object.schemaUrl) : "",
      type: isSet(object.type) ? globalThis.String(object.type) : "",
      idKeys: globalThis.Array.isArray(object?.idKeys) ? object.idKeys.map((e: any) => globalThis.String(e)) : [],
      descriptionKeys: globalThis.Array.isArray(object?.descriptionKeys)
        ? object.descriptionKeys.map((e: any) => globalThis.String(e))
        : [],
    };
  },

  toJSON(message: EntityRef): unknown {
    const obj: any = {};
    if (message.schemaUrl !== "") {
      obj.schemaUrl = message.schemaUrl;
    }
    if (message.type !== "") {
      obj.type = message.type;
    }
    if (message.idKeys?.length) {
      obj.idKeys = message.idKeys;
    }
    if (message.descriptionKeys?.length) {
      obj.descriptionKeys = message.descriptionKeys;
    }
    return obj;
  },

  create(base?: DeepPartial<EntityRef>): EntityRef {
    return EntityRef.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<EntityRef>): EntityRef {
    const message = createBaseEntityRef();
    message.schemaUrl = object.schemaUrl ?? "";
    message.type = object.type ?? "";
    message.idKeys = object.idKeys?.map((e) => e) || [];
    message.descriptionKeys = object.descriptionKeys?.map((e) => e) || [];
    return message;
  },
};

function bytesFromBase64(b64: string): Uint8Array {
  const bin = globalThis.atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  arr.forEach((byte) => {
    bin.push(globalThis.String.fromCharCode(byte));
  });
  return globalThis.btoa(bin.join(""));
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export interface MessageFns<T> {
  encode(message: T, writer?: BinaryWriter): BinaryWriter;
  decode(input: BinaryReader | Uint8Array, length?: number): T;
  fromJSON(object: any): T;
  toJSON(message: T): unknown;
  create(base?: DeepPartial<T>): T;
  fromPartial(object: DeepPartial<T>): T;
}
