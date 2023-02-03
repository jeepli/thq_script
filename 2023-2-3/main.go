package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	inputfile  = ""
	outputfile = ""
)

func init() {
	flag.StringVar(&inputfile, "input", "", "input file name")
	flag.StringVar(&outputfile, "output", "", "output file name")
	flag.Parse()
}

type atomItem struct {
	id    uint32
	mId   uint32
	iType uint32
	q     float64
	x     float64
	y     float64
	z     float64
	m     int32
	n     int32
	k     int32
}

func (i *atomItem) String() string {
	return fmt.Sprintf("%d %d %d %g %g %g %g %d %d %d",
		i.id, i.mId, i.iType, i.q, i.x, i.y, i.z, i.m, i.n, i.k)
}

type atoms struct {
	note  string
	items []*atomItem
}

func (a *atoms) String() string {
	str := "Atoms # full\n"
	str += "\n"
	str += a.note
	str += "\n"
	for _, item := range a.items {
		str += item.String()
		str += "\n"
	}
	return str
}

type velItem struct {
	aId uint32
	vx  float64
	vy  float64
	vz  float64
}

func (i *velItem) String() string {
	return fmt.Sprintf("%d %g %g %g", i.aId, i.vx, i.vy, i.vz)
}

type velocities struct {
	note  string
	items []*velItem
}

func (v *velocities) String() string {
	str := "Velocities\n"
	str += "\n"
	str += v.note
	str += "\n"
	for _, item := range v.items {
		str += item.String()
		str += "\n"
	}
	return str
}

type bondItem struct {
	id    uint32
	btype uint32
	aId1  uint32
	aId2  uint32
}

func (i *bondItem) String() string {
	return fmt.Sprintf("%d %d %d %d", i.id, i.btype, i.aId1, i.aId2)
}

type bonds struct {
	note  string
	items []*bondItem
}

func (b *bonds) String() string {
	str := "Bonds\n"
	str += "\n"
	str += b.note
	str += "\n"
	for _, item := range b.items {
		str += item.String()
		str += "\n"
	}
	return str
}

type angleItem struct {
	id    uint32
	atype uint32
	aId1  uint32
	aId2  uint32
	aId3  uint32
}

func (i *angleItem) String() string {
	return fmt.Sprintf("%d %d %d %d %d", i.id, i.atype, i.aId1, i.aId2, i.aId3)
}

type angles struct {
	note  string
	items []*angleItem
}

func (a *angles) String() string {
	str := "Angles\n"
	str += "\n"
	str += a.note
	str += "\n"
	for _, item := range a.items {
		str += item.String()
		str += "\n"
	}
	return str
}

type doc struct {
	header     string
	atoms      *atoms
	velocities *velocities
	bonds      *bonds
	angles     *angles
}

func (d *doc) String() string {
	str := d.header
	str += d.atoms.String()
	str += d.velocities.String()
	str += d.bonds.String()
	str += d.angles.String()
	return str
}

func checkPara() {
	if inputfile == "" {
		log.Fatal("input file not specified")
	}
	if outputfile == "" {
		log.Fatal("output file not specified")
	}
}

func readHeader(r *bufio.Reader) string {
	str := ""
	for {
		l, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(l)
		if strings.Contains(line, "Atoms") {
			break
		}

		str += line
		str += "\n"
	}

	return str
}

func parseAtomItem(line string) *atomItem {
	tokens := strings.Split(line, " ")
	if len(tokens) != 10 {
		return nil
	}
	return &atomItem{
		id:    toUint32(tokens[0]),
		mId:   toUint32(tokens[1]),
		iType: toUint32(tokens[2]),
		q:     toFloat64(tokens[3]),
		x:     toFloat64(tokens[4]),
		y:     toFloat64(tokens[5]),
		z:     toFloat64(tokens[6]),
		m:     toInt32(tokens[7]),
		n:     toInt32(tokens[8]),
		k:     toInt32(tokens[9]),
	}
}

func readAtoms(r *bufio.Reader) *atoms {
	a := &atoms{}
	for {
		l, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(l)
		if line == "" {
			continue
		}
		if strings.Contains(line, "#") {
			a.note = line
			continue
		}
		if strings.Contains(line, "Velocities") {
			break
		}
		ai := parseAtomItem(line)
		if ai != nil {
			a.items = append(a.items, ai)
		}
	}

	return a
}

func parsevelItem(line string) *velItem {
	tokens := strings.Split(line, " ")
	if len(tokens) != 4 {
		return nil
	}
	return &velItem{
		aId: toUint32(tokens[0]),
		vx:  toFloat64(tokens[1]),
		vy:  toFloat64(tokens[2]),
		vz:  toFloat64(tokens[3]),
	}
}

func readVelocities(r *bufio.Reader) *velocities {
	v := &velocities{}
	for {
		l, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(l)
		if strings.Contains(line, "#") {
			v.note = line
			continue
		}
		if line == "" {
			continue
		}

		if strings.Contains(line, "Bonds") {
			break
		}
		vi := parsevelItem(line)
		if vi != nil {
			v.items = append(v.items, vi)
		}
	}

	return v
}

func parseBondItem(line string) *bondItem {
	tokens := strings.Split(line, " ")
	if len(tokens) != 4 {
		return nil
	}
	return &bondItem{
		id:    toUint32(tokens[0]),
		btype: toUint32(tokens[1]),
		aId1:  toUint32(tokens[2]),
		aId2:  toUint32(tokens[3]),
	}
}

func readBonds(r *bufio.Reader) *bonds {
	b := &bonds{}
	for {
		l, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(l)
		if strings.Contains(line, "#") {
			b.note = line
			continue
		}
		if line == "" {
			continue
		}

		if strings.Contains(line, "Angles") {
			break
		}
		bi := parseBondItem(line)
		if bi != nil {
			b.items = append(b.items, bi)
		}
	}

	return b
}

func parseAngleItem(line string) *angleItem {
	tokens := strings.Split(line, " ")
	if len(tokens) != 5 {
		return nil
	}
	return &angleItem{
		id:    toUint32(tokens[0]),
		atype: toUint32(tokens[1]),
		aId1:  toUint32(tokens[2]),
		aId2:  toUint32(tokens[3]),
		aId3:  toUint32(tokens[4]),
	}
}

func readAngles(r *bufio.Reader) *angles {
	a := &angles{}
	for {
		l, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(l)
		if strings.Contains(line, "#") {
			a.note = line
			continue
		}
		if line == "" {
			continue
		}
		ai := parseAngleItem(line)
		if ai != nil {
			a.items = append(a.items, ai)
		}
	}

	return a
}

func readDoc(path string) (*doc, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	d := &doc{}
	d.header = readHeader(br)
	d.atoms = readAtoms(br)
	d.velocities = readVelocities(br)
	d.bonds = readBonds(br)
	d.angles = readAngles(br)
	return d, nil
}

func writeDocToFile(doc *doc, path string) error {
	fi, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer fi.Close()

	_, err = fi.WriteString(doc.String())
	return err
}

func proc(d *doc) {
	atoms := d.atoms.items
	sort.Slice(atoms, func(i, j int) bool {
		return atoms[i].id < atoms[j].id
	})
	m := map[uint32]uint32{}
	for i := 0; i < len(atoms); i++ {
		m[atoms[i].id] = uint32(i + 1)
		atoms[i].id = uint32(i + 1)
	}
	vels := d.velocities.items
	sort.Slice(vels, func(i, j int) bool {
		return vels[i].aId < vels[j].aId
	})
	for i := 0; i < len(vels); i++ {
		vels[i].aId = uint32(i + 1)
	}

	bonds := d.bonds.items
	for _, bond := range bonds {
		bond.aId1 = m[bond.aId1]
		bond.aId2 = m[bond.aId2]
	}

	angles := d.angles.items
	for _, angle := range angles {
		angle.aId1 = m[angle.aId1]
		angle.aId2 = m[angle.aId2]
		angle.aId3 = m[angle.aId3]
	}
}

func main() {
	checkPara()
	d, err := readDoc(inputfile)
	if err != nil {
		log.Fatal(err)
	}
	proc(d)
	writeDocToFile(d, outputfile)
}

func toUint32(s string) uint32 {
	res, _ := strconv.ParseUint(s, 10, 32)
	return uint32(res)
}

func toInt32(s string) int32 {
	res, _ := strconv.ParseInt(s, 10, 32)
	return int32(res)
}

func toFloat64(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}
