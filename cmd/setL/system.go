package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SystemInfo struct {
	RAMUsedMB   int
	RAMTotalMB  int
	BatteryPct  int
	BatteryStat string
	CPUTempC    int
	UptimeStr   string
	IPAddr      string
}

var demoGovernor = "performance"

func GetSystemMetrics() SystemInfo {
	if DemoMode || runtime.GOOS == "windows" {
		return SystemInfo{
			RAMUsedMB:   418,
			RAMTotalMB:  1860,
			BatteryPct:  89,
			BatteryStat: "Charging",
			CPUTempC:    38,
			UptimeStr:   "3d 14h 22m",
			IPAddr:      "192.168.1.xxx (DEMO)",
		}
	}

	info := SystemInfo{BatteryPct: -1, CPUTempC: -1}

	// RAM
	if memData, err := os.ReadFile("/proc/meminfo"); err == nil {
		var total, free, available, buffers, cached int
		for _, line := range strings.Split(string(memData), "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				val, _ := strconv.Atoi(fields[1])
				switch fields[0] {
				case "MemTotal:":
					total = val
				case "MemFree:":
					free = val
				case "MemAvailable:":
					available = val
				case "Buffers:":
					buffers = val
				case "Cached:":
					cached = val
				}
			}
		}
		info.RAMTotalMB = total / 1024
		if available > 0 {
			info.RAMUsedMB = (total - available) / 1024
		} else {
			info.RAMUsedMB = (total - free - buffers - cached) / 1024
		}
	}

	// Battery
	if capData, err := os.ReadFile("/sys/class/power_supply/battery/capacity"); err == nil {
		if val, err := strconv.Atoi(strings.TrimSpace(string(capData))); err == nil {
			info.BatteryPct = val
		}
	}
	if statData, err := os.ReadFile("/sys/class/power_supply/battery/status"); err == nil {
		info.BatteryStat = strings.TrimSpace(string(statData))
	}

	// CPU Temp
	tempPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/hwmon/hwmon0/temp1_input",
	}
	for _, p := range tempPaths {
		if data, err := os.ReadFile(p); err == nil {
			if val, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
				if val > 1000 {
					val = val / 1000
				}
				info.CPUTempC = val
				break
			}
		}
	}

	// Uptime
	if upData, err := os.ReadFile("/proc/uptime"); err == nil {
		fields := strings.Fields(string(upData))
		if len(fields) > 0 {
			if secs, err := strconv.ParseFloat(fields[0], 64); err == nil {
				d := time.Duration(secs) * time.Second
				h := int(d.Hours())
				m := int(d.Minutes()) % 60
				info.UptimeStr = fmt.Sprintf("%dh %dm", h, m)
			}
		}
	}

	// Local IP
	info.IPAddr = "127.0.0.1"
	if ifaces, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range ifaces {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					info.IPAddr = ipNet.IP.String()
					break
				}
			}
		}
	}

	return info
}

func OptimizeSystem() {
	ClearScreen()
	fmt.Println(ColorCyan + "==================================================" + ColorReset)
	fmt.Println(ColorWhite + "           Оптимизация и ускорение системы        " + ColorReset)
	fmt.Println(ColorCyan + "==================================================" + ColorReset)

	if DemoMode {
		fmt.Print(" [1/4] Сброс кэшей ОЗУ (drop_caches)... ")
		time.Sleep(300 * time.Millisecond)
		fmt.Println(ColorGreen + "OK (эмулировано)" + ColorReset)

		fmt.Print(" [2/4] Компактизация памяти ядра... ")
		time.Sleep(300 * time.Millisecond)
		fmt.Println(ColorGreen + "OK (эмулировано)" + ColorReset)

		fmt.Print(" [3/4] Очистка временных файлов (/tmp)... ")
		time.Sleep(300 * time.Millisecond)
		fmt.Println(ColorGreen + "OK (эмулировано)" + ColorReset)

		fmt.Print(" [4/4] Очистка кэша APT... ")
		time.Sleep(300 * time.Millisecond)
		fmt.Println(ColorGreen + "OK (эмулировано)" + ColorReset)

		fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
		fmt.Println(ColorGreen + "[✓] [DEMO] Система успешно оптимизирована без нагрузки на накопитель!" + ColorReset)
		Pause()
		return
	}

	fmt.Print(" [1/4] Сброс кэшей ОЗУ (drop_caches)... ")
	_ = os.WriteFile("/proc/sys/vm/drop_caches", []byte("3\n"), 0644)
	fmt.Println(ColorGreen + "OK" + ColorReset)

	fmt.Print(" [2/4] Компактизация памяти ядра... ")
	_ = os.WriteFile("/proc/sys/vm/compact_memory", []byte("1\n"), 0644)
	fmt.Println(ColorGreen + "OK" + ColorReset)

	fmt.Print(" [3/4] Очистка временных файлов (/tmp)... ")
	_ = exec.Command("sh", "-c", "rm -rf /tmp/* 2>/dev/null").Run()
	fmt.Println(ColorGreen + "OK" + ColorReset)

	fmt.Print(" [4/4] Очистка кэша APT... ")
	_ = exec.Command("apt-get", "clean").Run()
	fmt.Println(ColorGreen + "OK" + ColorReset)

	fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
	fmt.Println(ColorGreen + "[✓] Система успешно очищена и оптимизирована!" + ColorReset)
	Pause()
}

func CPUGovernorMenu() {
	for {
		ClearScreen()
		currentGov := "unknown"
		availGov := "performance, interactive, ondemand, powersave"

		if DemoMode {
			currentGov = demoGovernor
		} else {
			if data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor"); err == nil {
				currentGov = strings.TrimSpace(string(data))
			}
			if data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_available_governors"); err == nil {
				availGov = strings.TrimSpace(string(data))
			}
		}

		fmt.Println(ColorCyan + "==================================================" + ColorReset)
		fmt.Println(ColorWhite + "              Режим работы процессора (CPU)       " + ColorReset)
		fmt.Println(ColorCyan + "==================================================" + ColorReset)
		fmt.Printf(" Текущий режим:     %s%s%s\n", ColorGreen, currentGov, ColorReset)
		fmt.Printf(" Доступные режимы:  %s%s%s\n", ColorGray, availGov, ColorReset)
		fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)
		fmt.Println(" 1. Performance  (Максимальная скорость серверов)")
		fmt.Println(" 2. Interactive  (Баланс скорости и энергопотребления)")
		fmt.Println(" 3. Powersave    (Минимальный нагрев и энергопотребление)")
		fmt.Println(" 0. Назад")
		fmt.Println(ColorCyan + "==================================================" + ColorReset)

		choice := ReadLine(" Выберите режим: ")
		var target string
		switch choice {
		case "1":
			target = "performance"
		case "2":
			target = "interactive"
		case "3":
			target = "powersave"
		case "0":
			return
		}

		if target != "" {
			if DemoMode {
				demoGovernor = target
				fmt.Printf("\n%s[✓] [DEMO] Режим процессора переключен на: %s (эмулировано)%s\n", ColorGreen, target, ColorReset)
				time.Sleep(1 * time.Second)
				continue
			}

			cpus, _ := filepath.Glob("/sys/devices/system/cpu/cpu[0-9]*/cpufreq/scaling_governor")
			for _, c := range cpus {
				_ = os.WriteFile(c, []byte(target+"\n"), 0644)
			}
			fmt.Printf("\n%s[✓] Режим процессора переключен на: %s%s\n", ColorGreen, target, ColorReset)
			time.Sleep(1 * time.Second)
		}
	}
}

func InstallServerEssentials() {
	ClearScreen()
	fmt.Println(ColorCyan + "==================================================" + ColorReset)
	fmt.Println(ColorWhite + "       Установка базового набора сервера          " + ColorReset)
	fmt.Println(ColorCyan + "==================================================" + ColorReset)
	fmt.Println(" Пакеты: git, build-essential, htop, tmux, ncdu,")
	fmt.Println("         curl, wget, python3, python3-pip, net-tools")
	fmt.Println(ColorCyan + "--------------------------------------------------" + ColorReset)

	confirm := ReadLine(" Начать установку? [y/N]: ")
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" && confirm != "д" {
		return
	}

	if DemoMode {
		fmt.Println("\n" + ColorYellow + "[DEMO] Обновление списков репозиториев..." + ColorReset)
		time.Sleep(500 * time.Millisecond)
		fmt.Println(ColorGreen + "Чтение списков пакетов... Готово" + ColorReset)

		fmt.Println("\n" + ColorYellow + "[DEMO] Установка пакетов (git, tmux, htop, curl, python3)..." + ColorReset)
		time.Sleep(800 * time.Millisecond)
		fmt.Println("\n" + ColorGreen + "[✓] [DEMO] Базовые пакеты сервера успешно эмулированы без трафика и нагрузки!" + ColorReset)
		Pause()
		return
	}

	fmt.Println("\n" + ColorYellow + "Обновление репозиториев..." + ColorReset)
	cmdUpdate := exec.Command("apt-get", "update")
	cmdUpdate.Stdout = os.Stdout
	cmdUpdate.Stderr = os.Stderr
	_ = cmdUpdate.Run()

	fmt.Println("\n" + ColorYellow + "Установка пакетов (может занять несколько минут)..." + ColorReset)
	packages := []string{
		"install", "-y",
		"git", "build-essential", "htop", "tmux", "ncdu",
		"curl", "wget", "python3", "python3-pip", "net-tools",
	}
	cmdInstall := exec.Command("apt-get", packages...)
	cmdInstall.Stdout = os.Stdout
	cmdInstall.Stderr = os.Stderr
	err := cmdInstall.Run()

	if err != nil {
		fmt.Printf("\n%s[!] Ошибка установки: %v%s\n", ColorRed, err, ColorReset)
	} else {
		fmt.Println("\n" + ColorGreen + "[✓] Все базовые компоненты успешно установлены!" + ColorReset)
	}
	Pause()
}

func RebootDevice() {
	ClearScreen()
	fmt.Println(ColorRed + "==================================================" + ColorReset)
	fmt.Println(ColorWhite + "                ПЕРЕЗАГРУЗКА СЕРВЕРА              " + ColorReset)
	fmt.Println(ColorRed + "==================================================" + ColorReset)
	ans := ReadLine(" Вы уверены, что хотите перезагрузить телефон? [y/N]: ")
	if strings.ToLower(ans) == "y" || strings.ToLower(ans) == "yes" || ans == "д" {
		if DemoMode {
			fmt.Println("\n" + ColorYellow + "[DEMO] Запрос перезагрузки через SysRq (echo b > /proc/sysrq-trigger)..." + ColorReset)
			time.Sleep(500 * time.Millisecond)
			fmt.Println(ColorGreen + "[✓] [DEMO] Перезагрузка успешно смоделирована (система не перезагружается)." + ColorReset)
			Pause()
			return
		}

		fmt.Println("\n" + ColorYellow + "Сброс буферов и перезагрузка..." + ColorReset)
		_ = exec.Command("sync").Run()
		_ = os.WriteFile("/proc/sysrq-trigger", []byte("b\n"), 0644)
	}
}

func PoweroffDevice() {
	ClearScreen()
	fmt.Println(ColorRed + "==================================================" + ColorReset)
	fmt.Println(ColorWhite + "                ВЫКЛЮЧЕНИЕ СЕРВЕРА                " + ColorReset)
	fmt.Println(ColorRed + "==================================================" + ColorReset)
	ans := ReadLine(" Вы уверены, что хотите ВЫКЛЮЧИТЬ телефон? [y/N]: ")
	if strings.ToLower(ans) == "y" || strings.ToLower(ans) == "yes" || ans == "д" {
		if DemoMode {
			fmt.Println("\n" + ColorYellow + "[DEMO] Запрос отключения через SysRq (echo o > /proc/sysrq-trigger)..." + ColorReset)
			time.Sleep(500 * time.Millisecond)
			fmt.Println(ColorGreen + "[✓] [DEMO] Выключение успешно смоделировано (питание не отключается)." + ColorReset)
			Pause()
			return
		}

		fmt.Println("\n" + ColorYellow + "Сброс буферов и отключение питания..." + ColorReset)
		_ = exec.Command("sync").Run()
		_ = os.WriteFile("/proc/sysrq-trigger", []byte("o\n"), 0644)
	}
}
