package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

type net struct {
	n, k            int
	host            []host
	partCount       []uint8
	partOrder       []uint8
	partHosts       [][]uint8
	partSrc         []uint8 // хост для запроса каждой части
	totalPartsCount int
	roundCount      int
}

func (nt *net) init(n, k int) {

	if debugEnable {
		log.Println("=== init")
	}

	nt.n = n
	nt.k = k

	// создаем счетчики для каждой части
	nt.partCount = make([]uint8, k)

	// cоздаем массив для ранжирования частей
	nt.partOrder = make([]uint8, k)
	for partID := range nt.partOrder {
		nt.partOrder[partID] = uint8(partID)
	}

	// nt.partHosts = make([][]uint8, k)
	// for partID := range nt.partHosts {
	// 	nt.partHosts[partID] = make([]uint8, 0, n)
	// }

	{
		buf := make([]uint8, n*k)
		nt.partHosts = make([][]uint8, k)
		for i, j := 0, 0; i < len(nt.partHosts); i, j = i+1, j+n {
			nt.partHosts[i] = buf[j:j]
		}
	}

	// создаем массив под источники для каждой части
	nt.partSrc = make([]uint8, k)

	// создаем хосты
	nt.host = make([]host, n)
	for hostID := range nt.host {
		nt.host[hostID].init(uint8(hostID), nt)
	}

	// заполняем первый хост всеми частями
	nt.fillFirstHost()
}

func (nt *net) fillFirstHost() {
	h := &nt.host[0]
	for partID := 0; partID < nt.k; partID++ {
		h.receivePartFrom(0, uint8(partID))
	}
}

func (nt *net) doUpdate() {

	full := nt.n * nt.k
	if debugEnable {
		log.Printf("N.totalPartCount: %d/%d", nt.totalPartsCount, full)
	}

	for int(nt.totalPartsCount) < full {
		nt.doRound()

		if debugEnable {
			log.Printf("N.totalPartCount: %d/%d", nt.totalPartsCount, full)
		}
	}
}

func (nt *net) doRound() {

	nt.roundCount++
	if debugEnable {
		log.Println("=== round:", nt.roundCount)
	}

	// Каждое устройство выбирает отсутствующую на нем часть обновления, которая встречается в сети реже всего.
	// Если таких частей несколько, то выбирается отсутствующая на устройстве часть обновления с наименьшим номером.

	// Предварительно ранжируем части по размеру и id. Это должно уменьшить константу когда хост выберает часть.
	// В обычном случае ему нужно пройти весь список. Вслучае если части отсортированы, хост выбирает первую из
	// списка, которой у него нет.
	sort.Slice(nt.partOrder, func(i, j int) bool {
		i = int(nt.partOrder[i])
		j = int(nt.partOrder[j])
		return nt.partCount[i] < nt.partCount[j] || nt.partCount[i] == nt.partCount[j] && i < j
	})

	// После этого устройство делает запрос выбранной части обновления у одного из устройств, на котором такая часть
	// обновления уже скачана. Если таких устройств несколько — выбирается устройство, на котором скачано наименьшее
	// количество частей обновления. Если и таких устройств оказалось несколько — выбирается устройство с минимальным номером.

	for partID := 0; partID < nt.k; partID++ {
		hostID := 0     // на хосте 0 всегда
		minimum := nt.k // есть все части

		for _, id := range nt.partHosts[partID] {
			h := &nt.host[id]
			if int(h.partsCount) < minimum || int(h.partsCount) == minimum && int(id) < hostID {
				hostID = int(id)
				minimum = int(h.partsCount)
			}
		}

		nt.partSrc[partID] = uint8(hostID)
	}

	// хосты делают запросы
	for id := range nt.host {
		nt.host[id].doRequest()
	}

	// хосты делают выбор
	for id := range nt.host {
		nt.host[id].doChoice()
	}

	// хосты делают посылки
	for id := range nt.host {
		nt.host[id].sendPart()
	}
}

func (nt *net) hostReceivedPart(hostID, partID uint8) {

	if debugEnable {
		log.Println("N.hostReceivedPart:", hostID, partID)
	}

	nt.partCount[partID]++
	nt.totalPartsCount++
	nt.partHosts[partID] = append(nt.partHosts[partID], uint8(hostID))
}

type request struct {
	hostID uint8
	partID uint8
}

type host struct {
	net        *net
	partExists []bool
	request    []request
	peerValue  []uint8
	// поля переупорядочены, чтобы оптимизировать размер памяти занимаемой структурой
	id         uint8
	partsCount uint8
	finish     uint16
	choice     request
}

func (h *host) init(id uint8, net *net) {
	h.id = id
	h.net = net
	h.partExists = make([]bool, net.k)
	h.request = make([]request, 0, net.n)
	h.peerValue = make([]uint8, net.n)
}

func (h *host) doRequest() {

	if h.finish != 0 {
		return
	}

	// Каждое устройство выбирает отсутствующую на нем часть обновления, которая встречается в сети реже всего.
	// Если таких частей несколько, то выбирается отсутствующая на устройстве часть обновления с наименьшим номером.

	// Выбираем первую отсутствующую часть из предварительно ранжированого списка
	for _, partID := range h.net.partOrder {
		if !h.partExists[partID] {
			id := h.net.partSrc[partID]
			h.net.host[id].receiveRequestFrom(h.id, partID)
			return
		}
	}
}

func (h *host) receiveRequestFrom(fromID, partID uint8) {

	if debugEnable {
		log.Printf("%d.receiveRequestFrom: %d %d", h.id, fromID, partID)
	}

	h.request = append(h.request, request{fromID, partID})
}

func (h *host) doChoice() {
	
	if len(h.request) == 0 {
		return
	}

	// Устройство A удовлетворяет тот запрос, который поступил от наиболее ценного для A устройства.
	// Ценность устройства B для устройства A определяется как количество частей обновления, ранее
	// полученных устройством A от устройства B. Если на устройство A пришло несколько запросов от
	// одинаково ценных устройств, то удовлетворяется запрос того устройства, на котором меньше всего
	// скачанных частей обновления. Если и таких запросов несколько, то среди них выбирается устройство
	// с наименьшим номером.

	br := h.request[0]
	bh := &h.net.host[br.hostID]

	for _, cr := range h.request {
		ch := &h.net.host[cr.hostID]

		if h.peerValue[ch.id] > h.peerValue[bh.id] ||
			h.peerValue[ch.id] == h.peerValue[bh.id] && (ch.partsCount < bh.partsCount ||
				bh.partsCount == ch.partsCount && ch.id < bh.id) {
			br = cr
			bh = ch
		}
	}

	// NOTE: здесь не очищаем список запросов, т.к. непустой список будем использовать, как признак
	//  необходимости сделать посылку. Список запросов будет очищен в sendPart
	h.choice = br
}

func (h *host) sendPart() {

	if len(h.request) == 0 {
		return
	}
	h.request = h.request[:0]

	peer := &h.net.host[h.choice.hostID]
	peer.receivePartFrom(h.id, h.choice.partID)
}

func (h *host) receivePartFrom(fromID, partID uint8) {

	if debugEnable {
		log.Printf("%d.receivePartFrom: %d %d", h.id, fromID, partID)
	}

	h.peerValue[fromID]++
	h.partExists[partID] = true

	h.partsCount++
	if debugEnable {
		log.Printf("%d.partCount: %d/%d", h.id, h.partsCount, h.net.k)
	}

	if int(h.partsCount) == h.net.k {
		h.finish = uint16(h.net.roundCount)
	}

	// сообщаем, что получили часть
	h.net.hostReceivedPart(h.id, partID)
}

func solve(n, k int) []int {
	var nt net
	nt.init(n, k)
	nt.doUpdate()

	res := make([]int, 0, n)
	for id := 1; id < len(nt.host); id++ {
		res = append(res, int(nt.host[id].finish))
	}

	return res
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, k, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	res := solve(n, k)

	writeInts(bw, res, " ")
	bw.WriteByte('\n')

	return nil
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(unsafeString(sc.Bytes()))
}

func scanTwoInt(sc *bufio.Scanner) (v1, v2 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	return v1, v2, err
}

func scanThreeInt(sc *bufio.Scanner) (v1, v2, v3 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	return v1, v2, v3, err
}

func scanInts(sc *bufio.Scanner, a []int) error {
	for i := range a {
		v, err := scanInt(sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

func writeInt[I Int](bw *bufio.Writer, v I) error {
	var buf [32]byte
	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))
	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, sep string) error {
	if len(a) == 0 {
		return nil
	}

	var buf [32]byte

	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))
	for i := 1; err == nil && i < len(a); i++ {
		_, err = bw.WriteString(sep)
		if err == nil {
			_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
		}
	}

	return err
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
