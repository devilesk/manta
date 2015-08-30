package manta

import (
	"strconv"
)

// This is a list of encoding functions that are know to resolve to a correct layout.
// All of these have been verified
// ----------------------------------------------------------------------------------
//
// ID  NAME                           WEIGHT  ORIG  LEN  BITS
//  0  PlusOne                         36271          1   0
//  1  EncodingFinish                  25474          2   10
//  2  PlusTwo                         10334          4   1110
//  3  PlusN						    4128          5   11010
//  4  PlusThree                        1375          6   110010
//  5  PopAllButOnePlusOne              1837          6   110011
//  6  PushOneLeftDeltaOneRightZero      521          8   11011010
//  7  NonTopoComplexPack4Bits            99         10   1101100010
//  8  NonTopoComplex                     76         11   11011000111
//  9  PushOneLeftDeltaZeroRightZero      35         12   110110001101
// 10  PopOnePlusOne                       1     2   15   110110001100001
// 11  PopNPlusOne                         0         16   1101100011000110
// 12  PushTwoPack5LeftDeltaZero           0         17   11011000110010011
//
// Thanks to @spheenik for being resilient in his efforts to figure out the rest of the tree

// A single field to be read
type fieldpath_field struct {
	Name  string
	Path  string
	Field *dt_field
}

// A fieldpath, used to walk through the flattened table hierarchy
type fieldpath struct {
	parent   *dt
	fields   []*fieldpath_field
	index    []int32
	tree     *HuffmanTree
	finished bool
}

// Contains the weight and lookup function for a single operation
type fieldpathOp struct {
	Name     string
	Function func(*Reader, *fieldpath)
	Weight   int
}

// Global fieldpath lookup array
var fieldpathLookup = []fieldpathOp{
	{"PlusOne", PlusOne, 36271},
	{"PlusTwo", PlusTwo, 10334},
	{"PlusThree", PlusThree, 1375},
	{"PlusFour", PlusFour, 646},
	{"PlusN", PlusN, 4128},
	{"PushOneLeftDeltaZeroRightZero", PushOneLeftDeltaZeroRightZero, 35},
	{"PushOneLeftDeltaZeroRightNonZero", PushOneLeftDeltaZeroRightNonZero, 3},
	{"PushOneLeftDeltaOneRightZero", PushOneLeftDeltaOneRightZero, 521},
	{"PushOneLeftDeltaOneRightNonZero", PushOneLeftDeltaOneRightNonZero, 2942},
	{"PushOneLeftDeltaNRightZero", PushOneLeftDeltaNRightZero, 560},
	{"PushOneLeftDeltaNRightNonZero", PushOneLeftDeltaNRightNonZero, 471},
	{"PushOneLeftDeltaNRightNonZeroPack6Bits", PushOneLeftDeltaNRightNonZeroPack6Bits, 10530},
	{"PushOneLeftDeltaNRightNonZeroPack8Bits", PushOneLeftDeltaNRightNonZeroPack8Bits, 251},
	{"PushTwoLeftDeltaZero", PushTwoLeftDeltaZero, 0},
	{"PushTwoPack5LeftDeltaZero", PushTwoPack5LeftDeltaZero, 0},
	{"PushThreeLeftDeltaZero", PushThreeLeftDeltaZero, 0},
	{"PushThreePack5LeftDeltaZero", PushThreePack5LeftDeltaZero, 0},
	{"PushTwoLeftDeltaOne", PushTwoLeftDeltaOne, 0},
	{"PushTwoPack5LeftDeltaOne", PushTwoPack5LeftDeltaOne, 0},
	{"PushThreeLeftDeltaOne", PushThreeLeftDeltaOne, 0},
	{"PushThreePack5LeftDeltaOne", PushThreePack5LeftDeltaOne, 0},
	{"PushTwoLeftDeltaN", PushTwoLeftDeltaN, 0},
	{"PushTwoPack5LeftDeltaN", PushTwoPack5LeftDeltaN, 0},
	{"PushThreeLeftDeltaN", PushThreeLeftDeltaN, 0},
	{"PushThreePack5LeftDeltaN", PushThreePack5LeftDeltaN, 0},
	{"PushN", PushN, 0},
	{"PushNAndNonTopological", PushNAndNonTopological, 310},
	{"PopOnePlusOne", PopOnePlusOne, 2},
	{"PopOnePlusN", PopOnePlusN, 0},
	{"PopAllButOnePlusOne", PopAllButOnePlusOne, 1837},
	{"PopAllButOnePlusN", PopAllButOnePlusN, 149},
	{"PopAllButOnePlusNPack3Bits", PopAllButOnePlusNPack3Bits, 300},
	{"PopAllButOnePlusNPack6Bits", PopAllButOnePlusNPack6Bits, 634},
	{"PopNPlusOne", PopNPlusOne, 0},
	{"PopNPlusN", PopNPlusN, 0},
	{"PopNAndNonTopographical", PopNAndNonTopographical, 1},
	{"NonTopoComplex", NonTopoComplex, 76},
	{"NonTopoPenultimatePlusOne", NonTopoPenultimatePlusOne, 271},
	{"NonTopoComplexPack4Bits", NonTopoComplexPack4Bits, 99},
	{"FieldPathEncodeFinish", FieldPathEncodeFinish, 25474},
}

// Initialize a fieldpath object
func newFieldpath(parentTbl *dt, huf *HuffmanTree) *fieldpath {
	fp := &fieldpath{
		parent:   parentTbl,
		fields:   make([]*fieldpath_field, 0),
		index:    make([]int32, 0),
		tree:     huf,
		finished: false,
	}

	fp.index = append(fp.index, -1) // Always start at -1

	return fp
}

// Walk an encoded fieldpath based on a huffman tree
func (fp *fieldpath) walk(r *Reader) {
	cnt := 0
	root := fp.tree
	node := root

	for fp.finished == false {
		cnt++
		if r.readBits(1) == 1 {
			if i := (*node).Right(); i.IsLeaf() {
				node = root
				fieldpathLookup[i.Value()].Function(r, fp)
				fp.addField()

				_debugfl(6, "Reached in %d bits, %s, %d", cnt, fp.fields[len(fp.fields)-1].Name, r.pos)
				cnt = 0
			} else {
				node = &i
			}
		} else {
			if i := (*node).Left(); i.IsLeaf() {
				node = root
				fieldpathLookup[i.Value()].Function(r, fp)
				fp.addField()

				_debugfl(6, "Reached in %d bits, %s, %d", cnt, fp.fields[len(fp.fields)-1].Name, r.pos)
				cnt = 0
			} else {
				node = &i
			}
		}
	}

	// Will always add one additional field for the finishEncoding operation, remove it
	fp.fields = fp.fields[:len(fp.fields)-1]
}

// Adds a field based on the current index
func (fp *fieldpath) addField() {
	cDt := fp.parent

	var path string
	var name string
	i := 0

	for i = 0; i < len(fp.index)-1; i++ {
		path += strconv.Itoa(int(fp.index[i])) + "/"
	}

	_debugfl(10, "Adding field with path: %s%d", path, fp.index[len(fp.index)-1])

	for i = 0; i < len(fp.index)-1; i++ {
		if cDt.Properties[fp.index[i]].Table != nil {
			cDt = cDt.Properties[fp.index[i]].Table
			name += cDt.Name + "."
		} else {
			_panicf("expected table in fp properties")
		}
	}

	fp.fields = append(fp.fields, &fieldpath_field{name + cDt.Properties[fp.index[i]].Field.Name, path, cDt.Properties[fp.index[i]].Field})
}

// Returns a huffman tree based on the operation weights
func newFieldpathHuffman() HuffmanTree {
	// Generate feq map
	huffmanlist := make([]int, 40)
	for i, fpo := range fieldpathLookup {
		huffmanlist[i] = fpo.Weight
	}

	return buildTree(huffmanlist)
}

func PlusOne(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// Increment the index
	fp.index[len(fp.index)-1] += 1
}

func PlusTwo(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// Increment the index
	fp.index[len(fp.index)-1] += 2
}

func PlusThree(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// Increment the index
	fp.index[len(fp.index)-1] += 3
}

func PlusFour(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// Increment the index
	fp.index[len(fp.index)-1] += 4
}

func PlusN(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// This one reads a variable-length header encoding the number of bits
	// to read for N. Five is always added to the result.

	fp.index[len(fp.index)-1] += int32(r.readUBitVarFP()) + 5
}

func PushOneLeftDeltaZeroRightZero(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// Get current field and index
	fp.index = append(fp.index, 0)
}

func PushOneLeftDeltaZeroRightNonZero(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)

	// should be correct, not encountered however
	/*rBits := []int{2, 4, 10, 17, 30}

	for _, bits := range rBits {
		if r.readBits(1) == 1 {
			fp.index = append(fp.index, int32(r.readBits(bits)))
			_debugf("Index: %v, BitsL %v", fp.index, bits)
			return
		}
	}*/
}

func PushOneLeftDeltaOneRightZero(r *Reader, fp *fieldpath) {
	_debugf("Name: %s", fp.parent.Name)

	// Push +1, set index to 0
	fp.index[len(fp.index)-1] += 1
	fp.index = append(fp.index, 0)
}

func PushOneLeftDeltaOneRightNonZero(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushOneLeftDeltaNRightZero(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushOneLeftDeltaNRightNonZero(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)

}

func PushOneLeftDeltaNRightNonZeroPack6Bits(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushOneLeftDeltaNRightNonZeroPack8Bits(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushTwoLeftDeltaZero(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)

	// wrong
	//fp.index = append(fp.index, 0)
	//fp.index = append(fp.index, 0)
}

func PushTwoLeftDeltaOne(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushTwoLeftDeltaN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushTwoPack5LeftDeltaZero(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	fp.index = append(fp.index, int32(r.readBits(5)))
	fp.index = append(fp.index, int32(r.readBits(5)))
}

func PushTwoPack5LeftDeltaOne(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushTwoPack5LeftDeltaN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushThreeLeftDeltaZero(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushThreeLeftDeltaOne(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushThreeLeftDeltaN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushThreePack5LeftDeltaZero(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushThreePack5LeftDeltaOne(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushThreePack5LeftDeltaN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PushNAndNonTopological(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PopOnePlusOne(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// Check if we can pop an element
	if len(fp.index) <= 1 {
		_panicf("Trying to pop last element")
	}

	fp.index = fp.index[:len(fp.index)-1]
	fp.index[len(fp.index)-1] += 1
}

func PopOnePlusN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PopAllButOnePlusOne(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// Remove all hierarchy and index element
	fp.index = fp.index[:1]
	fp.index[len(fp.index)-1] += 1
}

func PopAllButOnePlusN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PopAllButOnePlusNPackN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PopAllButOnePlusNPack3Bits(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PopAllButOnePlusNPack6Bits(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PopNPlusOne(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	fp.index = fp.index[:len(fp.index)-(int(r.readUBitVarFP()))]
	fp.index[len(fp.index)-1] += 1
}

func PopNPlusN(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func PopNAndNonTopographical(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	fp.index = fp.index[:len(fp.index)-(int(r.readUBitVarFP()))]

	for i := 0; i < len(fp.index); i++ {
		if r.readBoolean() {
			fp.index[i] += r.readVarInt32()
		}
	}
}

func NonTopoComplex(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// See NonTopoComplexPack4Bits

	for i := 0; i < len(fp.index); i++ {
		if r.readBoolean() {
			fp.index[i] += r.readVarInt32()
		}
	}
}

func NonTopoPenultimatePlusOne(r *Reader, fp *fieldpath) {
	_panicf("Name: %s", fp.parent.Name)
}

func NonTopoComplexPack4Bits(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	// NonTopological = Disregard the hierarchy, work directly on the field
	// indizies for now
	//
	// Variables:
	// v4 = 0; // Incremented by 1 each loop
	// v3 = CFieldPath;
	//
	// Assumptions:
	// - Path data (array with MaxDepth) is first element of CFieldPath
	// - Current depth has an offset of 8 from CFieldPath
	//
	// Each loop does the following:
	// - Read 1 bit, if it is set, break
	// - Read 4 bits, substract 7 = v5
	// - Apply the data read to the v4'th index: v3[v4] += v5
	//
	// End condition:
	// - r.readBits(1) == 1
	// - Reached current depth (see assumption)

	for i := 0; i < len(fp.index); i++ {
		if r.readBoolean() {
			fp.index[i] += int32(r.readBits(4)) - 7
		}
	}
}

func FieldPathEncodeFinish(r *Reader, fp *fieldpath) {
	_debugfl(10, "Name: %s", fp.parent.Name)

	fp.finished = true
}
