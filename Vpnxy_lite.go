package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "math/rand"
        "net"
        "strings"
        "sync"
        "time"
)

// ============ PROXY LIST ============

var proxies = []struct {
        addr     string
        country  string
        flag     string
        code     string
        timezone string
        myth     string
        protocol string
}{
        // --- USA ---
        {"45.79.203.254:48388", "USA", "🇺🇸", "USA", "America/New_York", "CIA Surveillance Van 🚐", "SOCKS5"},
        {"104.219.236.127:1080", "USA", "🇺🇸", "USA", "America/Chicago", "Area 51 Secret Tunnel 🛸", "SOCKS5"},
        {"165.22.110.253:1080", "USA", "🇺🇸", "USA", "America/Los_Angeles", "Hollywood Green Screen 🎬", "SOCKS5"},
        {"192.241.230.75:1080", "USA", "🇺🇸", "USA", "America/New_York", "NYC Subway Node 🚇", "SOCKS5"},

        // --- GERMANY ---
        {"185.133.239.244:16299", "Germany", "🇩🇪", "GER", "Europe/Berlin", "Bratwurst Security 🌭", "SOCKS5"},
        {"185.194.217.97:1080", "Germany", "🇩🇪", "GER", "Europe/Berlin", "Octobersfest Hidden Lager 🍺", "SOCKS5"},
        {"84.200.125.162:1080", "Germany", "🇩🇪", "GER", "Europe/Berlin", "AutoBahn High Speed 🏎️", "SOCKS5"},
        {"46.4.53.115:1080", "Germany", "🇩🇪", "GER", "Europe/Berlin", "Black Forest Gateway 🌲", "SOCKS5"},

        // --- JAPAN ---
        {"20.210.113.32:8123", "Japan", "🇯🇵", "JPN", "Asia/Tokyo", "Akihabara Glitch 🤖", "HTTP"},
        {"89.116.88.19:80", "Japan", "🇯🇵", "JPN", "Asia/Tokyo", "Shibuya Crossing Ghost 👻", "HTTP"},
        {"153.122.100.18:1080", "Japan", "🇯🇵", "JPN", "Asia/Tokyo", "Mount Fuji Uplink 🏔️", "SOCKS5"},
        {"45.125.44.118:1080", "Japan", "🇯🇵", "JPN", "Asia/Tokyo", "Bullet Train Tunnel 🚄", "SOCKS5"},

        // --- BRAZIL ---
        {"186.26.95.249:61445", "Brazil", "🇧🇷", "BRA", "America/Sao_Paulo", "Amazon Rain Forest Wi-Fi 🌳", "SOCKS5"},
        {"187.17.201.203:38737", "Brazil", "🇧🇷", "BRA", "America/Sao_Paulo", "Maracana Stadium Node ⚽", "SOCKS5"},
        {"177.136.124.47:56113", "Brazil", "🇧🇷", "BRA", "America/Bahia", "Rio Carnival Mask 🎭", "SOCKS5"},
        {"191.252.62.147:1080", "Brazil", "🇧🇷", "BRA", "America/Sao_Paulo", "Samba Beat Rhythm 🥁", "SOCKS5"},

        // --- INDIA ---
        {"110.235.246.62:1080", "India", "🇮🇳", "IND", "Asia/Kolkata", "Taj Mahal Mirror 🕌", "SOCKS5"},
        {"64.227.131.240:1080", "India", "🇮🇳", "IND", "Asia/Kolkata", "Bangalore Tech Spirit 🧘", "SOCKS5"},
        {"139.59.24.173:1080", "India", "🇮🇳", "IND", "Asia/Kolkata", "Curry Powered Server 🍛", "SOCKS5"},
        {"103.149.162.194:1080", "India", "🇮🇳", "IND", "Asia/Kolkata", "Bollywood Dance Number 💃", "SOCKS5"},

        // --- SINGAPORE ---
        {"165.22.80.17:1080", "Singapore", "🇸🇬", "SGP", "Asia/Singapore", "Merlion Water Cannon 🌊", "SOCKS5"},
        {"167.172.112.65:1080", "Singapore", "🇸🇬", "SGP", "Asia/Singapore", "Marina Bay Sands Node 🏨", "SOCKS5"},
        {"139.59.125.101:1080", "Singapore", "🇸🇬", "SGP", "Asia/Singapore", "Satay by the Bay BBQ 🍢", "SOCKS5"},

        // --- NETHERLANDS ---
        {"46.101.11.45:1080", "Netherlands", "🇳🇱", "NLD", "Europe/Amsterdam", "Tulip Field Server 🌷", "SOCKS5"},
        {"188.166.98.210:1080", "Netherlands", "🇳🇱", "NLD", "Europe/Amsterdam", "Canal Boat Connection 🛶", "SOCKS5"},
        {"95.179.175.62:1080", "Netherlands", "🇳🇱", "NLD", "Europe/Amsterdam", "Windmill Power ⚡", "SOCKS5"},

        // --- CANADA ---
        {"167.71.205.251:1080", "Canada", "🇨🇦", "CAN", "America/Toronto", "Maple Syrup Router 🍁", "SOCKS5"},
        {"159.89.192.73:1080", "Canada", "🇨🇦", "CAN", "America/Montreal", "Poutine Protocol 🍟", "SOCKS5"},
        {"138.197.199.102:1080", "Canada", "🇨🇦", "CAN", "America/Vancouver", "Whistler Ski Lift ⛷️", "SOCKS5"},
}

// ============ CONFIG ============

type Config struct {
        SelectedCountries []string `json:"selected_countries"`
        RotateInterval    int      `json:"rotate_interval"`
        ProxyPort         int      `json:"proxy_port"`
}

var config Config
var configFile = "proxy-config.json"

// ============ GLOBAL STATE ============

var activeProxies []int
var proxyHealth = make(map[int]bool)
var proxyLatency = make(map[int]time.Duration)
var healthMu sync.RWMutex
var currentIndex int
var mu sync.Mutex
var patternCount int
var history []string
var usedPatterns []int
var lastSuccessIndex int = -1
var startTime time.Time
var running = false

// ============ MYSTICAL PATTERNS ============

var mysticalPatterns = []struct {
        name   string
        symbol string
}{
        {"Bermuda Triangle", "🔺"},
        {"Proxy Pentagram", "⭐"},
        {"Illuminati Butterfly", "🦋"},
        {"Ouroboros Snake", "🐍"},
        {"Ninja Run", "�"},
        {"Area 51 Raid", "🛸"},
        {"Mount Fuji", "🗻"},
        {"Oktoberfest", "🍺"},
        {"Carnival Chaos", "🎭"},
        {"Bollywood Dance", "💃"},
        {"Merlion Splash", "🌊"},
        {"Tulip Mania", "🌷"},
        {"Maple Syrup", "🍁"},
}

// ============ MAIN ============

func main() {
        rand.Seed(time.Now().UnixNano())
        clearScreen()
        loadConfig()
        showMainMenu()
}

// ============ UI ============

func clearScreen() {
        fmt.Print("\033[H\033[2J")
}

func printHeader(title string) {
        fmt.Println("╔════════════════════════════════════════════════════╗")
        fmt.Printf("║  %-46s ║\n", title)
        fmt.Println("╚════════════════════════════════════════════════════╝")
        fmt.Println()
}

func printSuccess(msg string) {
        fmt.Printf("✅ %s\n", msg)
}

func printError(msg string) {
        fmt.Printf("❌ %s\n", msg)
}

func printInfo(msg string) {
        fmt.Printf("ℹ️  %s\n", msg)
}

// ============ CONFIG ============

func loadConfig() {
        config = Config{
                SelectedCountries: []string{"USA", "Germany", "Japan", "Brazil", "India"},
                RotateInterval:    15,
                ProxyPort:         1080,
        }

        data, err := ioutil.ReadFile(configFile)
        if err == nil {
                json.Unmarshal(data, &config)
        }
}

func saveConfig() {
        data, _ := json.MarshalIndent(config, "", "  ")
        ioutil.WriteFile(configFile, data, 0644)
}

// ============ MENU ============

func showMainMenu() {
        for {
                clearScreen()
                printHeader("🔥 PROXY ROTATOR 🔥")

                fmt.Println("  [1] 🚀 START ROTATOR")
                fmt.Println("  [2] ⚙️  CONFIG")
                fmt.Println("  [3] 🌍 SELECT COUNTRIES")
                fmt.Println("  [4] 📊 STATUS")
                fmt.Println("  [5] ❌ EXIT")
                fmt.Println()
                fmt.Print("👉 Choose: ")

                var choice int
                fmt.Scanln(&choice)

                switch choice {
                case 1:
                        startRotator()
                case 2:
                        showConfigMenu()
                case 3:
                        showCountrySelector()
                case 4:
                        showStatus()
                case 5:
                        fmt.Println("\n👋 Goodbye!")
                        return
                default:
                        printError("Invalid option!")
                        time.Sleep(1 * time.Second)
                }
        }
}

func showConfigMenu() {
        clearScreen()
        printHeader("⚙️ CONFIG")

        fmt.Printf("  Interval: %ds\n", config.RotateInterval)
        fmt.Printf("  Port: %d\n", config.ProxyPort)
        fmt.Println()
        fmt.Println("  [1] Change interval")
        fmt.Println("  [2] Change port")
        fmt.Println("  [3] Back")
        fmt.Println()
        fmt.Print("👉 Choose: ")

        var choice int
        fmt.Scanln(&choice)

        switch choice {
        case 1:
                fmt.Print("Enter interval (5-60): ")
                var val int
                fmt.Scanln(&val)
                if val >= 5 && val <= 60 {
                        config.RotateInterval = val
                        saveConfig()
                        printSuccess(fmt.Sprintf("Interval = %d", val))
                }
        case 2:
                fmt.Print("Enter port (1024-65535): ")
                var val int
                fmt.Scanln(&val)
                if val >= 1024 && val <= 65535 {
                        config.ProxyPort = val
                        saveConfig()
                        printSuccess(fmt.Sprintf("Port = %d", val))
                }
        case 3:
                return
        }

        time.Sleep(2 * time.Second)
}

func showCountrySelector() {
        clearScreen()
        printHeader("🌍 SELECT COUNTRIES")

        countryMap := make(map[string]string)
        countryFlags := make(map[string]string)
        for _, p := range proxies {
                countryMap[p.country] = p.code
                countryFlags[p.country] = p.flag
        }

        var countries []string
        for country := range countryMap {
                countries = append(countries, country)
        }

        for i, country := range countries {
                selected := " "
                for _, c := range config.SelectedCountries {
                        if c == country {
                                selected = "✓"
                                break
                        }
                }
                fmt.Printf("  [%d] %s %-12s %s\n", i+1, countryFlags[country], country, selected)
        }

        fmt.Println()
        fmt.Println("  [A] All  [N] None  [D] Default  [B] Back")
        fmt.Print("👉 Choose: ")

        var input string
        fmt.Scanln(&input)

        switch strings.ToUpper(input) {
        case "A":
                config.SelectedCountries = countries
                printSuccess("All countries selected!")
        case "N":
                config.SelectedCountries = []string{}
                printSuccess("No countries selected!")
        case "D":
                config.SelectedCountries = []string{"USA", "Germany", "Japan", "Brazil", "India"}
                printSuccess("Default countries selected!")
        case "B":
                return
        default:
                var idx int
                fmt.Sscan(input, &idx)
                if idx >= 1 && idx <= len(countries) {
                        country := countries[idx-1]
                        found := -1
                        for i, c := range config.SelectedCountries {
                                if c == country {
                                        found = i
                                        break
                                }
                        }
                        if found >= 0 {
                                config.SelectedCountries = append(config.SelectedCountries[:found], config.SelectedCountries[found+1:]...)
                        } else {
                                config.SelectedCountries = append(config.SelectedCountries, country)
                        }
                }
        }

        saveConfig()
        time.Sleep(1 * time.Second)
        showCountrySelector()
}

func showStatus() {
        clearScreen()
        printHeader("📊 STATUS")

        countryCount := make(map[string]int)
        for _, p := range proxies {
                countryCount[p.country]++
        }

        fmt.Println("  Available proxies:")
        for country, count := range countryCount {
                flag := ""
                for _, p := range proxies {
                        if p.country == country {
                                flag = p.flag
                                break
                        }
                }
                fmt.Printf("    %s %s: %d\n", flag, country, count)
        }

        fmt.Println()
        fmt.Printf("  Total: %d proxies\n", len(proxies))

        selectedCount := 0
        for _, p := range proxies {
                for _, c := range config.SelectedCountries {
                        if p.country == c {
                                selectedCount++
                                break
                        }
                }
        }
        fmt.Printf("  Selected: %d proxies\n", selectedCount)

        fmt.Println()
        fmt.Println("  Press Enter to continue...")
        fmt.Scanln()
}

// ============ ROTATOR ============

func startRotator() {
        activeProxies = []int{}
        for i, p := range proxies {
                for _, c := range config.SelectedCountries {
                        if p.country == c {
                                activeProxies = append(activeProxies, i)
                                break
                        }
                }
        }

        if len(activeProxies) == 0 {
                printError("No countries selected!")
                time.Sleep(2 * time.Second)
                return
        }

        clearScreen()
        printHeader("🚀 ROTATOR RUNNING")

        fmt.Printf("  🌍 %d proxies\n", len(activeProxies))
        fmt.Printf("  ⏱️  Every %ds\n", config.RotateInterval)
        fmt.Printf("  🔌 Port %d\n", config.ProxyPort)
        fmt.Println()
        fmt.Println(strings.Repeat("─", 50))
        fmt.Println()

        startTime = time.Now()
        currentIndex = activeProxies[rand.Intn(len(activeProxies))]
        patternCount = 0
        history = []string{}
        running = true

        go healthChecker()
        go rotatorLoop()
        startLocalProxy()
}

func healthChecker() {
        for running {
                var wg sync.WaitGroup
                for _, idx := range activeProxies {
                        wg.Add(1)
                        go func(i int) {
                                defer wg.Done()
                                checkProxy(i)
                        }(idx)
                }
                wg.Wait()

                healthMu.RLock()
                healthy := 0
                for _, idx := range activeProxies {
                        if proxyHealth[idx] {
                                healthy++
                        }
                }
                healthMu.RUnlock()

                fmt.Printf("  📊 Healthy: %d/%d\n", healthy, len(activeProxies))

                time.Sleep(30 * time.Second)
        }
}

func checkProxy(index int) {
        addr := proxies[index].addr
        p := proxies[index]

        start := time.Now()
        conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
        latency := time.Since(start)

        if err != nil {
                healthMu.Lock()
                proxyHealth[index] = false
                proxyLatency[index] = 0
                healthMu.Unlock()
                return
        }
        conn.Close()

        // Simple SOCKS5 test
        alive := false
        if p.protocol == "SOCKS5" {
                conn, err = net.DialTimeout("tcp", addr, 3*time.Second)
                if err == nil {
                        conn.Write([]byte{0x05, 0x01, 0x00})
                        buf := make([]byte, 2)
                        conn.SetReadDeadline(time.Now().Add(2 * time.Second))
                        n, _ := conn.Read(buf)
                        conn.Close()
                        alive = (n == 2 && buf[0] == 0x05)
                }
        } else {
                alive = true // Assume HTTP proxies are alive if TCP works
        }

        healthMu.Lock()
        if alive {
                proxyHealth[index] = true
                proxyLatency[index] = latency
        } else {
                proxyHealth[index] = false
                proxyLatency[index] = 0
        }
        healthMu.Unlock()
}

func rotatorLoop() {
        time.Sleep(3 * time.Second)

        for running {
                time.Sleep(time.Duration(config.RotateInterval) * time.Second)

                mu.Lock()

                // Get new index
                healthMu.RLock()
                var healthy []int
                for _, idx := range activeProxies {
                        if proxyHealth[idx] {
                                healthy = append(healthy, idx)
                        }
                }
                healthMu.RUnlock()

                if len(healthy) > 0 {
                        newIndex := healthy[rand.Intn(len(healthy))]
                        if newIndex != currentIndex {
                                currentIndex = newIndex
                                lastSuccessIndex = currentIndex

                                p := proxies[currentIndex]
                                pattern := getRandomPattern()
                                patternCount++

                                fromCountry := "🌐 Start"
                                if len(history) > 0 {
                                        fromCountry = history[len(history)-1]
                                }
                                history = append(history, fmt.Sprintf("%s %s", p.flag, p.country))

                                healthMu.RLock()
                                latency := proxyLatency[currentIndex]
                                healthMu.RUnlock()

                                fmt.Println()
                                fmt.Printf("  🌀 [%d] %s %s\n", patternCount, pattern.name, pattern.symbol)
                                fmt.Printf("     ✈️  %s → %s %s\n", fromCountry, p.flag, p.country)
                                fmt.Printf("     🌍 %s | ⏱️  %v\n", p.addr, latency)
                                fmt.Printf("     🎭 %s\n", p.myth)
                                fmt.Println()
                        }
                }

                mu.Unlock()
        }
}

func startLocalProxy() {
        addr := fmt.Sprintf("127.0.0.1:%d", config.ProxyPort)
        listener, err := net.Listen("tcp", addr)
        if err != nil {
                fmt.Printf("❌ Failed to create listener: %v\n", err)
                return
        }
        defer listener.Close()

        fmt.Printf("  ✅ Local proxy: %s\n", addr)
        fmt.Printf("  🌐 Set browser to SOCKS5 %s\n", addr)
        fmt.Println()
        fmt.Println(strings.Repeat("─", 50))
        fmt.Println()

        for running {
                conn, err := listener.Accept()
                if err != nil {
                        continue
                }
                go handleConnection(conn)
        }
}

func handleConnection(client net.Conn) {
        defer client.Close()

        mu.Lock()
        targetAddr := proxies[currentIndex].addr
        mu.Unlock()

        backend, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
        if err != nil {
                return
        }
        defer backend.Close()

        go relay(backend, client)
        relay(client, backend)
}

func relay(dst, src net.Conn) {
        buf := make([]byte, 32768)
        for {
                src.SetReadDeadline(time.Now().Add(30 * time.Second))
                n, err := src.Read(buf)
                if err != nil {
                        return
                }
                dst.SetWriteDeadline(time.Now().Add(30 * time.Second))
                dst.Write(buf[:n])
        }
}

func getRandomPattern() struct {
        name   string
        symbol string
} {
        if len(usedPatterns) == 0 {
                usedPatterns = rand.Perm(len(mysticalPatterns))
        }

        idx := usedPatterns[0]
        usedPatterns = usedPatterns[1:]
        return mysticalPatterns[idx]
}
