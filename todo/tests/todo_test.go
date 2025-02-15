package tests

//Introduction to testing.  Note that testing is built into go and we will be using
//it extensively in this class. Below is a starter for your testing code.  In
//addition to what is built into go, we will be using a few third party packages
//that improve the testing experience.  The first is testify.  This package brings
//asserts to the table, that is much better than directly interacting with the
//testing.T object.  Second is gofakeit.  This package provides a significant number
//of helper functions to generate random data to make testing easier.

import (
	"fmt"
	"os"
	"testing"

	"drexel.edu/todo/db"
	fake "github.com/brianvoe/gofakeit/v6" //aliasing package name
	"github.com/stretchr/testify/assert"
)

// Note the default file path is relative to the test package location.  The
// project has a /tests path where you are at and a /data path where the
// database file sits.  So to get there we need to back up a directory and
// then go into the /data directory.  Thus this is why we are setting the
// default file name to "../data/todo.json"
const (
	DEFAULT_DB_FILE_NAME = "../data/todo.json"
)

var (
	DB *db.ToDo
)

// note init() is a helpful function in golang.  If it exists in a package
// such as we are doing here with the testing package, it will be called
// exactly once.  This is a great place to do setup work for your tests.
func init() {
	//Below we are setting up the gloabal DB variable that we can use in
	//all of our testing functions to make life easier
	testdb, err := db.New(DEFAULT_DB_FILE_NAME)
	if err != nil {
		fmt.Print("ERROR CREATING DB:", err)
		os.Exit(1)
	}

	DB = testdb //setup the global DB variable to support test cases

	//Now lets start with a fresh DB with the sample test data
	testdb.RestoreDB()
}

// Sample Test, will always pass, comparing the second parameter to true, which
// is hard coded as true
func TestTrue(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestAddHardCodedItem(t *testing.T) {

	item := db.ToDoItem{
		Id:     999,
		Title:  "This is a test case item",
		IsDone: false,
	}
	t.Log("Testing adding a hard coded item: ", item)

	//TODO: finish this test, add an item to the database and then
	//check that it was added correctly by looking it back up
	//use assert.NoError() to ensure errors are not returned.
	//explore other useful asserts in the testify package, see
	//https://github.com/stretchr/testify.  Specifically look
	//at things like assert.Equal() and assert.Condition()

	//I will get you started, uncomment the lines below to add to the DB
	//and ensure no errors:
	//---------------------------------------------------------------
	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	//TODO: Now finish the test case by looking up the item in the DB
	//and making sure it matches the item that you put in the DB above

	retrievedItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error retrieving item from database")
	assert.Equal(t, item, retrievedItem, "Retrieved item did not match hard coded item")
}

func TestAddRandomStructItem(t *testing.T) {
	//You can also use the Stuct() fake function to create a random struct
	//Not going to do anyting
	item := db.ToDoItem{}
	err := fake.Struct(&item)
	t.Log("Testing adding a randomly generated struct: ", item)

	assert.NoError(t, err, "Created fake item OK")

	//TODO: Complete the test

	err = DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	retrievedItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error retrieving item from database")
	assert.Equal(t, item, retrievedItem, "retrieved item did not match hard coded item")
}

func TestAddRandomItem(t *testing.T) {
	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	t.Log("Testing adding an item with random fields: ", item)

	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	retrievedItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error retrieving item from database")
	assert.Equal(t, item, retrievedItem, "Retrieved item did not match hard coded item")

}

//TODO: Create additional tests to showcase the correct operation of your program
//for example getting an item, getting all items, updating items, and so on. Be
//creative here.

func TestRestoreDB(t *testing.T) {

	//Should overwrite with a blank file
	file, err := os.OpenFile(DEFAULT_DB_FILE_NAME, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	assert.NoError(t, err, "Couldn't create a blank db file")
	defer file.Close()

	err = DB.RestoreDB()
	assert.NoError(t, err, "Error while restoring database")

	areFilesEqual, err := areFilesEqual(t, DEFAULT_DB_FILE_NAME, DEFAULT_DB_FILE_NAME+".bak")
	assert.NoError(t, err, "Error occurred comparing files")
	assert.Equal(t, true, areFilesEqual)
}

// component of TestRestoreDB.
// Not exported so it won't be confused as a standalone test
func areFilesEqual(t *testing.T, file1, file2 string) (bool, error) {
	contentFromFile1, err := os.Open(file1)
	assert.NoError(t, err, "Could not read from file1")

	contentFromFile2, err := os.Open(file2)
	assert.NoError(t, err, "Could not read from file2")

	bufferFile1 := make([]byte, 1024)
	bufferFile2 := make([]byte, 1024)

	for {
		bytesReadFromFile1, err := contentFromFile1.Read(bufferFile1)
		if err != nil && err.Error() != "EOF" {
			return false, err
		}

		bytesReadFromFile2, err := contentFromFile2.Read(bufferFile2)
		if err != nil && err.Error() != "EOF" {
			return false, err
		}

		if bytesReadFromFile1 != bytesReadFromFile2 {
			return false, nil

		}

		if bytesReadFromFile1 == 0 {
			break
		}

		if len(bufferFile1) != len(bufferFile2) {
			return false, nil
		}

		for i := range bufferFile1 {
			if bufferFile1[i] != bufferFile2[i] {
				return false, nil
			}
		}
	}
	return true, nil
}

func TestDeleteItem(t *testing.T) {
	item := db.ToDoItem{
		Id:     1002,
		Title:  "This item should not be inside of the final object",
		IsDone: false,
	}
	t.Log("Testing adding a hard coded item: ", item)

	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	retrievedItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error retrieving item from database")
	assert.Equal(t, item, retrievedItem, "Retrieved item did not match hard coded item")

	err = DB.DeleteItem(item.Id)
	assert.NoError(t, err, "Error deleting item from database")

	retrievedItem, err = DB.GetItem(item.Id)
	assert.Error(t, err, "There should have been an error retrieving data from the database")
	assert.Equal(t, db.ToDoItem{}, retrievedItem, "retrieved item was not empty")

}

func TestUpdateItem(t *testing.T) {
	item := db.ToDoItem{
		Id:     1003,
		Title:  "This is a test case item",
		IsDone: false,
	}
	t.Log("Testing adding a hard coded item: ", item)

	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	retrievedItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error retrieving item from database")
	assert.Equal(t, item, retrievedItem, "Retrieved item did not match hard coded item")

	updatedItem := db.ToDoItem{
		Id:     item.Id,
		Title:  item.Title,
		IsDone: !item.IsDone,
	}

	err = DB.UpdateItem(updatedItem)
	assert.NoError(t, err, "Error updating item in database")

	retrievedItem, err = DB.GetItem(item.Id)
	assert.NoError(t, err, "There was an error retrieving data from the database")
	assert.Equal(t, updatedItem, retrievedItem, "Retrieved item was not updated.")
}

func TestGetItem(t *testing.T) {

	item := db.ToDoItem{
		Id:     1058,
		Title:  "This is a test case item",
		IsDone: false,
	}

	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	retrievedItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error retrieving item from database")

	assert.Equal(t, item.Id, retrievedItem.Id, "Retrieved item Id did not match hard coded item")
	assert.Equal(t, item.Title, retrievedItem.Title, "Retrieved item Title did not match hard coded item")
	assert.Equal(t, item.IsDone, retrievedItem.IsDone, "Retrieved item IdDone status did not match hard coded item")
}

func TestGetAllItems(t *testing.T) {

	item := db.ToDoItem{
		Id:     1056,
		Title:  "This should be in the array",
		IsDone: false,
	}
	//	t.Log("Testing Adding a Hard Coded Item: ", item)

	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	retrievedItems, err := DB.GetAllItems()
	assert.NoError(t, err, "Error retrieving item from database")
	assert.GreaterOrEqual(t, len(retrievedItems), 1)

	counter := 0

	for _, value := range retrievedItems {
		if value.Id == item.Id {
			counter += 1
		}
	}

	assert.GreaterOrEqual(t, 1, counter, " The item we added was not in the list of retrieved items")
}

func TestChangeItemDoneStatus(t *testing.T) {
	item := db.ToDoItem{
		Id:     1006,
		Title:  "This is a test case item",
		IsDone: false,
	}
	t.Log("Testing adding a hard coded item: ", item)

	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to database")

	retrievedItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error retrieving item from database")
	assert.Equal(t, item, retrievedItem, "Retrieved item did not match hard coded item")

	err = DB.ChangeItemDoneStatus(item.Id, !item.IsDone)
	assert.NoError(t, err, "Error updating item in database")

	retrievedItem, err = DB.GetItem(item.Id)
	assert.NoError(t, err, "There was an error retrieving data from the database")
	assert.Equal(t, true, retrievedItem.IsDone, "Retrieved item was not updated.")
}
