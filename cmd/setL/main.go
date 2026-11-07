package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func ShowSSHStatus() {
	ClearScreen()
	fmt.Println(ColorCyan + "==================================================" + ColorReset)
	fmt.Println(ColorWhite + "               Статус службы OpenSSH              " + ColorReset)
	fmt.Println(ColorCyan + "==================================================" + ColorReset)

	if DemoMode {
		fmt.Printf(" Служба sshd: %sРАБОТАЕТ (PID: 1248) [DEMO]%s\n", ColorGreen, ColorReset)
		fmt.Printf(" Порт:        %s22 (LISTEN)%s\n", ColorGreen, ColorReset)
		fmt.Printf(" Ключей SSH:  %s1 (/root/.ssh/authorized_keys)%s\n", ColorWhite, ColorReset)
		fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
		fmt.Println(" Команда для подключения с ПК:")
		fmt.Printf(" %sssh root@192.168.1.xxx%s\n", ColorGreen, ColorReset)
		fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
		Pause()
		return
	}

	out, _ := exec.Command("sh", "-c", "ps -ef | grep [s]shd").Output()
	if len(out) > 0 {
		fmt.Printf(" Служба sshd: %sРАБОТАЕТ (PID: %s)%s\n", ColorGreen, strings.Fields(string(out))[1], ColorReset)
	} else {
		fmt.Printf(" Служба sshd: %sОСТАНОВЛЕНА%s\n", ColorRed, ColorReset)
	}

	netOut, _ := exec.Command("sh", "-c", "netstat -tlpn 2>/dev/null | grep sshd").Output()
	if len(netOut) > 0 {
		fmt.Printf(" Порт:        %s22 (LISTEN)%s\n", ColorGreen, ColorReset)
	}

	keysOut, _ := exec.Command("sh", "-c", "wc -l < /root/.ssh/authorized_keys 2>/dev/null").Output()
	fmt.Printf(" Ключей SSH:  %s%s%s\n", ColorWhite, strings.TrimSpace(string(keysOut)), ColorReset)

	fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
	fmt.Println(" Команда для подключения с ПК:")
	fmt.Printf(" %sssh root@%s%s\n", ColorGreen, GetSystemMetrics().IPAddr, ColorReset)
	fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)

	Pause()
}

func main() {
	for {
		metrics := GetSystemMetrics()
		pCfg := LoadProxyConfig()
		dnsList := GetCurrentDNS()
		dnsStr := "Не задан"
		if len(dnsList) > 0 {
			dnsStr = dnsList[0]
		}

		proxyStr := ColorRed + "ВЫКЛ" + ColorReset
		if pCfg.Enabled {
			proxyStr = ColorGreen + "ВКЛ" + ColorReset
		}

		ClearScreen()
		fmt.Println(ColorCyan + "================================================================" + ColorReset)
		fmt.Println(ColorWhite + "           ⚡ setup-loader (setL) — Панель управления           " + ColorReset)
		if DemoMode {
			fmt.Println(ColorYellow + "       🎭 [РЕЖИМ ДЕМО: Симуляция включена — без нагрузки]       " + ColorReset)
		}
		fmt.Println(ColorCyan + "================================================================" + ColorReset)

		batteryStr := "N/A"
		if metrics.BatteryPct >= 0 {
			batteryStr = fmt.Sprintf("%d%%", metrics.BatteryPct)
		}
		tempStr := "N/A"
		if metrics.CPUTempC >= 0 {
			tempStr = fmt.Sprintf("%d°C", metrics.CPUTempC)
		}

		fmt.Printf(" IP: %-15s | DNS: %-13s | Прокси: %s\n",
			metrics.IPAddr, dnsStr, proxyStr)
		fmt.Printf(" RAM: %dM / %dM     | Батарея: %-10s | CPU Темп: %s\n",
			metrics.RAMUsedMB, metrics.RAMTotalMB, batteryStr, tempStr)
		if metrics.UptimeStr != "" {
			fmt.Printf(" Uptime: %-12s\n", metrics.UptimeStr)
		}
		fmt.Println(ColorCyan + "================================================================" + ColorReset)
		fmt.Println(" 1. 🌐 Управление Прокси (Вкл / Выкл / Настройка)")
		fmt.Println(" 2. 📡 Настройка DNS (Quad9, Cloudflare, Google, AdGuard, Свой)")
		fmt.Println(" 3. 🚀 Быстрая оптимизация системы (Очистка ОЗУ, кэша)")
		fmt.Println(" 4. ⚡ Режим работы процессора (CPU Governor)")
		fmt.Println(" 5. 📦 Установка базовых пакетов сервера (Essentials Pack)")
		fmt.Println(" 6. 🔑 Проверка и статус службы SSH")
		fmt.Println(" 7. 🔄 Перезагрузка сервера (Reboot)")
		fmt.Println(" 8. ⏹  Выключение сервера (Poweroff)")

		demoStatus := ColorRed + "ВЫКЛ" + ColorReset
		if DemoMode {
			demoStatus = ColorGreen + "ВКЛ (симуляция)" + ColorReset
		}
		fmt.Printf(" 9. 🎭 Демо-режим: %s\n", demoStatus)
		fmt.Println(" 0. 🚪 Выйти")
		fmt.Println(ColorCyan + "================================================================" + ColorReset)

		choice := ReadLine(" Выберите действие: ")
		switch choice {
		case "1":
			ProxyMenu()
		case "2":
			DNSMenu()
		case "3":
			OptimizeSystem()
		case "4":
			CPUGovernorMenu()
		case "5":
			InstallServerEssentials()
		case "6":
			ShowSSHStatus()
		case "7":
			RebootDevice()
		case "8":
			PoweroffDevice()
		case "9":
			DemoMode = !DemoMode
		case "0":
			ClearScreen()
			return
		}
	}
}

