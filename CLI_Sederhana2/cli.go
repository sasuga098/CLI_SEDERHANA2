package main

import (
	"fmt"
	"time"
	"bufio"
	"os"
    "errors"
	"strings"
	"strconv"
	"encoding/csv"
    "path/filepath"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
	"cli_sederhana2/config"
	"cli_sederhana2/model"
	"gorm.io/gorm"
	"github.com/joho/godotenv"

)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("env not found, using global env")
	}
	config.OpenDB()
}

func TambahBuku() {
    Init()
	inputanUser := bufio.NewReader(os.Stdin)
	var buku model.Buku

	fmt.Println("=================================")
	fmt.Println("Tambah Buku Baru")
	fmt.Println("=================================")

	fmt.Print("ISBN: ")
	isbn, _ := inputanUser.ReadString('\n')
	buku.ISBN = strings.TrimSpace(isbn)

	fmt.Print("Penulis: ")
	penulis, _ := inputanUser.ReadString('\n')
	buku.Penulis = strings.TrimSpace(penulis)

	fmt.Print("Tahun Terbit: ")
	fmt.Scanln(&buku.Tahun)

	fmt.Print("Judul: ")
	judul, _ := inputanUser.ReadString('\n')
	buku.Judul = strings.TrimSpace(judul)

	fmt.Print("Gambar: ")
	gambar, _ := inputanUser.ReadString('\n')
	buku.Gambar = strings.TrimSpace(gambar)

	fmt.Print("Stok: ")
	fmt.Scanln(&buku.Stok)

	// Simpan buku ke database
	err := buku.Create(config.Mysql.DB)
	if err != nil {
		fmt.Println("Error creating book in database:", err)
		return
	}

	fmt.Println("Buku berhasil ditambahkan ke dalam database!")
}
func LihatSemuaBuku() {
    Init()
	fmt.Println("=================================")
	fmt.Println("Daftar Buku")
	fmt.Println("=================================")

	var buku []model.Buku
	result := config.Mysql.DB.Find(&buku)
	if result.Error != nil {
		fmt.Println("Error fetching books from database:", result.Error)
		return
	}

	fmt.Printf("%-10s | %-30s | %-20s | %-10s | %-5s\n", "Kode", "Judul", "Penulis", "Tahun", "Stok")
	fmt.Println("---------------------------------------------------------------")
	for _, b := range buku {
		fmt.Printf("%-10s | %-30s | %-20s | %-10d | %-5d\n", b.ISBN, b.Judul, b.Penulis, b.Tahun, b.Stok)
	}
	fmt.Println("---------------------------------------------------------------")
}

func HapusBuku() {
    Init()
	fmt.Println("=================================")
	fmt.Println("Hapus Buku")
	fmt.Println("=================================")

	LihatSemuaBuku()

	fmt.Print("Masukkan Kode Buku yang akan dihapus: ")
	inputanUser := bufio.NewReader(os.Stdin)
	kodeBuku, _ := inputanUser.ReadString('\n')
	kodeBuku = strings.TrimSpace(kodeBuku)

	var buku model.Buku
	result := config.Mysql.DB.Where("kode_buku = ?", kodeBuku).First(&buku)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Buku dengan Kode Buku tersebut tidak ditemukan.")
		return
	}

	// Soft-delete data buku dari database
	err := config.Mysql.DB.Delete(&buku).Error
	if err != nil {
		fmt.Println("Error deleting book in database:", err)
		return
	}

	fmt.Println("Buku berhasil dihapus dari database secara soft-delete!")
}

func EditBuku() {
    Init()
	fmt.Println("=================================")
	fmt.Println("Edit Buku")
	fmt.Println("=================================")

	LihatSemuaBuku()

	fmt.Print("Masukkan Kode Buku yang akan diubah: ")
	inputanUser := bufio.NewReader(os.Stdin)
	kodeBuku, _ := inputanUser.ReadString('\n')
	kodeBuku = strings.TrimSpace(kodeBuku)

	var buku model.Buku
	result := config.Mysql.DB.Where("kode_buku = ?", kodeBuku).First(&buku)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Buku dengan Kode Buku tersebut tidak ditemukan.")
		return
	}

	fmt.Println("Masukkan informasi baru:")
	fmt.Print("Judul Buku: ")
	judul, _ := inputanUser.ReadString('\n')
	buku.Judul = strings.TrimSpace(judul)

	fmt.Print("Penulis: ")
	penulis, _ := inputanUser.ReadString('\n')
	buku.Penulis = strings.TrimSpace(penulis)

	fmt.Print("Tahun Terbit: ")
	fmt.Scanln(&buku.Tahun)

	fmt.Print("Gambar: ")
	gambar, _ := inputanUser.ReadString('\n')
	buku.Gambar = strings.TrimSpace(gambar)

	fmt.Print("Stok: ")
	fmt.Scanln(&buku.Stok)

	// Update data buku di database
	err := config.Mysql.DB.Save(&buku).Error
	if err != nil {
		fmt.Println("Error updating book in database:", err)
		return
	}

	fmt.Println("Buku berhasil diubah di dalam database!")
}

func MuatDataBukuDariFile() {
    Init()
	files, err := os.ReadDir("books")
	if err != nil {
		fmt.Println("Error reading books directory:", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join("books", file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

			var buku model.Buku
			err = json.Unmarshal(data, &buku)
			if err != nil {
				fmt.Println("Error unmarshalling JSON data:", err)
				continue
			}

			// Simpan buku ke database
			err = buku.Create(config.Mysql.DB)
			if err != nil {
				fmt.Println("Error creating book in database:", err)
				continue
			}

			fmt.Println("Buku berhasil dimuat dari file:", file.Name())
		}
	}
}

func ImporDataDariCSV() {
    Init()
	fmt.Println("=================================")
	fmt.Println("Impor Data dari File CSV (Upsert)")
	fmt.Println("=================================")

	fmt.Print("Masukkan path/lokasi file CSV: ")
	inputanUser := bufio.NewReader(os.Stdin)
	pathCSV, _ := inputanUser.ReadString('\n')
	pathCSV = strings.TrimSpace(pathCSV)

	file, err := os.Open(pathCSV)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV records:", err)
		return
	}

	for _, row := range records {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			fmt.Println("Error converting id to int:", err)
			continue
		}

		tahun, err := strconv.Atoi(row[3])
		if err != nil {
			fmt.Println("Error converting tahun to int:", err)
			continue
		}

		stok, err := strconv.Atoi(row[6])
		if err != nil {
			fmt.Println("Error converting stok to int:", err)
			continue
		}

		buku := model.Buku{
			Model: model.Model{ID: uint(id)}, // Menginisialisasi ID dari Model
			ISBN:    row[1],
			Penulis: row[2],
			Tahun:   uint(tahun),
			Judul:   row[4],
			Gambar:  row[5],
			Stok:    uint(stok),
		}

		// Cek apakah buku sudah ada berdasarkan ID
		var existingBook model.Buku
		result := config.Mysql.DB.Where("id = ?", buku.ID).First(&existingBook)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Jika belum ada, tambahkan buku baru
			err := buku.Create(config.Mysql.DB)
			if err != nil {
				fmt.Println("Error importing book to database:", err)
				return
			}
		} else {
			// Jika sudah ada, update data buku
			existingBook.ISBN = buku.ISBN
			existingBook.Penulis = buku.Penulis
			existingBook.Tahun = buku.Tahun
			existingBook.Judul = buku.Judul
			existingBook.Gambar = buku.Gambar
			existingBook.Stok = buku.Stok

			err := existingBook.UpdateByID(config.Mysql.DB)
			if err != nil {
				fmt.Println("Error updating book in database:", err)
				return
			}
		}
	}

	fmt.Println("Data berhasil diimpor dari file CSV ke database dengan mekanisme upsert!")
}

func GeneratePdfBuku(buku model.Buku) {
    Init()
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)
	pdf.SetLeftMargin(10)
	pdf.SetRightMargin(10)

	bukuText := fmt.Sprintf(
		"Kode Buku: %s\nJudul: %s\nPenulis: %s\nTahun Terbit: %d",
		buku.ISBN, buku.Judul, buku.Penulis, buku.Tahun)

	pdf.MultiCell(0, 10, bukuText, "0", "L", false)
	pdf.Ln(5)

	waktuCetak := time.Now().Format("2006-01-02 15:04:05")
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(0, 10, "Waktu Cetak: "+waktuCetak, "", 0, "R", false, 0, "")

	err := pdf.OutputFileAndClose(fmt.Sprintf("pdf/%s.pdf", buku.ISBN))
	if err != nil {
		fmt.Println("Error creating PDF file:", err)
	} else {
		fmt.Println("Buku berhasil dicetak ke file PDF.")
	}
}

func main() {
    Init()
	config.OpenDB() // Buka koneksi database

	var pilihanMenu int

	for {
		fmt.Println("=================================")
		fmt.Println("Aplikasi Manajemen Daftar Buku Perpustakaan")
		fmt.Println("=================================")
		fmt.Println("Silahkan Pilih : ")
		fmt.Println("1. Tambah Buku")
		fmt.Println("2. Lihat Semua Buku")
		fmt.Println("3. Hapus Buku")
		fmt.Println("4. Edit Buku")
		fmt.Println("5. Impor Data dari File CSV (Upsert)")
		fmt.Println("6. Generate PDF Buku")
		fmt.Println("7. Keluar")
		fmt.Println("=================================")
		fmt.Print("Masukkan Pilihan : ")
		fmt.Scanln(&pilihanMenu)

		switch pilihanMenu {
		case 1:
			TambahBuku()
		case 2:
			LihatSemuaBuku()
		case 3:
			HapusBuku()
		case 4:
			EditBuku()
		case 5:
			ImporDataDariCSV()
		case 6:
			// Meminta input kode buku
			fmt.Print("Masukkan Kode Buku untuk dicetak ke PDF: ")
			inputanUser := bufio.NewReader(os.Stdin)
			kodeBuku, _ := inputanUser.ReadString('\n')
			kodeBuku = strings.TrimSpace(kodeBuku)

			// Mengambil informasi buku dari database berdasarkan kode buku
			var buku model.Buku
			result := config.Mysql.DB.Where("isbn = ?", kodeBuku).First(&buku)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				fmt.Println("Buku dengan Kode Buku tersebut tidak ditemukan.")
			} else {
				// Memanggil fungsi GeneratePdfBuku untuk membuat file PDF dari informasi buku
				GeneratePdfBuku(buku)
			}
		case 7:
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

