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
	fmt.Fprintln(writer, "import { VERSION } from \"ckb-js-toolkit\";")
	fmt.Fprintln(writer, "if (parseInt(VERSION.split(\".\")[1]) < 5) {")
	fmt.Fprintln(writer, "  throw new Error(\"moleculec-es requires at least ckb-js-toolkit 0.5.0!\");")
	fmt.Fprintln(writer, "}")
	fmt.Fprintln(writer)

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
				fmt.Fprintln(writer, "  view() {")
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
				fmt.Fprintln(writer, "  view() {")
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
		default:
			return fmt.Errorf("Invalid declaration type: %s", declaration.Type)
		}
		fmt.Fprintln(writer, "}")
		fmt.Fprintln(writer)
	}

	return nil
}
