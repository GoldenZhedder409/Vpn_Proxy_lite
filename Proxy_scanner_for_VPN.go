package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// ============ PROXY LIST DARI SCRIPT UTAMA ============

var knownProxies = []struct {
	addr     string
	country  string
	flag     string
	code     string
	protocol string
}{
	// --- USA ---
	{"45.79.203.254:48388", "USA", "🇺🇸", "USA", "SOCKS5"},
	{"104.219.236.127:1080", "USA", "🇺🇸", "USA", "SOCKS5"},
	{"165.22.110.253:1080", "USA", "🇺🇸", "USA", "SOCKS5"},
	{"192.241.230.75:1080", "USA", "🇺🇸", "USA", "SOCKS5"},

	// --- GERMANY ---
	{"185.133.239.244:16299", "Germany", "🇩🇪", "GER", "SOCKS5"},
	{"185.194.217.97:1080", "Germany", "🇩🇪", "GER", "SOCKS5"},
	{"84.200.125.162:1080", "Germany", "🇩🇪", "GER", "SOCKS5"},
	{"46.4.53.115:1080", "Germany", "🇩🇪", "GER", "SOCKS5"},

	// --- JAPAN ---
	{"20.210.113.32:8123", "Japan", "🇯🇵", "JPN", "HTTP"},
	{"89.116.88.19:80", "Japan", "🇯🇵", "JPN", "HTTP"},
	{"153.122.100.18:1080", "Japan", "🇯🇵", "JPN", "SOCKS5"},
	{"45.125.44.118:1080", "Japan", "🇯🇵", "JPN", "SOCKS5"},

	// --- BRAZIL ---
	{"186.26.95.249:61445", "Brazil", "🇧🇷", "BRA", "SOCKS5"},
	{"187.17.201.203:38737", "Brazil", "🇧🇷", "BRA", "SOCKS5"},
	{"177.136.124.47:56113", "Brazil", "🇧🇷", "BRA", "SOCKS5"},
	{"191.252.62.147:1080", "Brazil", "🇧🇷", "BRA", "SOCKS5"},

	// --- INDIA ---
	{"110.235.246.62:1080", "India", "🇮🇳", "IND", "SOCKS5"},
	{"64.227.131.240:1080", "India", "🇮🇳", "IND", "SOCKS5"},
	{"139.59.24.173:1080", "India", "🇮🇳", "IND", "SOCKS5"},
	{"103.149.162.194:1080", "India", "🇮🇳", "IND", "SOCKS5"},

	// --- SINGAPORE ---
	{"165.22.80.17:1080", "Singapore", "🇸🇬", "SGP", "SOCKS5"},
	{"167.172.112.65:1080", "Singapore", "🇸🇬", "SGP", "SOCKS5"},
	{"139.59.125.101:1080", "Singapore", "🇸🇬", "SGP", "SOCKS5"},

	// --- NETHERLANDS ---
	{"46.101.11.45:1080", "Netherlands", "🇳🇱", "NLD", "SOCKS5"},
	{"188.166.98.210:1080", "Netherlands", "🇳🇱", "NLD", "SOCKS5"},
	{"95.179.175.62:1080", "Netherlands", "🇳🇱", "NLD", "SOCKS5"},

	// --- CANADA ---
	{"167.71.205.251:1080", "Canada", "🇨🇦", "CAN", "SOCKS5"},
	{"159.89.192.73:1080", "Canada", "🇨🇦", "CAN", "SOCKS5"},
	{"138.197.199.102:1080", "Canada", "🇨🇦", "CAN", "SOCKS5"},
}

// ============ IP GEOLOCATION CACHE ============

type GeoInfo struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
}

var geoCache = make(map[string]*GeoInfo)

// ============ MAIN ============

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		clearScreen()
		printBanner()
		printMenu()
		
		fmt.Print("👉 Choose: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			analyzeSingleProxy(scanner)
		case "2":
			analyzeAllKnownProxies(scanner)
		case "3":
			analyzeCustomProxy(scanner)
		case "4":
			showStats()
			pause(scanner)
		case "5":
			fmt.Println("\n👋 Goodbye!")
			return
		default:
			fmt.Println("\n❌ Invalid option!")
			time.Sleep(1 * time.Second)
		}
	}
}

// ============ UI FUNCTIONS ============

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBanner() {
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║        🕵️  PROXY ANALYZER TOOL  🕵️               ║")
	fmt.Println("║    Check SOCKS5/SOCKS4/HTTP Proxies & Country    ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println()
}

func printMenu() {
	fmt.Println("  [1] 🔍 Analyze Single Proxy")
	fmt.Println("  [2] 📋 Analyze ALL Known Proxies (from script)")
	fmt.Println("  [3] ✏️  Analyze Custom Proxy (manual input)")
	fmt.Println("  [4] 📊 Statistics")
	fmt.Println("  [5] ❌ Exit")
	fmt.Println()
}

func printSeparator() {
	fmt.Println(strings.Repeat("─", 60))
}

func printSuccess(msg string) {
	fmt.Printf("✅ %s\n", msg)
}

func printError(msg string) {
	fmt.Printf("❌ %s\n", msg)
}

func printWarning(msg string) {
	fmt.Printf("⚠️  %s\n", msg)
}

func printInfo(msg string) {
	fmt.Printf("ℹ️  %s\n", msg)
}

func pause(scanner *bufio.Scanner) {
	fmt.Println("\nPress Enter to continue...")
	scanner.Scan()
}

// ============ PROXY ANALYSIS ============

func analyzeSingleProxy(scanner *bufio.Scanner) {
	clearScreen()
	printHeader("🔍 ANALYZE SINGLE PROXY")
	
	fmt.Print("Enter proxy (IP:PORT): ")
	scanner.Scan()
	proxyAddr := strings.TrimSpace(scanner.Text())
	
	if !isValidProxyFormat(proxyAddr) {
		printError("Invalid format! Use IP:PORT (e.g., 45.79.203.254:48388)")
		pause(scanner)
		return
	}
	
	fmt.Println()
	printInfo("Analyzing proxy...")
	fmt.Println()
	
	// Check if in known list
	known := findInKnownProxies(proxyAddr)
	if known != nil {
		printSuccess(fmt.Sprintf("✅ Found in known list: %s %s", known.flag, known.country))
		fmt.Printf("   Expected protocol: %s\n", known.protocol)
	} else {
		printWarning("⚠️  Not in known proxy list")
	}
	fmt.Println()
	
	// Check if alive
	alive, latency := isAlive(proxyAddr)
	if alive {
		printSuccess(fmt.Sprintf("✅ Proxy is ALIVE (latency: %v)", latency))
	} else {
		printError("❌ Proxy is DEAD (connection failed)")
		pause(scanner)
		return
	}
	fmt.Println()
	
	// Detect protocol
	protocol := detectProtocol(proxyAddr)
	if protocol != "UNKNOWN" {
		printSuccess(fmt.Sprintf("✅ Detected protocol: %s", protocol))
	} else {
		printError("❌ Could not detect protocol (not SOCKS5/SOCKS4/HTTP)")
	}
	fmt.Println()
	
	// Get country from IP
	country, flag := getCountryFromIP(proxyAddr)
	if country != "Unknown" {
		printSuccess(fmt.Sprintf("🌍 IP Location: %s %s", flag, country))
		
		// Check if matches known
		if known != nil && known.country != country {
			printWarning(fmt.Sprintf("⚠️  Country MISMATCH! Known: %s %s, Actual: %s %s", 
				known.flag, known.country, flag, country))
		}
	} else {
		printError("❌ Could not detect country (API limit or error)")
	}
	
	printSeparator()
	pause(scanner)
}

func analyzeAllKnownProxies(scanner *bufio.Scanner) {
	clearScreen()
	printHeader("📋 ANALYZE ALL KNOWN PROXIES")
	
	fmt.Printf("Total proxies in list: %d\n", len(knownProxies))
	fmt.Println()
	printInfo("Testing all proxies (this may take a minute)...")
	fmt.Println()
	
	results := struct {
		alive     int
		socks5    int
		socks4    int
		http      int
		unknown   int
		mismatch  int
		total     int
	}{}
	
	for i, p := range knownProxies {
		fmt.Printf("  [%d/%d] Testing %s %s %s... ", 
			i+1, len(knownProxies), p.flag, p.country, p.addr)
		
		// Check alive
		alive, latency := isAlive(p.addr)
		if !alive {
			fmt.Printf("❌ DEAD\n")
			continue
		}
		
		results.alive++
		
		// Detect protocol
		protocol := detectProtocol(p.addr)
		switch protocol {
		case "SOCKS5":
			results.socks5++
		case "SOCKS4":
			results.socks4++
		case "HTTP":
			results.http++
		default:
			results.unknown++
		}
		
		// Get country
		country, flag := getCountryFromIP(p.addr)
		if country != "Unknown" && country != p.country {
			results.mismatch++
			fmt.Printf("🌍 %s %s | ⏱️  %v | %s | COUNTRY MISMATCH!\n", 
				flag, country, latency.Round(time.Millisecond), protocol)
		} else {
			fmt.Printf("✅ %s %s | ⏱️  %v | %s\n", 
				flag, country, latency.Round(time.Millisecond), protocol)
		}
		
		results.total++
		
		// Small delay to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}
	
	fmt.Println()
	printSeparator()
	fmt.Println("📊 SUMMARY:")
	fmt.Printf("  Total proxies: %d\n", len(knownProxies))
	fmt.Printf("  Alive: %d\n", results.alive)
	fmt.Printf("  SOCKS5: %d\n", results.socks5)
	fmt.Printf("  SOCKS4: %d\n", results.socks4)
	fmt.Printf("  HTTP: %d\n", results.http)
	fmt.Printf("  Unknown protocol: %d\n", results.unknown)
	fmt.Printf("  Country mismatch: %d\n", results.mismatch)
	printSeparator()
	
	pause(scanner)
}

func analyzeCustomProxy(scanner *bufio.Scanner) {
	clearScreen()
	printHeader("✏️  ANALYZE CUSTOM PROXY")
	
	fmt.Print("Enter proxy (IP:PORT): ")
	scanner.Scan()
	proxyAddr := strings.TrimSpace(scanner.Text())
	
	if !isValidProxyFormat(proxyAddr) {
		printError("Invalid format! Use IP:PORT")
		pause(scanner)
		return
	}
	
	fmt.Println()
	printInfo("Analyzing custom proxy...")
	fmt.Println()
	
	// Check if alive
	alive, latency := isAlive(proxyAddr)
	if alive {
		printSuccess(fmt.Sprintf("✅ Proxy is ALIVE (latency: %v)", latency))
	} else {
		printError("❌ Proxy is DEAD")
		pause(scanner)
		return
	}
	fmt.Println()
	
	// Detect protocol
	protocol := detectProtocol(proxyAddr)
	if protocol != "UNKNOWN" {
		printSuccess(fmt.Sprintf("✅ Detected protocol: %s", protocol))
	} else {
		printError("❌ Could not detect protocol")
	}
	fmt.Println()
	
	// Get country
	country, flag := getCountryFromIP(proxyAddr)
	if country != "Unknown" {
		printSuccess(fmt.Sprintf("🌍 IP Location: %s %s", flag, country))
	} else {
		printError("❌ Could not detect country")
	}
	
	// Check if in known list
	known := findInKnownProxies(proxyAddr)
	if known != nil {
		fmt.Println()
		printWarning("⚠️  This proxy is IN THE KNOWN LIST!")
		fmt.Printf("   Known as: %s %s (%s)\n", known.flag, known.country, known.protocol)
	}
	
	printSeparator()
	pause(scanner)
}

func showStats() {
	clearScreen()
	printHeader("📊 PROXY STATISTICS")
	
	// Group by country
	countryCount := make(map[string]int)
	protocolCount := make(map[string]int)
	
	for _, p := range knownProxies {
		countryCount[p.country]++
		protocolCount[p.protocol]++
	}
	
	fmt.Println("  By Country:")
	for country, count := range countryCount {
		flag := ""
		for _, p := range knownProxies {
			if p.country == country {
				flag = p.flag
				break
			}
		}
		fmt.Printf("    %s %s: %d proxies\n", flag, country, count)
	}
	
	fmt.Println()
	fmt.Println("  By Protocol:")
	for protocol, count := range protocolCount {
		fmt.Printf("    %s: %d\n", protocol, count)
	}
	
	fmt.Println()
	fmt.Printf("  Total proxies: %d\n", len(knownProxies))
	
	printSeparator()
}

// ============ HELPER FUNCTIONS ============

func printHeader(title string) {
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Printf("║  %-46s ║\n", title)
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println()
}

func isValidProxyFormat(addr string) bool {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return false
	}
	ip := net.ParseIP(parts[0])
	return ip != nil
}

func findInKnownProxies(addr string) *struct {
	addr     string
	country  string
	flag     string
	code     string
	protocol string
} {
	for _, p := range knownProxies {
		if p.addr == addr {
			return &p
		}
	}
	return nil
}

func isAlive(addr string) (bool, time.Duration) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	latency := time.Since(start)
	
	if err != nil {
		return false, 0
	}
	conn.Close()
	return true, latency
}

func detectProtocol(addr string) string {
	// Test SOCKS5
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return "UNKNOWN"
	}
	
	// SOCKS5 test
	conn.Write([]byte{0x05, 0x01, 0x00})
	buf := make([]byte, 2)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _ := conn.Read(buf)
	if n == 2 && buf[0] == 0x05 {
		conn.Close()
		return "SOCKS5"
	}
	conn.Close()
	
	// SOCKS4 test
	conn, err = net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return "UNKNOWN"
	}
	conn.Write([]byte{0x04, 0x01, 0x00, 0x50, 0x00, 0x00, 0x00, 0x01, 0x00})
	buf = make([]byte, 8)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _ = conn.Read(buf)
	if n >= 8 && buf[0] == 0x00 && buf[1] == 0x5A {
		conn.Close()
		return "SOCKS4"
	}
	conn.Close()
	
	// HTTP test
	conn, err = net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return "UNKNOWN"
	}
	conn.Write([]byte("CONNECT www.google.com:80 HTTP/1.0\r\n\r\n"))
	buf = make([]byte, 13)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _ = conn.Read(buf)
	if n >= 13 && strings.Contains(string(buf[:13]), "200") {
		conn.Close()
		return "HTTP"
	}
	conn.Close()
	
	return "UNKNOWN"
}

func getCountryFromIP(addr string) (string, string) {
	ip := strings.Split(addr, ":")[0]
	
	// Check cache
	if cached, ok := geoCache[ip]; ok {
		return cached.Country, getFlag(cached.CountryCode)
	}
	
	// Try free IP API
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip + "?fields=country,countryCode")
	if err != nil {
		return "Unknown", "🌍"
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	
	var geo GeoInfo
	err = json.Unmarshal(body, &geo)
	if err != nil || geo.Country == "" {
		return "Unknown", "🌍"
	}
	
	// Save to cache
	geoCache[ip] = &geo
	
	return geo.Country, getFlag(geo.CountryCode)
}

func getFlag(countryCode string) string {
	// Convert country code to flag emoji
	if len(countryCode) != 2 {
		return "🌍"
	}
	
	// Regional indicator symbols: A=0x1F1E6, B=0x1F1E7, etc.
	first := 0x1F1E6 + rune(countryCode[0]-'A')
	second := 0x1F1E6 + rune(countryCode[1]-'A')
	
	return string(first) + string(second)
}
