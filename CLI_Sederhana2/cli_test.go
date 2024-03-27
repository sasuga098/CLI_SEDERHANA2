package main

import (
	"fmt"
	"testing"
	// "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"cli_sederhana2/config"
	"cli_sederhana2/model"
)

// func Init() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		fmt.Println("env not found, using global env")
// 	}
// 	config.OpenDB()
// }

func TestCreateBuku(t *testing.T) {
	Init()

	bukuData := model.Buku{
		ISBN:   "123456789",
		Penulis: "John Doe",
		Tahun:  2022,
		Judul:  "Sample Book",
		Gambar: "sample.jpg",
		Stok:   10,
	}

	err := bukuData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	fmt.Println(bukuData.ID)
}

func TestGetBukuByID(t *testing.T) {
	Init()

	bukuData := model.Buku{
		Model: model.Model{
			ID: 1,
		},
	}

	data, err := bukuData.GetByID(config.Mysql.DB)
	assert.Nil(t, err)

	fmt.Println(data)
}

func TestGetAllBuku(t *testing.T) {
	Init()

	bukuData := model.Buku{
		ISBN:   "805073663",
		Penulis: "Doug Stanton",
		Tahun:  2001,
		Judul:  "In Harm's Way: The Sinking of the U.S.S. Indianapolis and the Extraordinary Story of Its Survivors",
		Gambar: "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
		Stok:   24,
	}

	err := bukuData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	res, err := bukuData.GetAll(config.Mysql.DB)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(res), 1)

	fmt.Println(res)
}

func TestDeleteBukuByID(t *testing.T) {
	Init()

	bukuData := model.Buku{
		Model: model.Model{
			ID: 1,
		},
	}

err := bukuData.DeleteByID(config.Mysql.DB)
	assert.Nil(t, err)
}