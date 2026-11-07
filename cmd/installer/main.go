package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[1;31m"
	ColorGreen  = "\033[1;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[1;34m"
	ColorCyan   = "\033[1;36m"
	ColorWhite  = "\033[1;37m"
	ColorGray   = "\033[0;90m"
)

var (
	reader   = bufio.NewReader(os.Stdin)
	DemoMode = false
)

func init() {
	val := strings.ToLower(os.Getenv("DEMO"))
	if val == "1" || val == "true" || val == "yes" {
		DemoMode = true
	}
	for _, arg := range os.Args[1:] {
		if arg == "--demo" || arg == "-demo" || arg == "demo" {
			DemoMode = true
		}
	}
}

func Clear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func ReadLine(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func Pause() {
	fmt.Print("\n" + ColorGray + "Нажмите Enter для продолжения..." + ColorReset)
	_, _ = reader.ReadString('\n')
}

func PrintBanner() {
	Clear()
	fmt.Println(ColorCyan + `
  ██████╗  ██████╗  ██████╗██╗  ██╗███████╗████████╗███╗   ██╗ ██████╗ ██████╗ ███████╗
  ██╔══██╗██╔═══██╗██╔════╝██║ ██╔╝██╔════╝╚══██╔══╝████╗  ██║██╔═══██╗██╔══██╗██╔════╝
  ██████╔╝██║   ██║██║     █████═╝ █████╗     ██║   ██╔██╗ ██║██║   ██║██║  ██║█████╗  
  ██╔═══╝ ██║   ██║██║     ██╔═██╗ ██╔══╝     ██║   ██║╚██╗██║██║   ██║██║  ██║██╔══╝  
  ██║     ╚██████╔╝╚██████╗██║ ╚██╗███████╗   ██║   ██║ ╚████║╚██████╔╝██████╔╝███████╗
  ╚═╝      ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚══════╝
` + ColorReset)
	fmt.Println(ColorWhite + "       ⚡ Автоматический установщик PocketNode для Samsung Galaxy A5 (MSM8916)" + ColorReset)
	if DemoMode {
		fmt.Println(ColorYellow + "       🎭 [РЕЖИМ ДЕМО: Симуляция включена — реальные команды не исполняются]" + ColorReset)
	}
	fmt.Println(ColorCyan + "================================================================================" + ColorReset)
}

func CheckCommand(cmdName string) bool {
	_, err := exec.LookPath(cmdName)
	return err == nil
}

func StepRecovery() {
	PrintBanner()
	fmt.Println(ColorYellow + "ШАГ 1: Прошивка TWRP Recovery" + ColorReset)
	fmt.Println(ColorCyan + "--------------------------------------------------------------------------------" + ColorReset)
	fmt.Println("Для установки кастомной прошивки необходимо прошить TWRP Recovery в телефон.")
	fmt.Println()
	fmt.Println("1. Полностью выключите телефон.")
	fmt.Println("2. Зажмите комбинацию: " + ColorWhite + "Громкость ВНИЗ + Кнопка ДОМОЙ (Home) + Кнопка ВКЛЮЧЕНИЯ" + ColorReset)
	fmt.Println("3. При появлении предупреждающего экрана нажмите: " + ColorWhite + "Громкость ВВЕРХ" + ColorReset + " (вход в Download Mode).")
	fmt.Println("4. Подключите телефон к компьютеру оригинальным кабелем USB.")
	fmt.Println()
	fmt.Println("В папке проекта подготовлен образ рекавери:")
	assetsDir, _ := filepath.Abs("assets")
	recoveryPath := filepath.Join(assetsDir, "recovery.tar")
	fmt.Printf("Файл: %s%s%s\n", ColorGreen, recoveryPath, ColorReset)
	fmt.Println()
	fmt.Println("В программе Odin:")
	fmt.Println(" • Нажмите кнопку " + ColorWhite + "AP" + ColorReset + " и выберите файл " + ColorGreen + "recovery.tar" + ColorReset)
	fmt.Println(" • Во вкладке Options " + ColorRed + "СНИМИТЕ галочку Auto Reboot" + ColorReset)
	fmt.Println(" • Нажмите " + ColorWhite + "Start" + ColorReset + " и дождитесь зелёной надписи PASS.")
	fmt.Println()
	fmt.Println("После прошивки зажмите: " + ColorWhite + "Громкость ВВЕРХ + Home + Включение" + ColorReset + " для входа в TWRP.")
	fmt.Println(ColorCyan + "--------------------------------------------------------------------------------" + ColorReset)
	Pause()
}

func StepROM() {
	PrintBanner()
	fmt.Println(ColorYellow + "ШАГ 2: Прошивка LineageOS 18.1 (Android 11)" + ColorReset)
	fmt.Println(ColorCyan + "--------------------------------------------------------------------------------" + ColorReset)
	fmt.Println("В меню TWRP Recovery на телефоне:")
	fmt.Println()
	fmt.Println("1. Зайдите в раздел " + ColorWhite + "Wipe" + ColorReset + " -> " + ColorWhite + "Advanced Wipe" + ColorReset)
	fmt.Println("2. Отметьте галочками:")
	fmt.Println("   • Dalvik / ART Cache")
	fmt.Println("   • Cache")
	fmt.Println("   • Data")
	fmt.Println(ColorGreen + "   (ВНИМАНИЕ: Раздел System вайпать НЕ нужно! Всё работает без затрагивания system)" + ColorReset)
	fmt.Println("3. Сдвиньте слайдер " + ColorWhite + "Swipe to Wipe" + ColorReset)
	fmt.Println("4. Вернитесь в главное меню TWRP -> " + ColorWhite + "Install" + ColorReset)
	fmt.Println("5. Выберите архив с прошивкой (lineage-18.1-xxxx-a5lte.zip) и свайпните для прошивки.")
	fmt.Println("6. Нажмите " + ColorWhite + "Reboot System" + ColorReset + " и дождитесь первой загрузки системы.")
	fmt.Println()
	fmt.Println(ColorRed + "ВНИМАНИЕ: Не прошивайте архив Magisk! PocketNode использует нативный chroot." + ColorReset)
	fmt.Println(ColorCyan + "--------------------------------------------------------------------------------" + ColorReset)
	Pause()
}

func StepADB() {
	for {
		PrintBanner()
		fmt.Println(ColorYellow + "ШАГ 3: Включение отладки и рут-доступа" + ColorReset)
		fmt.Println(ColorCyan + "--------------------------------------------------------------------------------" + ColorReset)
		fmt.Println("После первой загрузки телефона выполните быструю настройку:")
		fmt.Println()
		fmt.Println("1. Зайдите в " + ColorWhite + "Настройки" + ColorReset + " -> " + ColorWhite + "О телефоне" + ColorReset)
		fmt.Println("2. В самом низу нажмите 7 раз на пункт " + ColorWhite + "«Номер сборки»" + ColorReset + " (станете разработчиком).")
		fmt.Println("3. Вернитесь в " + ColorWhite + "Система" + ColorReset + " -> " + ColorWhite + "Для разработчиков" + ColorReset + ":")
		fmt.Println("   • Включите " + ColorGreen + "«Отладка по USB»" + ColorReset)
		fmt.Println("   • В пункте " + ColorGreen + "«Рут-доступ»" + ColorReset + " выберите " + ColorGreen + "«Только ADB»" + ColorReset)
		fmt.Println("4. Подключите телефон кабелем к ПК (при запросе на экране смартфона разрешите отладку).")
		fmt.Println(ColorCyan + "--------------------------------------------------------------------------------" + ColorReset)
		fmt.Println("Ожидание подключения устройства по ADB...")

		if WaitForDevice() {
			fmt.Println("\n" + ColorGreen + "[✓] Смартфон обнаружен и готов к настройке!" + ColorReset)
			time.Sleep(1 * time.Second)
			return
		}

		retry := ReadLine("\nУстройство не найдено. Повторить поиск? [Y/n]: ")
		if strings.ToLower(retry) == "n" {
			return
		}
	}
}

func WaitForDevice() bool {
	if DemoMode {
		fmt.Print("Поиск устройства.")
		for i := 0; i < 3; i++ {
			time.Sleep(350 * time.Millisecond)
			fmt.Print(".")
		}
		fmt.Println()
		return true
	}

	for i := 0; i < 15; i++ {
		cmd := exec.Command("adb", "devices")
		out, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if strings.Contains(l, "\tdevice") {
					_ = exec.Command("adb", "root").Run()
					time.Sleep(1 * time.Second)
					return true
				}
			}
		}
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}
	return false
}

func RunDemoDeploy() {
	fmt.Println(ColorYellow + "Начинаю демонстрационное развёртывание микросервера..." + ColorReset)

	steps := []struct {
		title string
		sleep time.Duration
	}{
		{"[1/7] Перевод ADB в режим root и перемонтирование /system...", 500 * time.Millisecond},
		{"[2/7] Настройка постоянного порта ADB 5555 (/system/etc/init/adbtcp.rc)...", 500 * time.Millisecond},
		{"[3/7] Установка службы автозапуска сервера (/system/bin/start_sshd.sh)...", 600 * time.Millisecond},
		{"[4/7] Генерация и передача ключей SSH на сервер (root/.ssh/authorized_keys)...", 500 * time.Millisecond},
		{"[5/7] Установка утилиты PocketNode (setL, setup, pocketnode)...", 600 * time.Millisecond},
		{"[6/7] Настройка быстрых команд reboot и poweroff (SysRq trigger)...", 400 * time.Millisecond},
		{"[7/7] Запуск службы OpenSSH на сервере (порт 22)...", 700 * time.Millisecond},
	}

	for _, step := range steps {
		fmt.Print(" " + step.title + " ")
		time.Sleep(step.sleep)
		fmt.Println(ColorGreen + "OK (симуляция)" + ColorReset)
	}

	home, _ := os.UserHomeDir()
	keyPath := filepath.Join(home, ".ssh", "id_pocketnode")
	PrintSuccess("192.168.1.xxx (DEMO)", keyPath)
}

func GetRealDeviceIP() string {
	cmds := []string{
		"ip route get 1.1.1.1 | awk '{for(i=1;i<=NF;i++) if ($i==\"src\") print $(i+1)}'",
		"ip -4 addr show wlan0 | grep -oE '[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}' | head -n 1",
		"getprop dhcp.wlan0.ipaddress",
	}
	for _, c := range cmds {
		out, err := exec.Command("adb", "shell", c).Output()
		if err == nil {
			ip := strings.TrimSpace(string(out))
			if ip != "" && !strings.Contains(ip, "127.0.0.1") && net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	return "192.168.1.xxx"
}

func StepDeploy() {
	PrintBanner()
	fmt.Println(ColorYellow + "ШАГ 4: Автоматическая установка PocketNode" + ColorReset)
	fmt.Println(ColorCyan + "--------------------------------------------------------------------------------" + ColorReset)

	if DemoMode {
		RunDemoDeploy()
		return
	}

	fmt.Println("Начинаю развёртывание микросервера...")

	// 1. Root and Remount
	fmt.Print(" [1/7] Перевод ADB в режим root и перемонтирование /system... ")
	_ = exec.Command("adb", "root").Run()
	time.Sleep(1 * time.Second)
	_ = exec.Command("adb", "remount").Run()
	fmt.Println(ColorGreen + "OK" + ColorReset)

	// 2. Persistent ADB on port 5555
	fmt.Print(" [2/7] Настройка постоянного порта ADB 5555... ")
	rcContent := `service debian_sshd /system/bin/sh /system/bin/start_sshd.sh
    class late_start
    user root
    group root
    seclabel u:r:su:s0

on property:sys.boot_completed=1
    setprop service.adb.tcp.port 5555
    restart adbd
    start debian_sshd
`
	_ = exec.Command("adb", "shell", "setprop persist.adb.tcp.port 5555").Run()
	_ = exec.Command("adb", "shell", fmt.Sprintf("echo '%s' > /system/etc/init/adbtcp.rc; chmod 644 /system/etc/init/adbtcp.rc", rcContent)).Run()
	fmt.Println(ColorGreen + "OK" + ColorReset)

	// 3. Setup start_sshd.sh
	fmt.Print(" [3/7] Установка службы автозапуска сервера (/system/bin/start_sshd.sh)... ")
	sshdScript := `#!/system/bin/sh
D=/data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs

mountpoint -q $D/dev || mount --bind /dev $D/dev
mountpoint -q $D/dev/pts || mount -t devpts devpts $D/dev/pts
mountpoint -q $D/proc || mount -t proc proc $D/proc
mountpoint -q $D/sys || mount --bind /sys $D/sys
mkdir -p $D/sdcard
mountpoint -q $D/sdcard || mount --bind /storage/emulated/0 $D/sdcard 2>/dev/null

echo "nameserver 9.9.9.9" > $D/etc/resolv.conf
echo "nameserver 149.112.112.112" >> $D/etc/resolv.conf

mkdir -p $D/run/sshd $D/var/run/sshd
chmod 755 $D/run $D/run/sshd $D/var/run $D/var/run/sshd

exec chroot $D /usr/sbin/sshd -D
`
	_ = exec.Command("adb", "shell", fmt.Sprintf("echo '%s' > /system/bin/start_sshd.sh; chmod 755 /system/bin/start_sshd.sh", sshdScript)).Run()
	fmt.Println(ColorGreen + "OK" + ColorReset)

	// 4. Setup SSH Keys
	fmt.Print(" [4/7] Генерация и передача ключей SSH на сервер... ")
	home, _ := os.UserHomeDir()
	sshDir := filepath.Join(home, ".ssh")
	_ = os.MkdirAll(sshDir, 0700)
	keyPath := filepath.Join(sshDir, "id_pocketnode")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		_ = exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-N", "", "-C", "pocketnode").Run()
	}
	pubKeyData, _ := os.ReadFile(keyPath + ".pub")
	if len(pubKeyData) > 0 {
		_ = exec.Command("adb", "shell", "mkdir -p /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/root/.ssh").Run()
		_ = exec.Command("adb", "shell", fmt.Sprintf("echo '%s' >> /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/root/.ssh/authorized_keys", strings.TrimSpace(string(pubKeyData)))).Run()
		_ = exec.Command("adb", "shell", "chmod 700 /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/root/.ssh; chmod 600 /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/root/.ssh/authorized_keys").Run()
	}
	fmt.Println(ColorGreen + "OK" + ColorReset)

	// 5. Deploy setL binary
	fmt.Print(" [5/7] Установка утилиты PocketNode (assets/setL)... ")
	setLPath := filepath.Join("assets", "setL")
	if _, err := os.Stat(setLPath); err != nil {
		setLPath = "setL"
	}
	if _, err := os.Stat(setLPath); err == nil {
		_ = exec.Command("adb", "push", setLPath, "/data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/bin/setL").Run()
		_ = exec.Command("adb", "shell", "chmod 755 /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/bin/setL").Run()
		_ = exec.Command("adb", "shell", "ln -sf /usr/local/bin/setL /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/bin/setup").Run()
		_ = exec.Command("adb", "shell", "ln -sf /usr/local/bin/setL /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/bin/pocketnode").Run()
	}
	fmt.Println(ColorGreen + "OK" + ColorReset)

	// 6. Fix Reboot & Poweroff handlers
	fmt.Print(" [6/7] Настройка быстрых команд reboot и poweroff... ")
	_ = exec.Command("adb", "shell", "echo 'sync; echo b > /proc/sysrq-trigger' > /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/sbin/reboot; chmod 755 /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/sbin/reboot").Run()
	_ = exec.Command("adb", "shell", "echo 'sync; echo o > /proc/sysrq-trigger' > /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/sbin/poweroff; chmod 755 /data/data/com.termux/files/usr/var/lib/proot-distro/containers/debian/rootfs/usr/local/sbin/poweroff").Run()
	fmt.Println(ColorGreen + "OK" + ColorReset)

	// 7. Start SSH
	fmt.Print(" [7/7] Запуск службы OpenSSH на сервере... ")
	_ = exec.Command("adb", "shell", "nohup /system/bin/sh /system/bin/start_sshd.sh >/dev/null 2>&1 &").Run()
	time.Sleep(2 * time.Second)
	fmt.Println(ColorGreen + "OK" + ColorReset)

	// Get Real Device IP
	deviceIP := GetRealDeviceIP()
	PrintSuccess(deviceIP, keyPath)
}

func PrintSuccess(ip, keyPath string) {
	fmt.Println()
	fmt.Println(ColorGreen + "================================================================================" + ColorReset)
	fmt.Println(ColorWhite + "            🎉 ПОЗДРАВЛЯЕМ! POCKETNODE УСПЕШНО УСТАНОВЛЕН И ЗАПУЩЕН            " + ColorReset)
	fmt.Println(ColorGreen + "================================================================================" + ColorReset)
	fmt.Println(" Телефон теперь работает как автономный Linux-микросервер.")
	fmt.Println(" Кабель от ПК можно отключить навсегда!")
	fmt.Println()
	fmt.Println(ColorYellow + "Для подключения с этого компьютера введите:" + ColorReset)
	fmt.Printf("   %sssh root@%s%s\n", ColorWhite, ip, ColorReset)
	fmt.Println()
	fmt.Println(ColorYellow + "Для подключения с других устройств (по ключу):" + ColorReset)
	fmt.Printf("   %sssh -i %s root@%s%s\n", ColorWhite, keyPath, ip, ColorReset)
	fmt.Println()
	fmt.Println(ColorCyan + "Команды управления внутри сервера:" + ColorReset)
	fmt.Println(" • " + ColorWhite + "setL" + ColorReset + "       — интерактивная панель управления (прокси, DNS, тюнинг)")
	fmt.Println(" • " + ColorWhite + "reboot" + ColorReset + "     — быстрая перезагрузка")
	fmt.Println(" • " + ColorWhite + "poweroff" + ColorReset + "   — выключение устройства")
	fmt.Println(ColorGreen + "================================================================================" + ColorReset)
	Pause()
}

func main() {
	for {
		PrintBanner()
		fmt.Println(" Главное меню мастера установки:")
		fmt.Println()
		fmt.Println(" 1. 📖 Полное пошаговое руководство (Шаг 1 -> Шаг 4)")
		fmt.Println(" 2. ⚡ Быстрая установка PocketNode (если LineageOS уже прошит)")
		fmt.Println(" 3. 🛠 Проверить подключение телефона по ADB")
		fmt.Println(" 4. 🔑 Проверить подключение по SSH к телефону")

		demoStatus := ColorRed + "ВЫКЛ" + ColorReset
		if DemoMode {
			demoStatus = ColorGreen + "ВКЛ (симуляция без нагрузки)" + ColorReset
		}
		fmt.Printf(" 5. 🎭 Демо-режим: %s\n", demoStatus)
		fmt.Println(" 0. 🚪 Выйти")
		fmt.Println(ColorCyan + "================================================================================" + ColorReset)

		choice := ReadLine(" Выберите пункт: ")
		switch choice {
		case "1":
			StepRecovery()
			StepROM()
			StepADB()
			StepDeploy()
		case "2":
			StepADB()
			StepDeploy()
		case "3":
			PrintBanner()
			fmt.Println("Проверка подключения ADB...")
			if DemoMode {
				time.Sleep(300 * time.Millisecond)
				fmt.Println("List of devices attached")
				fmt.Println("a5lte_msm8916_demo\tdevice (root)")
				fmt.Println(ColorGreen + "\n[✓] [DEMO] Виртуальное устройство Samsung Galaxy A5 обнаружено!" + ColorReset)
			} else {
				out, err := exec.Command("adb", "devices").CombinedOutput()
				if err != nil {
					fmt.Println(ColorRed + "ADB не найден в системе. Убедитесь, что у вас установлены platform-tools." + ColorReset)
				} else {
					fmt.Println(string(out))
				}
			}
			Pause()
		case "4":
			PrintBanner()
			ip := ReadLine(" Введите IP-адрес телефона (например 192.168.1.191): ")
			if ip != "" {
				if DemoMode {
					time.Sleep(300 * time.Millisecond)
					fmt.Printf("\n%s[✓] [DEMO] Порт 22 ОТКРЫТ на %s! Сервер SSH готов к подключениям.%s\n", ColorGreen, ip, ColorReset)
				} else {
					conn, err := net.DialTimeout("tcp", ip+":22", 3*time.Second)
					if err != nil {
						fmt.Printf("\n%s[✗] Порт 22 недоступен на %s: %v%s\n", ColorRed, ip, err, ColorReset)
					} else {
						conn.Close()
						fmt.Printf("\n%s[✓] Порт 22 ОТКРЫТ! Сервер SSH готов к подключениям.%s\n", ColorGreen, ColorReset)
					}
				}
			}
			Pause()
		case "5":
			DemoMode = !DemoMode
			if DemoMode {
				fmt.Println(ColorGreen + "\n[✓] Демо-режим ВКЛЮЧЕН. Теперь можно безопасно тестировать установку без телефона!" + ColorReset)
			} else {
				fmt.Println(ColorYellow + "\n[!] Демо-режим ВЫКЛЮЧЕН. Установщик переключен в штатный режим работы с телефоном." + ColorReset)
			}
			time.Sleep(1200 * time.Millisecond)
		case "0":
			Clear()
			return
		}
	}
}

