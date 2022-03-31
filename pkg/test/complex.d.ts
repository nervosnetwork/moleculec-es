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

export type Uint64Type = CanCastToArrayBuffer;

export type Uint128Type = CanCastToArrayBuffer;

export type Byte32Type = CanCastToArrayBuffer;

export type Uint256Type = CanCastToArrayBuffer;

export type BytesType = CanCastToArrayBuffer;

export type BytesOptType = BytesType | undefined;

export type BytesVecType = BytesType[];

export type Byte32VecType = Byte32Type[];

export type BaseDataType =
	|{ type: "Bytes", value: BytesType }
	|{ type: "Uint32", value: Uint32Type }
	|{ type: "Uint64", value: Uint64Type };

export type BigNumberType =
	|{ type: "Uint64", value: Uint64Type }
	|{ type: "Uint128", value: Uint128Type };

export type AllRoadType =
	|{ type: "BaseData", value: BaseDataType }
	|{ type: "BigNumber", value: BigNumberType };

export type BaseDataOptType = BaseDataType | undefined;

export type BigNumberOptType = BigNumberType | undefined;

export type BaseDataOptVecType = BaseDataOptType[];

export type BigNumberOptVecType = BigNumberOptType[];

export interface VehicleType {
  distance: BaseDataOptVecType;
  gas: BigNumberOptVecType;
}

export interface GarageType {
  car: VehicleType;
  monitor: AllRoadType;
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

export function SerializeUint64(value: CanCastToArrayBuffer): ArrayBuffer;
export class Uint64 {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): number;
  raw(): ArrayBuffer;
  static size(): Number;
}

export function SerializeUint128(value: CanCastToArrayBuffer): ArrayBuffer;
export class Uint128 {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): number;
  raw(): ArrayBuffer;
  static size(): Number;
}

export function SerializeByte32(value: CanCastToArrayBuffer): ArrayBuffer;
export class Byte32 {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): number;
  raw(): ArrayBuffer;
  static size(): Number;
}

export function SerializeUint256(value: CanCastToArrayBuffer): ArrayBuffer;
export class Uint256 {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): number;
  raw(): ArrayBuffer;
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

export function SerializeByte32Vec(value: Array<Byte32Type>): ArrayBuffer;
export class Byte32Vec {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): Byte32;
  length(): number;
}

export function SerializeBaseData(value: BaseDataType): ArrayBuffer;
export class BaseData {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  unionType(): string;
  value(): any;
}

export function SerializeBigNumber(value: BigNumberType): ArrayBuffer;
export class BigNumber {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  unionType(): string;
  value(): any;
}

export function SerializeAllRoad(value: AllRoadType): ArrayBuffer;
export class AllRoad {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  unionType(): string;
  value(): any;
}

export function SerializeBaseDataOpt(value: BaseDataType | null): ArrayBuffer;
export class BaseDataOpt {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  value(): BaseData;
  hasValue(): boolean;
}

export function SerializeBigNumberOpt(value: BigNumberType | null): ArrayBuffer;
export class BigNumberOpt {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  value(): BigNumber;
  hasValue(): boolean;
}

export function SerializeBaseDataOptVec(value: Array<BaseDataOptType>): ArrayBuffer;
export class BaseDataOptVec {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): BaseDataOpt;
  length(): number;
}

export function SerializeBigNumberOptVec(value: Array<BigNumberOptType>): ArrayBuffer;
export class BigNumberOptVec {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  indexAt(i: number): BigNumberOpt;
  length(): number;
}

export function SerializeVehicle(value: VehicleType): ArrayBuffer;
export class Vehicle {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  getDistance(): BaseDataOptVec;
  getGas(): BigNumberOptVec;
}

export function SerializeGarage(value: GarageType): ArrayBuffer;
export class Garage {
  constructor(reader: CanCastToArrayBuffer, options?: CreateOptions);
  validate(compatible?: boolean): void;
  getCar(): Vehicle;
  getMonitor(): AllRoad;
}

