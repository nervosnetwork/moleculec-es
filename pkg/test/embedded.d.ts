export interface CastToArrayBuffer {
  toArrayBuffer(): ArrayBuffer;
}

export type CanCastToArrayBuffer = ArrayBuffer | CastToArrayBuffer;

export interface CreateOptions {
  validate?: boolean;
}

export interface UnionType {
  type: string;
  value: any;
}

export type Uint32Type = CanCastToArrayBuffer;

export type BytesType = CanCastToArrayBuffer;

export type BytesOptType = BytesType | undefined;

export type BytesVecType = BytesType[];

export type BytesMatrixType = BytesVecType[];

export interface BytesTableType {
  bytes_matrix: BytesMatrixType;
}

export type BytesTableOptType = BytesTableType | undefined;

export type BytesTableOptVecType = BytesTableOptType[];

export interface BlockType {
  data: BytesTableOptVecType;
}

export function SerializeUint32(value: CanCastToArrayBuffer): ArrayBuffer;
export class Uint32 {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): number;
  raw(): ArrayBuffer;
  toBigEndianUint32(): number;
  toLittleEndianUint32(): number;
  static size(): Number;
}

export function SerializeBytes(value: CanCastToArrayBuffer): ArrayBuffer;
export class Bytes {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): number;
  raw(): ArrayBuffer;
  length(): number;
}

export function SerializeBytesOpt(value: BytesType | null): ArrayBuffer;
export class BytesOpt {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  value(): Bytes;
  hasValue(): boolean;
}

export function SerializeBytesVec(value: Array<BytesType>): ArrayBuffer;
export class BytesVec {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): Bytes;
  length(): number;
}

export function SerializeBytesMatrix(value: Array<BytesVecType>): ArrayBuffer;
export class BytesMatrix {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): BytesVec;
  length(): number;
}

export function SerializeBytesTable(value: BytesTableType): ArrayBuffer;
export class BytesTable {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  getBytesMatrix(): BytesMatrix;
}

export function SerializeBytesTableOpt(value: BytesTableType | null): ArrayBuffer;
export class BytesTableOpt {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  value(): BytesTable;
  hasValue(): boolean;
}

export function SerializeBytesTableOptVec(value: Array<BytesTableOptType>): ArrayBuffer;
export class BytesTableOptVec {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): BytesTableOpt;
  length(): number;
}

export function SerializeBlock(value: BlockType): ArrayBuffer;
export class Block {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  getData(): BytesTableOptVec;
}

