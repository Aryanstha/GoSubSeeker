package lib

import (
	"bufio"
	"context"
	"fmt"
	"github.com/aryanstha/GoSubSeeker/lib/dns"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Host string
	Addr []net.IP
}

type Scanner struct {
	opts        *Options    // Options
	preWordChan chan string //
	wordChan    chan string
	found       int
	count       int
	issued      int
	log         *os.File // log
	mu          *sync.RWMutex
	timeStart   time.Time
	BlackList   map[string]string
}

func NewScanner(opts *Options) *Scanner {
	var this Scanner
	this.opts = opts

	this.wordChan = make(chan string, this.opts.Threads)
	this.preWordChan = make(chan string)
	this.mu = new(sync.RWMutex)
	this.BlackList = make(map[string]string)
	//this.LoadBlackListFile()

	f, err := os.Create(opts.Log)
	this.log = f

	if err != nil {
		log.Fatalln(err)
	}
	return &this
}

func (this *Scanner) WildcardsDomain(s string) bool {
	//log.Printf("[+] Validate wildcard domain *.%v exists", this.opts.Domain)
	if ip, ok := this.IsWildcardsDomain(s); ok {
		//log.Printf("[+] Domain %v is wildcard,*.%v ip is %v", this.opts.Domain, this.opts.Domain, ip)
		//if ! this.opts.WildcardDomain {
		//	return true
		//}

		for _, v := range ip {
			this.BlackList[v.String()] = fmt.Sprintf("*.%s", this.opts.Domain)
		}

		return true
	}

	return false
}

func (this *Scanner) Start() {
	if ok := this.WildcardsDomain(this.opts.Domain); ok && !this.opts.WildcardDomain {
		return
	}

	this.timeStart = time.Now()
	var wg sync.WaitGroup
	//wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go this.progressPrint(ctx, &wg)
	go this.preWord(&wg)
	for i := 0; i < this.opts.Threads; i++ {
		go this.worker(&wg)
	}

	// 读取字典
	f, err := os.Open(this.opts.Dict)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	this.count = getCountLine(f)

	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		wg.Add(1)
		word := strings.TrimSpace(scanner.Text())
		this.wordChan <- word
	}

	wg.Wait()

	format := "All Done. %d found, %.4f/s, %d scanned in %.2f seconds\n"
	this.mu.RLock()
	this.progressClean()
	log.Printf(format,
		this.found,
		float64(this.issued)/time.Since(this.timeStart).Seconds(),
		this.issued,
		time.Since(this.timeStart).Seconds(),
	)
	log.Printf("The output result file is %s\n", this.opts.Log)
	this.mu.RUnlock()

}

func (this *Scanner) incr() {
	this.mu.Lock()
	this.issued++
	this.mu.Unlock()
}

// goroutine >= 1

func (this *Scanner) worker(wg *sync.WaitGroup) {
	for v := range this.wordChan {
		this.incr()
		host := fmt.Sprintf("%s.%s", v, this.opts.Domain)
		ip, err := this.LookupHost(host)
		if err == nil {
			this.result(Result{host, ip}, wg)
		}

		wg.Done()
	}
}
func (this *Scanner) preWord(wg *sync.WaitGroup) {
	// TODO
	f, err := os.Open("dict/next_sub.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	// TODO 深度扫描暂时不支持泛域名扫描
	for v := range this.preWordChan {
		if ok := this.WildcardsDomain(v); !ok {
			f.Seek(0, 0)
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				this.mu.Lock()
				this.count++
				this.mu.Unlock()

				wg.Add(1)
				word := strings.TrimSpace(scanner.Text())
				word = fmt.Sprintf("%s.%s", word, v[0:strings.LastIndex(v, this.opts.Domain)-1])
				this.wordChan <- word
			}
		}

		wg.Done()
	}
}

func (this *Scanner) result(re Result, wg *sync.WaitGroup) {
	// 如果没有一个可用ip存在,则不记录
	if this.IsBlackIPs(re.Addr) {
		return
	}

	if this.opts.Depth > 1 {
		wg.Add(1)
		go func() {
			this.preWordChan <- re.Host
		}()
	}

	this.progressClean()

	fmt.Printf("[+] %v\n", re)

	this.mu.Lock()
	this.found++
	this.log.WriteString(fmt.Sprintf("%v\t%v\n", re.Host, re.Addr))
	this.mu.Unlock()
}

// 清除光标所在行的所有字符
func (this *Scanner) progressClean() {
	fmt.Fprint(os.Stderr, "\r\x1b[2K")
}

// goroutine = 1
// 启动后该方法负责打印进度
// 直到进度到100%跳出死循环
func (this *Scanner) progressPrint(c context.Context, wg *sync.WaitGroup) {
	tick := time.NewTicker(1 * time.Second)
	format := "\r%d|%.4f%%|%.4f/s|%d scanned in %.2f seconds"
	log.Println("Starting")

	for {
		select {
		case <-tick.C:
			this.mu.RLock()
			fmt.Fprintf(os.Stderr, format,
				this.count,
				float64(this.issued)/float64(this.count)*100,
				float64(this.issued)/time.Since(this.timeStart).Seconds(),
				this.issued,
				time.Since(this.timeStart).Seconds(),
			)
			this.mu.RUnlock()
			// Force quit
			//if this.issued == this.count {
			//	break Loop;
			//}

		case <-c.Done():
			return
		}
	}
}

// 获取泛域名ip地址
func (this *Scanner) IsWildcardsDomain(s string) (ip []net.IP, ok bool) {
	// Go package net exists bug?
	// @link https://github.com/golang/go/issues/28947
	// Nonsupport RFC 4592
	// net.LookupHost("*.qzone.qq.com") //  --> lookup *.qzone.qq.com: no such host

	// md5(random string)
	// byte := md5.Sum([]byte(time.Now().String()))
	// randSub:=hex.EncodeToString(byte[:])
	// host := fmt.Sprintf("%s.%s", randSub, this.opts.Domain)
	// addrs, err := net.LookupHost(host)

	host := fmt.Sprintf("*.%s", s)
	addrs, err := this.LookupHost(host)

	if err != nil {
		return addrs, false
	}

	return addrs, true
}

// 验证DNS域传送
func (this *Scanner) TestAXFR(domain string) (results []string, err error) {
	server, err := this.LookupNS(domain)

	if results, err = dns.Axrf(domain, server); err == nil {
		for _, v := range results {
			this.mu.Lock()
			this.log.WriteString(fmt.Sprintf("%s\n", v))
			this.mu.Unlock()
		}
	}
	return
}

// 验证DNS服务器是否稳定
func (this *Scanner) TestDNSServer() bool {
	ipaddr, err := this.LookupHost("google-public-dns-a.google.com") // test lookup an existed domain

	if err != nil {
		//log.Println(err)
		return false
	}
	// Validate dns pollution
	if ipaddr[0].String() != "8.8.8.8" {
		// Non-existed domain test
		_, err := this.LookupHost("test.bad.dns.fengdingbo.com")
		// Bad DNS Server
		if err == nil {
			return false
		}
	}

	return true
}

func getCountLine(f *os.File) int {
	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i++
	}

	return i
}
