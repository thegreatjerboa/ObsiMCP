package api

import (
	"fmt"
	"obsimcp/src/config"
	"os"
	"testing"
)

var api LocalRestApi

func TestMain(m *testing.M) {
	// 初始化配置
	fmt.Println("Initializing environment...")
	config.InitConfig()
	// 初始化 LocalRestApi 实例
	api = NewLocalRestApi()
	// 运行测试
	code := m.Run()

	// 清理环境
	fmt.Println("Cleaning up environment...")

	// 退出测试
	os.Exit(code)
}

func TestLocalRestApi_GetVaultFile(t *testing.T) {
	filename := "Inbox/test.md"
	details, err := api.GetVaultFile(filename)
	if err != nil {
		t.Errorf("Error getting vault file: %v", err)
		return
	}
	fmt.Printf("Vault file details: %+v\n", details)
}

func TestLocalRestApi_DeleteVaultFile(t *testing.T) {
	filename := "Inbox/xsdsada.md"
	err := api.DeleteVaultFile(filename)
	if err != nil {
		t.Errorf("Error deleting vault file: %v", err)
		return
	}
	fmt.Println("Vault file deleted successfully.")
}

func TestLocalRestApi_AppendVaultFile(t *testing.T) {
	filename := "Inbox/test.md"
	content := "## test\n\n- [ ] test"
	err := api.AppendVaultFile(filename, content)
	if err != nil {
		t.Errorf("Error appending vault file: %v", err)
		return
	}
	fmt.Println("Vault file appended successfully.")
}

func TestLocalRestApi_CreateOrUpdateVaultFile(t *testing.T) {
	filename := "Inbox/test.md"
	content := "## test\n\n- [ ] test1212"
	err := api.CreateOrUpdateVaultFile(filename, content)
	if err != nil {
		t.Errorf("Error creating or updating vault file: %v", err)
		return
	}
	fmt.Println("Vault file created or updated successfully.")
}

func TestLocalRestApi_GetDirectory(t *testing.T) {
	path := "Inbox"
	directory, err := api.ListDirectory(path)
	if err != nil {
		t.Errorf("Error getting directory: %v", err)
		return
	}
	fmt.Printf("Directory structure: %+v\n", directory)
}

func TestLocalRestApi_SimpleSearch(t *testing.T) {
	query := "test"
	contextLength := 100
	results, err := api.SimpleSearch(query, contextLength)
	if err != nil {
		t.Errorf("Error performing simple search: %v", err)
		return
	}
	for i, result := range results {
		fmt.Printf("Result %d: Filename: %s\n", i+1, result.Filename)
		for j, match := range result.Matches {
			fmt.Printf("Match %d: Context: %s\n", j+1, match.Context)
		}
	}
}
