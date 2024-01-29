/*
Copyright © 2024 NAME HERE md3852@drexel.edu
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"drexel.edu/todo/db"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Global variables to hold the command line flags to drive the todo CLI
// application
var (
	dbFileNameFlag string
	restoreDbFlag  bool
	listFlag       bool
	itemStatusFlag bool
	queryFlag      int
	addFlag        string
	updateFlag     string
	deleteFlag     int
)

type AppOptType int

// To make the code a little more clean we will use the following
// constants as basically enumerations for different options.  This
// allows us to use a switch statement in main to process the command line
// flags effectively
const (
	LIST_DB_ITEM AppOptType = iota
	RESTORE_DB_ITEM
	QUERY_DB_ITEM
	ADD_DB_ITEM
	UPDATE_DB_ITEM
	DELETE_DB_ITEM
	CHANGE_ITEM_STATUS
	NOT_IMPLEMENTED
	INVALID_APP_OPT
)

// processCmdLineFlags parses the command line flags for our CLI
//
// TODO: This function uses the flag package to parse the command line
//		 flags.  The flag package is not very flexible and can lead to
//		 some confusing code.

//			 REQUIRED:     Study the code below, and make sure you understand
//						   how it works.  Go online and readup on how the
//						   flag package works.  Then, write a nice comment
//				  		   block to document this function that highights that
//						   you understand how it works.
//
//			 EXTRA CREDIT: The best CLI and command line processor for
//						   go is called Cobra.  Refactor this function to
//						   use it.  See github.com/spf13/cobra for information
//						   on how to use it.
//
//	 YOUR ANSWER:
/*
The original processCmdLineFlags function used the flag package to parse
command-line arguments. The flag package has built-in functions that allow the
developer to define the type of argument expected (i.e., strings, boolean values,
and more). After the flags are specified, the flags.parse() function
populates the associated global variables with the command-line flags.
(i.e., -db= dbFileNameFlag, -restore=restoreDbFlag, -l=listFlag, and more). Once
the values are parsed, they are accessible via the global variables. The
processCmdLineFlags function also accesses the flag package's Visit function to
loop over the flags to determine which operation to execute.I have modified
processCmdLineFlags to use Cobra.

Cobra uses pFlags to achieve simialar behavior to flags. Cobra handles parsing the
command-line arguments in the init function. It uses the root CMD Object, the
persistantFlags() function, and various parsers to define the flag in the command
line and assign the value provided into a package global variable. The global
variables are already set by the time we start processCmdLineFlags function. the
processCmdLineFlags function checks how many arguments were provided. If there are
none, it shows a help menu. Lastly, the processCmdLineFlags function uses the visit
function to loop through every flag Cobra identified. Each flag is submitted to a
switch statement which assigns an operation to the variable that the function will
return (i.e., appOpt). If no operation is assigned, then a help menu is printed.

One not of interest is the case for multiple flags. Most flags are meant to
be used by themselves; however, the -s (i.e., status) flag needs the -q (query)
flag. The order of operations is important because we can only return one operation,
and the last operation overrides previous ones. Flags are processed in
lexicographical order. This happens to work out because status is after query.
Suppose we added other flags in the future that require multiple arguments. In that
case, we may want to refactor the operations into separate sub-commands instead of doing
it in the root. Otherwise, we need to implement a mechanism to ensure that combinations
of flags invoke the correct operation.

Once the processCmdLineFlags function returns the operation, the Run function will
submit the operation to a switch statement, which will execute the operation.
*/
func processCmdLineFlags(cmd *cobra.Command) (AppOptType, error) {

	// Note: We used Cobra to attempt the extra credit assignment.

	// The flags are assigned in the init() function instead of this
	// processCmdLineFlags function. We also use pFlags instaed of flags
	// pFlags is a third party flag package, while flags is a standard
	// library. pFlags builds on flags by extending its functionality by
	// adding shorthand flags, POSIX/GNU style flags, subcommands,
	// default values for flags, and more advanced flag parsing options.

	var appOpt AppOptType = INVALID_APP_OPT

	//show help if no flags are set
	if len(os.Args) == 1 {
		cmd.Help()
		return appOpt, errors.New("no flags were set")
	}

	// Loop over the flags and check which ones are set, set appOpt
	// accordingly
	cmd.Flags().Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "list":
			appOpt = LIST_DB_ITEM
		case "restore":
			appOpt = RESTORE_DB_ITEM
		case "query":
			appOpt = QUERY_DB_ITEM
		case "add":
			appOpt = ADD_DB_ITEM
		case "update":
			appOpt = UPDATE_DB_ITEM
		case "delete":
			appOpt = DELETE_DB_ITEM

		//TODO: EXTRA CREDIT - Implment the -s flag that changes the
		//done status of an item in the database.  For example -s=true
		//will set the done status for a particular item to true, and
		//-s=false will set the done states for a particular item to
		//false.
		//
		//HINT FOR EXTRA CREDIT
		//Note the -s option also requires an id for the item to that
		//you want to change.  I recommend you use the -q option to
		//specify the item id.  Therefore, the -s option is only valid
		//if the -q option is also set
		case "status":
			//For extra credit you will need to change some things here
			//and also in main under the CHANGE_ITEM_STATUS case
			appOpt = CHANGE_ITEM_STATUS
		default:
			appOpt = INVALID_APP_OPT
		}
	})

	if appOpt == INVALID_APP_OPT || appOpt == NOT_IMPLEMENTED {
		fmt.Println("Invalid option set or the desired option is not currently implemented")
		cmd.Help()
		return appOpt, errors.New("no flags or unimplemented were set")
	}

	return appOpt, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A CLI tool to manage to do items",
	Long: `This is a go language CLI tool to manage a list of todo items.
	
	This application will be driven by a simple text file based database. 
	See the todo.json file. Notice that this file is structured as a JSON array, 
	with collections of individual JSON objects. Each object contains an id, 
	description and a done flag. For example:
	
	[
	  {
		"id": 1,
		"title": "Learn Go / GoLang",
		"done": false
	  },
	  {
		"id": 2,
		"title": "Learn Kubernetes",
		"done": false
	  }
	]

	By default our program uses ./data/todo.json as the default database. 
	You can override the database name from the command line via the -db flag 
	providing a new database name. For example -db ./data/my_new_database.db.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		//Process the command line flags
		opts, err := processCmdLineFlags(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//Create a new db object
		todo, err := db.New(dbFileNameFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//Switch over the command line flags and call the appropriate
		//function in the db package
		switch opts {
		case RESTORE_DB_ITEM:
			fmt.Println("Running RESTORE_DB_ITEM...")
			if err := todo.RestoreDB(); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Database restored from backup file")
		case LIST_DB_ITEM:
			fmt.Println("Running QUERY_DB_ITEM...")
			todoList, err := todo.GetAllItems()
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			for _, item := range todoList {
				todo.PrintItem(item)
			}
			fmt.Println("THERE ARE", len(todoList), "ITEMS IN THE DB")
			fmt.Println("Ok")

		case QUERY_DB_ITEM:
			fmt.Println("Running QUERY_DB_ITEM...")
			item, err := todo.GetItem(queryFlag)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			todo.PrintItem(item)
			fmt.Println("Ok")
		case ADD_DB_ITEM:
			fmt.Println("Running ADD_DB_ITEM...")
			item, err := todo.JsonToItem(addFlag)
			if err != nil {
				fmt.Println("Add option requires a valid JSON todo item string")
				fmt.Println("Error: ", err)
				break
			}
			if err := todo.AddItem(item); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case UPDATE_DB_ITEM:
			fmt.Println("Running UPDATE_DB_ITEM...")
			item, err := todo.JsonToItem(updateFlag)
			if err != nil {
				fmt.Println("Update option requires a valid JSON todo item string")
				fmt.Println("Error: ", err)
				break
			}
			if err := todo.UpdateItem(item); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case DELETE_DB_ITEM:
			fmt.Println("Running DELETE_DB_ITEM...")
			err := todo.DeleteItem(deleteFlag)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case CHANGE_ITEM_STATUS:
			//For the CHANGE_ITEM_STATUS extra credit you will also
			//need to add some code here
			fmt.Println("Running CHANGE_ITEM_STATUS...")
			if queryFlag == 0 {
				fmt.Println("Warning, a zero value for id was detected. this can happen if you forget to define the -q flag. The usage for CHANGE_ITEM_STATUS is -q <targetId:int> -s=<isDone:bool>. if your intent was to target an item with an id of 0, please ignore this message")
			}
			err := todo.ChangeItemDoneStatus(queryFlag, itemStatusFlag)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		default:
			fmt.Println("INVALID_APP_OPT")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&addFlag, "add", "a", "", "Adds an item to the database")
	//Note two letter shorthand flags are not allowed. we have to use two  dashes with db. i.e., --db
	rootCmd.PersistentFlags().StringVar(&dbFileNameFlag, "db", "./data/todo.json", "Name of the target database file (default: \"./data/todo.json\")")
	rootCmd.PersistentFlags().IntVarP(&deleteFlag, "delete", "d", 0, "Deletes an item from the database")
	rootCmd.PersistentFlags().BoolVarP(&listFlag, "list", "l", false, "List all the items in the database")
	rootCmd.PersistentFlags().IntVarP(&queryFlag, "query", "q", 0, "Query an item in the database")
	//Note two letter shorthand flags are not allowed. we have to use two  dashes with restore. i.e., --restore or -r for short.
	rootCmd.PersistentFlags().BoolVarP(&restoreDbFlag, "restore", "r", false, "Restore the database from the backup file")
	rootCmd.PersistentFlags().BoolVarP(&itemStatusFlag, "status", "s", false, "Change item 'done' status to true or false. must be used alongside -q flag. ex: -q<target:int> -s=<isDone:bool>")
	rootCmd.PersistentFlags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&updateFlag, "update", "u", "", "Updates an item in the database")
}
