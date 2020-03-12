package generator

import (
	"fmt"
	"io"
	"strings"

	"github.com/iancoleman/strcase"
)

type innerWriter struct {
	err    error
	writer io.Writer
}

func (w *innerWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return len(p), nil
	}
	var n int
	n, w.err = w.writer.Write(p)
	return n, nil
}

func Generate(schema Schema, writer io.Writer) error {
	iw := &innerWriter{
		writer: writer,
	}
	err := doGenerate(schema, iw)
	if err != nil {
		return err
	}
	return iw.err
}

func doGenerate(schema Schema, writer *innerWriter) error {
	fmt.Fprintln(writer, `function verifyAndExtractOffsets(view, expectedFieldCount, compatible) {
  if (view.byteLength < 4) {
    throw new Error(`+"`"+`Data should at least be 4 bytes long! Actual: ${view.byteLength}`+"`"+`);
  }
  const requiredByteLength = view.getUint32(0, true);
  if (requiredByteLength !== view.byteLength) {
    throw new Error(`+"`"+`Invalid data length! Required: ${requiredByteLength}, actual: ${view.byteLength}`+"`"+`);
  }
  if (requiredByteLength === 4) {
    return [requiredByteLength];
  }
  if (requiredByteLength < 8) {
    throw new Error(`+"`"+`Non empty data should at least be of length 8! Actual: ${view.byteLength}`+"`"+`);
  }
  const firstOffset = view.getUint32(4, true);
  if (firstOffset % 4 !== 0 || firstOffset < 8) {
    throw new Error(`+"`"+`Invalid first offset: ${firstOffset}`+"`"+`);
  }
  const itemCount := firstOffset / 4 - 1;
  if (itemCount < expectedFieldCount) {
    throw new Error(`+"`"+`Item count not enough! Required: ${expectedFieldCount}, actual: ${itemCount}`+"`"+`);
  } else if ((!compatible) && itemCount > expectedFieldCount) {
    throw new Error(`+"`"+`Item count is more than required! Required: ${expectedFieldCount}, actual: ${itemCount}`+"`"+`);
  }
  if requiredByteLength < firstOffset {
    throw new Error(`+"`"+`First offset is larger than byte length: ${firstOffset}`+"`"+`);
  }
  const offsets = [];
  for (let i = 0; i < itemCount; i++) {
    const start = 4 + i * 4;
    offsets.push(view.getUint32(start, true));
  }
  offsets.push(requiredByteLength);
  for (let i = 0; i < offsets.length - 1; i++) {
    if (offsets[i] > offsets[i + 1]) {
      throw new Error(`+"`"+`Offset index ${i}: ${offsets[i]} is larger than offset index ${i + 1}: ${offsets[i + 1]}`+"`"+`);
    }
  }
  return offsets;
}

function serializeTable(buffers) {
  const itemCount = buffers.length;
  let totalSize = 4 * (itemCount + 1);
  const offsets = [];

  for (let i = 0; i < itemCount; i++) {
    offsets.push(totalSize);
    totalSize += buffers[i].byteLength;
  }

  const buffer = new ArrayBuffer(totalSize);
  const array = new Uint8Array(buffer);
  const view = new DataView(buffer);

  view.setUint32(0, totalSize, true);
  for (let i = 0; i < itemCount; i++) {
    view.setUint32(4 + i * 4, offsets[i], true);
    array.set(new Uint8Array(buffers[i]), offsets[i]);
  }
  return buffer;
}`)
	for _, declaration := range schema.Declarations {
		fmt.Fprintf(writer, "export class %s {\n", declaration.Name)
		fmt.Fprintln(writer, "  constructor(reader, { validate = true } = {}) {")
		fmt.Fprintln(writer, "    if (reader instanceof Object && reader.toArrayBuffer instanceof Function) {")
		fmt.Fprintln(writer, "      reader = reader.toArrayBuffer();")
		fmt.Fprintln(writer, "    }")
		fmt.Fprintln(writer, "    if (!(reader instanceof ArrayBuffer)) {")
		fmt.Fprintln(writer, "      throw new Error(\"Provided value must be an ArrayBuffer or can be transformed into ArrayBuffer!\")")
		fmt.Fprintln(writer, "    }")
		fmt.Fprintln(writer, "    this.view = new DataView(reader);")
		fmt.Fprintln(writer, "    if (validate) {")
		fmt.Fprintln(writer, "      this.validate();")
		fmt.Fprintln(writer, "    }")
		fmt.Fprintln(writer, "  }")
		fmt.Fprintln(writer)
		switch declaration.Type {
		case "array":
			if declaration.Item == "byte" {
				fmt.Fprintln(writer, "  validate(compatible = false) {")
				fmt.Fprintf(writer, "    if (this.view.byteLength !== %d) {\n", declaration.ItemCount)
				fmt.Fprintf(writer, "      throw new Error(`Invalid data length! Required: %d, actual: ${this.view.byteLength}`);\n", declaration.ItemCount)
				fmt.Fprintln(writer, "    }")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, "  indexAt(i) {")
				fmt.Fprintln(writer, "    return this.view.getUint8(i);")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, "  raw() {")
				fmt.Fprintln(writer, "    return this.view;")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				switch declaration.ItemCount {
				case 2:
					fmt.Fprintln(writer, "  toBigEndianUint16() {")
					fmt.Fprintln(writer, "    return this.view.getUint16(0, false);")
					fmt.Fprintln(writer, "  }")
					fmt.Fprintln(writer)
					fmt.Fprintln(writer, "  toLittleEndianUint16() {")
					fmt.Fprintln(writer, "    return this.view.getUint16(0, true);")
					fmt.Fprintln(writer, "  }")
					fmt.Fprintln(writer)
				case 4:
					fmt.Fprintln(writer, "  toBigEndianUint32() {")
					fmt.Fprintln(writer, "    return this.view.getUint32(0, false);")
					fmt.Fprintln(writer, "  }")
					fmt.Fprintln(writer)
					fmt.Fprintln(writer, "  toLittleEndianUint32() {")
					fmt.Fprintln(writer, "    return this.view.getUint32(0, true);")
					fmt.Fprintln(writer, "  }")
					fmt.Fprintln(writer)
				case 8:
					fmt.Fprintln(writer, "  toBigEndianBigUint64() {")
					fmt.Fprintln(writer, "    return this.view.getBigUint64(0, false);")
					fmt.Fprintln(writer, "  }")
					fmt.Fprintln(writer)
					fmt.Fprintln(writer, "  toLittleEndianBigUint64() {")
					fmt.Fprintln(writer, "    return this.view.getUint64(0, true);")
					fmt.Fprintln(writer, "  }")
					fmt.Fprintln(writer)
				}
				fmt.Fprintln(writer, "  size() {")
				fmt.Fprintf(writer, "    return %d;\n", declaration.ItemCount)
				fmt.Fprintln(writer, "  }")
			} else {
				fmt.Fprintln(writer, "  validate(compatible = false) {")
				fmt.Fprintf(writer, "    if (this.view.byteLength !== %s.size() * %d) {\n", declaration.Item, declaration.ItemCount)
				fmt.Fprintf(writer, "      throw new Error(`Invalid data length! Required: ${%s.size() * %d}, actual: ${this.view.byteLength}`);\n", declaration.Item, declaration.ItemCount)
				fmt.Fprintln(writer, "    }")
				fmt.Fprintf(writer, "    for (let i = 0; i < %d; i++) {\n", declaration.ItemCount)
				fmt.Fprintln(writer, "      const item = this.indexAt(i);")
				fmt.Fprintln(writer, "      item.validate(compatible);")
				fmt.Fprintln(writer, "    }")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, "  indexAt(i) {")
				fmt.Fprintf(writer, "    return new %s(this.view.buffer.slice(i * %s.size(), (i + 1) * %s.size(), { validate: false });\n", declaration.Item, declaration.Item, declaration.Item)
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, "  size() {")
				fmt.Fprintf(writer, "    return %s.size() * %d;\n", declaration.Item, declaration.ItemCount)
				fmt.Fprintln(writer, "  }")
			}
		case "fixvec":
			if declaration.Item == "byte" {
				fmt.Fprintln(writer, "  validate(compatible = false) {")
				fmt.Fprintln(writer, "    if (this.view.byteLength < 4) {")
				fmt.Fprintln(writer, "      throw new Error(`Data should at least be 4 bytes long! Actual: ${this.view.byteLength}`);")
				fmt.Fprintln(writer, "    }")
				fmt.Fprintln(writer, "    const requiredByteLength = this.length() + 4;")
				fmt.Fprintln(writer, "    if (this.view.byteLength !== requiredByteLength) {")
				fmt.Fprintln(writer, "      throw new Error(`Invalid data length! Required: ${requiredByteLength}, actual: ${this.view.byteLength}`);")
				fmt.Fprintln(writer, "    }")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, "  raw() {")
				fmt.Fprintln(writer, "    return new DataView(this.view.buffer, 4);")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, "  indexAt(i) {")
				fmt.Fprintln(writer, "    return this.view.getUint8(4 + i);")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
			} else {
				fmt.Fprintln(writer, "  validate(compatible = false) {")
				fmt.Fprintln(writer, "    if (this.view.byteLength < 4) {")
				fmt.Fprintln(writer, "      throw new Error(`Data should at least be 4 bytes long! Actual: ${this.view.byteLength}`);")
				fmt.Fprintln(writer, "    }")
				fmt.Fprintf(writer, "    const requiredByteLength = this.length() * %s.size() + 4;\n", declaration.Item)
				fmt.Fprintln(writer, "    if (this.view.byteLength !== requiredByteLength) {")
				fmt.Fprintln(writer, "      throw new Error(`Invalid data length! Required: ${requiredByteLength}, actual: ${this.view.byteLength}`);")
				fmt.Fprintln(writer, "    }")
				fmt.Fprintf(writer, "    for (let i = 0; i < %d; i++) {\n", declaration.ItemCount)
				fmt.Fprintln(writer, "      const item = this.indexAt(i);")
				fmt.Fprintln(writer, "      item.validate(compatible);")
				fmt.Fprintln(writer, "    }")
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, "  indexAt(i) {")
				fmt.Fprintf(writer, "    return new %s(this.view.buffer.slice(4 + i * %s.size(), 4 + (i + 1) * %s.size(), { validate: false });\n", declaration.Item, declaration.Item, declaration.Item)
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer)
			}
			fmt.Fprintln(writer, "  length() {")
			fmt.Fprintln(writer, "    return this.view.getUint32(0, true);")
			fmt.Fprintln(writer, "  }")
		case "struct":
			sizes := []string{"0"}
			for _, field := range declaration.Fields {
				if field.Type == "byte" {
					fmt.Fprintf(writer, `  %s() {
    return this.view.getUint8(%s);
  }`+"\n\n",
						strcase.ToLowerCamel(fmt.Sprintf("get_%s", field.Name)),
						strings.Join(sizes, " + "))
					sizes = append(sizes, "1")
				} else {
					fmt.Fprintf(writer, `  %s() {
    return new %s(this.view.buffer.slice(%s, %s.size()), { validate: false });
  }`+"\n\n",
						strcase.ToLowerCamel(fmt.Sprintf("get_%s", field.Name)),
						field.Type,
						strings.Join(sizes, " + "),
						field.Type)
					sizes = append(sizes, fmt.Sprintf("%s.size()", field.Type))
				}
			}
			fmt.Fprintf(writer, `  validate(compatible = false) {
    const requiredByteLength = %s;
    if (this.view.byteLength !== requiredByteLength) {
      throw new Error(`+"`"+`Invalid data length! Required: ${requiredByteLength}, actual: ${this.view.byteLength}`+"`"+`);
    }`+"\n", strings.Join(sizes, " + "))
			for _, field := range declaration.Fields {
				if field.Type != "byte" {
					fmt.Fprintf(writer, "    this.%s().validate(compatible);\n",
						strcase.ToLowerCamel(fmt.Sprintf("get_%s", field.Name)))
				}
			}
			fmt.Fprintln(writer, "  }")
		case "dynvec":
			fmt.Fprintf(writer, `  validate(compatible = false) {
    const offsets = verifyAndExtractOffsets(this.view, 0, true);
    for (let i = 0; i < len(offsets) - 1; i++) {
      new %s(this.view.buffer.slice(offsets[i], offsets[i + 1]), { validate: false }).validate();
    }
  }

  length() {
    if (this.view.byteLength < 8) {
      return 0;
    } else {
      return this.view.getUint32(4, true) / 4 - 1;
    }
  }

  indexAt(i) {
    const start = 4 + i * 4;
    const offset = this.view.getUint32(start, true);
    let offset_end = this.view.byteLength;
    if (i + 1 < this.length()) {
      offset_end = this.view.getUint32(start + 4, true);
    }
    return new %s(this.view.buffer.slice(offset, offset_end), { validate: false };)
  }`+"\n", declaration.Item, declaration.Item)
		case "table":
			fmt.Fprintln(writer, `  validate(compatible = false) {
    const offsets = verifyAndExtractOffsets(this.view, 0, true);`)
			for i, field := range declaration.Fields {
				if field.Type == "byte" {
					fmt.Fprintf(writer, `    if (offset[%d] - offset[%d] !== 1) {
      throw new Error(`+"`"+`Invalid offset for %s: ${offset[%d]} - ${offset[%d]}`+"`"+`)
    }`+"\n", i+1, i, field.Name, i, i+1)
				} else {
					fmt.Fprintf(writer, "    new %s(this.view.buffer.slice(offset[%d], offset[%d]), { validate: false }).validate();\n", field.Type, i, i+1)
				}
			}
			fmt.Fprintln(writer, "  }")
			fmt.Fprintln(writer)
			for i, field := range declaration.Fields {
				last := i == len(declaration.Fields)-1
				fmt.Fprintf(writer, "  get%s() {\n", strcase.ToCamel(field.Name))
				if last {
					fmt.Fprintf(writer, `    const start = %d;
    const offset = this.view.getUint32(start, true);
    const offset_end = this.view.byteLength;`+"\n", 4+i*4)
				} else {
					fmt.Fprintf(writer, `    const start = %d;
    const offset = this.view.getUint32(start, true);
    const offset_end = this.view.getUint32(start + 4, true);`+"\n", 4+i*4)
				}
				if field.Type == "byte" {
					fmt.Fprintln(writer, "    return new DataView(this.view.buffer.slice(offset, offset_end)).getUint8(0);")
				} else {
					fmt.Fprintf(writer, "    return new %s(this.view.buffer.slice(offset, offset_end), { validate: false });\n", field.Type)
				}
				fmt.Fprintln(writer, "  }")
				if !last {
					fmt.Fprintln(writer)
				}
			}
		case "option":
			if declaration.Item == "byte" {
				fmt.Fprintln(writer, `  validate(compatible = false) {
    if (this.view.byteLength !== 0 && this.view.byteLength !== 1) {
      throw new Error(`+"`"+`Option that stores byte can only be of length 0 or 1! Actual: ${this.view.byteLength}`+"`"+`);
    }
  }`)
				fmt.Fprintln(writer)
				fmt.Fprintln(writer, `  value() {
    return this.view.getUint8(0);
  }`)
				fmt.Fprintln(writer)
			} else {
				fmt.Fprintln(writer, `  validate(compatible = false) {
    if (this.hasValue()) {
      this.value().validate(compatible);
    }
  }`)
				fmt.Fprintln(writer)
				fmt.Fprintf(writer, `  value() {
    return new %s(this.view.buffer, { validate: false });
  }`+"\n", declaration.Item)
				fmt.Fprintln(writer)
			}
			fmt.Fprintln(writer, `  hasValue() {
    return this.view.byteLength > 0;
  }`)
		case "union":
			fmt.Fprintln(writer, `  validate(compatible = false) {
    if (this.view.byteLength < 4) {
      throw new Error(`+"`"+`Data should at least be 4 bytes long! Actual: ${this.view.byteLength}`+"`"+`);
    }
    const t = this.view.getUint32(0, true);
    switch (t) {`)
			for i, item := range declaration.Items {
				if item == "byte" {
					fmt.Fprintf(writer, `    case %d:
      if (this.view.byteLength !== 5) {
        throw new Error(`+"`"+`Invalid data length! Required: 5, actual: ${this.view.byteLength}`+"`"+`);
      }
      break;`+"\n", i)
				} else {
					fmt.Fprintf(writer, `    case %d:
      new %s(this.view.buffer.slice(4), { validate: false }).validate();
      break;`+"\n",
						i, item)
				}
			}
			fmt.Fprintln(writer, `    default:
      throw new Error(`+"`"+`Invalid type: ${t}`+"`"+`);
    }
  }`)
			fmt.Fprintln(writer)
			fmt.Fprintln(writer, `  unionType() {
    const t = this.view.getUint32(0, true);
    switch (t) {`)
			for i, item := range declaration.Items {
				fmt.Fprintf(writer, `    case %d:
      return "%s";`+"\n", i, item)
			}
			fmt.Fprintln(writer, `    default:
      throw new Error(`+"`"+`Invalid type: ${t}`+"`"+`);
    }
  }`)
			fmt.Fprintln(writer)
			fmt.Fprintln(writer, `  value() {
    const t = this.view.getUint32(0, true);
    switch (t) {`)
			for i, item := range declaration.Items {
				if item == "byte" {
					fmt.Fprintf(writer, `    case %d:
      return this.view.buffer.getUint8(4);`+"\n", i)
				} else {
					fmt.Fprintf(writer, `    case %d:
      return new %s(this.view.buffer.slice(4), { validate: false });`+"\n", i, item)
				}
			}
			fmt.Fprintln(writer, `    default:
      throw new Error(`+"`"+`Invalid type: ${t}`+"`"+`);
    }
  }`)
		default:
			return fmt.Errorf("Invalid declaration type: %s", declaration.Type)
		}
		fmt.Fprintln(writer, "}")
		fmt.Fprintln(writer)

		fmt.Fprintf(writer, "export function Serialize%s(value) {\n", declaration.Name)
		switch declaration.Type {
		case "array":
			if declaration.Item == "byte" {
				fmt.Fprintln(writer, "  return new Reader(value).toArrayBuffer();")
			} else {
				fmt.Fprintf(writer, "  const array = new Uint8Array(%s.size() * value.length);\n", declaration.Item)
				fmt.Fprintln(writer, "  for (let i = 0; i < value.length; i++) {")
				fmt.Fprintf(writer, "    const itemBuffer = Serialize%s(value[i]);\n", declaration.Item)
				fmt.Fprintf(writer, "    array.set(new Uint8Array(itemBuffer), i * %s.size());\n", declaration.Item)
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer, "  return array.buffer;")
			}
		case "fixvec":
			if declaration.Item == "byte" {
				fmt.Fprintln(writer, "  const reader = new Reader(value);")
				fmt.Fprintln(writer, "  const array = new Uint8Array(4 + reader.length());")
				fmt.Fprintln(writer, "  (new DataView(array.buffer)).setUint32(0, reader.length(), true);")
				fmt.Fprintln(writer, "  array.set(new Uint8Array(reader.toArrayBuffer()), 4);")
				fmt.Fprintln(writer, "  return array.buffer;")
			} else {
				fmt.Fprintf(writer, "  const array = new Uint8Array(4 + %s.size() * value.length);\n", declaration.Item)
				fmt.Fprintln(writer, "  (new DataView(array.buffer)).setUint32(0, value.length, true);")
				fmt.Fprintln(writer, "  for (let i = 0; i < value.length; i++) {")
				fmt.Fprintf(writer, "    const itemBuffer = Serialize%s(value[i]);\n", declaration.Item)
				fmt.Fprintf(writer, "    array.set(new Uint8Array(itemBuffer), 4 + i * %s.size());\n", declaration.Item)
				fmt.Fprintln(writer, "  }")
				fmt.Fprintln(writer, "  return array.buffer;")
			}
		case "struct":
			sizes := []string{"0"}
			for _, field := range declaration.Fields {
				if field.Type == "byte" {
					sizes = append(sizes, "1")
				} else {
					sizes = append(sizes, fmt.Sprintf("%s.size()", field.Type))
				}
			}
			fmt.Fprintf(writer, "  const array = new Uint8Array(%s);\n", strings.Join(sizes, " + "))
			fmt.Fprintln(writer, "  const view = new DataView(array.buffer);")
			sizes = []string{"0"}
			for _, field := range declaration.Fields {
				if field.Type == "byte" {
					fmt.Fprintf(writer, "  view.setUint8(%s, value.%s);\n", strings.Join(sizes, " + "), field.Name)
					sizes = append(sizes, "1")
				} else {
					fmt.Fprintf(writer, "  const itemBuffer = Serialize%s(value.%s);\n", field.Type, field.Name)
					fmt.Fprintf(writer, "  array.set(new Uint8Array(itemBuffer), %s);\n", strings.Join(sizes, " + "))
					sizes = append(sizes, fmt.Sprintf("%s.size()", field.Type))
				}
			}
			fmt.Fprintln(writer, "  return array.buffer;")
		case "dynvec":
			fmt.Fprintf(writer, "  return serializeTable(value.map(item => Serialize%s(item)));\n", declaration.Item)
		case "table":
			fmt.Fprintln(writer, "  const buffers = [];")
			for _, field := range declaration.Fields {
				camelCaseName := strcase.ToLowerCamel(field.Name)
				if field.Type == "byte" {
					fmt.Fprintf(writer, `  const %sView = new DataView(new ArrayBuffer(1));
  %sView.setUint8(0, value.%s);
  buffers.push(%sView.buffer)`+"\n", camelCaseName, camelCaseName, field.Name, camelCaseName)
				} else {
					fmt.Fprintf(writer, "  buffers.push(Serialize%s(value.%s));\n", field.Type, field.Name)
				}
			}
			fmt.Fprintln(writer, "  return serializeTable(buffers);")
		case "option":
			if declaration.Item == "byte" {
				fmt.Fprintln(writer, `  if (value) {
    const buffer = new ArrayBuffer(1);
    const view = new DataView(buffer);
    view.setUint8(0, value);
    return buffer;
  } else {
    return new ArrayBuffer(0);
  }`)
			} else {
				fmt.Fprintf(writer, `  if (value) {
    return Serialize%s(value);
  } else {
    return new ArrayBuffer(0);
  }`+"\n", declaration.Item)
			}
		case "union":
			fmt.Fprintln(writer, "  switch (value.type) {")
			for i, item := range declaration.Items {
				fmt.Fprintf(writer, "  case \"%s\":\n", item)
				if item == "byte" {
					fmt.Fprintf(writer, `    const view = new DataView(new ArrayBuffer(5));
    view.setUint32(0, %d, true);
    view.setUint8(4, value.value);
    return view.buffer;`+"\n", i)
				} else {
					fmt.Fprintf(writer, `    const itemBuffer = Serialize%s(value.value);
    const array = new Uint8Array(4 + itemBuffer.byteLength);
    const view = new DataView(array.buffer);
    view.setUint32(0, %d, true);
    array.set(new Uint8Array(itemBuffer), 4);
    return array.buffer;`+"\n", item, i)
				}
			}
			fmt.Fprintln(writer, `  default:
    throw new Error(`+"`"+`Invalid type: ${value.type}`+"`"+`);
  }`+"\n")
		default:
			return fmt.Errorf("Invalid declaration type: %s", declaration.Type)
		}
		fmt.Fprintln(writer, "}")
		fmt.Fprintln(writer)
	}

	return nil
}
