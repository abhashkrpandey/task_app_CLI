package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Items struct {
	Id          int
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func readFromFile(filePath string, structData *[]Items) {
	file, fileOpenErr := os.OpenFile(filePath, os.O_RDWR, 0750)
	if fileOpenErr != nil {
		fmt.Println(fileOpenErr)
		os.Exit(1)
	} else {
		fileInfo, fileStatErr := file.Stat()
		if fileStatErr != nil {
			fmt.Println(fileStatErr)
			os.Exit(1)
		} else {
			dataArray := make([]byte, fileInfo.Size())
			file.Seek(0, 0)
			_, fileReadErr := file.Read(dataArray)
			if fileReadErr != nil && fileReadErr != io.EOF {
				fmt.Println(fileReadErr)
			} else {
				json.Unmarshal(dataArray, &structData)
				file.Close()
			}
		}
	}
	return
}
func writeInFile(filePath string, structData []Items) {
	file, fileOpenErr := os.OpenFile(filePath, os.O_RDWR, 0750)
	if fileOpenErr != nil {
		fmt.Println(fileOpenErr)
		os.Exit(1)
	} else {
		marshedData, dataConversionErr := json.Marshal(structData)
		if dataConversionErr != nil {
			fmt.Println(dataConversionErr)
			os.Exit(1)
		} else {
			file.Truncate(0)
			file.Seek(0, 0)
			_, fileWriteErr := file.Write(marshedData)
			if fileWriteErr != nil {
				fmt.Println(fileWriteErr)
				os.Exit(1)
			}
			return
		}
	}
}
func addInList(filePath string, description string) {
	structData := make([]Items, 0, 5)
	readFromFile(filePath, &structData)
	element := Items{
		Id:          len(structData) + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	structData = append(structData, element)
	writeInFile(filePath, structData)
	fmt.Println("item added in list")
	os.Exit(0)
}
func updateById(id int, filePath string, description string) {
	var itemFound bool = false
	structData := make([]Items, 0, 5)
	readFromFile(filePath, &structData)
	if len(structData) == 0 {
		fmt.Println("First enter the item in list")
		os.Exit(1)
	} else {
		for i := range structData {
			if structData[i].Id == id {
				structData[i].Description = description
				structData[i].UpdatedAt = time.Now()
				itemFound = true
			}
		}
		if itemFound == false {
			fmt.Println("id not found")
			os.Exit(1)
		} else {
			writeInFile(filePath, structData)
		}
	}
}
func deleteItemById(id int, filePath string) {
	itemFound := false
	structData := make([]Items, 0, 5)
	readFromFile(filePath, &structData)
	for i := range structData {
		if structData[i].Id == id {
			structData = append(structData[0:i], structData[i+1:]...)
			itemFound = true
		}
	}
	if itemFound == false {
		fmt.Println("id not found")
		os.Exit(1)
	} else {
		writeInFile(filePath, structData)
		fmt.Println("item deleted")
		os.Exit(0)
	}
}
func mark(filePath string, itemID int, indicator int) {
	itemFound := false
	structData := make([]Items, 0, 5)
	readFromFile(filePath, &structData)
	for i := range structData {
		if structData[i].Id == itemID {
			if indicator == 1 {
				structData[i].Status = "in progress"
			} else if indicator == 2 {
				structData[i].Status = "done"
			}
			structData[i].UpdatedAt = time.Now()
			itemFound = true
		}
	}
	if itemFound == false {
		fmt.Println("id  not found")
		os.Exit(1)
	} else {
		writeInFile(filePath, structData)
		fmt.Println("Status updated")
		os.Exit(0)
	}
}
func list(filePath string, indicator int) {
	structData := make([]Items, 0, 5)
	readFromFile(filePath, &structData)
	elementsList := make([]Items, 0, 5)
	for i := range structData {
		if indicator == 0 && structData[i].Status == "todo" {
			elementsList = append(elementsList, structData[i])
		} else if indicator == 1 && structData[i].Status == "in progress" {
			elementsList = append(elementsList, structData[i])
		} else if indicator == 2 && structData[i].Status == "done" {
			elementsList = append(elementsList, structData[i])
		}
	}
	for i := range elementsList {
		fmt.Printf("item id:%v\n", elementsList[i].Id)
		fmt.Printf("item description:%v\n", elementsList[i].Description)
		fmt.Printf("item status:%v\n", elementsList[i].Status)
		fmt.Printf("item created at:%v\n", elementsList[i].CreatedAt)
		fmt.Printf("item updated at:%v\n", elementsList[i].UpdatedAt)
	}
	os.Exit(0)
}
func main() {

	//task-app createFolder -name=folder1

	argumentsWithoutProgram := os.Args[1:]
	if len(argumentsWithoutProgram) < 1 {
		fmt.Println("Incomplete command")
		os.Exit(1)
	}
	switch argumentsWithoutProgram[0] {
	case "start":
		fmt.Println("----- 👋 Welcome to task_app CLI: Manage Your Daily Tasks with Ease -----")

		fmt.Println("📘 Quick Command Guide:")

		fmt.Println("📁 Create a folder in the current location:")
		fmt.Println("   task_app createFolder <folder_name>")

		fmt.Println("📝 Create a new task list inside a folder:")
		fmt.Println("   task_app addList <folder_name> <list_name.json>")

		fmt.Println("➕ Add a task to any list:")
		fmt.Println("   task_app addInList <folder_name>/<list_name.json> \"task description\"")

		fmt.Println("✏️ Update a task description by ID:")
		fmt.Println("   task_app updateById <folder_name>/<list_name.json> <task_id> \"new description\"")

		fmt.Println("❌ Delete a task by ID:")
		fmt.Println("   task_app deleteItemById <folder_name>/<list_name.json> <task_id>")

		fmt.Println("🗑️ Delete a list:")
		fmt.Println("   task_app deleteList <folder_name>/<list_name.json>")

		fmt.Println("🚧 Mark a task as In Progress:")
		fmt.Println("   task_app markAsInProgress <folder_name>/<list_name.json> <task_id>")

		fmt.Println("✅ Mark a task as Done:")
		fmt.Println("   task_app markAsDone <folder_name>/<list_name.json> <task_id>")

		fmt.Println("📋 View tasks categorized as To Do:")
		fmt.Println("   task_app listByTodo <folder_name>/<list_name.json>")

		fmt.Println("📋 View tasks categorized as In Progress:")
		fmt.Println("   task_app listByInProgess <folder_name>/<list_name.json>")

		fmt.Println("📋 View tasks categorized as Done:")
		fmt.Println("   task_app listByDone <folder_name>/<list_name.json>")

	case "createFolder": // createFolder dir
		if len(argumentsWithoutProgram) < 2 {
			fmt.Println("correct command -> task_app createFolder dirname")
			os.Exit(1)
		} else {
			dirName := os.Args[2]
			err := os.Mkdir(dirName, 0750)
			if err != nil {
				fmt.Println("No Directory created")
				os.Exit(1)
			} else {
				fmt.Println("folder:" + dirName + " created")
			}
		}
	case "addList": // addList dir filename
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app addList dirname filename")
			os.Exit(1)
		} else {
			filePath := os.Args[2] + "/" + os.Args[3] + ".json"
			_, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0754)
			if err != nil {
				fmt.Println("Error occured while creating a list")
				os.Exit(1)
			} else {
				fmt.Println("List added to the folder")
			}
		}
	case "addInList": // addInList dirname/filename.json "description"
		if len(argumentsWithoutProgram) < 4 {
			fmt.Println("correct command ->task_app addInList dirname/filename.json 'description'")
			os.Exit(1)
		} else {
			description := os.Args[3]
			filePath := os.Args[2]
			addInList(filePath, description)
		}
	case "updateById": // updateById dir/filename id "change_description"
		if len(argumentsWithoutProgram) < 5 {
			fmt.Println("correct command ->task_app updateById dirname/filename.json itemId 'change_description' ")
			os.Exit(1)
		} else {
			itemId, _ := strconv.Atoi(os.Args[3])
			filePath := os.Args[2]
			newDescription := os.Args[4]
			updateById(itemId, filePath, newDescription)
		}
	case "deleteItemById": // deleteItemById dir/filename id
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app deleteItemById dirname/filename.json itemId ")
			os.Exit(1)
		} else {
			itemId, _ := strconv.Atoi(os.Args[3])
			fileName := os.Args[2]
			deleteItemById(itemId, fileName)
		}
	case "deleteList": // deleteList dir/filename
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app deleteList dirname/filename.json")
			os.Exit(1)
		} else {
			filePath := os.Args[2]
			fileDeletionErr := os.Remove(filePath)
			if fileDeletionErr == nil {
				fmt.Println("file Deleted")
				os.Exit(0)
			} else {
				fmt.Println(fileDeletionErr)
				os.Exit(1)
			}
		}
	case "markAsInProgressById": // markAsInProgress dir/filename id
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app markAsInProgress dirname/filename.json id")
			os.Exit(1)
		} else {
			filepath := os.Args[2]
			itemId, _ := strconv.Atoi(os.Args[3])
			mark(filepath, itemId, 1)
		}
	case "markAsDoneById": // markAsDone dir/filename id
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app markAsDone dirname/filename.json id")
			os.Exit(1)
		} else {
			filepath := os.Args[2]
			itemId, _ := strconv.Atoi(os.Args[3])
			mark(filepath, itemId, 2)
		}
	case "listByTodo": // listByTodo dir/filename
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app listByTodo dirname/filename.json ")
			os.Exit(1)
		} else {
			filePath := os.Args[2]
			list(filePath, 0)
		}
	case "listByDone":
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app listByDone dirname/filename.json ")
			os.Exit(1)
		} else {
			filePath := os.Args[2]
			list(filePath, 1)
		}
	case "listByInProgress":
		if len(argumentsWithoutProgram) < 3 {
			fmt.Println("correct command ->task_app listByInProgess dirname/filename.json ")
			os.Exit(1)
		} else {
			filePath := os.Args[2]
			list(filePath, 2)
		}
	default:
		fmt.Println("Enter proper command")
	}
}
