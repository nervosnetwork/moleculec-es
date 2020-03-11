package generator

import (
	"fmt"
	"io"
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
	fmt.Fprintln(writer, "}\n")

	for _, declaration := range schema.Declarations {
		fmt.Fprintf(writer, "export class %s {\n", declaration.Name)
		fmt.Fprintln(writer, "  constructor(reader) {")
		fmt.Fprintln(writer, "    if (reader instanceof Object && reader.toArrayBuffer instanceof Function) {")
		fmt.Fprintln(writer, "      reader = reader.toArrayBuffer();")
		fmt.Fprintln(writer, "    }")
		fmt.Fprintln(writer, "    if (!(reader instanceof ArrayBuffer)) {")
		fmt.Fprintln(writer, "      throw new Error(\"Provided value must be an ArrayBuffer or can be transformed into ArrayBuffer!\")")
		fmt.Fprintln(writer, "    }")
		fmt.Fprintln(writer, "    this.view = new DataView(reader);")
		fmt.Fprintln(writer, "  }\n")
		switch declaration.Type {
		case "array":
			if declaration.Item == "byte" {
				fmt.Fprintln(writer, "  indexAt(i) {")
				fmt.Fprintln(writer, "    return this.view.getUint8(i);")
				fmt.Fprintln(writer, "  }\n")
				fmt.Fprintln(writer, "  view() {")
				fmt.Fprintln(writer, "    return this.view;")
				fmt.Fprintln(writer, "  }\n")
				switch declaration.ItemCount {
				case 2:
					fmt.Fprintln(writer, "  toBigEndianUint16() {")
					fmt.Fprintln(writer, "    return this.view.getUint16(0, false);")
					fmt.Fprintln(writer, "  }\n")
					fmt.Fprintln(writer, "  toLittleEndianUint16() {")
					fmt.Fprintln(writer, "    return this.view.getUint16(0, true);")
					fmt.Fprintln(writer, "  }\n")
				case 4:
					fmt.Fprintln(writer, "  toBigEndianUint32() {")
					fmt.Fprintln(writer, "    return this.view.getUint32(0, false);")
					fmt.Fprintln(writer, "  }\n")
					fmt.Fprintln(writer, "  toLittleEndianUint32() {")
					fmt.Fprintln(writer, "    return this.view.getUint32(0, true);")
					fmt.Fprintln(writer, "  }\n")
				case 8:
					fmt.Fprintln(writer, "  toBigEndianBigUint64() {")
					fmt.Fprintln(writer, "    return this.view.getBigUint64(0, false);")
					fmt.Fprintln(writer, "  }\n")
					fmt.Fprintln(writer, "  toLittleEndianBigUint64() {")
					fmt.Fprintln(writer, "    return this.view.getUint64(0, true);")
					fmt.Fprintln(writer, "  }\n")
				}
				fmt.Fprintln(writer, "  size() {")
				fmt.Fprintf(writer, "    return %d;\n", declaration.ItemCount)
				fmt.Fprintln(writer, "  }")
			} else {
				fmt.Fprintln(writer, "  indexAt(i) {")
				fmt.Fprintf(writer, "    return new %s(this.view.buffer.slice(i * %s.size(), (i + 1) * %s.size());\n", declaration.Item, declaration.Item, declaration.Item)
				fmt.Fprintln(writer, "  }\n")
				fmt.Fprintln(writer, "  size() {")
				fmt.Fprintf(writer, "    return %s.size() * %d;\n", declaration.Item, declaration.ItemCount)
				fmt.Fprintln(writer, "  }")
			}
		default:
			return fmt.Errorf("Invalid declaration type: %s", declaration.Type)
		}
		fmt.Fprintln(writer, "}\n")

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
		default:
			return fmt.Errorf("Invalid declaration type: %s", declaration.Type)
		}
		fmt.Fprintln(writer, "}\n")
	}

	return nil
}
