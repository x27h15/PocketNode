package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	ProxyConfFile  = "/etc/proxy.conf"
	ProxyEnvFile   = "/etc/profile.d/99proxy.sh"
	ProxyAptFile   = "/etc/apt/apt.conf.d/80proxy"
	ProxyCheckURL  = "https://api.ipify.org"
)

type ProxyConfig struct {
	Enabled bool
	Addr    string
	User    string
	Pass    string
}

var demoProxy = ProxyConfig{
	Enabled: false,
	Addr:    "127.0.0.1:1080",
	User:    "",
	Pass:    "",
}

func LoadProxyConfig() ProxyConfig {
	if DemoMode {
		return demoProxy
	}
	cfg := ProxyConfig{}
	data, err := os.ReadFile(ProxyConfFile)
	if err != nil {
		return cfg
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "PROXY_ENABLED=") {
			val := strings.TrimPrefix(line, "PROXY_ENABLED=")
			cfg.Enabled = (val == "1" || val == "true")
		} else if strings.HasPrefix(line, "PROXY_ADDR=") {
			cfg.Addr = strings.Trim(strings.TrimPrefix(line, "PROXY_ADDR="), "\"")
		} else if strings.HasPrefix(line, "PROXY_USER=") {
			cfg.User = strings.Trim(strings.TrimPrefix(line, "PROXY_USER="), "\"")
		} else if strings.HasPrefix(line, "PROXY_PASS=") {
			cfg.Pass = strings.Trim(strings.TrimPrefix(line, "PROXY_PASS="), "\"")
		}
	}
	return cfg
}

func SaveProxyConfig(cfg ProxyConfig) error {
	if DemoMode {
		demoProxy = cfg
		return nil
	}

	enabledStr := "0"
	if cfg.Enabled {
		enabledStr = "1"
	}
	content := fmt.Sprintf("PROXY_ENABLED=%s\nPROXY_ADDR=\"%s\"\nPROXY_USER=\"%s\"\nPROXY_PASS=\"%s\"\n",
		enabledStr, cfg.Addr, cfg.User, cfg.Pass)
	return os.WriteFile(ProxyConfFile, []byte(content), 0600)
}

func BuildProxyURL(addr, user, pass string) string {
	if addr == "" {
		return ""
	}
	proto := "http"
	hostport := addr
	if strings.HasPrefix(addr, "socks5://") {
		proto = "socks5"
		hostport = strings.TrimPrefix(addr, "socks5://")
	} else if strings.HasPrefix(addr, "socks4://") {
		proto = "socks4"
		hostport = strings.TrimPrefix(addr, "socks4://")
	} else if strings.HasPrefix(addr, "http://") {
		proto = "http"
		hostport = strings.TrimPrefix(addr, "http://")
	} else if strings.HasPrefix(addr, "https://") {
		proto = "http"
		hostport = strings.TrimPrefix(addr, "https://")
	}

	if user != "" {
		if pass != "" {
			return fmt.Sprintf("%s://%s:%s@%s", proto, url.QueryEscape(user), url.QueryEscape(pass), hostport)
		}
		return fmt.Sprintf("%s://%s@%s", proto, url.QueryEscape(user), hostport)
	}
	return fmt.Sprintf("%s://%s", proto, hostport)
}

func ApplyProxy(cfg *ProxyConfig) error {
	fullURL := BuildProxyURL(cfg.Addr, cfg.User, cfg.Pass)
	if fullURL == "" {
		return fmt.Errorf("адрес прокси не заполнен")
	}

	if DemoMode {
		cfg.Enabled = true
		return SaveProxyConfig(*cfg)
	}

	envScript := fmt.Sprintf(`export http_proxy="%s"
export https_proxy="%s"
export ftp_proxy="%s"
export all_proxy="%s"
export HTTP_PROXY="%s"
export HTTPS_PROXY="%s"
export FTP_PROXY="%s"
export ALL_PROXY="%s"
export no_proxy="localhost,127.0.0.1,localaddress,.localdomain.com,192.168.0.0/16,10.0.0.0/8"
export NO_PROXY="localhost,127.0.0.1,localaddress,.localdomain.com,192.168.0.0/16,10.0.0.0/8"
`, fullURL, fullURL, fullURL, fullURL, fullURL, fullURL, fullURL, fullURL)

	_ = os.WriteFile(ProxyEnvFile, []byte(envScript), 0644)

	aptConfig := fmt.Sprintf("Acquire::http::Proxy \"%s\";\nAcquire::https::Proxy \"%s\";\n", fullURL, fullURL)
	_ = os.WriteFile(ProxyAptFile, []byte(aptConfig), 0644)

	os.Setenv("http_proxy", fullURL)
	os.Setenv("https_proxy", fullURL)
	os.Setenv("all_proxy", fullURL)

	cfg.Enabled = true
	return SaveProxyConfig(*cfg)
}

func DisableProxy(cfg *ProxyConfig) error {
	if DemoMode {
		cfg.Enabled = false
		return SaveProxyConfig(*cfg)
	}

	_ = os.Remove(ProxyEnvFile)
	_ = os.Remove(ProxyAptFile)

	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
	os.Unsetenv("all_proxy")

	cfg.Enabled = false
	return SaveProxyConfig(*cfg)
}

func TestProxyConnection(proxyURL string) (string, error) {
	if DemoMode {
		time.Sleep(300 * time.Millisecond)
		return "198.51.100.42", nil
	}

	client := &http.Client{Timeout: 5 * time.Second}
	if proxyURL != "" {
		pURL, err := url.Parse(proxyURL)
		if err != nil {
			return "", err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(pURL)}
	}
	resp, err := client.Get(ProxyCheckURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func MaskPassword(pass string) string {
	if pass == "" {
		return "<Не задан>"
	}
	return strings.Repeat("*", len(pass))
}

func ProxyMenu() {
	for {
		cfg := LoadProxyConfig()
		ClearScreen()

		statusStr := ColorRed + "ВЫКЛЮЧЕН" + ColorReset
		currentIP := ""
		if cfg.Enabled {
			fullURL := BuildProxyURL(cfg.Addr, cfg.User, cfg.Pass)
			ip, err := TestProxyConnection(fullURL)
			if err == nil && ip != "" {
				statusStr = ColorGreen + "ВКЛЮЧЕН (Работает, IP: " + ip + ")" + ColorReset
				currentIP = ip
			} else {
				statusStr = ColorYellow + "ВКЛЮЧЕН (Сервер не отвечает)" + ColorReset
			}
		}

		configStatus := ColorRed + "Не заполнены" + ColorReset
		if cfg.Addr != "" {
			configStatus = ColorGreen + "Заполнены (" + cfg.Addr + ")" + ColorReset
		}

		directIP, _ := TestProxyConnection("")
		if directIP == "" {
			directIP = "Нет сети"
		}

		fmt.Println(ColorCyan + "==================================================" + ColorReset)
		fmt.Println(ColorWhite + "              Управление Прокси (VPN)             " + ColorReset)
		fmt.Println(ColorCyan + "==================================================" + ColorReset)
		fmt.Printf(" Статус:    %s\n", statusStr)
		fmt.Printf(" Настройки: %s\n", configStatus)
		if currentIP != "" {
			fmt.Printf(" Прокси IP: %s%s%s\n", ColorGreen, currentIP, ColorReset)
		}
		fmt.Printf(" Прямой IP: %s%s%s\n", ColorWhite, directIP, ColorReset)
		fmt.Println(ColorCyan + "==================================================" + ColorReset)
		fmt.Println(" 1. Включить прокси")
		fmt.Println(" 2. Выключить прокси")
		fmt.Println(" 3. Настроить прокси")
		fmt.Println(" 0. Выйти")
		fmt.Println(ColorCyan + "==================================================" + ColorReset)

		choice := ReadLine(" Выберите действие: ")
		switch choice {
		case "1":
			if cfg.Addr == "" {
				fmt.Println("\n" + ColorRed + "[!] Ошибка: данные для подключения не заполнены!" + ColorReset)
				fmt.Println("Сначала выберите пункт 3 и настройте адрес сервера.")
				Pause()
			} else {
				fmt.Println("\n" + ColorYellow + "Применение настроек прокси..." + ColorReset)
				if err := ApplyProxy(&cfg); err != nil {
					fmt.Printf("\n%s[!] Ошибка: %s%s\n", ColorRed, err, ColorReset)
				} else {
					fmt.Println(ColorGreen + "[✓] Прокси успешно включен!" + ColorReset)
				}
				time.Sleep(1 * time.Second)
			}
		case "2":
			fmt.Println("\n" + ColorYellow + "Выключение прокси..." + ColorReset)
			_ = DisableProxy(&cfg)
			fmt.Println(ColorGreen + "[✓] Прокси выключен." + ColorReset)
			time.Sleep(1 * time.Second)
		case "3":
			configureProxySubmenu(&cfg)
		case "0":
			return
		}
	}
}

func configureProxySubmenu(cfg *ProxyConfig) {
	tmpAddr := cfg.Addr
	tmpUser := cfg.User
	tmpPass := cfg.Pass

	for {
		ClearScreen()
		fmt.Println(ColorCyan + "==================================================" + ColorReset)
		fmt.Println(ColorWhite + "                 Настройка Прокси                 " + ColorReset)
		fmt.Println(ColorCyan + "==================================================" + ColorReset)
		fmt.Printf(" Адрес сервера: %s%s%s\n", ColorWhite, defaultStr(tmpAddr, "<Не задан>"), ColorReset)
		fmt.Printf(" Логин:         %s%s%s\n", ColorWhite, defaultStr(tmpUser, "<Не задан>"), ColorReset)
		fmt.Printf(" Пароль:        %s%s%s\n", ColorWhite, MaskPassword(tmpPass), ColorReset)
		fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
		fmt.Println(" 1. Изменить")
		fmt.Println(" 2. Сохранить")
		fmt.Println(" 0. Выйти")
		fmt.Println(ColorCyan + "==================================================" + ColorReset)

		subChoice := ReadLine(" Выберите действие: ")
		switch subChoice {
		case "1":
			fmt.Println()
			tmpAddr = ReadLine(" Адрес сервера (например, 1.2.3.4:8080): ")
			tmpUser = ReadLine(" Логин (Enter если нет): ")
			tmpPass = ReadLine(" Пароль (Enter если нет): ")

			for {
				ClearScreen()
				fmt.Println(ColorCyan + "==================================================" + ColorReset)
				fmt.Println(ColorWhite + "                   Новые данные                   " + ColorReset)
				fmt.Println(ColorCyan + "==================================================" + ColorReset)
				fmt.Printf(" Адрес сервера: %s%s%s\n", ColorWhite, defaultStr(tmpAddr, "<Не задан>"), ColorReset)
				fmt.Printf(" Логин:         %s%s%s\n", ColorWhite, defaultStr(tmpUser, "<Не задан>"), ColorReset)
				fmt.Printf(" Пароль:        %s%s%s\n", ColorWhite, MaskPassword(tmpPass), ColorReset)
				fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
				fmt.Println(" 1. Сохранить и выйти на главную")
				fmt.Println(" 2. Сохранить")
				fmt.Println(" 3. Тест подключения")
				fmt.Println(" 0. Выйти без сохранения")
				fmt.Println(ColorCyan + "==================================================" + ColorReset)

				postChoice := ReadLine(" Выберите действие: ")
				switch postChoice {
				case "1":
					cfg.Addr = tmpAddr
					cfg.User = tmpUser
					cfg.Pass = tmpPass
					_ = SaveProxyConfig(*cfg)
					if cfg.Enabled {
						_ = ApplyProxy(cfg)
					}
					fmt.Println("\n" + ColorGreen + "[✓] Настройки сохранены!" + ColorReset)
					time.Sleep(1 * time.Second)
					return
				case "2":
					cfg.Addr = tmpAddr
					cfg.User = tmpUser
					cfg.Pass = tmpPass
					_ = SaveProxyConfig(*cfg)
					if cfg.Enabled {
						_ = ApplyProxy(cfg)
					}
					fmt.Println("\n" + ColorGreen + "[✓] Настройки сохранены!" + ColorReset)
					time.Sleep(1 * time.Second)
					goto SubmenuLoop
				case "3":
					fmt.Println("\n" + ColorYellow + "Тестирование соединения через указанный прокси..." + ColorReset)
					testURL := BuildProxyURL(tmpAddr, tmpUser, tmpPass)
					ip, err := TestProxyConnection(testURL)
					if err == nil && ip != "" {
						fmt.Printf("%s[✓] Соединение успешно! Внешний IP: %s%s\n", ColorGreen, ip, ColorReset)
					} else {
						fmt.Printf("%s[✗] Ошибка соединения через прокси: %v%s\n", ColorRed, err, ColorReset)
					}
					Pause()
				case "0":
					tmpAddr = cfg.Addr
					tmpUser = cfg.User
					tmpPass = cfg.Pass
					goto SubmenuLoop
				}
			}
		case "2":
			cfg.Addr = tmpAddr
			cfg.User = tmpUser
			cfg.Pass = tmpPass
			_ = SaveProxyConfig(*cfg)
			if cfg.Enabled {
				_ = ApplyProxy(cfg)
			}
			fmt.Println("\n" + ColorGreen + "[✓] Настройки сохранены!" + ColorReset)
			time.Sleep(1 * time.Second)
			return
		case "0":
			return
		}
	SubmenuLoop:
	}
}

func defaultStr(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
